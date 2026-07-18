// Package tui is the full-screen interactive homepage: a boot screen that
// runs preflight checks, then a picker you navigate with arrow keys.
package tui

import (
	"errors"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/palette"
	"github.com/JacobTDang/Ballroom/internal/preflight"
	"github.com/JacobTDang/Ballroom/internal/tutor"
)

const ollamaHost = "http://localhost:11434"

// maxStepLogLines caps how many of a single docker-build step's own
// lines stay on screen — the step entry itself persists for the whole
// build (so progress doesn't disappear), only its own scrollback within
// that step is windowed.
const maxStepLogLines = 3

// maxOutputLines caps how many lines of a resolved check's real command
// output are shown — same windowing idea as maxStepLogLines, applied to
// preflight.Check.Output instead of docker build's step lines.
const maxOutputLines = 3

// checkStartDelay paces checks so they visibly run one at a time instead
// of all resolving within the same frame — a var (not a const) so tests
// can zero it out instead of actually sleeping.
var checkStartDelay = 250 * time.Millisecond

var (
	checkOKStyle   = lipgloss.NewStyle().Foreground(palette.Lip(palette.Teal)).Bold(true)
	checkFailStyle = lipgloss.NewStyle().Foreground(palette.Lip(palette.Red)).Bold(true)
	checkDimStyle  = lipgloss.NewStyle().Foreground(palette.Lip(palette.PaleGray))
	hintStyle      = lipgloss.NewStyle().Foreground(palette.Lip(palette.Gold)).Bold(true)
	buildLogStyle  = lipgloss.NewStyle().Foreground(palette.Lip(palette.PaleGray)).Faint(true)
)

type checkDoneMsg preflight.Check
type buildLineMsg string
type buildDoneMsg struct{ err error }
type pullLineMsg string
type pullDoneMsg struct{ err error }

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

// lastBuildStepOutput carries the last docker-build step's lines into
// the resolved Image check's Output, so it doesn't lose all its visible
// output the moment the live build panel collapses into a plain check.
func lastBuildStepOutput(steps []buildStepLog) string {
	if len(steps) == 0 {
		return ""
	}
	return strings.Join(steps[len(steps)-1].lines, "\n")
}

// lastBuildLines flattens every step's lines (in order) and returns at
// most the last n across the whole build — a real build has far more
// steps than fit on screen, so the live panel needs one rolling window
// over the whole stream, not a per-step cap that grows with step count.
func lastBuildLines(steps []buildStepLog, n int) []string {
	var all []string
	for _, step := range steps {
		all = append(all, step.lines...)
	}
	if len(all) > n {
		all = all[len(all)-n:]
	}
	return all
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

// pullModelFn is a var (not a direct call) so tests can substitute a fake
// pull stream instead of making a real HTTP call to Ollama.
var pullModelFn = preflight.PullModel

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

	// pullingModel etc. mirror the building/buildSteps fields above, for
	// the same live-panel treatment applied to a fallback `ollama pull`
	// instead of `docker build` — see checkDoneMsg's handling of
	// preflight.CheckNameModel.
	pullingModel     bool
	pullLines        []string
	pullLineCh       <-chan string
	pullErrCh        <-chan error
	pullFallbackFrom string // the originally configured model, captured when the fallback pull starts

	width, height int

	// checkNames mirrors pending, one entry per slot, so a check's name
	// can be shown while it's still queued — before it has actually run
	// and produced a Check with a real Command/Output to display.
	checkNames []string
}

func newBootModel(cfg config.Config) bootModel {
	// modelCheck defaults to the real local-Ollama lookup, but an
	// OpenRouter-prefixed model was never a candidate for that at all —
	// preflight.CheckModel only ever queries Ollama's own /api/tags, so
	// it always reported an OpenRouter model as "not pulled". That
	// always tripped checkDoneMsg's pull-fallback path (Ollama itself
	// being reachable was the only other condition), which pulled
	// config.DefaultTutorModel and silently overwrote cfg.TutorModel for
	// the session — a real bug found live: the worker model kept
	// reverting to llama on every restart, and the corrupted value then
	// got persisted to settings.json the next time any setting was
	// saved. Trivially OK here instead of ever making that call.
	modelCheck := func() preflight.Check { return preflight.CheckModel(ollamaHost, cfg.TutorModel) }
	if strings.HasPrefix(cfg.TutorModel, tutor.OpenRouterModelPrefix) {
		modelCheck = func() preflight.Check {
			return preflight.Check{Name: preflight.CheckNameModel, OK: true, Detail: cfg.TutorModel + " (OpenRouter, not a local Ollama model)"}
		}
	}
	return bootModel{
		cfg: cfg,
		pending: []func() preflight.Check{
			preflight.CheckDocker,
			// The image slot's own result is never used — checkDoneMsg
			// always runs a real `docker build` here instead (see
			// imageCheckIndex), so there's no point spending a
			// `docker image inspect` call first just to discard it.
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameImage} },
			func() preflight.Check { return preflight.CheckOllama(ollamaHost) },
			modelCheck,
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

// pullCommand is the actual Ollama request preflight.PullModel makes to
// fetch the fallback default model, shown next to the live pull panel for
// the same reason every other check shows its command.
func (m bootModel) pullCommand() string {
	return fmt.Sprintf("POST %s/api/pull {\"model\":%q}", ollamaHost, config.DefaultTutorModel)
}

// checkOK reports whether checks contains a check with the given name
// that passed — used to gate the model-pull fallback on Ollama itself
// actually being reachable, so it doesn't attempt (and immediately fail)
// a pull when the real problem is that Ollama isn't running at all.
func checkOK(checks []preflight.Check, name string) bool {
	for _, c := range checks {
		if c.Name == name {
			return c.OK
		}
	}
	return false
}

// requiredCheckNames are the checks that must pass before the boot
// screen will let you continue: without a running Docker daemon and a
// built practice image, `docker run` at launch just fails outright, so
// there's nothing "Press Enter to continue" would actually do (issue
// #230). Ollama and the tutor model are deliberately not required --
// without them the tutor pane won't work, but you can still code, run
// tests, and submit inside a session (see this package's own doc
// comment on checks being informational).
var requiredCheckNames = map[string]bool{
	preflight.CheckNameDocker: true,
	preflight.CheckNameImage:  true,
}

// allRequiredChecksPassed reports whether every required check (see
// requiredCheckNames) present in checks passed.
func allRequiredChecksPassed(checks []preflight.Check) bool {
	for _, c := range checks {
		if requiredCheckNames[c.Name] && !c.OK {
			return false
		}
	}
	return true
}

// blockedReason names every failed required check and its remedy (the
// check's own Detail, e.g. "not running — start Docker Desktop") --
// e.g. "Docker daemon: not running — start Docker Desktop". A required
// check with an empty Detail was never actually attempted (the image
// slot's placeholder result when Docker itself failed -- see
// checkDoneMsg) rather than a real failure with its own remedy, so it's
// skipped: Docker's own failure above is already the actionable reason.
func blockedReason(checks []preflight.Check) string {
	var lines []string
	for _, c := range checks {
		if requiredCheckNames[c.Name] && !c.OK && c.Detail != "" {
			lines = append(lines, c.Name+": "+c.Detail)
		}
	}
	return strings.Join(lines, "\n")
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

// waitForPullLine is waitForBuildLine's counterpart for a live model pull.
func waitForPullLine(lineCh <-chan string, errCh <-chan error) tea.Cmd {
	return func() tea.Msg {
		line, ok := <-lineCh
		if ok {
			return pullLineMsg(line)
		}
		return pullDoneMsg{err: <-errCh}
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
	// Only a required check (Docker, the practice image) blocks
	// continuing -- Ollama/the tutor model failing still resolves the
	// sequence as "ready", same as before this field existed, so you're
	// never stuck over a check that isn't required (issue #230).
	m.ready = allRequiredChecksPassed(m.checks)
	return m, nil
}

// blocked reports whether the boot sequence has finished (every pending
// check has resolved) but ready is still false — a required check (see
// requiredCheckNames) failed, and nothing changes until the user
// re-runs the checks (the "r" key in Update) and they pass.
func (m bootModel) blocked() bool {
	return len(m.pending) > 0 && len(m.checks) == len(m.pending) && !m.ready
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
		// The image check always runs a real `docker build` live right
		// here instead of just inspecting whether the tag exists —
		// Docker's own layer cache makes this fast when nothing
		// changed (every step shows CACHED), and it's the only way to
		// show real build output either way: downloading fresh layers,
		// or confirming the cache was used. Only skipped if Docker
		// itself isn't reachable. Gated on the check's Name (not its
		// position in pending) so this only ever fires for the actual
		// image check.
		if check.Name == preflight.CheckNameImage && len(m.checks) > 0 && m.checks[0].OK {
			lineCh, errCh := buildImageFn(m.cfg)
			m.building = true
			m.buildLineCh = lineCh
			m.buildErrCh = errCh
			return m, waitForBuildLine(lineCh, errCh)
		}
		// The configured tutor model isn't pulled — fall back to pulling
		// the default (config.DefaultTutorModel) live right here, same
		// "don't leave you stuck, actually fix it" treatment as the
		// image-build step
		// above. Gated on Ollama itself being reachable so this doesn't
		// also attempt (and immediately fail) a pull when the real
		// problem is Ollama not running at all.
		if check.Name == preflight.CheckNameModel && !check.OK && checkOK(m.checks, preflight.CheckNameOllama) {
			m.pullFallbackFrom = m.cfg.TutorModel
			lineCh, errCh := pullModelFn(ollamaHost, config.DefaultTutorModel)
			m.pullingModel = true
			m.pullLineCh = lineCh
			m.pullErrCh = errCh
			return m, waitForPullLine(lineCh, errCh)
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

	case pullLineMsg:
		line := string(msg)
		m.pullLines = append(m.pullLines, line)
		if len(m.pullLines) > maxOutputLines {
			m.pullLines = m.pullLines[len(m.pullLines)-maxOutputLines:]
		}
		return m, waitForPullLine(m.pullLineCh, m.pullErrCh)

	case pullDoneMsg:
		m.pullingModel = false
		result := preflight.Check{
			Name:    preflight.CheckNameModel,
			Command: m.pullCommand(),
			Output:  strings.Join(m.pullLines, "\n"),
		}
		if msg.err != nil {
			result.OK = false
			result.Detail = fmt.Sprintf("%s not pulled, and falling back to %s failed: %v", m.pullFallbackFrom, config.DefaultTutorModel, msg.err)
		} else {
			result.OK = true
			result.Detail = fmt.Sprintf("%s not pulled — downloaded default %s instead", m.pullFallbackFrom, config.DefaultTutorModel)
			m.cfg.TutorModel = config.DefaultTutorModel
		}
		m.pullLines = nil
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
		case "r":
			// Only meaningful once blocked -- re-running mid-sequence or
			// after a successful boot has nothing to fix.
			if m.blocked() {
				m.checks = nil
				m.building = false
				m.buildSteps = nil
				m.pullingModel = false
				m.pullLines = nil
				return m, m.runCheck(0)
			}
		}
	}
	return m, nil
}

// RunBoot shows the boot screen and blocks until the user presses Enter
// (proceed=true) or quits (proceed=false). The returned Config may differ
// from the one passed in: if the configured tutor model wasn't pulled,
// boot falls back to pulling the default (config.DefaultTutorModel) live
// and switches to it for this run only — the persisted setting is left
// untouched, so a future launch still tries the real pick first.
func RunBoot(cfg config.Config) (result config.Config, proceed bool, err error) {
	final, err := tea.NewProgram(newBootModel(cfg), tea.WithAltScreen()).Run()
	if err != nil {
		return cfg, false, err
	}
	fm := final.(bootModel)
	return fm.cfg, !fm.quit, nil
}

// expandedCheckRow shows a check with its real invoked command and up to
// maxOutputLines of its real output, indented underneath. Every check
// stays in this form for the whole boot sequence — nothing collapses
// away or gets replaced by a summary once it's been shown.
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
// live checks, each with the real command it ran and its real output /
// build log / continue prompt. The animated Ballroom banner above this
// is added by renderDashboardPanel, not here.
func (m bootModel) renderRightColumn() string {
	var b strings.Builder

	// The checks already stream in one at a time (checkStartDelay), so
	// framing them as a titled sequence costs nothing and makes the
	// wait read as a machine coming up rather than a stalled screen.
	b.WriteString(hintStyle.Render(heading("System check")))
	b.WriteString("\n\n")

	for _, c := range m.checks {
		mark := checkOKStyle.Render("✓")
		if !c.OK {
			mark = checkFailStyle.Render("✗")
		}
		b.WriteString(expandedCheckRow(mark, c.Name, c.Command, c.Output))
	}

	startIdx := len(m.checks)
	if m.building {
		b.WriteString(expandedCheckRow(hintStyle.Render("▾"), preflight.CheckNameImage, m.buildCommand(), ""))
		for _, line := range lastBuildLines(m.buildSteps, maxOutputLines) {
			fmt.Fprintf(&b, "      %s\n", buildLogStyle.Render(truncateTitle(line, 90)))
		}
		startIdx++ // the image slot is shown above, not in the queued loop
	} else if m.pullingModel {
		b.WriteString(expandedCheckRow(hintStyle.Render("▾"), preflight.CheckNameModel, m.pullCommand(), ""))
		for _, line := range m.pullLines {
			fmt.Fprintf(&b, "      %s\n", buildLogStyle.Render(truncateTitle(line, 90)))
		}
		startIdx++ // the model slot is shown above, not in the queued loop
	}
	for i := startIdx; i < len(m.pending); i++ {
		fmt.Fprintf(&b, "  %s %s\n", checkDimStyle.Render("…"), m.checkNames[i])
	}

	if m.ready {
		b.WriteString("\n")
		b.WriteString("  " + passStyle.Render("READY.") + "\n\n")
		b.WriteString("  " + hintStyle.Render("Press Enter to continue") + checkDimStyle.Render("  (q to quit)") + "\n")
	} else if m.blocked() {
		b.WriteString("\n")
		b.WriteString(renderFriendlyError("can't continue", errors.New(blockedReason(m.checks))))
		b.WriteString("\n\n")
		b.WriteString("  " + hintStyle.Render("r") + checkDimStyle.Render(" to re-run checks  ·  ") + hintStyle.Render("q") + checkDimStyle.Render(" to quit") + "\n")
	}

	return b.String()
}

func (m bootModel) View() string {
	right := m.renderRightColumn()
	if m.width == 0 || m.height == 0 {
		return right
	}
	panel := renderDashboardPanel(m.width, m.height, m.phase, right, layoutTop, "")
	return placeBlock(m.width, m.height, panel)
}
