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
	"github.com/JacobTDang/Ballroom/internal/preflight"
)

const (
	ollamaHost = "http://localhost:11434"
	tutorModel = "qwen2.5-coder:7b"
)

var (
	checkOKStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("120")).Bold(true)
	checkFailStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("210")).Bold(true)
	checkDimStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	hintStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("87")).Bold(true)
)

type checkDoneMsg preflight.Check

// bootModel runs preflight checks one at a time (so they visibly appear in
// sequence, like a real boot log) and then waits for the user to continue.
type bootModel struct {
	pending []func() preflight.Check
	checks  []preflight.Check
	ready   bool
	quit    bool
}

func newBootModel(cfg config.Config) bootModel {
	image := cfg.DockerImage
	return bootModel{
		pending: []func() preflight.Check{
			preflight.CheckDocker,
			func() preflight.Check { return preflight.CheckImage(image) },
			func() preflight.Check { return preflight.CheckOllama(ollamaHost) },
			func() preflight.Check { return preflight.CheckModel(ollamaHost, tutorModel) },
		},
	}
}

func (m bootModel) Init() tea.Cmd {
	return m.runCheck(0)
}

func (m bootModel) runCheck(i int) tea.Cmd {
	if i >= len(m.pending) {
		return nil
	}
	fn := m.pending[i]
	return func() tea.Msg { return checkDoneMsg(fn()) }
}

func (m bootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case checkDoneMsg:
		m.checks = append(m.checks, preflight.Check(msg))
		if len(m.checks) < len(m.pending) {
			return m, m.runCheck(len(m.checks))
		}
		m.ready = true
		return m, nil

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
	final, err := tea.NewProgram(newBootModel(cfg)).Run()
	if err != nil {
		return false, err
	}
	return !final.(bootModel).quit, nil
}

func (m bootModel) View() string {
	var b strings.Builder
	b.WriteString(catalog.Banner())
	b.WriteString("\n")

	for _, c := range m.checks {
		mark := checkOKStyle.Render("✓")
		if !c.OK {
			mark = checkFailStyle.Render("✗")
		}
		fmt.Fprintf(&b, "  %s %-16s %s\n", mark, c.Name, checkDimStyle.Render(c.Detail))
	}
	for i := len(m.checks); i < len(m.pending); i++ {
		fmt.Fprintf(&b, "  %s\n", checkDimStyle.Render("… checking"))
	}

	if m.ready {
		b.WriteString("\n")
		b.WriteString("  " + hintStyle.Render("Press Enter to continue") + checkDimStyle.Render("  (q to quit)") + "\n")
	}
	return b.String()
}
