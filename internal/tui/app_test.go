package tui

import (
	"errors"
	"fmt"
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

// fakeCheckToolCalling substitutes checkToolCallingFn in tests so
// selectModel's background check never makes a real Ollama round-trip.
// Returns a restore func to defer.
func fakeCheckToolCalling(supported bool, err error) func() {
	orig := checkToolCallingFn
	checkToolCallingFn = func(string, string, string) (bool, error) { return supported, err }
	return func() { checkToolCallingFn = orig }
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
		fakeStatusIn("two-pointers", "two-pointers-01"),
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

func TestAppModel_Categories_EnterOnDSAGoesToDSASubcategories(t *testing.T) {
	m := categoriesFixture(t)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // dsa is first
	if cmd != nil {
		t.Error("expected no external command")
	}
	got := newM.(appModel)
	if got.stage != stageDSACategories {
		t.Fatalf("stage = %v, want stageDSACategories", got.stage)
	}
	if len(got.dsaCategories) != 1 || got.dsaCategories[0] != "two-pointers" {
		t.Errorf("dsaCategories = %v, want [two-pointers]", got.dsaCategories)
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

// --- stageDSACategories ---

func dsaCategoriesFixture(t *testing.T) appModel {
	t.Helper()
	m := categoriesFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // dsa -> stageDSACategories
	return newM.(appModel)
}

func TestAppModel_DSACategories_DownStopsAtLast(t *testing.T) {
	m := dsaCategoriesFixture(t)
	for i := 0; i < 10; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(appModel)
	}
	if m.dsaCategoryCursor != len(m.dsaCategories)-1 {
		t.Errorf("dsaCategoryCursor = %d, want %d", m.dsaCategoryCursor, len(m.dsaCategories)-1)
	}
}

func TestAppModel_DSACategories_EnterFiltersProblemsAndAdvances(t *testing.T) {
	m := dsaCategoriesFixture(t)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // two-pointers is the only subcategory
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

func TestAppModel_DSACategories_QGoesBackToCategories(t *testing.T) {
	m := dsaCategoriesFixture(t)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd != nil {
		t.Error("expected no external command")
	}
	if newM.(appModel).stage != stageCategories {
		t.Error("expected q to return to stageCategories")
	}
}

func TestAppModel_RenderDSACategories_ShowsSubcategoryName(t *testing.T) {
	m := dsaCategoriesFixture(t)
	out := stripAnsiTUI(m.View())
	if !strings.Contains(out, "Two Pointers") {
		t.Errorf("expected the subcategory list to show %q, got:\n%s", "Two Pointers", out)
	}
}

// --- stageProblems ---

func problemsFixture(t *testing.T) appModel {
	t.Helper()
	m := dsaCategoriesFixture(t)
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

func TestAppModel_Problems_QGoesBackToDSACategoriesWhenGrouped(t *testing.T) {
	m := problemsFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if newM.(appModel).stage != stageDSACategories {
		t.Error("expected q to return to stageDSACategories for a grouped (NeetCode) category")
	}
}

func TestAppModel_Problems_QGoesBackToCategoriesWhenUngrouped(t *testing.T) {
	m := categoriesFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown}) // move to "debug" (ungrouped)
	m = newM.(appModel)
	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // debug -> stageProblems directly
	m = newM2.(appModel)
	if m.stage != stageProblems {
		t.Fatalf("stage = %v, want stageProblems", m.stage)
	}

	newM3, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if newM3.(appModel).stage != stageCategories {
		t.Error("expected q to return directly to stageCategories for an ungrouped category")
	}
}

func TestAppModel_RenderProblems_ShowsDisplayCategoryAsHeader(t *testing.T) {
	m := problemsFixture(t)
	out := stripAnsiTUI(m.View())
	if !strings.Contains(out, "Two Pointers") {
		t.Errorf("expected the problems header to show %q, got:\n%s", "Two Pointers", out)
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
	defer fakeCheckToolCalling(true, nil)()

	dir := t.TempDir()
	cfg := config.Config{DataDir: dir, TutorModel: "old-model"}
	m := appModel{cfg: cfg, stage: stageModelPicker, modelLoading: true}
	newM, _ := m.Update(modelsLoadedMsg{models: []string{"qwen2.5-coder:7b", "llama3:8b"}})
	m = newM.(appModel)

	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newM2.(appModel)

	newM3, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected a background command kicking off the tool-calling check")
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
	if got.modelPendingDownloadTag != "custom:tag" {
		t.Errorf("modelPendingDownloadTag = %q, want %q", got.modelPendingDownloadTag, "custom:tag")
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
		t.Fatal("expected enter on a suggested-but-unpulled model to NOT quit, just prompt")
	}
	if got.stage != stageModelPicker {
		t.Error("expected to stay in the model picker while the download prompt shows")
	}
	if got.modelWarning == "" {
		t.Error("expected a non-empty warning message")
	}
	if got.modelPendingDownloadTag != config.DeepSeekCoderV2LiteModel {
		t.Errorf("modelPendingDownloadTag = %q, want %q", got.modelPendingDownloadTag, config.DeepSeekCoderV2LiteModel)
	}
}

func TestAppModel_ModelPicker_YOnDownloadPromptStartsLivePullAndSelectsOnSuccess(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: config.Qwen25Coder14BModel + " not pulled — ollama pull " + config.Qwen25Coder14BModel})()
	defer fakeCheckToolCalling(true, nil)()

	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	defer fakePullModel(lineCh, errCh)()

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
	if m.modelPendingDownloadTag == "" {
		t.Fatal("expected a pending download prompt before answering y/n")
	}

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	if cmd == nil {
		t.Fatal("expected 'y' to start waiting for pull output")
	}
	m = newM2.(appModel)
	if !m.modelDownloading {
		t.Fatal("expected modelDownloading=true after answering 'y'")
	}
	if m.modelPendingDownloadTag != "" {
		t.Error("expected the download prompt to clear once 'y' starts the pull")
	}

	newM3, cmd3 := m.Update(pullDoneMsg{err: nil})
	if cmd3 == nil {
		t.Error("expected a background command kicking off the tool-calling check once the pull succeeds")
	}
	got := newM3.(appModel)
	if got.stage != stageMain {
		t.Errorf("stage = %v, want stageMain", got.stage)
	}
	if got.cfg.TutorModel != config.Qwen25Coder14BModel {
		t.Errorf("cfg.TutorModel = %q, want %q after a successful download", got.cfg.TutorModel, config.Qwen25Coder14BModel)
	}
	if got.modelDownloading {
		t.Error("expected modelDownloading=false once the pull resolves")
	}
}

func TestAppModel_ModelPicker_NOnDownloadPromptCancelsWithoutSelecting(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: "custom:tag not pulled — ollama pull custom:tag"})()

	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("custom:tag")})
	m = newM.(appModel)
	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM2.(appModel)
	if m.modelPendingDownloadTag == "" {
		t.Fatal("expected a pending download prompt")
	}

	newM3, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if cmd != nil {
		t.Fatal("expected 'n' to just cancel, not quit or leave the picker")
	}
	got := newM3.(appModel)
	if got.stage != stageModelPicker {
		t.Error("expected to stay in the model picker after declining the download")
	}
	if got.cfg.TutorModel != "" {
		t.Error("expected no selection after declining the download")
	}
	if got.modelPendingDownloadTag != "" {
		t.Error("expected the download prompt to clear after 'n'")
	}
	if got.modelDownloading {
		t.Error("expected no download to have started after 'n'")
	}
}

func TestAppModel_ModelPicker_PullFailureShowsWarningWithoutSelecting(t *testing.T) {
	defer fakeCheckModel(preflight.Check{OK: false, Detail: "custom:tag not pulled — ollama pull custom:tag"})()

	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	defer fakePullModel(lineCh, errCh)()

	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("custom:tag")})
	m = newM.(appModel)
	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM2.(appModel)
	newM3, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})
	m = newM3.(appModel)

	newM4, cmd := m.Update(pullDoneMsg{err: errors.New("boom")})
	if cmd != nil {
		t.Fatal("expected no quit/external command when the download fails")
	}
	got := newM4.(appModel)
	if got.stage != stageModelPicker {
		t.Error("expected to stay in the model picker after a failed download")
	}
	if got.cfg.TutorModel != "" {
		t.Error("expected no selection when the download fails")
	}
	if got.modelDownloading {
		t.Error("expected modelDownloading=false once the failed pull resolves")
	}
	if got.modelWarning == "" {
		t.Error("expected a non-empty warning explaining the failure")
	}
}

func TestAppModel_ModelPicker_DownloadLineCapsAtThreeLinesTotal(t *testing.T) {
	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	m := appModel{stage: stageModelPicker, modelDownloading: true, modelDownloadLineCh: lineCh, modelDownloadErrCh: errCh}

	for i := 0; i < maxOutputLines+5; i++ {
		newM, cmd := m.Update(pullLineMsg(fmt.Sprintf("line %d", i)))
		if cmd == nil {
			t.Fatal("expected pullLineMsg to keep listening for more output")
		}
		m = newM.(appModel)
	}
	if len(m.modelDownloadLines) != maxOutputLines {
		t.Errorf("modelDownloadLines = %d, want capped at %d", len(m.modelDownloadLines), maxOutputLines)
	}
}

func TestAppModel_RenderModelPicker_ShowsDownloadPromptAndLiveOutput(t *testing.T) {
	m := appModel{
		stage:                   stageModelPicker,
		modelPendingDownloadTag: "custom:tag",
		modelWarning:            "custom:tag not pulled — ollama pull custom:tag",
	}
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, "download custom:tag? (y/n)") {
		t.Errorf("expected a y/n download prompt, got:\n%s", view)
	}

	m2 := appModel{
		stage:               stageModelPicker,
		modelDownloading:    true,
		modelDownloadTarget: "custom:tag",
		modelDownloadLines:  []string{"pulling manifest", "downloading (42%)"},
	}
	view2 := stripAnsiTUI(m2.View())
	if !strings.Contains(view2, "downloading custom:tag") {
		t.Errorf("expected a downloading-in-progress notice, got:\n%s", view2)
	}
	if !strings.Contains(view2, "pulling manifest") || !strings.Contains(view2, "downloading (42%)") {
		t.Errorf("expected live pull output visible, got:\n%s", view2)
	}
}

func TestAppModel_ModelPicker_SelectingLocalModelNeverCallsCheckModel(t *testing.T) {
	// No fakeCheckModel set up here on purpose — if selecting an
	// already-pulled local model called checkModelFn at all, this test
	// would hit the real network-calling default and could hang/flake in
	// CI. Selecting a genuinely local entry must short-circuit before
	// ever reaching that call. checkToolCallingFn is faked separately —
	// selecting a model always kicks that check off, local or not.
	defer fakeCheckToolCalling(true, nil)()

	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b", "llama3:8b"})
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected a background command kicking off the tool-calling check")
	}
	got := newM.(appModel)
	if got.stage != stageMain {
		t.Errorf("stage = %v, want stageMain", got.stage)
	}
	if got.cfg.TutorModel != "qwen2.5-coder:7b" {
		t.Errorf("cfg.TutorModel = %q, want %q", got.cfg.TutorModel, "qwen2.5-coder:7b")
	}
}

// --- OpenRouter model selection (handleModelEnter routing + key entry) ---

func TestAppModel_ModelPicker_TypingOpenRouterTagWithNoKeySetGoesToKeyEntry(t *testing.T) {
	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("openrouter:anthropic/claude-3.5-sonnet")})
	m = newM.(appModel)

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command — entering the key-entry stage is an internal stage change")
	}
	got := newM2.(appModel)
	if got.stage != stageOpenRouterKeyEntry {
		t.Errorf("stage = %v, want stageOpenRouterKeyEntry", got.stage)
	}
	if got.openRouterPendingModel != "openrouter:anthropic/claude-3.5-sonnet" {
		t.Errorf("openRouterPendingModel = %q, want %q", got.openRouterPendingModel, "openrouter:anthropic/claude-3.5-sonnet")
	}
	// checkModelFn must never be reached for an openrouter: tag — it's
	// meaningless against Ollama's /api/tags and would misreport it as
	// "not pulled".
}

func TestAppModel_ModelPicker_TypingOpenRouterTagWithKeyAlreadySetSelectsImmediately(t *testing.T) {
	defer fakeCheckToolCalling(true, nil)()

	dir := t.TempDir()
	m := appModel{cfg: config.Config{DataDir: dir, OpenRouterAPIKey: "sk-existing"}, stage: stageModelPicker, modelLoading: true}
	newM, _ := m.Update(modelsLoadedMsg{models: []string{"qwen2.5-coder:7b"}})
	m = newM.(appModel)
	newM2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("openrouter:anthropic/claude-3.5-sonnet")})
	m = newM2.(appModel)

	newM3, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected a background command kicking off the tool-calling check")
	}
	got := newM3.(appModel)
	if got.stage != stageMain {
		t.Errorf("stage = %v, want stageMain (key already available, should select immediately)", got.stage)
	}
	if got.cfg.TutorModel != "openrouter:anthropic/claude-3.5-sonnet" {
		t.Errorf("cfg.TutorModel = %q, want %q", got.cfg.TutorModel, "openrouter:anthropic/claude-3.5-sonnet")
	}
}

func TestAppModel_OpenRouterKeyEntry_EnterSavesKeyAndSelectsPendingModel(t *testing.T) {
	defer fakeCheckToolCalling(true, nil)()

	dir := t.TempDir()
	cfg := config.Config{DataDir: dir, TutorModel: "old-model"}
	m := appModel{cfg: cfg, stage: stageOpenRouterKeyEntry, openRouterPendingModel: "openrouter:anthropic/claude-3.5-sonnet"}

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("sk-typed-key")})
	m = newM.(appModel)

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected a background command kicking off the tool-calling check")
	}
	got := newM2.(appModel)
	if got.stage != stageMain {
		t.Fatalf("stage = %v, want stageMain", got.stage)
	}
	if got.cfg.OpenRouterAPIKey != "sk-typed-key" {
		t.Errorf("cfg.OpenRouterAPIKey = %q, want %q", got.cfg.OpenRouterAPIKey, "sk-typed-key")
	}
	if got.cfg.TutorModel != "openrouter:anthropic/claude-3.5-sonnet" {
		t.Errorf("cfg.TutorModel = %q, want %q", got.cfg.TutorModel, "openrouter:anthropic/claude-3.5-sonnet")
	}

	saved, err := config.LoadSettings(cfg.SettingsPath())
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if saved.OpenRouterAPIKey != "sk-typed-key" {
		t.Errorf("persisted OpenRouterAPIKey = %q, want %q", saved.OpenRouterAPIKey, "sk-typed-key")
	}
	if saved.TutorModel != "openrouter:anthropic/claude-3.5-sonnet" {
		t.Errorf("persisted TutorModel = %q, want %q", saved.TutorModel, "openrouter:anthropic/claude-3.5-sonnet")
	}
}

func TestAppModel_OpenRouterKeyEntry_EnterWithEmptyKeyDoesNothing(t *testing.T) {
	m := appModel{cfg: config.Config{DataDir: t.TempDir()}, stage: stageOpenRouterKeyEntry, openRouterPendingModel: "openrouter:x/y"}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no command when the key input is empty")
	}
	got := newM.(appModel)
	if got.stage != stageOpenRouterKeyEntry {
		t.Errorf("stage = %v, want to stay at stageOpenRouterKeyEntry with nothing typed", got.stage)
	}
}

func TestAppModel_OpenRouterKeyEntry_EscCancelsBackToModelPickerWithoutSelecting(t *testing.T) {
	m := appModel{cfg: config.Config{DataDir: t.TempDir(), TutorModel: "old-model"}, stage: stageOpenRouterKeyEntry, openRouterPendingModel: "openrouter:x/y", openRouterKeyInput: "sk-partial"}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Error("expected no external command on cancel")
	}
	got := newM.(appModel)
	if got.stage != stageModelPicker {
		t.Errorf("stage = %v, want stageModelPicker", got.stage)
	}
	if got.cfg.TutorModel != "old-model" {
		t.Errorf("cfg.TutorModel = %q, want unchanged %q", got.cfg.TutorModel, "old-model")
	}
	if got.cfg.OpenRouterAPIKey != "" {
		t.Error("expected the partially-typed key to be discarded, not saved")
	}
}

func TestAppModel_OpenRouterKeyEntry_BackspaceRemovesLastCharacter(t *testing.T) {
	m := appModel{cfg: config.Config{DataDir: t.TempDir()}, stage: stageOpenRouterKeyEntry, openRouterKeyInput: "sk-abc"}

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	got := newM.(appModel)
	if got.openRouterKeyInput != "sk-ab" {
		t.Errorf("openRouterKeyInput = %q, want %q", got.openRouterKeyInput, "sk-ab")
	}
}

func TestAppModel_SelectModel_PreservesOpenRouterAPIKeyAcrossAnOllamaPick(t *testing.T) {
	defer fakeCheckToolCalling(true, nil)()

	dir := t.TempDir()
	cfg := config.Config{DataDir: dir, OpenRouterAPIKey: "sk-preserve-me"}
	m := appModel{cfg: cfg}

	newM, _ := m.selectModel("llama3:8b")
	got := newM.(appModel)
	if got.cfg.OpenRouterAPIKey != "sk-preserve-me" {
		t.Errorf("cfg.OpenRouterAPIKey = %q, want it preserved across an unrelated Ollama pick", got.cfg.OpenRouterAPIKey)
	}

	saved, err := config.LoadSettings(cfg.SettingsPath())
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if saved.OpenRouterAPIKey != "sk-preserve-me" {
		t.Errorf("persisted OpenRouterAPIKey = %q, want it preserved, not wiped by picking an unrelated Ollama model", saved.OpenRouterAPIKey)
	}
}

func TestAppModel_RenderOpenRouterKeyEntry_MasksTypedKey(t *testing.T) {
	m := appModel{cfg: config.Config{}, stage: stageOpenRouterKeyEntry, openRouterPendingModel: "openrouter:x/y", openRouterKeyInput: "sk-secret-value"}
	view := stripAnsiTUI(m.View())
	if strings.Contains(view, "sk-secret-value") {
		t.Error("expected the typed API key to be masked, not shown in plain text")
	}
}

// --- tool-calling validation (selectModel's checkToolCallingCmd) ---

func TestAppModel_SelectModel_RunningItsCmdReportsUnsupportedToolCalling(t *testing.T) {
	defer fakeCheckToolCalling(false, nil)()

	m := appModel{cfg: config.Config{DataDir: t.TempDir()}}
	newM, cmd := m.selectModel("qwen2.5-coder:7b")
	if cmd == nil {
		t.Fatal("expected selectModel to return the tool-calling check command")
	}

	msg := cmd()
	checkMsg, ok := msg.(toolCallingCheckMsg)
	if !ok {
		t.Fatalf("cmd() returned %T, want toolCallingCheckMsg", msg)
	}
	if checkMsg.model != "qwen2.5-coder:7b" || checkMsg.supported {
		t.Errorf("toolCallingCheckMsg = %+v, want model=qwen2.5-coder:7b supported=false", checkMsg)
	}

	got, _ := newM.Update(checkMsg)
	final := got.(appModel)
	if final.toolCallingWarning == "" {
		t.Error("expected toolCallingWarning to be set for an unsupported model")
	}
	if !strings.Contains(final.toolCallingWarning, "qwen2.5-coder:7b") {
		t.Errorf("toolCallingWarning = %q, want it to name the model", final.toolCallingWarning)
	}
}

func TestAppModel_ToolCallingCheckMsg_SupportedLeavesNoWarning(t *testing.T) {
	m := appModel{cfg: config.Config{TutorModel: "llama3.1:8b"}}
	newM, cmd := m.Update(toolCallingCheckMsg{model: "llama3.1:8b", supported: true})
	if cmd != nil {
		t.Error("expected no further command once the check resolves")
	}
	got := newM.(appModel)
	if got.toolCallingWarning != "" {
		t.Errorf("toolCallingWarning = %q, want empty for a supported model", got.toolCallingWarning)
	}
}

func TestAppModel_ToolCallingCheckMsg_ErrorSetsWarningWithRealDetail(t *testing.T) {
	// A real bug found live: this used to leave toolCallingWarning empty
	// on any error, on the reasoning that a network blip checking a
	// perfectly fine model shouldn't scare the user into thinking that
	// model is broken. In practice this swallowed a much more
	// significant case: Ollama hard-rejecting a request with 400 "does
	// not support tools" for a model picked without real tool-calling
	// support IS an error here (CheckToolCalling returns it as such,
	// not just supported=false) — so the picker went completely silent,
	// and the problem only surfaced once inside a live tutor session as
	// a confusing "could not reach" message. An error must still warn,
	// just phrased as the check having failed rather than flatly
	// blaming the model, since it could still be a transient issue.
	m := appModel{cfg: config.Config{TutorModel: "llama3.1:8b"}}
	newM, _ := m.Update(toolCallingCheckMsg{model: "llama3.1:8b", supported: false, err: errors.New("does not support tools")})
	got := newM.(appModel)
	if got.toolCallingWarning == "" {
		t.Fatal("expected toolCallingWarning to be set when the check itself errored")
	}
	if !strings.Contains(got.toolCallingWarning, "llama3.1:8b") {
		t.Errorf("toolCallingWarning = %q, want it to name the model", got.toolCallingWarning)
	}
	if !strings.Contains(got.toolCallingWarning, "does not support tools") {
		t.Errorf("toolCallingWarning = %q, want the real underlying error detail included", got.toolCallingWarning)
	}
}

func TestAppModel_ToolCallingCheckMsg_StaleResultForReplacedModelIgnored(t *testing.T) {
	// The user picked "llama3.1:8b" again (or something else) before the
	// check for the earlier "qwen2.5-coder:7b" pick resolved — that
	// stale result must not set a warning about a model that isn't even
	// selected anymore.
	m := appModel{cfg: config.Config{TutorModel: "llama3.1:8b"}, toolCallingWarning: ""}
	newM, _ := m.Update(toolCallingCheckMsg{model: "qwen2.5-coder:7b", supported: false})
	got := newM.(appModel)
	if got.toolCallingWarning != "" {
		t.Errorf("toolCallingWarning = %q, want empty — result was for a model that is no longer selected", got.toolCallingWarning)
	}
}

func TestAppModel_RenderMain_ShowsToolCallingWarning(t *testing.T) {
	m := appModel{stage: stageMain, toolCallingWarning: "qwen2.5-coder:7b may not support real tool calling"}
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, "may not support real tool calling") {
		t.Errorf("expected the tool-calling warning in the main menu view, got:\n%s", view)
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

	m := newAppModel(config.Config{}, appResume{stage: stageProblems, category: "two-pointers"})
	if m.stage != stageProblems {
		t.Fatalf("stage = %v, want stageProblems", m.stage)
	}
	if len(m.categoryProblems) != 1 || m.categoryProblems[0].ProblemID != "two-pointers-01" {
		t.Errorf("categoryProblems = %+v, want just two-pointers-01", m.categoryProblems)
	}
}
