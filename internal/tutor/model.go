package tutor

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// minTextareaRows/minViewportRows are floors, not targets — the textarea
// always shows at least one row even when empty, and the viewport always
// keeps at least a few rows of conversation visible even when the
// textarea has grown to its cap (see recomputeLayout).
const (
	minTextareaRows    = 1
	minViewportRows    = 3
	textareaBorderRows = 2 // top + bottom border, added by textareaBoxStyle
	textareaBorderCols = 2 // left + right border, same style
)

// textareaBoxStyle replaces the old hand-rolled box borders
// (internal/tutor/scrollbox.go, deleted alongside this rewrite) with a
// lipgloss border — same teal accent used elsewhere in this project's
// palette (docker/tmux.conf, internal/tui/styles.go).
var textareaBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#2FA6A6"))

// viewportContentStyle left-pads the scrolling conversation area — a
// real readability complaint found live: text printed flush against the
// pane's own left edge read as cramped. A lipgloss primitive instead of
// prefixing every print call site by hand (the old architecture's only
// option).
var viewportContentStyle = lipgloss.NewStyle().PaddingLeft(2)

// tutorModel is the tutor pane's bubbletea Model — replaces
// internal/tutor/scrollbox.go's hand-rolled ANSI positioning entirely.
// Mirrors internal/tui/app.go's own Model/Update/View architecture,
// which has needed zero manual cursor/escape-sequence code all session,
// unlike the hand-rolled approach this replaces.
type tutorModel struct {
	viewport viewport.Model
	textarea textarea.Model
	width    int
	height   int

	// history accumulates submitted lines for display. Stage 1 only:
	// echoes what was typed, no real model/tool wiring yet — that lands
	// in a later stage of this same rewrite (see the plan for this
	// feature).
	history []string
}

// newTutorModel builds the initial model: an empty, focused textarea
// with Enter reserved for submit (not newline-insert — see
// estimatedTextareaRows's doc comment for why "Enter submits" needs a
// width-aware growth estimate to work correctly for a single long
// wrapped line, not just explicit paragraph breaks), and an empty
// viewport for the scrolling conversation.
func newTutorModel() tutorModel {
	ta := textarea.New()
	ta.Placeholder = "Ask a question..."
	ta.Prompt = "> "
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.Focus()

	vp := viewport.New(0, minViewportRows)
	vp.Style = viewportContentStyle

	return tutorModel{viewport: vp, textarea: ta}
}

func (m tutorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m tutorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.recomputeLayout()
		return m, tea.ClearScreen

	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			return m.submit()
		}
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		m.recomputeLayout()
		return m, cmd
	}
	return m, nil
}

// recomputeLayout resizes the textarea and viewport to fit the current
// terminal dimensions and the textarea's current content — called on
// every resize AND every keystroke (per bubbles/textarea and
// bubbles/viewport's own design: neither auto-grows or reacts to resize
// itself, the composing model always owns this). This per-keystroke
// recompute is the actual mechanism that makes input growth dynamic,
// not something built from scratch.
//
// The textarea's height is capped at half the terminal height (a floor
// of minTextareaRows, a ceiling that scales with the real terminal size
// rather than a fixed row count) so a pathological single huge message
// can't starve the viewport entirely — beyond the cap, bubbles/textarea
// scrolls *within its own bounded region*, which is the structural fix
// over the old hand-rolled box: View() always renders exactly as many
// rows as it's told to, so there is no way for overflow to land outside
// the box and corrupt anything else, unlike raw cooked-mode wrapping.
func (m *tutorModel) recomputeLayout() {
	m.textarea.SetWidth(max(m.width-textareaBorderCols, 1))

	maxTaRows := max(m.height/2, minTextareaRows)
	taRows := estimatedTextareaRows(m.textarea.Value(), m.textarea.Width())
	taRows = min(max(taRows, minTextareaRows), maxTaRows)
	m.textarea.SetHeight(taRows)

	m.viewport.Width = m.width
	m.viewport.Height = max(m.height-taRows-textareaBorderRows, minViewportRows)
}

// estimatedTextareaRows estimates how many visual rows value will wrap
// to at the given width — deliberately from the raw value text, not
// textarea.LineCount() (which only counts explicit newlines). A real bug
// this is fixing: a single long line that wraps purely from exceeding
// the box's width never grew LineCount() at all, so a height computed
// from it never grew either — exactly the scenario that used to overflow
// past the old hand-rolled box's last row and corrupt the terminal. Pure
// and testable without a real terminal. width <= 0 is treated as 1 to
// keep the division defined during startup/resize edge cases before a
// real size is known.
func estimatedTextareaRows(value string, width int) int {
	if width <= 0 {
		width = 1
	}
	lines := strings.Split(value, "\n")
	rows := 0
	for _, line := range lines {
		w := utf8.RuneCountInString(line)
		rows += max(1, (w+width-1)/width) // ceil division
	}
	return max(rows, 1)
}

// submit handles Enter: empty input is a no-op (nothing to send), a real
// message resets the textarea (growth immediately collapses back down —
// recomputeLayout runs again right after) and echoes into the viewport.
// Stage 1 only echoes; real turn-loop wiring (comprehension check,
// routing, the actual model call) lands in a later stage of this
// rewrite.
func (m tutorModel) submit() (tea.Model, tea.Cmd) {
	line := strings.TrimSpace(m.textarea.Value())
	if line == "" {
		return m, nil
	}
	m.textarea.Reset()
	m.history = append(m.history, "> "+line)
	m.viewport.SetContent(strings.Join(m.history, "\n"))
	m.viewport.GotoBottom()
	m.recomputeLayout()
	return m, nil
}

func (m tutorModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), textareaBoxStyle.Render(m.textarea.View()))
}
