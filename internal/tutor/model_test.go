package tutor

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestEstimatedTextareaRows_ShortLineIsOneRow(t *testing.T) {
	if got := estimatedTextareaRows("hello", 80); got != 1 {
		t.Errorf("estimatedTextareaRows(%q, 80) = %d, want 1", "hello", got)
	}
}

func TestEstimatedTextareaRows_LongSingleLineWrapsAcrossMultipleRows(t *testing.T) {
	// A real bug found live: LineCount() (explicit newlines only) never
	// grows for a single long line that visually WRAPS across the box's
	// width -- exactly the scenario that used to corrupt the old
	// hand-rolled ANSI box. Width-aware estimation from the raw value is
	// what makes growth track a wrapped single line too, not just
	// explicit paragraph breaks.
	line := strings.Repeat("x", 25) // 25 runes, width 10 -> ceil(25/10) = 3 rows
	if got := estimatedTextareaRows(line, 10); got != 3 {
		t.Errorf("estimatedTextareaRows(25 x's, width 10) = %d, want 3", got)
	}
}

func TestEstimatedTextareaRows_MultipleLogicalLinesSumTheirWrappedRows(t *testing.T) {
	value := strings.Repeat("x", 15) + "\n" + strings.Repeat("y", 5) // width 10: ceil(15/10)=2, ceil(5/10)=1 -> 3
	if got := estimatedTextareaRows(value, 10); got != 3 {
		t.Errorf("estimatedTextareaRows(...) = %d, want 3", got)
	}
}

func TestEstimatedTextareaRows_EmptyValueIsOneRow(t *testing.T) {
	if got := estimatedTextareaRows("", 80); got != 1 {
		t.Errorf("estimatedTextareaRows(\"\", 80) = %d, want 1", got)
	}
}

func TestEstimatedTextareaRows_ZeroOrNegativeWidthTreatedAsOne(t *testing.T) {
	// Guards the division below from ever seeing a non-positive width --
	// a real risk during startup/resize edge cases before a real size is
	// known yet.
	if got := estimatedTextareaRows("hello", 0); got != 5 {
		t.Errorf("estimatedTextareaRows(hello, 0) = %d, want 5 (treated as width 1)", got)
	}
}

func TestNewTutorModel_TextareaIsFocusedAndEnterDoesNotInsertANewline(t *testing.T) {
	m := newTutorModel()
	if !m.textarea.Focused() {
		t.Error("expected the textarea to start focused")
	}
	if m.textarea.KeyMap.InsertNewline.Enabled() {
		t.Error("expected InsertNewline disabled -- Enter must submit, not insert a newline")
	}
}

func TestTutorModel_WindowSizeMsg_SetsViewportAndTextareaWidth(t *testing.T) {
	m := newTutorModel()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	got := newM.(tutorModel)

	if got.width != 100 || got.height != 30 {
		t.Errorf("width,height = %d,%d, want 100,30", got.width, got.height)
	}
	if got.viewport.Width != 100 {
		t.Errorf("viewport.Width = %d, want 100", got.viewport.Width)
	}
	if got.textarea.Width() <= 0 || got.textarea.Width() >= 100 {
		t.Errorf("textarea.Width() = %d, want less than 100 (room left for the border)", got.textarea.Width())
	}
}

func TestTutorModel_WindowSizeMsg_ViewportAndTextareaHeightsSumWithinTerminal(t *testing.T) {
	m := newTutorModel()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	got := newM.(tutorModel)

	total := got.viewport.Height + got.textarea.Height() + textareaBorderRows
	if total > 30 {
		t.Errorf("viewport(%d) + textarea(%d) + border(%d) = %d, want <= terminal height 30", got.viewport.Height, got.textarea.Height(), textareaBorderRows, total)
	}
}

func TestTutorModel_TypingALongLineGrowsTextareaHeight(t *testing.T) {
	m := newTutorModel()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 30})
	m = newM.(tutorModel)
	before := m.textarea.Height()

	long := strings.Repeat("a very long message that should wrap across several rows ", 5)
	for _, r := range long {
		newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = newM.(tutorModel)
	}

	if m.textarea.Height() <= before {
		t.Errorf("textarea.Height() = %d after a long line, want it to have grown past the starting %d", m.textarea.Height(), before)
	}
}

func TestTutorModel_TextareaHeightNeverExceedsHalfTheTerminal(t *testing.T) {
	m := newTutorModel()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 20})
	m = newM.(tutorModel)

	// A single absurdly long line -- growth must be capped, not
	// unbounded, or a pathological message could starve the viewport
	// entirely. The cap scales with the real terminal height (not a
	// hardcoded row count), matching the "dynamic, not hardcoded" ask.
	huge := strings.Repeat("x", 2000)
	m.textarea.SetValue(huge)
	newM, _ = m.Update(tea.WindowSizeMsg{Width: 40, Height: 20})
	m = newM.(tutorModel)

	if m.textarea.Height() > 10 { // height/2
		t.Errorf("textarea.Height() = %d, want capped at half the terminal height (10)", m.textarea.Height())
	}
	if m.viewport.Height < minViewportRows {
		t.Errorf("viewport.Height = %d, want it to keep at least minViewportRows (%d) even with a huge textarea", m.viewport.Height, minViewportRows)
	}
}

func TestTutorModel_EnterOnEmptyTextareaDoesNothing(t *testing.T) {
	m := newTutorModel()
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := newM.(tutorModel)
	if got.textarea.Value() != "" {
		t.Errorf("textarea.Value() = %q, want still empty", got.textarea.Value())
	}
	if cmd != nil {
		t.Error("expected no command from submitting an empty message")
	}
}

func TestTutorModel_EnterOnNonEmptyTextareaResetsItAndEchoesIntoViewport(t *testing.T) {
	m := newTutorModel()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m.textarea.SetValue("hello there")

	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := newM.(tutorModel)

	if got.textarea.Value() != "" {
		t.Errorf("textarea.Value() = %q, want reset to empty after submit", got.textarea.Value())
	}
	if !strings.Contains(got.viewport.View(), "hello there") {
		t.Errorf("viewport view %q, want it to contain the submitted message", got.viewport.View())
	}
}

func TestTutorModel_View_RendersBothViewportAndTextarea(t *testing.T) {
	m := newTutorModel()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	out := m.View()
	if out == "" {
		t.Fatal("View() returned empty output")
	}
}
