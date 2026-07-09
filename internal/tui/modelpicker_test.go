package tui

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/config"
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
	// 2 local + 2 suggested (DeepSeek-Coder-V2-Lite-Instruct,
	// Qwen2.5-Coder-14B-Instruct) — see TestModelPickerModel_Suggested*
	// for the discoverability behavior itself.
	if len(m.filtered) != 4 {
		t.Errorf("filtered = %v, want 4 entries (2 local + 2 suggested)", m.filtered)
	}
}

func TestModelPickerModel_SuggestedModelsAppearEvenWhenNotPulledLocally(t *testing.T) {
	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})

	foundDeepSeek, foundQwen14B := false, false
	for _, name := range m.filtered {
		if name == config.DeepSeekCoderV2LiteModel {
			foundDeepSeek = true
		}
		if name == config.Qwen25Coder14BModel {
			foundQwen14B = true
		}
	}
	if !foundDeepSeek {
		t.Errorf("expected %s to be listed even though it isn't pulled, filtered = %v", config.DeepSeekCoderV2LiteModel, m.filtered)
	}
	if !foundQwen14B {
		t.Errorf("expected %s to be listed even though it isn't pulled, filtered = %v", config.Qwen25Coder14BModel, m.filtered)
	}
}

func TestModelPickerModel_SuggestedModelAlreadyPulledDoesNotAppearTwice(t *testing.T) {
	m := loadedModelPicker([]string{"qwen2.5-coder:7b", config.DeepSeekCoderV2LiteModel})

	count := 0
	for _, name := range m.filtered {
		if name == config.DeepSeekCoderV2LiteModel {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected %s to appear exactly once, appeared %d times in %v", config.DeepSeekCoderV2LiteModel, count, m.filtered)
	}
}

func TestModelPickerModel_ViewMarksSuggestedNotPulledModelsDistinctly(t *testing.T) {
	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, config.DeepSeekCoderV2LiteModel) {
		t.Fatalf("expected view to list %s, got:\n%s", config.DeepSeekCoderV2LiteModel, view)
	}
	if !strings.Contains(view, "not pulled") {
		t.Errorf("expected a not-pulled marker in the view, got:\n%s", view)
	}
}

func TestModelPickerModel_SelectingSuggestedNotPulledModelWarnsFirst(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: config.DeepSeekCoderV2LiteModel + " not pulled — ollama pull " + config.DeepSeekCoderV2LiteModel})()

	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	idx := -1
	for i, name := range m.filtered {
		if name == config.DeepSeekCoderV2LiteModel {
			idx = i
		}
	}
	if idx < 0 {
		t.Fatalf("expected %s in filtered list, got %v", config.DeepSeekCoderV2LiteModel, m.filtered)
	}
	for i := 0; i < idx; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(modelPickerModel)
	}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	mm := newM.(modelPickerModel)
	if cmd != nil {
		t.Fatal("expected enter on a suggested-but-unpulled model to NOT quit, just prompt")
	}
	if mm.selected != nil {
		t.Error("expected no selection while the download prompt is showing")
	}
	if mm.warning == "" {
		t.Error("expected a non-empty warning message")
	}
	if mm.pendingDownloadTag != config.DeepSeekCoderV2LiteModel {
		t.Errorf("pendingDownloadTag = %q, want %q", mm.pendingDownloadTag, config.DeepSeekCoderV2LiteModel)
	}
}

func TestModelPickerModel_YOnDownloadPromptStartsLivePullAndSelectsOnSuccess(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: config.Qwen25Coder14BModel + " not pulled — ollama pull " + config.Qwen25Coder14BModel})()

	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	defer fakePullModel(lineCh, errCh)()

	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	idx := -1
	for i, name := range m.filtered {
		if name == config.Qwen25Coder14BModel {
			idx = i
		}
	}
	if idx < 0 {
		t.Fatalf("expected %s in filtered list, got %v", config.Qwen25Coder14BModel, m.filtered)
	}
	for i := 0; i < idx; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(modelPickerModel)
	}

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(modelPickerModel)
	if m.pendingDownloadTag == "" {
		t.Fatal("expected a pending download prompt before answering y/n")
	}

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if cmd == nil {
		t.Fatal("expected 'y' to start waiting for pull output")
	}
	m = newM2.(modelPickerModel)
	if !m.downloading {
		t.Fatal("expected downloading=true after answering 'y'")
	}
	if m.pendingDownloadTag != "" {
		t.Error("expected the download prompt to clear once 'y' starts the pull")
	}

	newM3, cmd3 := m.Update(pullDoneMsg{err: nil})
	if cmd3 == nil {
		t.Fatal("expected quit once the pull succeeds")
	}
	mm := newM3.(modelPickerModel)
	if mm.selected == nil || *mm.selected != config.Qwen25Coder14BModel {
		t.Errorf("expected %s selected after a successful download, got %+v", config.Qwen25Coder14BModel, mm.selected)
	}
	if mm.downloading {
		t.Error("expected downloading=false once the pull resolves")
	}
}

func TestModelPickerModel_NOnDownloadPromptCancelsWithoutSelecting(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: "custom:tag not pulled — ollama pull custom:tag"})()

	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("custom:tag")})
	m = newM.(modelPickerModel)
	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM2.(modelPickerModel)
	if m.pendingDownloadTag == "" {
		t.Fatal("expected a pending download prompt")
	}

	newM3, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if cmd != nil {
		t.Fatal("expected 'n' to just cancel, not quit")
	}
	mm := newM3.(modelPickerModel)
	if mm.selected != nil {
		t.Error("expected no selection after declining the download")
	}
	if mm.pendingDownloadTag != "" {
		t.Error("expected the download prompt to clear after 'n'")
	}
	if mm.downloading {
		t.Error("expected no download to have started after 'n'")
	}
}

func TestModelPickerModel_PullFailureShowsWarningWithoutSelecting(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: "custom:tag not pulled — ollama pull custom:tag"})()

	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	defer fakePullModel(lineCh, errCh)()

	m := loadedModelPicker([]string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("custom:tag")})
	m = newM.(modelPickerModel)
	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM2.(modelPickerModel)
	newM3, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	m = newM3.(modelPickerModel)

	newM4, cmd := m.Update(pullDoneMsg{err: errors.New("boom")})
	if cmd != nil {
		t.Fatal("expected no quit command when the download fails")
	}
	mm := newM4.(modelPickerModel)
	if mm.selected != nil {
		t.Error("expected no selection when the download fails")
	}
	if mm.downloading {
		t.Error("expected downloading=false once the failed pull resolves")
	}
	if mm.warning == "" {
		t.Error("expected a non-empty warning explaining the failure")
	}
}

func TestModelPickerModel_DownloadLineCapsAtThreeLinesTotal(t *testing.T) {
	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	m := modelPickerModel{downloading: true, downloadLineCh: lineCh, downloadErrCh: errCh}

	for i := 0; i < maxOutputLines+5; i++ {
		newM, cmd := m.Update(pullLineMsg(fmt.Sprintf("line %d", i)))
		if cmd == nil {
			t.Fatal("expected pullLineMsg to keep listening for more output")
		}
		m = newM.(modelPickerModel)
	}
	if len(m.downloadLines) != maxOutputLines {
		t.Errorf("downloadLines = %d, want capped at %d", len(m.downloadLines), maxOutputLines)
	}
}

func TestModelPickerModel_ViewShowsDownloadPromptAndLiveOutput(t *testing.T) {
	m := modelPickerModel{
		pendingDownloadTag: "custom:tag",
		warning:            "custom:tag not pulled — ollama pull custom:tag",
	}
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, "download custom:tag? (y/n)") {
		t.Errorf("expected a y/n download prompt, got:\n%s", view)
	}

	m2 := modelPickerModel{
		downloading:    true,
		downloadTarget: "custom:tag",
		downloadLines:  []string{"pulling manifest", "downloading (42%)"},
	}
	view2 := stripAnsiTUI(m2.View())
	if !strings.Contains(view2, "downloading custom:tag") {
		t.Errorf("expected a downloading-in-progress notice, got:\n%s", view2)
	}
	if !strings.Contains(view2, "pulling manifest") || !strings.Contains(view2, "downloading (42%)") {
		t.Errorf("expected live pull output visible, got:\n%s", view2)
	}
}

func TestModelPickerModel_SelectingLocalModelNeverCallsCheckModel(t *testing.T) {
	// No fakeCheckModel set up here on purpose — if selecting an
	// already-pulled local model called checkModelFn at all, this test
	// would hit the real network-calling default and could hang/flake in
	// CI. Selecting a genuinely local entry must short-circuit before
	// ever reaching that call.
	m := loadedModelPicker([]string{"qwen2.5-coder:7b", "llama3:8b"})
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter to quit immediately for an already-pulled local model")
	}
	mm := newM.(modelPickerModel)
	if mm.selected == nil || *mm.selected != "qwen2.5-coder:7b" {
		t.Errorf("expected qwen2.5-coder:7b selected, got %+v", mm.selected)
	}
}

func TestModelPickerModel_TypingFiltersListCaseInsensitively(t *testing.T) {
	m := loadedModelPicker([]string{"qwen2.5-coder:7b", "llama3:8b", "qwen2.5-coder:1.5b"})

	for _, r := range "QWEN" {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = newM.(modelPickerModel)
	}

	// 2 local qwen models + the suggested Qwen2.5-Coder-14B-Instruct,
	// which also legitimately matches "qwen" and should surface here too.
	if len(m.filtered) != 3 {
		t.Fatalf("filtered = %v, want the 2 local qwen models plus the suggested 14B one", m.filtered)
	}
	for _, name := range m.filtered {
		if name != "qwen2.5-coder:7b" && name != "qwen2.5-coder:1.5b" && name != config.Qwen25Coder14BModel {
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
		t.Fatal("expected enter on an unpulled tag to NOT quit, just prompt")
	}
	if mm.selected != nil {
		t.Error("expected no selection while the download prompt is showing")
	}
	if mm.warning == "" {
		t.Error("expected a non-empty warning message")
	}
	if mm.pendingDownloadTag != "custom:tag" {
		t.Errorf("pendingDownloadTag = %q, want %q", mm.pendingDownloadTag, "custom:tag")
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
