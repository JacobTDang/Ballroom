package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
}

// TestLoad_DefaultsToCWDWhenCWDLooksLikeCheckout covers the ordinary
// "ballroom run from inside the checkout" case: no PRACTICE_ROOT, cwd
// itself has docker/Dockerfile, so Root/ExercisesDir/TestsDir/DBPath
// all resolve straight off it, same as before ResolveRoot existed.
func TestLoad_DefaultsToCWDWhenCWDLooksLikeCheckout(t *testing.T) {
	withRootCache(t)
	dir := checkoutDir(t)
	// resolve symlinks (macOS TempDir can be under /var -> /private/var)
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	chdir(t, dir)
	t.Setenv("PRACTICE_ROOT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if cfg.Root != resolved {
		t.Errorf("Root = %q, want %q", cfg.Root, resolved)
	}
	if want := filepath.Join(resolved, "exercises"); cfg.ExercisesDir != want {
		t.Errorf("ExercisesDir = %q, want %q", cfg.ExercisesDir, want)
	}
	if want := filepath.Join(resolved, "tests"); cfg.TestsDir != want {
		t.Errorf("TestsDir = %q, want %q", cfg.TestsDir, want)
	}
	if want := filepath.Join(resolved, "data", "tracker.db"); cfg.DBPath != want {
		t.Errorf("DBPath = %q, want %q", cfg.DBPath, want)
	}
}

// TestLoad_FallsBackToCachedRootWhenCWDIsNotACheckout is issue #255's
// core fix: `go install ./cmd/ballroom` then running `ballroom` from
// $HOME (or anywhere else outside the checkout) must still populate
// real exercises/tests/data dirs, not silently resolve them under
// $HOME. It does, via the same cached-checkout-root fallback the
// docker build root already used (see ResolveRoot).
func TestLoad_FallsBackToCachedRootWhenCWDIsNotACheckout(t *testing.T) {
	cachePath := withRootCache(t)
	realCheckout := checkoutDir(t)
	cacheRoot(cachePath, realCheckout)

	chdir(t, t.TempDir()) // stand-in for $HOME: not a checkout
	t.Setenv("PRACTICE_ROOT", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Root != realCheckout {
		t.Errorf("Root = %q, want the cached checkout %q", cfg.Root, realCheckout)
	}
	if want := filepath.Join(realCheckout, "exercises"); cfg.ExercisesDir != want {
		t.Errorf("ExercisesDir = %q, want %q", cfg.ExercisesDir, want)
	}
	if want := filepath.Join(realCheckout, "data", "tracker.db"); cfg.DBPath != want {
		t.Errorf("DBPath = %q, want %q", cfg.DBPath, want)
	}
}

// TestLoad_ErrorsClearlyWhenCWDIsNotACheckoutAndNoCacheExists is the
// other half of the fallback: a genuinely fresh machine (never run
// ballroom from inside the checkout, so nothing is cached) must fail
// loud with an explanation, not silently boot into an empty picker.
func TestLoad_ErrorsClearlyWhenCWDIsNotACheckoutAndNoCacheExists(t *testing.T) {
	withRootCache(t) // isolated, empty cache
	chdir(t, t.TempDir())
	t.Setenv("PRACTICE_ROOT", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected Load to error when cwd isn't a checkout and nothing is cached")
	}
	if !strings.Contains(err.Error(), "docker/Dockerfile") {
		t.Errorf("Load err = %v, want it to explain the missing checkout", err)
	}
}

func TestLoad_RespectsRootEnvOverride(t *testing.T) {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Root != resolved {
		t.Errorf("Root = %q, want %q", cfg.Root, resolved)
	}
}

// TestLoad_RootEnvOverrideBypassesCheckoutValidation: PRACTICE_ROOT is
// trusted as-is even when it doesn't look like a real checkout (no
// docker/Dockerfile) -- the explicit escape hatch this codebase's own
// tests rely on to sandbox Load into a throwaway temp dir, distinct
// from the no-override fallback-to-cache path.
func TestLoad_RootEnvOverrideBypassesCheckoutValidation(t *testing.T) {
	withRootCache(t) // empty cache -- if this were consulted, Load would error
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v, want PRACTICE_ROOT trusted without checkout validation", err)
	}
	if cfg.Root != resolved {
		t.Errorf("Root = %q, want %q", cfg.Root, resolved)
	}
}

func TestLoad_DefaultDockerImage(t *testing.T) {
	t.Setenv("PRACTICE_ROOT", t.TempDir())
	t.Setenv("PRACTICE_DOCKER_IMAGE", "")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DockerImage != "ballroom-practice" {
		t.Errorf("DockerImage = %q, want %q", cfg.DockerImage, "ballroom-practice")
	}
}

func TestLoad_DockerImageOverride(t *testing.T) {
	t.Setenv("PRACTICE_ROOT", t.TempDir())
	t.Setenv("PRACTICE_DOCKER_IMAGE", "custom-image")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DockerImage != "custom-image" {
		t.Errorf("DockerImage = %q, want %q", cfg.DockerImage, "custom-image")
	}
}

func TestExercisePath(t *testing.T) {
	cfg := Config{ExercisesDir: "/root/exercises"}
	want := filepath.Join("/root/exercises", "two-pointers-01", "exercise.json")
	if got := cfg.ExercisePath("two-pointers-01"); got != want {
		t.Errorf("ExercisePath = %q, want %q", got, want)
	}
}

func TestTestsPath(t *testing.T) {
	cfg := Config{TestsDir: "/root/tests"}
	want := filepath.Join("/root/tests", "two-pointers-01")
	if got := cfg.TestsPath("two-pointers-01"); got != want {
		t.Errorf("TestsPath = %q, want %q", got, want)
	}
}

func TestReferencePath(t *testing.T) {
	cfg := Config{ExercisesDir: "/root/exercises"}
	want := filepath.Join("/root/exercises", "two-pointers-01", ".reference")
	if got := cfg.ReferencePath("two-pointers-01"); got != want {
		t.Errorf("ReferencePath = %q, want %q", got, want)
	}
}

func TestSettingsPath(t *testing.T) {
	cfg := Config{DataDir: "/root/data"}
	want := filepath.Join("/root/data", "settings.json")
	if got := cfg.SettingsPath(); got != want {
		t.Errorf("SettingsPath = %q, want %q", got, want)
	}
}

func TestLoadSettings_MissingFileReturnsZeroValueNotError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	s, err := LoadSettings(path)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if s.TutorModel != "" {
		t.Errorf("TutorModel = %q, want empty for a missing file", s.TutorModel)
	}
}

func TestLoadSettings_MalformedFileReturnsError(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if _, err := LoadSettings(path); err == nil {
		t.Fatal("expected an error for malformed settings JSON, got nil")
	}
}

func TestSaveSettings_ThenLoadRoundTrips(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	want := Settings{TutorModel: "llama3:8b"}
	if err := SaveSettings(path, want); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	got, err := LoadSettings(path)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if got != want {
		t.Errorf("LoadSettings = %+v, want %+v", got, want)
	}
}

func TestSaveSettings_ThenLoadRoundTripsOpenRouterAPIKey(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	want := Settings{TutorModel: "openrouter:anthropic/claude-3.5-sonnet", OpenRouterAPIKey: "sk-abc123"}
	if err := SaveSettings(path, want); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	got, err := LoadSettings(path)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if got != want {
		t.Errorf("LoadSettings = %+v, want %+v", got, want)
	}
}

func TestSaveSettings_CreatesParentDirIfMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "data", "settings.json")

	if err := SaveSettings(path, Settings{TutorModel: "x"}); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected settings file to exist at %q: %v", path, err)
	}
}

func TestLoad_DefaultsTutorModelWhenNoSettingsFile(t *testing.T) {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.TutorModel != DefaultTutorModel {
		t.Errorf("TutorModel = %q, want default %q", cfg.TutorModel, DefaultTutorModel)
	}
}

func TestQwen25Coder14BModel_IsValidOllamaTag(t *testing.T) {
	want := "qwen2.5-coder:14b-instruct"
	if Qwen25Coder14BModel != want {
		t.Errorf("Qwen25Coder14BModel = %q, want %q", Qwen25Coder14BModel, want)
	}
}

func TestDeepSeekCoderV2LiteModel_IsValidOllamaTag(t *testing.T) {
	want := "deepseek-coder-v2:16b-lite-instruct-q4_K_M"
	if DeepSeekCoderV2LiteModel != want {
		t.Errorf("DeepSeekCoderV2LiteModel = %q, want %q", DeepSeekCoderV2LiteModel, want)
	}
}

func TestLoad_ReadsPersistedTutorModel(t *testing.T) {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)

	settingsPath := filepath.Join(resolved, "data", "settings.json")
	if err := SaveSettings(settingsPath, Settings{TutorModel: "llama3:8b"}); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.TutorModel != "llama3:8b" {
		t.Errorf("TutorModel = %q, want %q", cfg.TutorModel, "llama3:8b")
	}
}

func TestLoad_ReadsPersistedOpenRouterAPIKey(t *testing.T) {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)
	t.Setenv("OPENROUTER_API_KEY", "")

	settingsPath := filepath.Join(resolved, "data", "settings.json")
	if err := SaveSettings(settingsPath, Settings{TutorModel: "llama3:8b", OpenRouterAPIKey: "sk-from-settings"}); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.OpenRouterAPIKey != "sk-from-settings" {
		t.Errorf("OpenRouterAPIKey = %q, want %q", cfg.OpenRouterAPIKey, "sk-from-settings")
	}
}

func TestLoad_FallsBackToOpenRouterAPIKeyEnvVarWhenNotInSettings(t *testing.T) {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)
	t.Setenv("OPENROUTER_API_KEY", "sk-from-env")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.OpenRouterAPIKey != "sk-from-env" {
		t.Errorf("OpenRouterAPIKey = %q, want %q", cfg.OpenRouterAPIKey, "sk-from-env")
	}
}

func TestLoad_SettingsOpenRouterAPIKeyTakesPrecedenceOverEnvVar(t *testing.T) {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)
	t.Setenv("OPENROUTER_API_KEY", "sk-from-env")

	settingsPath := filepath.Join(resolved, "data", "settings.json")
	if err := SaveSettings(settingsPath, Settings{TutorModel: "llama3:8b", OpenRouterAPIKey: "sk-from-settings"}); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.OpenRouterAPIKey != "sk-from-settings" {
		t.Errorf("OpenRouterAPIKey = %q, want the settings.json value %q to win over the env var", cfg.OpenRouterAPIKey, "sk-from-settings")
	}
}

func TestLoad_OpenRouterAPIKeyEmptyWhenNeitherSourceHasIt(t *testing.T) {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)
	t.Setenv("OPENROUTER_API_KEY", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.OpenRouterAPIKey != "" {
		t.Errorf("OpenRouterAPIKey = %q, want empty", cfg.OpenRouterAPIKey)
	}
}

func TestSaveSettings_ThenLoadRoundTripsGraderModel(t *testing.T) {
	dir := t.TempDir()
	settingsPath := filepath.Join(dir, "settings.json")

	want := Settings{TutorModel: "llama3:8b", GraderModel: "openrouter:tencent/hy3:free"}
	if err := SaveSettings(settingsPath, want); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}
	got, err := LoadSettings(settingsPath)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if got.GraderModel != want.GraderModel {
		t.Errorf("GraderModel = %q, want %q", got.GraderModel, want.GraderModel)
	}
}

func TestSaveSettings_ThenLoadRoundTripsOrchestratorModel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")

	want := Settings{TutorModel: "openrouter:openai/gpt-oss-120b:free", OrchestratorModel: "openrouter:nvidia/nemotron-3-nano-30b-a3b:free"}
	if err := SaveSettings(path, want); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	got, err := LoadSettings(path)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if got != want {
		t.Errorf("LoadSettings = %+v, want %+v", got, want)
	}
}

func TestLoad_ReadsPersistedOrchestratorModel(t *testing.T) {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)

	settingsPath := filepath.Join(resolved, "data", "settings.json")
	if err := SaveSettings(settingsPath, Settings{TutorModel: "llama3:8b", OrchestratorModel: "openrouter:nvidia/nemotron-3-nano-30b-a3b:free"}); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.OrchestratorModel != "openrouter:nvidia/nemotron-3-nano-30b-a3b:free" {
		t.Errorf("OrchestratorModel = %q, want %q", cfg.OrchestratorModel, "openrouter:nvidia/nemotron-3-nano-30b-a3b:free")
	}
}

func TestLoad_OrchestratorModelEmptyByDefault(t *testing.T) {
	dir := t.TempDir()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("EvalSymlinks: %v", err)
	}
	t.Setenv("PRACTICE_ROOT", resolved)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.OrchestratorModel != "" {
		t.Errorf("OrchestratorModel = %q, want empty (routing off by default)", cfg.OrchestratorModel)
	}
}

func TestSaveSettings_ThenLoadRoundTripsDefaultLanguageAndNotesToggle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	in := Settings{TutorModel: "m", DefaultLanguage: "go", DisableTutorNotes: true}
	if err := SaveSettings(path, in); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}
	got, err := LoadSettings(path)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if got.DefaultLanguage != "go" || !got.DisableTutorNotes {
		t.Errorf("round-trip = %+v, want DefaultLanguage go and DisableTutorNotes true", got)
	}
}

func TestLoad_ReadsPersistedDefaultLanguageAndNotesToggle(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PRACTICE_ROOT", dir)
	if err := SaveSettings(filepath.Join(dir, "data", "settings.json"), Settings{DefaultLanguage: "cpp", DisableTutorNotes: true}); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DefaultLanguage != "cpp" || !cfg.DisableTutorNotes {
		t.Errorf("cfg = DefaultLanguage %q DisableTutorNotes %v, want cpp/true", cfg.DefaultLanguage, cfg.DisableTutorNotes)
	}
}

// TestLoad_InvalidDefaultLanguageFailsLoud: a hand-edited settings.json
// with an unsupported language must fail Load rather than silently
// behaving as "ask every time" — the user asked for a default and
// would otherwise never learn it isn't being honored.
func TestLoad_InvalidDefaultLanguageFailsLoud(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PRACTICE_ROOT", dir)
	if err := SaveSettings(filepath.Join(dir, "data", "settings.json"), Settings{DefaultLanguage: "java"}); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "default_language") {
		t.Fatalf("Load err = %v, want an error naming default_language", err)
	}
}

// TestSaveSettings_ThenLoadRoundTripsTutorModeOverride and its Load
// sibling below mirror DefaultLanguage's exact round-trip/validation
// shape (issue #255's per-session tutor-mode override) -- same "" =
// exercise default, hand-edited garbage fails loud rather than silently
// falling back.
func TestSaveSettings_ThenLoadRoundTripsTutorModeOverride(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	in := Settings{TutorModel: "m", TutorModeOverride: "syntax-only"}
	if err := SaveSettings(path, in); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}
	got, err := LoadSettings(path)
	if err != nil {
		t.Fatalf("LoadSettings: %v", err)
	}
	if got.TutorModeOverride != "syntax-only" {
		t.Errorf("TutorModeOverride = %q, want %q", got.TutorModeOverride, "syntax-only")
	}
}

func TestLoad_ReadsPersistedTutorModeOverride(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PRACTICE_ROOT", dir)
	if err := SaveSettings(filepath.Join(dir, "data", "settings.json"), Settings{TutorModeOverride: "hints-first"}); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.TutorModeOverride != "hints-first" {
		t.Errorf("TutorModeOverride = %q, want %q", cfg.TutorModeOverride, "hints-first")
	}
}

func TestLoad_TutorModeOverrideEmptyByDefault(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PRACTICE_ROOT", dir)
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.TutorModeOverride != "" {
		t.Errorf("TutorModeOverride = %q, want empty (exercise default) when nothing persisted", cfg.TutorModeOverride)
	}
}

// TestLoad_InvalidTutorModeOverrideFailsLoud: same fail-loud contract
// as TestLoad_InvalidDefaultLanguageFailsLoud -- a hand-edited
// settings.json with an unsupported value must not silently be treated
// as "exercise default".
func TestLoad_InvalidTutorModeOverrideFailsLoud(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PRACTICE_ROOT", dir)
	if err := SaveSettings(filepath.Join(dir, "data", "settings.json"), Settings{TutorModeOverride: "godmode"}); err != nil {
		t.Fatalf("SaveSettings: %v", err)
	}
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "tutor_mode_override") {
		t.Fatalf("Load err = %v, want an error naming tutor_mode_override", err)
	}
}
