// Package config holds shared paths and settings (exercises dir, tests dir,
// data dir, Docker image name) used across the launcher.
package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const defaultDockerImage = "ballroom-practice"

// Config holds resolved filesystem paths and settings for one invocation
// of the launcher.
type Config struct {
	Root         string // repo root
	ExercisesDir string // Root/exercises
	TestsDir     string // Root/tests (hidden tests, never mounted until submit)
	DataDir      string // Root/data
	DBPath       string // DataDir/tracker.db
	DockerImage  string
}

// Load resolves Config from the environment. Root defaults to the current
// working directory, overridable via PRACTICE_ROOT (paths are resolved
// through EvalSymlinks so tests/comparisons aren't tripped up by symlinked
// temp dirs, e.g. /tmp -> /private/tmp on macOS).
func Load() (Config, error) {
	root := os.Getenv("PRACTICE_ROOT")
	if root == "" {
		wd, err := os.Getwd()
		if err != nil {
			return Config{}, fmt.Errorf("config: getwd: %w", err)
		}
		root = wd
	}
	resolved, err := filepath.EvalSymlinks(root)
	if err != nil {
		return Config{}, fmt.Errorf("config: resolve root %s: %w", root, err)
	}
	root = resolved

	image := os.Getenv("PRACTICE_DOCKER_IMAGE")
	if image == "" {
		image = defaultDockerImage
	}

	return Config{
		Root:         root,
		ExercisesDir: filepath.Join(root, "exercises"),
		TestsDir:     filepath.Join(root, "tests"),
		DataDir:      filepath.Join(root, "data"),
		DBPath:       filepath.Join(root, "data", "tracker.db"),
		DockerImage:  image,
	}, nil
}

// ExercisePath returns the path to exercise <id>'s definition file.
func (c Config) ExercisePath(id string) string {
	return filepath.Join(c.ExercisesDir, id, "exercise.json")
}

// TestsPath returns the path to exercise <id>'s hidden test directory.
func (c Config) TestsPath(id string) string {
	return filepath.Join(c.TestsDir, id)
}
