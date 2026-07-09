package tui

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/preflight"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

var ansiPatternTUI = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripAnsiTUI(s string) string {
	return ansiPatternTUI.ReplaceAllString(s, "")
}

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

// fakeCatalogList substitutes catalogListFn in tests so no real exercises
// dir / sqlite db is touched — same indirection pattern as
// listModelsFn/checkModelFn.
func fakeCatalogList(statuses []catalog.ExerciseStatus, err error) func() {
	orig := catalogListFn
	catalogListFn = func(config.Config) ([]catalog.ExerciseStatus, error) { return statuses, err }
	return func() { catalogListFn = orig }
}

func fakeRecentAttempts(attempts []tracker.Attempt, err error) func() {
	orig := recentAttemptsFn
	recentAttemptsFn = func(config.Config, int) ([]tracker.Attempt, error) { return attempts, err }
	return func() { recentAttemptsFn = orig }
}

func practiceFixture() []catalog.ExerciseStatus {
	return []catalog.ExerciseStatus{
		fakeStatusIn("dsa", "two-pointers-01"),
		fakeStatusIn("debug", "off-by-one-01"),
	}
}

// --- stageMain ---

func TestAppModel_CursorStartsAtZero(t *testing.T) {
	m := appModel{}
	if m.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m.cursor)
	}
	if m.stage != stageMain {
		t.Errorf("stage = %v, want stageMain", m.stage)
	}
}

func TestAppModel_Main_UpStaysAtTop(t *testing.T) {
	m := appModel{}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if newM.(appModel).cursor != 0 {
		t.Error("cursor should stay at 0 when already at the top")
	}
}

func TestAppModel_Main_DownStopsAtLastOption(t *testing.T) {
	m := appModel{}
	for i := 0; i < 10; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(appModel)
	}
	if m.cursor != len(menuLabels)-1 {
		t.Errorf("cursor = %d, want %d (last option)", m.cursor, len(menuLabels)-1)
	}
}

func TestAppModel_Main_NumberKeysJumpDirectly(t *testing.T) {
	m := appModel{}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("3")})
	if newM.(appModel).cursor != 2 {
		t.Errorf("pressing 3 should jump cursor to index 2, got %d", newM.(appModel).cursor)
	}
}

func TestAppModel_Main_QQuitsTheWholeProgram(t *testing.T) {
	m := appModel{}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command")
	}
	if !newM.(appModel).quit {
		t.Error("expected quit=true")
	}
}

func TestAppModel_Main_TickAdvancesPhaseAndReschedules(t *testing.T) {
	m := appModel{}
	newM, cmd := m.Update(tickMsg{})
	if cmd == nil {
		t.Fatal("expected tick to reschedule another tick command")
	}
	if newM.(appModel).phase != 1 {
		t.Errorf("phase = %d, want 1", newM.(appModel).phase)
	}
}

// --- stageMain -> stageCategories (Practice) ---

func TestAppModel_EnterOnPractice_LoadsCategoriesInline(t *testing.T) {
	defer fakeCatalogList(practiceFixture(), nil)()

	m := appModel{cursor: 0}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command — this should stay inside the same program")
	}
	got := newM.(appModel)
	if got.stage != stageCategories {
		t.Fatalf("stage = %v, want stageCategories", got.stage)
	}
	if len(got.categories) != 2 {
		t.Fatalf("categories = %v, want 2 (dsa, debug)", got.categories)
	}
	if got.categories[0] != "dsa" {
		t.Errorf("categories[0] = %q, want %q (dsa sorts first, matching categoryOrder)", got.categories[0], "dsa")
	}
}

func TestAppModel_EnterOnPractice_CatalogErrorStaysOnMainAndSetsErr(t *testing.T) {
	defer fakeCatalogList(nil, errors.New("boom"))()

	m := appModel{cursor: 0}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := newM.(appModel)
	if got.stage != stageMain {
		t.Errorf("stage = %v, want stageMain (stay put on error)", got.stage)
	}
	if got.err == nil {
		t.Error("expected err to be set, not silently swallowed")
	}
}

// --- stageMain -> outcome (Sandbox) ---

func TestAppModel_EnterOnSandbox_SetsOutcomeAndQuits(t *testing.T) {
	m := appModel{cursor: 1}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter on Sandbox to return a quit command — it hands the terminal to docker")
	}
	got := newM.(appModel)
	if got.outcome != outcomeRunSandbox {
		t.Errorf("outcome = %v, want outcomeRunSandbox", got.outcome)
	}
}

// --- stageMain -> stageStats ---

func TestAppModel_EnterOnStats_LoadsStatsInline(t *testing.T) {
	defer fakeCatalogList(practiceFixture(), nil)()
	recent := []tracker.Attempt{{ExerciseID: "two-pointers-01", Result: tracker.ResultPass}}
	defer fakeRecentAttempts(recent, nil)()

	m := appModel{cursor: 2}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command for Stats")
	}
	got := newM.(appModel)
	if got.stage != stageStats {
		t.Fatalf("stage = %v, want stageStats", got.stage)
	}
	if len(got.statsRecent) != 1 {
		t.Errorf("statsRecent = %v, want 1 entry", got.statsRecent)
	}
}

func TestAppModel_StatsAnyKeyGoesBackToMain(t *testing.T) {
	m := appModel{stage: stageStats}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if cmd != nil {
		t.Error("expected no external command — back to main is an internal stage change")
	}
	if newM.(appModel).stage != stageMain {
		t.Error("expected any key to return to stageMain")
	}
}

// --- stageMain -> stageModelPicker ---

func TestAppModel_EnterOnModel_KicksOffAsyncLoad(t *testing.T) {
	defer fakeListModels([]string{"a", "b"}, nil)()

	m := appModel{cursor: 3}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter on Model to kick off an async models-loading command")
	}
	got := newM.(appModel)
	if got.stage != stageModelPicker {
		t.Fatalf("stage = %v, want stageModelPicker", got.stage)
	}
	if !got.modelLoading {
		t.Error("expected modelLoading=true while the command is in flight")
	}

	msg := cmd()
	loaded, ok := msg.(modelsLoadedMsg)
	if !ok {
		t.Fatalf("expected modelsLoadedMsg, got %T", msg)
	}
	newM2, _ := got.Update(loaded)
	got2 := newM2.(appModel)
	if got2.modelLoading {
		t.Error("expected modelLoading=false once loaded")
	}
	// 2 local ("a", "b") + 2 suggested (DeepSeek-Coder-V2-Lite,
	// Qwen2.5-Coder-14B) — see the SuggestedModels* tests below for the
	// discoverability behavior itself.
	if len(got2.modelFiltered) != 4 {
		t.Errorf("modelFiltered = %v, want 4 entries (2 local + 2 suggested)", got2.modelFiltered)
	}
}

// --- stageCategories ---

func categoriesFixture(t *testing.T) appModel {
	t.Helper()
	defer fakeCatalogList(practiceFixture(), nil)()
	m := appModel{cursor: 0}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	return newM.(appModel)
}

func TestAppModel_Categories_DownStopsAtLast(t *testing.T) {
	m := categoriesFixture(t)
	for i := 0; i < 10; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(appModel)
	}
	if m.categoryCursor != len(m.categories)-1 {
		t.Errorf("categoryCursor = %d, want %d", m.categoryCursor, len(m.categories)-1)
	}
}

func TestAppModel_Categories_EnterFiltersProblemsAndAdvances(t *testing.T) {
	m := categoriesFixture(t)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // dsa is first
	if cmd != nil {
		t.Error("expected no external command")
	}
	got := newM.(appModel)
	if got.stage != stageProblems {
		t.Fatalf("stage = %v, want stageProblems", got.stage)
	}
	if len(got.categoryProblems) != 1 || got.categoryProblems[0].ProblemID != "two-pointers-01" {
		t.Errorf("categoryProblems = %+v, want just two-pointers-01", got.categoryProblems)
	}
}

func TestAppModel_Categories_QGoesBackToMain(t *testing.T) {
	m := categoriesFixture(t)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd != nil {
		t.Error("expected no external command")
	}
	if newM.(appModel).stage != stageMain {
		t.Error("expected q to return to stageMain")
	}
}

func TestAppModel_RenderCategories_ShowsDSAUppercase(t *testing.T) {
	m := categoriesFixture(t)
	out := stripAnsiTUI(m.View())
	if !strings.Contains(out, "DSA") {
		t.Errorf("expected category list to show %q, got:\n%s", "DSA", out)
	}
	if strings.Contains(out, "dsa ") {
		t.Errorf("expected the raw id not to leak into the view:\n%s", out)
	}
}

// --- stageProblems ---

func problemsFixture(t *testing.T) appModel {
	t.Helper()
	m := categoriesFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	return newM.(appModel)
}

func TestAppModel_Problems_EnterAdvancesToLanguage(t *testing.T) {
	m := problemsFixture(t)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command")
	}
	got := newM.(appModel)
	if got.stage != stageLanguage {
		t.Fatalf("stage = %v, want stageLanguage", got.stage)
	}
	if got.selectedProblem.ProblemID != "two-pointers-01" {
		t.Errorf("selectedProblem = %+v, want two-pointers-01", got.selectedProblem)
	}
}

func TestAppModel_Problems_QGoesBackToCategories(t *testing.T) {
	m := problemsFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if newM.(appModel).stage != stageCategories {
		t.Error("expected q to return to stageCategories")
	}
}

func TestAppModel_RenderProblems_ShowsDisplayCategoryAsHeader(t *testing.T) {
	m := problemsFixture(t)
	out := stripAnsiTUI(m.View())
	if !strings.Contains(out, "DSA") {
		t.Errorf("expected the problems header to show %q, got:\n%s", "DSA", out)
	}
}

// --- stageLanguage ---

func languageFixture(t *testing.T) appModel {
	t.Helper()
	m := problemsFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	return newM.(appModel)
}

func TestAppModel_Language_EnterSetsOutcomeAndQuits(t *testing.T) {
	m := languageFixture(t)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter on a language variant to quit — it hands the terminal to docker")
	}
	got := newM.(appModel)
	if got.outcome != outcomeRunExercise {
		t.Errorf("outcome = %v, want outcomeRunExercise", got.outcome)
	}
	if got.exerciseToRun.ID != "two-pointers-01" {
		t.Errorf("exerciseToRun.ID = %q, want %q", got.exerciseToRun.ID, "two-pointers-01")
	}
}

func TestAppModel_Language_QGoesBackToProblems(t *testing.T) {
	m := languageFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if newM.(appModel).stage != stageProblems {
		t.Error("expected q to return to stageProblems")
	}
}

// --- stageModelPicker ---

// modelPickerFixture gives cfg a temp DataDir — selecting a model calls
// config.SaveSettings for real (see appModel.selectModel), so every test
// built on this fixture needs an isolated settings.json rather than
// writing into the repo's actual working directory.
func modelPickerFixture(t *testing.T, models []string) appModel {
	t.Helper()
	cfg := config.Config{DataDir: t.TempDir()}
	m := appModel{cfg: cfg, stage: stageModelPicker, modelLoading: true}
	newM, _ := m.Update(modelsLoadedMsg{models: models})
	return newM.(appModel)
}

func TestAppModel_ModelPicker_TypingFilters(t *testing.T) {
	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b", "llama3:8b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("llama")})
	got := newM.(appModel)
	if len(got.modelFiltered) != 1 || got.modelFiltered[0] != "llama3:8b" {
		t.Errorf("modelFiltered = %v, want [llama3:8b]", got.modelFiltered)
	}
}

func TestAppModel_ModelPicker_QWithNoFilterGoesBackToMain(t *testing.T) {
	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd != nil {
		t.Error("expected no external command — back to main is an internal stage change")
	}
	if newM.(appModel).stage != stageMain {
		t.Error("expected q with no filter to return to stageMain")
	}
}

func TestAppModel_ModelPicker_EnterOnLocalModelPersistsAndReturnsToMain(t *testing.T) {
	dir := t.TempDir()
	cfg := config.Config{DataDir: dir, TutorModel: "old-model"}
	m := appModel{cfg: cfg, stage: stageModelPicker, modelLoading: true}
	newM, _ := m.Update(modelsLoadedMsg{models: []string{"qwen2.5-coder:7b", "llama3:8b"}})
	m = newM.(appModel)

	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newM2.(appModel)

	newM3, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command — selecting a model doesn't leave the program")
	}
	got := newM3.(appModel)
	if got.stage != stageMain {
		t.Fatalf("stage = %v, want stageMain", got.stage)
	}
	if got.cfg.TutorModel != "llama3:8b" {
		t.Errorf("cfg.TutorModel = %q, want %q", got.cfg.TutorModel, "llama3:8b")
	}

	saved, err := config.LoadSettings(cfg.SettingsPath())
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if saved.TutorModel != "llama3:8b" {
		t.Errorf("persisted TutorModel = %q, want %q", saved.TutorModel, "llama3:8b")
	}
}

func TestAppModel_ModelPicker_TypingArbitraryTagNotPulledShowsWarningWithoutSelecting(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: "custom:tag not pulled"})()

	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("custom:tag")})
	m = newM.(appModel)

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected enter on an unpulled tag to not quit or leave the picker")
	}
	got := newM2.(appModel)
	if got.stage != stageModelPicker {
		t.Error("expected to stay in the model picker while the warning shows")
	}
	if got.modelWarning == "" {
		t.Error("expected a non-empty warning message")
	}
}

func TestAppModel_ModelPicker_SuggestedModelsAppearEvenWhenNotPulledLocally(t *testing.T) {
	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})

	foundDeepSeek, foundQwen14B := false, false
	for _, name := range m.modelFiltered {
		if name == config.DeepSeekCoderV2LiteModel {
			foundDeepSeek = true
		}
		if name == config.Qwen25Coder14BModel {
			foundQwen14B = true
		}
	}
	if !foundDeepSeek {
		t.Errorf("expected %s to be listed even though it isn't pulled, filtered = %v", config.DeepSeekCoderV2LiteModel, m.modelFiltered)
	}
	if !foundQwen14B {
		t.Errorf("expected %s to be listed even though it isn't pulled, filtered = %v", config.Qwen25Coder14BModel, m.modelFiltered)
	}
}

func TestAppModel_ModelPicker_SuggestedModelAlreadyPulledDoesNotAppearTwice(t *testing.T) {
	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b", config.DeepSeekCoderV2LiteModel})

	count := 0
	for _, name := range m.modelFiltered {
		if name == config.DeepSeekCoderV2LiteModel {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected %s to appear exactly once, appeared %d times in %v", config.DeepSeekCoderV2LiteModel, count, m.modelFiltered)
	}
}

func TestAppModel_RenderModelPicker_MarksSuggestedNotPulledModelsDistinctly(t *testing.T) {
	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, config.DeepSeekCoderV2LiteModel) {
		t.Fatalf("expected view to list %s, got:\n%s", config.DeepSeekCoderV2LiteModel, view)
	}
	if !strings.Contains(view, "not pulled") {
		t.Errorf("expected a not-pulled marker in the view, got:\n%s", view)
	}
}

func TestAppModel_ModelPicker_SelectingSuggestedNotPulledModelWarnsFirst(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: config.DeepSeekCoderV2LiteModel + " not pulled — ollama pull " + config.DeepSeekCoderV2LiteModel})()

	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})
	idx := -1
	for i, name := range m.modelFiltered {
		if name == config.DeepSeekCoderV2LiteModel {
			idx = i
		}
	}
	if idx < 0 {
		t.Fatalf("expected %s in filtered list, got %v", config.DeepSeekCoderV2LiteModel, m.modelFiltered)
	}
	for i := 0; i < idx; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(appModel)
	}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := newM.(appModel)
	if cmd != nil {
		t.Fatal("expected enter on a suggested-but-unpulled model to NOT quit, just warn")
	}
	if got.stage != stageModelPicker {
		t.Error("expected to stay in the model picker while the warning shows")
	}
	if got.modelWarning == "" {
		t.Error("expected a non-empty warning message")
	}
}

func TestAppModel_ModelPicker_SecondEnterOnSuggestedNotPulledModelConfirms(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: config.Qwen25Coder14BModel + " not pulled — ollama pull " + config.Qwen25Coder14BModel})()

	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})
	idx := -1
	for i, name := range m.modelFiltered {
		if name == config.Qwen25Coder14BModel {
			idx = i
		}
	}
	if idx < 0 {
		t.Fatalf("expected %s in filtered list, got %v", config.Qwen25Coder14BModel, m.modelFiltered)
	}
	for i := 0; i < idx; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(appModel)
	}

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(appModel)

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected the second enter to confirm without leaving the program (selecting a model is an internal stage change)")
	}
	got := newM2.(appModel)
	if got.stage != stageMain {
		t.Errorf("stage = %v, want stageMain", got.stage)
	}
	if got.cfg.TutorModel != config.Qwen25Coder14BModel {
		t.Errorf("cfg.TutorModel = %q, want %q after confirming", got.cfg.TutorModel, config.Qwen25Coder14BModel)
	}
}

func TestAppModel_ModelPicker_SelectingLocalModelNeverCallsCheckModel(t *testing.T) {
	// No fakeCheckModel set up here on purpose — if selecting an
	// already-pulled local model called checkModelFn at all, this test
	// would hit the real network-calling default and could hang/flake in
	// CI. Selecting a genuinely local entry must short-circuit before
	// ever reaching that call.
	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b", "llama3:8b"})
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected enter to select immediately for an already-pulled local model, no external command")
	}
	got := newM.(appModel)
	if got.stage != stageMain {
		t.Errorf("stage = %v, want stageMain", got.stage)
	}
	if got.cfg.TutorModel != "qwen2.5-coder:7b" {
		t.Errorf("cfg.TutorModel = %q, want %q", got.cfg.TutorModel, "qwen2.5-coder:7b")
	}
}

// --- newAppModel / resume ---

func TestNewAppModel_DefaultsToStageMain(t *testing.T) {
	m := newAppModel(config.Config{}, appResume{})
	if m.stage != stageMain {
		t.Errorf("stage = %v, want stageMain", m.stage)
	}
}

func TestNewAppModel_ResumeAtStageProblems_PreloadsCategory(t *testing.T) {
	defer fakeCatalogList(practiceFixture(), nil)()

	m := newAppModel(config.Config{}, appResume{stage: stageProblems, category: "dsa"})
	if m.stage != stageProblems {
		t.Fatalf("stage = %v, want stageProblems", m.stage)
	}
	if len(m.categoryProblems) != 1 || m.categoryProblems[0].ProblemID != "two-pointers-01" {
		t.Errorf("categoryProblems = %+v, want just two-pointers-01", m.categoryProblems)
	}
}
