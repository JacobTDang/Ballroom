// Package tui is the full-screen interactive homepage: a boot screen that
// runs preflight checks, then a picker you navigate with arrow keys.
package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/preflight"
)

const (
	ollamaHost = "http://localhost:11434"
	tutorModel = "qwen2.5-coder:7b"
)

// maxStepLogLines caps how many of a single docker-build step's own
// lines stay on screen — the step entry itself persists for the whole
// build (so progress doesn't disappear), only its own scrollback within
// that step is windowed.
const maxStepLogLines = 3

// maxOutputLines caps how many lines of a resolved check's real command
// output are shown — same windowing idea as maxStepLogLines, applied to
// preflight.Check.Output instead of docker build's step lines.
const maxOutputLines = 3

// recentCheckWindow is how many of the most recently invoked checks stay
// expanded (showing their real command + output); older ones collapse to
// a single summary line so the screen doesn't get more cluttered as more
// checks finish.
const recentCheckWindow = 3

// checkStartDelay paces checks so they visibly run one at a time instead
// of all resolving within the same frame — a var (not a const) so tests
// can zero it out instead of actually sleeping.
var checkStartDelay = 250 * time.Millisecond

// imageCheckIndex is the pending-checks slot for the image check — the
// one that gets special "build it live" treatment instead of just
// reporting a plain pass/fail. See newBootModel for the fixed check order.
const imageCheckIndex = 1

var (
	checkOKStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#2FA6A6")).Bold(true)
	checkFailStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F03C3C")).Bold(true)
	checkDimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#D9D3C4"))
	hintStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#E8A93C")).Bold(true)
	buildLogStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#D9D3C4")).Faint(true)
)

type checkDoneMsg preflight.Check
type buildLineMsg string
type buildDoneMsg struct{ err error }

// startCheckMsg fires after checkStartDelay to kick off the check at
// index — the delay between a check resolving and the next one starting
// is what makes them visibly run one at a time.
type startCheckMsg struct{ index int }

// delayedCheck schedules startCheckMsg for the given check index after
// checkStartDelay.
func delayedCheck(index int) tea.Cmd {
	return tea.Tick(checkStartDelay, func(time.Time) tea.Msg {
		return startCheckMsg{index: index}
	})
}

// lastLines splits raw command/response output into non-empty lines and
// returns at most the last n of them — the same "windowed scrollback"
// idea as a build step's lines, applied to a resolved check's real
// output instead.
func lastLines(output string, n int) []string {
	var lines []string
	for _, l := range strings.Split(strings.TrimSpace(output), "\n") {
		l = strings.TrimRight(l, "\r")
		if strings.TrimSpace(l) == "" {
			continue
		}
		lines = append(lines, l)
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return lines
}

// recentWindowStart returns the index of the first check that should
// stay expanded — everything before it collapses to a summary line.
func recentWindowStart(total, window int) int {
	start := total - window
	if start < 0 {
		return 0
	}
	return start
}

// lastBuildStepOutput carries the last docker-build step's lines into
// the resolved Image check's Output, so it doesn't lose all its visible
// output the moment the live build panel collapses into a plain check.
func lastBuildStepOutput(steps []buildStepLog) string {
	if len(steps) == 0 {
		return ""
	}
	return strings.Join(steps[len(steps)-1].lines, "\n")
}

// buildStepLog groups the docker-build output lines belonging to one
// step (identified by its leading "#NN" token) — the step itself stays
// on screen for the whole build, but only its most recent
// maxStepLogLines lines are kept.
type buildStepLog struct {
	id    string
	lines []string
}

// stepID extracts the leading "#NN" token docker buildkit's plain
// progress output prefixes every line with (e.g. "#5 [2/4] RUN ..." ->
// "#5"). Returns "" if the line doesn't start with one.
func stepID(line string) string {
	if len(line) < 2 || line[0] != '#' {
		return ""
	}
	i := 1
	for i < len(line) && line[i] >= '0' && line[i] <= '9' {
		i++
	}
	if i == 1 {
		return ""
	}
	return line[:i]
}

// buildImageFn is a var (not a direct call) so tests can substitute a
// fake build stream instead of shelling out to docker for real.
var buildImageFn = orchestrator.BuildImage

// bootModel runs preflight checks one at a time (so they visibly appear
// in sequence, like a real boot log) and then waits for the user to
// continue. If the practice image isn't built, it builds it right here —
// expanding into a live panel of `docker build` output — and "ready"
// doesn't become true until that finishes, instead of silently deferring
// the build to whenever you first launch something.
type bootModel struct {
	cfg     config.Config
	pending []func() preflight.Check
	checks  []preflight.Check
	ready   bool
	quit    bool
	phase   int

	building    bool
	buildSteps  []buildStepLog
	buildLineCh <-chan string
	buildErrCh  <-chan error

	width, height int

	// checkNames mirrors pending, one entry per slot, so a check's name
	// can be shown while it's still queued — before it has actually run
	// and produced a Check with a real Command/Output to display.
	checkNames []string
}

func newBootModel(cfg config.Config) bootModel {
	image := cfg.DockerImage
	return bootModel{
		cfg: cfg,
		pending: []func() preflight.Check{
			preflight.CheckDocker,
			func() preflight.Check { return preflight.CheckImage(image) },
			func() preflight.Check { return preflight.CheckOllama(ollamaHost) },
			func() preflight.Check { return preflight.CheckModel(ollamaHost, tutorModel) },
		},
		checkNames: []string{
			preflight.CheckNameDocker,
			preflight.CheckNameImage,
			preflight.CheckNameOllama,
			preflight.CheckNameModel,
		},
	}
}

// buildCommand is the actual docker build invocation orchestrator.BuildImage
// runs, shown next to the live build panel for the same reason every other
// check shows its command.
func (m bootModel) buildCommand() string {
	return fmt.Sprintf("docker build -f docker/Dockerfile -t %s .", m.cfg.DockerImage)
}

func (m bootModel) Init() tea.Cmd {
	return tea.Batch(m.runCheck(0), tickCmd())
}

func (m bootModel) runCheck(i int) tea.Cmd {
	if i >= len(m.pending) {
		return nil
	}
	fn := m.pending[i]
	return func() tea.Msg { return checkDoneMsg(fn()) }
}

// waitForBuildLine blocks on the build's line channel and forwards the
// next line; once that channel closes (build's output has ended) it
// waits on the error channel instead, for the final result.
func waitForBuildLine(lineCh <-chan string, errCh <-chan error) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-lineCh
		if ok {
			return buildLineMsg(line)
		}
		return buildDoneMsg{err: <-errCh}
	}
}

func (m bootModel) advance(check preflight.Check) (tea.Model, tea.Cmd) {
	m.checks = append(m.checks, check)
	if len(m.checks) < len(m.pending) {
		// Pace the next check behind a delay instead of firing it
		// immediately — otherwise fast checks all resolve within the
		// same frame and there's nothing to actually watch happen.
		return m, delayedCheck(len(m.checks))
	}
	m.ready = true
	return m, nil
}

func (m bootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// A resize (especially a big jump, e.g. small window -> full
		// screen) can leave stale content from the old, differently-
		// centered render behind — force a full repaint instead of
		// relying on line-by-line diffing to catch it.
		return m, tea.ClearScreen

	case tickMsg:
		m.phase++
		return m, tickCmd()

	case startCheckMsg:
		return m, m.runCheck(msg.index)

	case checkDoneMsg:
		check := preflight.Check(msg)
		// The image check gets special handling: if it's not OK but
		// Docker (already resolved, one slot earlier) is reachable,
		// build it live right here instead of just reporting failure.
		if len(m.checks) == imageCheckIndex && !check.OK && m.checks[0].OK {
			lineCh, errCh := orchestrator.BuildImage(m.cfg)
			m.building = true
			m.buildLineCh = lineCh
			m.buildErrCh = errCh
			return m, waitForBuildLine(lineCh, errCh)
		}
		return m.advance(check)

	case buildLineMsg:
		line := string(msg)
		id := stepID(line)
		if n := len(m.buildSteps); n > 0 && m.buildSteps[n-1].id == id {
			step := &m.buildSteps[n-1]
			step.lines = append(step.lines, line)
			if len(step.lines) > maxStepLogLines {
				step.lines = step.lines[len(step.lines)-maxStepLogLines:]
			}
		} else {
			m.buildSteps = append(m.buildSteps, buildStepLog{id: id, lines: []string{line}})
		}
		return m, waitForBuildLine(m.buildLineCh, m.buildErrCh)

	case buildDoneMsg:
		m.building = false
		result := preflight.Check{
			Name:    preflight.CheckNameImage,
			OK:      msg.err == nil,
			Detail:  "built",
			Command: m.buildCommand(),
			Output:  lastBuildStepOutput(m.buildSteps),
		}
		if msg.err != nil {
			result.Detail = fmt.Sprintf("build failed: %v", msg.err)
		}
		m.buildSteps = nil
		return m.advance(result)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "enter":
			if m.ready {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// RunBoot shows the boot screen and blocks until the user presses Enter
// (proceed=true) or quits (proceed=false).
func RunBoot(cfg config.Config) (proceed bool, err error) {
	final, err := tea.NewProgram(newBootModel(cfg), tea.WithAltScreen()).Run()
	if err != nil {
		return false, err
	}
	return !final.(bootModel).quit, nil
}

// collapsedCheckRow is a compact one-liner for a check that's aged out
// of the recent window: mark, name, and its final detail — no command or
// output, keeping settled history out of the way.
func collapsedCheckRow(mark, name, detail string) string {
	return fmt.Sprintf("  %s %-16s %s\n", mark, name, checkDimStyle.Render(detail))
}

// expandedCheckRow shows a check with its real invoked command and up to
// maxOutputLines of its real output, indented underneath — for checks
// still within the recent window (see recentCheckWindow).
func expandedCheckRow(mark, name, command, output string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "  %s %s\n", mark, name)
	fmt.Fprintf(&b, "      %s\n", buildLogStyle.Render("$ "+truncateTitle(command, 80)))
	for _, line := range lastLines(output, maxOutputLines) {
		fmt.Fprintf(&b, "      %s\n", buildLogStyle.Render(truncateTitle(line, 80)))
	}
	return b.String()
}

// renderRightColumn renders the boot screen's right-hand content: the
// live checks (the last few invoked ones expanded with their real
// command + output, older ones collapsed to a summary line) / build log
// / continue prompt. The animated Ballroom banner above this is added by
// renderDashboardPanel, not here.
func (m bootModel) renderRightColumn() string {
	var b strings.Builder

	expandFrom := recentWindowStart(len(m.checks), recentCheckWindow)
	for i, c := range m.checks {
		mark := checkOKStyle.Render("✓")
		if !c.OK {
			mark = checkFailStyle.Render("✗")
		}
		if i < expandFrom {
			b.WriteString(collapsedCheckRow(mark, c.Name, c.Detail))
			continue
		}
		b.WriteString(expandedCheckRow(mark, c.Name, c.Command, c.Output))
	}

	startIdx := len(m.checks)
	if m.building {
		b.WriteString(expandedCheckRow(hintStyle.Render("▾"), preflight.CheckNameImage, m.buildCommand(), ""))
		for _, step := range m.buildSteps {
			for _, line := range step.lines {
				fmt.Fprintf(&b, "      %s\n", buildLogStyle.Render(truncateTitle(line, 90)))
			}
		}
		startIdx++ // the image slot is shown above, not in the queued loop
	}
	for i := startIdx; i < len(m.pending); i++ {
		fmt.Fprintf(&b, "  %s %s\n", checkDimStyle.Render("…"), m.checkNames[i])
	}

	if m.ready {
		b.WriteString("\n")
		b.WriteString("  " + hintStyle.Render("Press Enter to continue") + checkDimStyle.Render("  (q to quit)") + "\n")
	}

	return b.String()
}

func (m bootModel) View() string {
	right := m.renderRightColumn()
	if m.width == 0 || m.height == 0 {
		return right
	}
	panel := renderDashboardPanel(m.width, m.height, m.phase, right)
	return placeBlock(m.width, m.height, panel)
}
