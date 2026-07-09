package config

import (
	"os"
	"path/filepath"
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

func TestLoad_DefaultsToCWD(t *testing.T) {
	dir := t.TempDir()
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

func TestLoad_DefaultDockerImage(t *testing.T) {
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
