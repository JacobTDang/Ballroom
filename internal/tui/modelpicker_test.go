package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/preflight"
)

// fakeListModels substitutes listModelsFn in tests so no real HTTP call to
// Ollama is made. Returns a restore func to defer.
func fakeListModels(models []string, err error) func() {
	orig := listModelsFn
	listModelsFn = func(string) ([]string, error) { return models, err }
	return func() { listModelsFn = orig }
}

// fakeCheckModel substitutes checkModelFn in tests so no real HTTP call to
// Ollama is made when the user types an arbitrary tag. Returns a restore
// func to defer.
func fakeCheckModel(result preflight.Check) func() {
	orig := checkModelFn
	checkModelFn = func(string, string) preflight.Check { return result }
	return func() { checkModelFn = orig }
}

func loadedModelPicker(models []string) modelPickerModel {
	m := newModelPickerModel("http://localhost:11434", "qwen2.5-coder:7b")
	newM, _ := m.Update(modelsLoadedMsg{models: models})
	return newM.(modelPickerModel)
}

func TestModelPickerModel_InitQueriesConfiguredHost(t *testing.T) {
	defer fakeListModels([]string{"a", "b"}, nil)()

	m := newModelPickerModel("http://localhost:11434", "a")
	msg := m.Init()()
	loaded, ok := msg.(modelsLoadedMsg)
	if !ok {
		t.Fatalf("expected modelsLoadedMsg, got %T", msg)
	}
	if len(loaded.models) != 2 {
		t.Errorf("models = %v, want 2 entries", loaded.models)
	}
}

func TestModelPickerModel_LoadedModelsPopulateFilteredList(t *testing.T) {
	m := loadedModelPicker([]string{"qwen2.5-coder:7b", "llama3:8b"})
	if m.loading {
		t.Error("expected loading=false once models are loaded")
	}
	if len(m.filtered) != 2 {
		t.Errorf("filtered = %v, want 2 entries", m.filtered)
	}
}

func TestModelPickerModel_TypingFiltersListCaseInsensitively(t *testing.T) {
	m := loadedModelPicker([]string{"qwen2.5-coder:7b", "llama3:8b", "qwen2.5-coder:1.5b"})

	for _, r := range "QWEN" {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = newM.(modelPickerModel)
	}

	if len(m.filtered) != 2 {
		t.Fatalf("filtered = %v, want the 2 qwen models", m.filtered)
	}
	for _, name := range m.filtered {
		if name != "qwen2.5-coder:7b" && name != "qwen2.5-coder:1.5b" {
			t.Errorf("unexpected model in filtered list: %q", name)
		}
	}
}

func TestModelPickerModel_BackspaceRemovesLastFilterChar(t *testing.T) {
	m := loadedModelPicker([]string{"llama3:8b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("llama")})
	m = newM.(modelPickerModel)
	if m.filter != "llama" {
		t.Fatalf("filter = %q, want %q", m.filter, "llama")
	}

	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = newM2.(modelPickerModel)
	if m.filter != "llam" {
		t.Errorf("filter = %q, want %q after backspace", m.filter, "llam")
	}
}

func TestModelPickerModel_EnterSelectsHighlightedLocalModel(t *testing.T) {
	m := loadedModelPicker([]string{"qwen2.5-coder:7b", "llama3:8b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newM.(modelPickerModel)

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter to return a quit command")
	}
	mm := newM2.(modelPickerModel)
	if mm.selected == nil || *mm.selected != "llama3:8b" {
		t.Errorf("expected llama3:8b selected, got %+v", mm.selected)
	}
}

func TestModelPickerModel_QWithNoFilterBacksOutWithoutSelecting(t *testing.T) {
	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command")
	}
	mm := newM.(modelPickerModel)
	if !mm.back {
		t.Error("expected back=true")
	}
	if mm.selected != nil {
		t.Error("expected no selection when backing out")
	}
}

func TestModelPickerModel_EscAlwaysBacksOut(t *testing.T) {
	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}) // typed as filter char first
	m = newM.(modelPickerModel)
	// "q" alone with an empty prior filter backs out (tested above); once
	// something is typed, esc is the reliable way back regardless of
	// filter content.
	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected esc to return a quit command")
	}
	mm := newM2.(modelPickerModel)
	if !mm.back {
		t.Error("expected back=true")
	}
}

func TestModelPickerModel_TypingArbitraryTagNotPulledShowsWarningWithoutSelecting(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: "custom:tag not pulled — ollama pull custom:tag"})()

	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("custom:tag")})
	m = newM.(modelPickerModel)
	if len(m.filtered) != 0 {
		t.Fatalf("expected no local matches for custom:tag, got %v", m.filtered)
	}

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	mm := newM2.(modelPickerModel)
	if cmd != nil {
		t.Fatal("expected enter on an unpulled tag to NOT quit, just warn")
	}
	if mm.selected != nil {
		t.Error("expected no selection while the not-pulled warning is showing")
	}
	if mm.warning == "" {
		t.Error("expected a non-empty warning message")
	}
}

func TestModelPickerModel_SecondEnterAfterWarningConfirmsTagAnyway(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: "custom:tag not pulled — ollama pull custom:tag"})()

	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("custom:tag")})
	m = newM.(modelPickerModel)
	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM2.(modelPickerModel)

	newM3, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected the second enter to confirm and quit")
	}
	mm := newM3.(modelPickerModel)
	if mm.selected == nil || *mm.selected != "custom:tag" {
		t.Errorf("expected custom:tag selected after confirming, got %+v", mm.selected)
	}
}

func TestModelPickerModel_TypingArbitraryTagThatIsPulledSelectsDirectlyOnEnter(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: true, Detail: "pulled:tag ready"})()

	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("pulled:tag")})
	m = newM.(modelPickerModel)

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter to quit once CheckModel confirms the tag is pulled")
	}
	mm := newM2.(modelPickerModel)
	if mm.selected == nil || *mm.selected != "pulled:tag" {
		t.Errorf("expected pulled:tag selected, got %+v", mm.selected)
	}
}

func TestFilterModels_EmptyFilterReturnsAllModels(t *testing.T) {
	got := filterModels([]string{"a", "b"}, "")
	if len(got) != 2 {
		t.Errorf("filterModels with empty filter = %v, want all models", got)
	}
}

func TestFilterModels_NoMatchesReturnsEmpty(t *testing.T) {
	got := filterModels([]string{"a", "b"}, "zzz")
	if len(got) != 0 {
		t.Errorf("filterModels = %v, want empty", got)
	}
}
