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
