package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/preflight"
)

// listModelsFn and checkModelFn are vars (not direct calls) so tests can
// substitute fakes instead of making real HTTP calls to Ollama — same
// indirection pattern as buildImageFn in boot.go.
var listModelsFn = preflight.ListModels
var checkModelFn = preflight.CheckModel

// modelsLoadedMsg carries the result of querying Ollama's /api/tags for
// the locally pulled model list.
type modelsLoadedMsg struct {
	models []string
	err    error
}

// modelPickerModel is a popup listing locally pulled Ollama models,
// reached from the main menu rather than re-shown before every exercise
// launch — the pick is persisted (see config.Settings) so it only needs
// asking once. Typing filters the list in place; when nothing local
// matches, the typed text is treated as a candidate arbitrary model tag
// instead — enter checks it against Ollama (reusing preflight.CheckModel's
// "not pulled" messaging) and warns rather than silently accepting a tag
// that isn't actually there. A second enter after the warning confirms the
// tag anyway, matching this codebase's "informational, not blocking"
// preflight-check philosophy.
type modelPickerModel struct {
	host    string
	current string

	loading bool
	loadErr error
	models  []string

	filter   string
	filtered []string
	cursor   int
	warning  string

	selected      *string
	back          bool
	width, height int
}

func newModelPickerModel(host, current string) modelPickerModel {
	return modelPickerModel{host: host, current: current, loading: true}
}

func (m modelPickerModel) Init() tea.Cmd {
	return func() tea.Msg {
		models, err := listModelsFn(m.host)
		return modelsLoadedMsg{models: models, err: err}
	}
}

// filterModels returns the models whose name contains filter
// (case-insensitive). An empty filter matches everything.
func filterModels(models []string, filter string) []string {
	if filter == "" {
		return models
	}
	lower := strings.ToLower(filter)
	var out []string
	for _, name := range models {
		if strings.Contains(strings.ToLower(name), lower) {
			out = append(out, name)
		}
	}
	return out
}

func (m modelPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, tea.ClearScreen

	case modelsLoadedMsg:
		m.loading = false
		m.loadErr = msg.err
		m.models = msg.models
		m.filtered = filterModels(m.models, m.filter)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case tea.KeyDown:
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
			return m, nil
		case tea.KeyBackspace:
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.filtered = filterModels(m.models, m.filter)
				m.cursor = 0
				m.warning = ""
			}
			return m, nil
		case tea.KeyEsc, tea.KeyCtrlC:
			m.back = true
			return m, tea.Quit
		case tea.KeyEnter:
			return m.handleEnter()
		case tea.KeyRunes:
			// "q" with nothing typed yet backs out, matching every other
			// picker in this package — once the user has started typing,
			// every rune (including "q") feeds the filter/custom tag
			// instead, since it might be part of a real model name.
			if m.filter == "" && string(msg.Runes) == "q" {
				m.back = true
				return m, tea.Quit
			}
			m.filter += string(msg.Runes)
			m.filtered = filterModels(m.models, m.filter)
			m.cursor = 0
			m.warning = ""
			return m, nil
		}
	}
	return m, nil
}

// handleEnter selects the highlighted local model, or — when the typed
// filter matches nothing local — treats the filter itself as a candidate
// model tag: the first enter checks it against Ollama and warns if it
// isn't pulled, and a second enter (with the warning already showing)
// confirms using it anyway.
func (m modelPickerModel) handleEnter() (tea.Model, tea.Cmd) {
	if len(m.filtered) > 0 {
		sel := m.filtered[m.cursor]
		m.selected = &sel
		return m, tea.Quit
	}

	tag := strings.TrimSpace(m.filter)
	if tag == "" {
		return m, nil
	}

	if m.warning != "" {
		sel := tag
		m.selected = &sel
		return m, tea.Quit
	}

	check := checkModelFn(m.host, tag)
	if !check.OK {
		m.warning = check.Detail
		return m, nil
	}
	sel := tag
	m.selected = &sel
	return m, tea.Quit
}

func (m modelPickerModel) View() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("Tutor model"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("type to search, or enter a custom tag"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("  %s%s", checkDimStyle.Render("› "), m.filter))
	b.WriteString("\n\n")

	switch {
	case m.loading:
		b.WriteString(checkDimStyle.Render("  loading models from " + m.host + "..."))
		b.WriteString("\n")
	case m.loadErr != nil:
		b.WriteString(failStyle.Render("  couldn't reach Ollama: " + m.loadErr.Error()))
		b.WriteString("\n")
		b.WriteString(checkDimStyle.Render("  you can still type a model tag directly"))
		b.WriteString("\n")
	case len(m.filtered) == 0:
		b.WriteString(checkDimStyle.Render("  no local matches"))
		b.WriteString("\n")
	default:
		for i, name := range m.filtered {
			label := name
			if name == m.current {
				label += "  " + hintStyle.Render("(current)")
			}
			if i == m.cursor {
				b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s", label)))
			} else {
				b.WriteString(fmt.Sprintf("  %s", label))
			}
			b.WriteString("\n")
		}
	}

	if m.warning != "" {
		b.WriteString("\n")
		b.WriteString(hintStyle.Render("  " + m.warning))
		b.WriteString("\n")
		b.WriteString(checkDimStyle.Render("  press enter again to use it anyway"))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · esc back"))

	content := popupBoxStyle.Render(b.String())
	if m.width > 0 && m.height > 0 {
		return placeBlock(m.width, m.height, content)
	}
	return content
}

// RunModelPicker shows the model popup and blocks until the user picks or
// types a model (ok=true) or backs out (ok=false).
func RunModelPicker(host, current string) (model string, ok bool, err error) {
	final, err := tea.NewProgram(newModelPickerModel(host, current), tea.WithAltScreen()).Run()
	if err != nil {
		return "", false, err
	}
	mm := final.(modelPickerModel)
	if mm.selected == nil {
		return "", false, nil
	}
	return *mm.selected, true, nil
}
