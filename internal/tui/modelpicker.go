package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/preflight"
)

// listModelsFn and checkModelFn are vars (not direct calls) so tests can
// substitute fakes instead of making real HTTP calls to Ollama — same
// indirection pattern as buildImageFn in boot.go.
var listModelsFn = preflight.ListModels
var checkModelFn = preflight.CheckModel

// suggestedModels are known-good tutor model tags surfaced in the picker
// even before they're pulled locally, so they're discoverable without
// already knowing their exact Ollama tag. Selecting one that isn't
// pulled yet goes through the same not-pulled warn-then-confirm flow as
// typing an arbitrary tag (see handleEnter).
var suggestedModels = []string{
	config.DeepSeekCoderV2LiteModel,
	config.Qwen25Coder14BModel,
}

// browsableModels merges local (locally pulled) with suggestedModels,
// deduplicated — a suggested tag already pulled just appears once, as a
// normal local entry.
func browsableModels(local []string) []string {
	seen := make(map[string]bool, len(local))
	out := make([]string, 0, len(local)+len(suggestedModels))
	for _, name := range local {
		seen[name] = true
		out = append(out, name)
	}
	for _, name := range suggestedModels {
		if !seen[name] {
			out = append(out, name)
		}
	}
	return out
}

// modelsLoadedMsg carries the result of querying Ollama's /api/tags for
// the locally pulled model list.
type modelsLoadedMsg struct {
	models []string
	err    error
}

// modelPickerModel is a popup listing locally pulled Ollama models plus
// suggestedModels (known-good tags surfaced for discoverability even
// before they're pulled), reached from the main menu rather than
// re-shown before every exercise launch — the pick is persisted (see
// config.Settings) so it only needs asking once. Typing filters the list
// in place; when nothing in that combined list matches, the typed text
// is treated as a candidate arbitrary model tag instead. Selecting any
// entry not actually pulled locally — whether a suggested one from the
// list or a freely typed tag — checks it against Ollama (reusing
// preflight.CheckModel's "not pulled" messaging) and warns rather than
// silently accepting a tag that isn't actually there. A second enter
// after the warning confirms the tag anyway, matching this codebase's
// "informational, not blocking" preflight-check philosophy.
//
// config.Qwen25Coder14BModel — a larger variant of the default 7B model —
// needs meaningfully more RAM/VRAM to run (roughly 12-16GB free vs. the
// 7B default's ~8GB). There's no per-entry hint UI here to surface that
// inline, so it's called out in the const's doc comment instead — check
// there before pulling/selecting it on constrained hardware.
type modelPickerModel struct {
	host    string
	current string

	loading     bool
	loadErr     error
	localModels []string // exactly what Ollama reports pulled, for isLocal
	models      []string // localModels + suggestedModels, deduplicated

	filter   string
	filtered []string
	cursor   int
	warning  string

	selected      *string
	back          bool
	width, height int
}

// isLocal reports whether name is actually pulled locally (as opposed to
// merely suggested) — selecting a local entry can skip straight to
// selection; anything else needs the not-pulled check first.
func (m modelPickerModel) isLocal(name string) bool {
	for _, local := range m.localModels {
		if local == name {
			return true
		}
	}
	return false
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
		m.localModels = msg.models
		m.models = browsableModels(msg.models)
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

// handleEnter selects the highlighted entry if it's already pulled
// locally, or — for a highlighted-but-unpulled suggested entry, or when
// the typed filter matches nothing — treats the tag as a candidate:
// confirmTag checks it against Ollama and warns if it isn't pulled, and a
// second enter (with the warning already showing) confirms using it
// anyway.
func (m modelPickerModel) handleEnter() (tea.Model, tea.Cmd) {
	if len(m.filtered) > 0 {
		sel := m.filtered[m.cursor]
		if m.isLocal(sel) {
			m.selected = &sel
			return m, tea.Quit
		}
		return m.confirmTag(sel)
	}

	tag := strings.TrimSpace(m.filter)
	if tag == "" {
		return m, nil
	}
	return m.confirmTag(tag)
}

// confirmTag is the not-pulled warn-then-confirm flow shared by both a
// freely typed tag and a highlighted suggested-but-unpulled entry.
func (m modelPickerModel) confirmTag(tag string) (tea.Model, tea.Cmd) {
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
		b.WriteString(checkDimStyle.Render("  no matches"))
		b.WriteString("\n")
	default:
		for i, name := range m.filtered {
			label := name
			if name == m.current {
				label += "  " + hintStyle.Render("(current)")
			} else if !m.isLocal(name) {
				label += "  " + checkDimStyle.Render("(not pulled)")
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
