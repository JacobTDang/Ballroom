// Package tui is the full-screen interactive homepage: a boot screen that
// runs preflight checks, then a picker you navigate with arrow keys.
package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/preflight"
)

const (
	ollamaHost = "http://localhost:11434"
	tutorModel = "qwen2.5-coder:7b"
)

// maxBuildLogLines caps how many recent `docker build` lines the boot
// screen's live panel keeps on screen — the point is to show it's making
// progress, not to be a full build log.
const maxBuildLogLines = 12

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
	buildLines  []string
	buildLineCh <-chan string
	buildErrCh  <-chan error

	width, height int
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
	}
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
		return m, m.runCheck(len(m.checks))
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
		m.buildLines = append(m.buildLines, string(msg))
		if len(m.buildLines) > maxBuildLogLines {
			m.buildLines = m.buildLines[len(m.buildLines)-maxBuildLogLines:]
		}
		return m, waitForBuildLine(m.buildLineCh, m.buildErrCh)

	case buildDoneMsg:
		m.building = false
		m.buildLines = nil
		result := preflight.Check{Name: preflight.CheckNameImage, OK: msg.err == nil, Detail: "built"}
		if msg.err != nil {
			result.Detail = fmt.Sprintf("build failed: %v", msg.err)
		}
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

func (m bootModel) View() string {
	var b strings.Builder
	b.WriteString(catalog.MosaicBanner(m.phase))
	b.WriteString("\n")

	for _, c := range m.checks {
		mark := checkOKStyle.Render("✓")
		if !c.OK {
			mark = checkFailStyle.Render("✗")
		}
		fmt.Fprintf(&b, "  %s %-16s %s\n", mark, c.Name, checkDimStyle.Render(c.Detail))
	}

	pendingCount := len(m.pending) - len(m.checks)
	if m.building {
		fmt.Fprintf(&b, "  %s %-16s %s\n", hintStyle.Render("▾"), preflight.CheckNameImage,
			checkDimStyle.Render("building — this can take a minute or two"))
		for _, line := range m.buildLines {
			fmt.Fprintf(&b, "      %s\n", buildLogStyle.Render(truncateTitle(line, 90)))
		}
		pendingCount-- // the image slot is shown above, not as "… checking"
	}
	for i := 0; i < pendingCount; i++ {
		fmt.Fprintf(&b, "  %s\n", checkDimStyle.Render("… checking"))
	}

	if m.ready {
		b.WriteString("\n")
		b.WriteString("  " + hintStyle.Render("Press Enter to continue") + checkDimStyle.Render("  (q to quit)") + "\n")
	}

	content := b.String()
	if m.width > 0 && m.height > 0 {
		return placeBlock(m.width, m.height, content)
	}
	return content
}
