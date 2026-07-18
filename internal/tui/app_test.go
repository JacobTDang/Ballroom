package tui

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
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
	checkToolCallingFn = func(context.Context, string, string, string) (bool, error) { return supported, err }
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

func fakeAllAttempts(attempts []tracker.Attempt, err error) func() {
	orig := allAttemptsFn
	allAttemptsFn = func(config.Config) ([]tracker.Attempt, error) { return attempts, err }
	return func() { allAttemptsFn = orig }
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
	m := appModel{cursor: int(menuSandbox)}
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

	m := appModel{cursor: int(menuStats)}
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

func TestAppModel_Stats_ShowsRubricWeakSpots(t *testing.T) {
	defer fakeCatalogList(nil, nil)()
	defer fakeRecentAttempts(nil, nil)()
	defer fakeAllAttempts([]tracker.Attempt{
		{Category: "system-design", GradeSummary: "1. Back-of-envelope estimates: missing. none"},
		{Category: "system-design", GradeSummary: "1. Back-of-envelope estimates: missing. again\n2. High-level design: strong. good"},
	}, nil)()

	m := appModel{stage: stageMain}
	m = m.loadStats()
	view := m.View()
	if !strings.Contains(view, "Rubric weak spots") {
		t.Fatalf("stats view missing the weak-spots section:\n%s", view)
	}
	if !strings.Contains(view, "Back-of-envelope estimates") || !strings.Contains(view, "missing 2/2") {
		t.Errorf("stats view should rank estimates as missing 2/2:\n%s", view)
	}
}

func TestAppModel_Stats_NoGradedAttemptsHidesWeakSpots(t *testing.T) {
	defer fakeCatalogList(nil, nil)()
	defer fakeRecentAttempts(nil, nil)()
	defer fakeAllAttempts([]tracker.Attempt{{Category: "dsa"}}, nil)()

	m := appModel{stage: stageMain}
	m = m.loadStats()
	if strings.Contains(m.View(), "Rubric weak spots") {
		t.Error("weak-spots section should be absent with no graded attempts")
	}
}

func TestAppModel_Stats_ShowsCodingWeakSpots(t *testing.T) {
	defer fakeCatalogList(nil, nil)()
	defer fakeRecentAttempts(nil, nil)()
	defer fakeAllAttempts([]tracker.Attempt{
		{Category: "trees", Result: tracker.ResultFail},
		{Category: "trees", Result: tracker.ResultFail},
		{Category: "trees", Result: tracker.ResultPass},
	}, nil)()

	m := appModel{stage: stageMain}
	m = m.loadStats()
	view := m.View()
	if !strings.Contains(view, "Coding weak spots") {
		t.Fatalf("stats view missing the coding weak-spots section:\n%s", view)
	}
	if !strings.Contains(view, "Trees") || !strings.Contains(view, "failed 2/3") {
		t.Errorf("stats view should rank Trees as failed 2/3:\n%s", view)
	}
}

func TestAppModel_Stats_TooFewCodingAttemptsHidesCodingWeakSpots(t *testing.T) {
	defer fakeCatalogList(nil, nil)()
	defer fakeRecentAttempts(nil, nil)()
	// One failure is under the evidence threshold -- a single bad day
	// must not brand a whole topic weak.
	defer fakeAllAttempts([]tracker.Attempt{{Category: "trees", Result: tracker.ResultFail}}, nil)()

	m := appModel{stage: stageMain}
	m = m.loadStats()
	if strings.Contains(m.View(), "Coding weak spots") {
		t.Error("coding weak-spots section should be absent under the attempts threshold")
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

func TestAppModel_EnterOnSettings_GoesToStageSettingsNoAsyncLoadYet(t *testing.T) {
	// Settings leads with the Worker/Orchestrator role choice, not
	// straight into a provider choice or the Ollama picker — no async
	// load should happen until Local is chosen several levels in (see
	// TestAppModel_ProviderChoice_SelectingLocal...).
	m := appModel{cursor: int(menuSettings)}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no command yet — the role choice itself needs no data")
	}
	got := newM.(appModel)
	if got.stage != stageSettings {
		t.Fatalf("stage = %v, want stageSettings", got.stage)
	}
}

func TestAppModel_Settings_SelectingWorkerGoesToProviderChoiceTargetingWorker(t *testing.T) {
	m := appModel{stage: stageSettings, settingsCursor: 0} // 0 = Worker model
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command — entering the provider choice is an internal stage change")
	}
	got := newM.(appModel)
	if got.stage != stageProviderChoice {
		t.Fatalf("stage = %v, want stageProviderChoice", got.stage)
	}
	if got.settingsEditing != settingsTargetWorker {
		t.Errorf("settingsEditing = %v, want settingsTargetWorker", got.settingsEditing)
	}
}

func TestAppModel_Settings_SelectingOrchestratorGoesToProviderChoiceTargetingOrchestrator(t *testing.T) {
	m := appModel{stage: stageSettings, settingsCursor: 1} // 1 = Orchestrator model
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command")
	}
	got := newM.(appModel)
	if got.stage != stageProviderChoice {
		t.Fatalf("stage = %v, want stageProviderChoice", got.stage)
	}
	if got.settingsEditing != settingsTargetOrchestrator {
		t.Errorf("settingsEditing = %v, want settingsTargetOrchestrator", got.settingsEditing)
	}
}

func TestAppModel_ProviderChoice_SelectingLocalKicksOffAsyncModelLoad(t *testing.T) {
	defer fakeListModels([]string{"a", "b"}, nil)()

	m := appModel{stage: stageProviderChoice, settingsCursor: 0, settingsEditing: settingsTargetWorker} // 0 = Local
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected enter on Local to kick off an async models-loading command")
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

func TestAppModel_ProviderChoice_SelectingAPIGoesToAPIModelEntry(t *testing.T) {
	m := appModel{stage: stageProviderChoice, settingsCursor: 1, settingsEditing: settingsTargetWorker} // 1 = API
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command — entering the slug-entry stage is an internal stage change")
	}
	got := newM.(appModel)
	if got.stage != stageAPIModelEntry {
		t.Fatalf("stage = %v, want stageAPIModelEntry", got.stage)
	}
	if len(got.apiModelFiltered) == 0 {
		t.Error("expected apiModelFiltered to be seeded with suggestedOpenRouterModels, got none")
	}
	if got.apiModelCursor != 0 {
		t.Errorf("apiModelCursor = %d, want 0", got.apiModelCursor)
	}
}

func TestAppModel_ProviderChoice_OrchestratorGetsAThirdNoneOption(t *testing.T) {
	m := appModel{
		cfg:             config.Config{DataDir: t.TempDir(), TutorModel: "worker-model", OrchestratorModel: "orchestrator-model"},
		stage:           stageProviderChoice,
		settingsCursor:  2, // None (disable routing) -- only present for the orchestrator target
		settingsEditing: settingsTargetOrchestrator,
	}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command")
	}
	got := newM.(appModel)
	if got.stage != stageSettings {
		t.Fatalf("stage = %v, want stageSettings (back out after clearing)", got.stage)
	}
	if got.cfg.OrchestratorModel != "" {
		t.Errorf("cfg.OrchestratorModel = %q, want empty after selecting None", got.cfg.OrchestratorModel)
	}
	if got.cfg.TutorModel != "worker-model" {
		t.Errorf("cfg.TutorModel = %q, want it preserved", got.cfg.TutorModel)
	}

	saved, err := config.LoadSettings(m.cfg.SettingsPath())
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if saved.OrchestratorModel != "" {
		t.Errorf("persisted OrchestratorModel = %q, want empty", saved.OrchestratorModel)
	}
	if saved.TutorModel != "worker-model" {
		t.Errorf("persisted TutorModel = %q, want it preserved", saved.TutorModel)
	}
}

func TestAppModel_ProviderChoice_WorkerTargetHasNoNoneOption(t *testing.T) {
	// The worker model is never optional -- there must always be one --
	// so its provider choice list stays at 2 items (Local/API), unlike
	// the orchestrator's 3.
	m := appModel{stage: stageProviderChoice, settingsCursor: 1, settingsEditing: settingsTargetWorker}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	got := newM.(appModel)
	if got.settingsCursor != 1 {
		t.Errorf("settingsCursor = %d, want to stay at 1 (only 2 items: Local, API)", got.settingsCursor)
	}
}

func TestAppModel_ProviderChoice_EscGoesBackToSettings(t *testing.T) {
	m := appModel{stage: stageProviderChoice}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Error("expected no external command")
	}
	if newM.(appModel).stage != stageSettings {
		t.Errorf("stage = %v, want stageSettings", newM.(appModel).stage)
	}
}

func TestAppModel_Settings_EscGoesBackToMain(t *testing.T) {
	m := appModel{stage: stageSettings}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Error("expected no external command")
	}
	if newM.(appModel).stage != stageMain {
		t.Errorf("stage = %v, want stageMain", newM.(appModel).stage)
	}
}

func TestAppModel_RenderSettings_ShowsCurrentModel(t *testing.T) {
	m := appModel{cfg: config.Config{TutorModel: "llama3.1:8b"}, stage: stageSettings}
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, "llama3.1:8b") {
		t.Errorf("expected the current model in the view, got:\n%s", view)
	}
}

func TestAppModel_APIModelEntry_EnterWithKeyAlreadySetSelectsImmediately(t *testing.T) {
	defer fakeCheckToolCalling(true, nil)()

	dir := t.TempDir()
	m := appModel{cfg: config.Config{DataDir: dir, OpenRouterAPIKey: "sk-existing"}, stage: stageAPIModelEntry, apiModelInput: "anthropic/claude-3.5-sonnet"}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected a background command kicking off the tool-calling check")
	}
	got := newM.(appModel)
	if got.stage != stageMain {
		t.Fatalf("stage = %v, want stageMain", got.stage)
	}
	if got.cfg.TutorModel != "openrouter:anthropic/claude-3.5-sonnet" {
		t.Errorf("cfg.TutorModel = %q, want %q", got.cfg.TutorModel, "openrouter:anthropic/claude-3.5-sonnet")
	}
}

func TestAppModel_APIModelEntry_EnterWithNoKeyGoesToKeyEntry(t *testing.T) {
	m := appModel{cfg: config.Config{DataDir: t.TempDir()}, stage: stageAPIModelEntry, apiModelInput: "anthropic/claude-3.5-sonnet"}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected no external command — no key yet, just an internal stage change")
	}
	got := newM.(appModel)
	if got.stage != stageOpenRouterKeyEntry {
		t.Fatalf("stage = %v, want stageOpenRouterKeyEntry", got.stage)
	}
	if got.openRouterPendingModel != "openrouter:anthropic/claude-3.5-sonnet" {
		t.Errorf("openRouterPendingModel = %q, want %q", got.openRouterPendingModel, "openrouter:anthropic/claude-3.5-sonnet")
	}
}

func TestAppModel_APIModelEntry_EscGoesBackToProviderChoice(t *testing.T) {
	m := appModel{stage: stageAPIModelEntry, apiModelInput: "partial"}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Error("expected no external command")
	}
	if newM.(appModel).stage != stageProviderChoice {
		t.Errorf("stage = %v, want stageProviderChoice", newM.(appModel).stage)
	}
}

func TestAppModel_APIModelEntry_BackspaceRemovesLastCharacter(t *testing.T) {
	m := appModel{stage: stageAPIModelEntry, apiModelInput: "anthropic/x"}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	got := newM.(appModel)
	if got.apiModelInput != "anthropic/" {
		t.Errorf("apiModelInput = %q, want %q", got.apiModelInput, "anthropic/")
	}
}

// apiModelEntryFixture enters stageAPIModelEntry the real way (via
// stageProviderChoice), so apiModelFiltered starts seeded with
// suggestedOpenRouterModels exactly like a real session.
func apiModelEntryFixture(t *testing.T) appModel {
	t.Helper()
	m := appModel{stage: stageProviderChoice, settingsCursor: 1, settingsEditing: settingsTargetWorker}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	return newM.(appModel)
}

func TestAppModel_APIModelEntry_TypingFiltersSuggestedModels(t *testing.T) {
	m := apiModelEntryFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("gpt-oss-20b")})
	got := newM.(appModel)
	if len(got.apiModelFiltered) != 1 || got.apiModelFiltered[0] != "openai/gpt-oss-20b:free" {
		t.Errorf("apiModelFiltered = %v, want [openai/gpt-oss-20b:free]", got.apiModelFiltered)
	}
}

func TestAppModel_APIModelEntry_DownMovesCursorWithinFilteredList(t *testing.T) {
	m := apiModelEntryFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	got := newM.(appModel)
	if got.apiModelCursor != 1 {
		t.Errorf("apiModelCursor = %d, want 1", got.apiModelCursor)
	}
}

func TestAppModel_APIModelEntry_DownStopsAtLastFilteredEntry(t *testing.T) {
	m := apiModelEntryFixture(t)
	for i := 0; i < len(m.apiModelFiltered)+5; i++ {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = newM.(appModel)
	}
	if m.apiModelCursor != len(m.apiModelFiltered)-1 {
		t.Errorf("apiModelCursor = %d, want %d (last entry)", m.apiModelCursor, len(m.apiModelFiltered)-1)
	}
}

func TestAppModel_APIModelEntry_EnterOnHighlightedSuggestionSelectsItWithOpenRouterPrefix(t *testing.T) {
	defer fakeCheckToolCalling(true, nil)()
	dir := t.TempDir()

	m := apiModelEntryFixture(t)
	m.cfg = config.Config{DataDir: dir, OpenRouterAPIKey: "sk-existing"}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newM.(appModel)
	want := "openrouter:" + m.apiModelFiltered[m.apiModelCursor]

	newM2, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Error("expected a background command kicking off the tool-calling check")
	}
	got := newM2.(appModel)
	if got.cfg.TutorModel != want {
		t.Errorf("cfg.TutorModel = %q, want %q", got.cfg.TutorModel, want)
	}
}

func TestAppModel_APIModelEntry_QWithNoFilterGoesBackToProviderChoice(t *testing.T) {
	m := apiModelEntryFixture(t)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd != nil {
		t.Error("expected no external command — back to the provider choice is an internal stage change")
	}
	if newM.(appModel).stage != stageProviderChoice {
		t.Error("expected q with no filter to return to stageProviderChoice")
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

// TestAppModel_DailyJumpsStraightToTheLanguagePicker drives the new
// main-menu entry end to end: "2" jumps the cursor to Daily, enter
// picks today's problem deterministically and lands on stageLanguage
// with the pick's category context set up for backing out.
func TestAppModel_DailyJumpsStraightToTheLanguagePicker(t *testing.T) {
	restore := fakeCatalogList([]catalog.ExerciseStatus{
		{Exercise: exercise.Exercise{ID: "two-pointers-01-go", ProblemID: "two-pointers-01", Title: "Two Sum II", Category: "two-pointers", Language: "go"}},
	}, nil)
	defer restore()

	m := newAppModel(config.Config{}, appResume{})
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("2")})
	m = newM.(appModel)
	if menuChoice(m.cursor) != menuDaily {
		t.Fatalf("cursor = %d after pressing 2, want the Daily entry", m.cursor)
	}
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(appModel)

	if m.err != nil {
		t.Fatalf("daily load error: %v", m.err)
	}
	if m.stage != stageLanguage {
		t.Fatalf("stage = %v, want stageLanguage -- Daily skips category/problem navigation", m.stage)
	}
	if m.selectedProblem.ProblemID != "two-pointers-01" {
		t.Errorf("selectedProblem = %q, want the only available problem", m.selectedProblem.ProblemID)
	}
	if m.category != "two-pointers" || len(m.categoryProblems) != 1 {
		t.Errorf("category context = (%q, %d problems), want the pick's category set up for backing out", m.category, len(m.categoryProblems))
	}

	// Backing out must walk populated picker screens, not empty ones --
	// a real gap found in review: the first version populated only
	// categoryProblems, so q -> q rendered an empty topic list and an
	// empty category list.
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	m = newM.(appModel)
	if m.stage != stageProblems {
		t.Fatalf("stage after q = %v, want stageProblems", m.stage)
	}
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	m = newM.(appModel)
	if m.stage != stageDSACategories || len(m.dsaCategories) == 0 {
		t.Fatalf("stage after q q = %v with %d topics, want a populated DSA topic list for a grouped-category pick", m.stage, len(m.dsaCategories))
	}
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	m = newM.(appModel)
	if m.stage != stageCategories || len(m.categories) == 0 {
		t.Fatalf("stage after q q q = %v with %d categories, want the populated category list", m.stage, len(m.categories))
	}
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

func TestFilterByCategory_ListsDueProblemsFirst(t *testing.T) {
	dueProblem := catalog.ProblemStatus{
		ProblemID: "search-kv-store-01", Title: "Design a key-value store", Category: exercise.CategorySystemDesign,
		Variants: []catalog.ExerciseStatus{
			{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}, Attempts: 1, LastResult: tracker.ResultPass},
			{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageInterviewer}},
		},
	}
	notDue := catalog.ProblemStatus{
		ProblemID: "mint-01", Title: "Design Mint.com", Category: exercise.CategorySystemDesign,
		Variants: []catalog.ExerciseStatus{
			{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}},
			{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageInterviewer}},
		},
	}
	otherCategory := catalog.ProblemStatus{ProblemID: "two-pointers-01", Category: exercise.CategoryDSA}

	// The due problem sits alphabetically after the non-due one -- the
	// exact shape that used to hide the "mock due" marker below the fold.
	got := filterByCategory([]catalog.ProblemStatus{notDue, dueProblem, otherCategory}, exercise.CategorySystemDesign)

	if len(got) != 2 {
		t.Fatalf("got %d problems, want 2 (the other-category problem filtered out)", len(got))
	}
	if got[0].ProblemID != "search-kv-store-01" || got[1].ProblemID != "mint-01" {
		t.Errorf("order = [%s, %s], want the due problem floated first", got[0].ProblemID, got[1].ProblemID)
	}
}

func codingProblem(id, title, difficulty string) catalog.ProblemStatus {
	return catalog.ProblemStatus{
		ProblemID: id, Title: title, Category: "two-pointers",
		Variants: []catalog.ExerciseStatus{
			{Exercise: exercise.Exercise{ID: id + "-go", ProblemID: id, Title: title, Category: "two-pointers", Language: "go", Difficulty: difficulty}},
		},
	}
}

func TestFilterProblems_CaseInsensitiveTitleMatch(t *testing.T) {
	problems := []catalog.ProblemStatus{
		codingProblem("p1", "Two Sum II", "easy"),
		codingProblem("p2", "Container With Most Water", "medium"),
		codingProblem("p3", "Trapping Rain Water", "hard"),
	}
	if got := filterProblems(problems, ""); len(got) != 3 {
		t.Errorf("empty filter matched %d, want all 3", len(got))
	}
	got := filterProblems(problems, "WATER")
	if len(got) != 2 || got[0].ProblemID != "p2" || got[1].ProblemID != "p3" {
		t.Errorf("filter WATER = %+v, want p2,p3 in order", got)
	}
	if got := filterProblems(problems, "zzz"); len(got) != 0 {
		t.Errorf("filter zzz matched %d, want none", len(got))
	}
}

func TestAppModel_Problems_TypingFiltersAndEnterSelectsFromFiltered(t *testing.T) {
	m := appModel{stage: stageProblems, category: "two-pointers", categoryProblems: []catalog.ProblemStatus{
		codingProblem("p1", "Two Sum II", "easy"),
		codingProblem("p2", "Container With Most Water", "medium"),
		codingProblem("p3", "Trapping Rain Water", "hard"),
	}}
	for _, r := range "rain" {
		newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = newM.(appModel)
	}
	if m.problemFilter != "rain" {
		t.Fatalf("problemFilter = %q, want %q", m.problemFilter, "rain")
	}
	if vis := m.visibleProblems(); len(vis) != 1 || vis[0].ProblemID != "p3" {
		t.Fatalf("visibleProblems = %+v, want just p3", vis)
	}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(appModel)
	if m.stage != stageLanguage || m.selectedProblem.ProblemID != "p3" {
		t.Errorf("enter selected %q (stage %v), want the filtered match p3 at stageLanguage", m.selectedProblem.ProblemID, m.stage)
	}
}

func TestAppModel_Problems_QFeedsFilterOnceTypingStarted(t *testing.T) {
	m := appModel{stage: stageProblems, category: "debug", categoryProblems: []catalog.ProblemStatus{
		codingProblem("p1", "Quick fix", "easy"),
	}}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	m = newM.(appModel)
	if m.stage != stageCategories {
		t.Fatalf("q with empty filter: stage = %v, want stageCategories (backs out)", m.stage)
	}

	m = appModel{stage: stageProblems, category: "debug", categoryProblems: []catalog.ProblemStatus{
		codingProblem("p1", "Quick fix", "easy"),
	}, problemFilter: "ui"}
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	m = newM.(appModel)
	if m.stage != stageProblems || m.problemFilter != "uiq" {
		t.Errorf("q mid-filter: stage=%v filter=%q, want it appended to the filter", m.stage, m.problemFilter)
	}

	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	m = newM.(appModel)
	if m.problemFilter != "ui" {
		t.Errorf("backspace: filter = %q, want %q", m.problemFilter, "ui")
	}
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newM.(appModel)
	if m.stage != stageCategories {
		t.Errorf("esc: stage = %v, want it to back out even mid-filter", m.stage)
	}
}

func TestAppModel_Problems_EntryResetsFilter(t *testing.T) {
	m := problemsFixture(t)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	m = newM.(appModel)
	if m.problemFilter != "x" {
		t.Fatalf("setup: filter = %q, want %q", m.problemFilter, "x")
	}
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc}) // back to DSA topics
	m = newM.(appModel)
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // re-enter the topic
	m = newM.(appModel)
	if m.stage != stageProblems || m.problemFilter != "" {
		t.Errorf("re-entry: stage=%v filter=%q, want a fresh empty filter", m.stage, m.problemFilter)
	}
}

func TestDifficultyBadge_LettersAndMixedRange(t *testing.T) {
	cases := []struct {
		p    catalog.ProblemStatus
		want string
	}{
		{codingProblem("p1", "A", "easy"), "E"},
		{codingProblem("p2", "B", "medium"), "M"},
		{codingProblem("p3", "C", "hard"), "H"},
		{codingProblem("p4", "D", ""), ""},
		{catalog.ProblemStatus{Variants: []catalog.ExerciseStatus{
			{Exercise: exercise.Exercise{Difficulty: "medium"}},
			{Exercise: exercise.Exercise{Difficulty: "hard"}},
		}}, "M/H"},
	}
	for _, c := range cases {
		if got := difficultyBadge(c.p); got != c.want {
			t.Errorf("difficultyBadge(%s) = %q, want %q", c.p.ProblemID, got, c.want)
		}
	}
}

func TestRenderProblems_ShowsDifficultyBadgesAndSearchPrompt(t *testing.T) {
	m := appModel{stage: stageProblems, category: "two-pointers", categoryProblems: []catalog.ProblemStatus{
		codingProblem("p1", "Two Sum II", "easy"),
		codingProblem("p2", "Container With Most Water", "medium"),
		codingProblem("p3", "Trapping Rain Water", "hard"),
	}}
	out := m.renderProblems()
	for _, want := range []string{"[E]", "[M]", "[H]", "type to search"} {
		if !strings.Contains(out, want) {
			t.Errorf("renderProblems missing %q", want)
		}
	}
}

func TestRenderProblems_WindowFollowsCursorOnLongLists(t *testing.T) {
	var problems []catalog.ProblemStatus
	for i := 0; i < 20; i++ {
		problems = append(problems, codingProblem(fmt.Sprintf("p%02d", i), fmt.Sprintf("Problem %02d", i), "easy"))
	}
	m := appModel{stage: stageProblems, category: "two-pointers", categoryProblems: problems}

	top := m.renderProblems()
	if !strings.Contains(top, "Problem 00") || strings.Contains(top, "Problem 19") {
		t.Errorf("cursor at top: want the first problem visible and the last one below the window")
	}
	if !strings.Contains(top, "more") {
		t.Errorf("cursor at top: want a more-below indicator, got:\n%s", top)
	}

	m.problemCursor = 19
	bottom := m.renderProblems()
	if !strings.Contains(bottom, "Problem 19") || strings.Contains(bottom, "Problem 00") {
		t.Errorf("cursor at bottom: want the last problem visible and the first one above the window")
	}
}

func TestAppModel_Problems_MockDueMarker(t *testing.T) {
	m := appModel{stage: stageProblems, category: exercise.CategorySystemDesign, categoryProblems: []catalog.ProblemStatus{
		{
			ProblemID: "url-shortener-01", Title: "Design Pastebin / Bit.ly", Category: exercise.CategorySystemDesign,
			Solved: true, Attempts: 1,
			Variants: []catalog.ExerciseStatus{
				{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}, Attempts: 1, LastResult: tracker.ResultPass},
				{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageInterviewer}},
			},
		},
		{
			ProblemID: "web-crawler-01", Title: "Design a Web Crawler", Category: exercise.CategorySystemDesign,
			Variants: []catalog.ExerciseStatus{
				{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}},
				{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageInterviewer}},
			},
		},
	}}
	view := m.View()
	if strings.Count(view, "mock due") != 1 {
		t.Errorf("want exactly one (mock due) marker (coach passed + interviewer untouched), view:\n%s", view)
	}
}

func TestAppModel_Problems_ReviewDueMarker(t *testing.T) {
	m := appModel{stage: stageProblems, category: exercise.CategoryDSA, categoryProblems: []catalog.ProblemStatus{
		{
			ProblemID: "two-pointers-01", Title: "Two Sum II", Category: exercise.CategoryDSA,
			Attempts: 1,
			Variants: []catalog.ExerciseStatus{
				// Failed long ago -- review-due regardless of when the test runs.
				{Exercise: exercise.Exercise{Language: "go"}, Attempts: 1, LastResult: tracker.ResultFail, LastAttemptDate: "2020-01-01"},
			},
		},
		{
			ProblemID: "off-by-one-01", Title: "Off by one", Category: exercise.CategoryDSA,
			Variants: []catalog.ExerciseStatus{
				{Exercise: exercise.Exercise{Language: "go"}},
			},
		},
	}}
	view := m.View()
	if strings.Count(view, "review due") != 1 {
		t.Errorf("want exactly one review-due marker (stale fail yes, untouched problem no), view:\n%s", view)
	}
	if strings.Contains(view, "mock due") {
		t.Errorf("a review-due coding problem must not claim to be mock due, view:\n%s", view)
	}
}

func TestAppModel_Language_DesignProblemSaysSessionStyle(t *testing.T) {
	// A design problem's "language" variants are session styles
	// (coach/interviewer) -- calling them a language on screen would be
	// wrong, this is the one cosmetic spot the variant trick shows.
	m := appModel{stage: stageLanguage, selectedProblem: catalog.ProblemStatus{
		ProblemID: "url-shortener-01",
		Title:     "Design Pastebin / Bit.ly",
		Category:  exercise.CategorySystemDesign,
		Variants: []catalog.ExerciseStatus{
			{Exercise: exercise.Exercise{ID: "url-shortener-01-coach", Kind: exercise.KindDesign, Language: exercise.LanguageCoach}},
			{Exercise: exercise.Exercise{ID: "url-shortener-01-interviewer", Kind: exercise.KindDesign, Language: exercise.LanguageInterviewer}},
		},
	}}
	view := m.View()
	if !strings.Contains(view, "choose a session style") {
		t.Errorf("design problem's variant picker should say \"choose a session style\", got:\n%s", view)
	}
	if strings.Contains(view, "choose a language") {
		t.Errorf("design problem's variant picker must not say language:\n%s", view)
	}
	for _, want := range []string{"coach", "interviewer"} {
		if !strings.Contains(view, want) {
			t.Errorf("variant picker missing %q:\n%s", want, view)
		}
	}
}

func TestAppModel_Language_CodingProblemStillSaysLanguage(t *testing.T) {
	m := languageFixture(t)
	if view := m.View(); !strings.Contains(view, "choose a language") {
		t.Errorf("coding problem's variant picker should still say \"choose a language\", got:\n%s", view)
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

func TestAppModel_ModelPicker_QWithNoFilterGoesBackToProviderChoice(t *testing.T) {
	// The model picker is now only reachable via Settings -> (Worker or
	// Orchestrator) -> Local, so backing out lands on the provider
	// choice, not all the way back to the main menu.
	m := modelPickerFixture(t, []string{"qwen2.5-coder:7b"})
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd != nil {
		t.Error("expected no external command — back to the provider choice is an internal stage change")
	}
	if newM.(appModel).stage != stageProviderChoice {
		t.Error("expected q with no filter to return to stageProviderChoice")
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

func TestAppModel_ModelPicker_MetadataDenialShowsNoToolsMarker(t *testing.T) {
	orig := claimsToolSupportFn
	claimsToolSupportFn = func(_ context.Context, _, model, _ string) (bool, bool) {
		return model != "gemma:2b", true // gemma denies tools, others claim them
	}
	t.Cleanup(func() { claimsToolSupportFn = orig })

	m := modelPickerFixture(t, nil)
	newM, cmd := m.Update(modelsLoadedMsg{models: []string{"llama3.1:8b", "gemma:2b"}})
	m = newM.(appModel)
	if cmd == nil {
		t.Fatal("modelsLoadedMsg produced no command, want the metadata check to fire")
	}
	newM, _ = m.Update(cmd())
	m = newM.(appModel)

	view := m.View()
	if !strings.Contains(view, "gemma:2b") || !strings.Contains(view, "(no tools)") {
		t.Errorf("view should mark gemma:2b with (no tools):\n%s", view)
	}
	if strings.Count(view, "(no tools)") != 1 {
		t.Errorf("exactly one row should carry the marker, view:\n%s", view)
	}
}

func TestAppModel_ModelPicker_MetadataFailureShowsNoMarker(t *testing.T) {
	orig := claimsToolSupportFn
	claimsToolSupportFn = func(_ context.Context, _, _, _ string) (bool, bool) {
		return false, false // fetch failed: no signal
	}
	t.Cleanup(func() { claimsToolSupportFn = orig })

	m := modelPickerFixture(t, nil)
	newM, cmd := m.Update(modelsLoadedMsg{models: []string{"llama3.1:8b"}})
	m = newM.(appModel)
	if cmd == nil {
		t.Fatal("modelsLoadedMsg produced no command")
	}
	newM, _ = m.Update(cmd())
	m = newM.(appModel)

	if strings.Contains(m.View(), "(no tools)") {
		t.Errorf("metadata failure must degrade to no marker, view:\n%s", m.View())
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

func TestAppModel_OpenRouterKeyEntry_EscCancelsBackToProviderChoiceWithoutSelecting(t *testing.T) {
	// A single consistent "cancel returns to the provider choice"
	// behavior, regardless of whether this stage was reached by typing
	// openrouter: directly in the local picker or via the newer
	// Settings -> API -> slug-entry path.
	m := appModel{cfg: config.Config{DataDir: t.TempDir(), TutorModel: "old-model"}, stage: stageOpenRouterKeyEntry, openRouterPendingModel: "openrouter:x/y", openRouterKeyInput: "sk-partial"}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd != nil {
		t.Error("expected no external command on cancel")
	}
	got := newM.(appModel)
	if got.stage != stageProviderChoice {
		t.Errorf("stage = %v, want stageProviderChoice", got.stage)
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

func TestAppModel_SelectModel_WritesOrchestratorModelWhenTargetingOrchestrator(t *testing.T) {
	defer fakeCheckToolCalling(true, nil)()

	dir := t.TempDir()
	cfg := config.Config{DataDir: dir, TutorModel: "worker-model"}
	m := appModel{cfg: cfg, settingsEditing: settingsTargetOrchestrator}

	newM, _ := m.selectModel("orchestrator-model")
	got := newM.(appModel)
	if got.cfg.OrchestratorModel != "orchestrator-model" {
		t.Errorf("cfg.OrchestratorModel = %q, want %q", got.cfg.OrchestratorModel, "orchestrator-model")
	}
	if got.cfg.TutorModel != "worker-model" {
		t.Errorf("cfg.TutorModel = %q, want it preserved (unaffected by an orchestrator pick)", got.cfg.TutorModel)
	}

	saved, err := config.LoadSettings(cfg.SettingsPath())
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if saved.OrchestratorModel != "orchestrator-model" {
		t.Errorf("persisted OrchestratorModel = %q, want %q", saved.OrchestratorModel, "orchestrator-model")
	}
	if saved.TutorModel != "worker-model" {
		t.Errorf("persisted TutorModel = %q, want it preserved", saved.TutorModel)
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

func TestUpdateProblems_DefaultLanguageRunsMatchingVariantDirectly(t *testing.T) {
	m := appModel{
		stage: stageProblems,
		cfg:   config.Config{DefaultLanguage: "go"},
		categoryProblems: []catalog.ProblemStatus{{
			ProblemID: "p1",
			Variants: []catalog.ExerciseStatus{
				{Exercise: exercise.Exercise{ID: "p1-python", Language: "python"}},
				{Exercise: exercise.Exercise{ID: "p1-go", Language: "go"}},
			},
		}},
	}
	newM, cmd := m.updateProblems(tea.KeyMsg{Type: tea.KeyEnter})
	got := newM.(appModel)
	if got.outcome != outcomeRunExercise || got.exerciseToRun.ID != "p1-go" {
		t.Fatalf("outcome=%v exerciseToRun=%q, want the go variant run directly", got.outcome, got.exerciseToRun.ID)
	}
	if cmd == nil {
		t.Fatal("want a tea.Quit cmd when the default language matches")
	}
}

func TestUpdateProblems_DefaultLanguageWithoutMatchingVariantStillAsks(t *testing.T) {
	// Design problems ride coach/interviewer in the language slot — a
	// python/go/cpp default must never match them.
	m := appModel{
		stage: stageProblems,
		cfg:   config.Config{DefaultLanguage: "go"},
		categoryProblems: []catalog.ProblemStatus{{
			ProblemID: "url-shortener",
			Variants: []catalog.ExerciseStatus{
				{Exercise: exercise.Exercise{ID: "u-coach", Language: "coach"}},
				{Exercise: exercise.Exercise{ID: "u-interviewer", Language: "interviewer"}},
			},
		}},
	}
	newM, cmd := m.updateProblems(tea.KeyMsg{Type: tea.KeyEnter})
	got := newM.(appModel)
	if got.stage != stageLanguage || cmd != nil {
		t.Fatalf("stage=%v cmd=%v, want the language picker kept when no variant matches", got.stage, cmd)
	}
}

func TestUpdateSettings_CycleDefaultLanguagePersists(t *testing.T) {
	m := appModel{stage: stageSettings, settingsCursor: 2, cfg: config.Config{DataDir: t.TempDir()}}
	for _, want := range []string{"python", "go", "cpp", ""} {
		newM, _ := m.updateSettings(tea.KeyMsg{Type: tea.KeyEnter})
		m = newM.(appModel)
		if m.cfg.DefaultLanguage != want {
			t.Fatalf("DefaultLanguage = %q, want %q in the cycle", m.cfg.DefaultLanguage, want)
		}
		s, err := config.LoadSettings(m.cfg.SettingsPath())
		if err != nil || s.DefaultLanguage != want {
			t.Fatalf("persisted DefaultLanguage = %q (err %v), want %q", s.DefaultLanguage, err, want)
		}
		if m.stage != stageSettings {
			t.Fatalf("stage = %v, want to stay on stageSettings", m.stage)
		}
	}
}

func TestUpdateSettings_ToggleTutorNotesPersists(t *testing.T) {
	m := appModel{stage: stageSettings, settingsCursor: 3, cfg: config.Config{DataDir: t.TempDir()}}
	newM, _ := m.updateSettings(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(appModel)
	if !m.cfg.DisableTutorNotes {
		t.Fatal("DisableTutorNotes = false after toggle, want true")
	}
	s, err := config.LoadSettings(m.cfg.SettingsPath())
	if err != nil || !s.DisableTutorNotes {
		t.Fatalf("persisted DisableTutorNotes = %v (err %v), want true", s.DisableTutorNotes, err)
	}
	newM, _ = m.updateSettings(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(appModel)
	if m.cfg.DisableTutorNotes {
		t.Fatal("DisableTutorNotes = true after second toggle, want false")
	}
}

func TestRenderMain_ErrorShowsHeadlineWithRawDetailKeptVisible(t *testing.T) {
	m := appModel{stage: stageMain, err: errors.New("sqlite: database is locked")}
	out := stripAnsiTUI(m.renderMain())
	if !strings.Contains(out, "something went wrong") {
		t.Errorf("renderMain = %q, want the human headline", out)
	}
	if !strings.Contains(out, "sqlite: database is locked") {
		t.Errorf("renderMain = %q, want the raw error detail still visible, never swallowed", out)
	}
}

func TestRenderModelPicker_OllamaErrorShowsHeadlineAndDetail(t *testing.T) {
	m := appModel{stage: stageModelPicker, modelLoadErr: errors.New("dial tcp: connection refused")}
	out := stripAnsiTUI(m.renderModelPicker())
	if !strings.Contains(out, "couldn't reach Ollama") || !strings.Contains(out, "connection refused") {
		t.Errorf("renderModelPicker = %q, want headline and raw detail", out)
	}
	if !strings.Contains(out, "type a model tag directly") {
		t.Errorf("renderModelPicker = %q, want the escape-hatch hint kept", out)
	}
}
