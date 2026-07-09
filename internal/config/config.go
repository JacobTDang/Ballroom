// Package config holds shared paths and settings (exercises dir, tests dir,
// data dir, Docker image name) used across the launcher.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const defaultDockerImage = "ballroom-practice"

// DefaultTutorModel is the Ollama model used when nothing has been
// persisted to settings.json yet (first run, before the user has ever
// picked a model). Must match tutor/chat.sh's own fallback so the two
// stay in sync.
const DefaultTutorModel = "qwen2.5-coder:7b"

// DeepSeekCoderV2LiteModel is a second tutor model confirmed to work
// end-to-end (preflight check, TUTOR_MODEL wiring, and a real chat
// round-trip) — deepseek-coder-v2:16b-lite-instruct-q4_K_M, verified
// against the tag list at ollama.com/library/deepseek-coder-v2/tags.
// It is not the default. There's no fixed "supported models" list in
// this codebase: the model picker (internal/tui/modelpicker.go) already
// accepts any locally pulled or freely typed Ollama tag, so selecting
// this one is just a matter of typing it there or running
// `ollama pull deepseek-coder-v2:16b-lite-instruct-q4_K_M` first. This
// const exists purely so the verified tag is documented and typo-proof
// (e.g. for scripting a pull, or referencing in tests) rather than
// re-typed from memory.
const DeepSeekCoderV2LiteModel = "deepseek-coder-v2:16b-lite-instruct-q4_K_M"

// Qwen25Coder14BModel is a second tutor model confirmed to work end-to-end
// (preflight check, TUTOR_MODEL wiring, and a real chat round-trip) —
// qwen2.5-coder:14b-instruct, verified against the tag list at
// ollama.com/library/qwen2.5-coder/tags (9.0GB, q4_K_M quantization,
// 32K context). It is not the default. There's no fixed "supported
// models" list in this codebase: the model picker
// (internal/tui/modelpicker.go) already accepts any locally pulled or
// freely typed Ollama tag, so selecting this one is just a matter of
// typing it there or running `ollama pull qwen2.5-coder:14b-instruct`
// first. This const exists purely so the verified tag is documented and
// typo-proof (e.g. for scripting a pull, or referencing in tests) rather
// than re-typed from memory.
//
// Hardware note: this is a meaningfully larger model than
// DefaultTutorModel (9.0GB on disk vs. ~4.7GB for the 7B default).
// Budget roughly 12-16GB of free RAM/VRAM for comfortable inference at
// this quantization (model weights plus KV cache headroom) — pulling or
// selecting it on a machine with less will be slow or may fail to load.
const Qwen25Coder14BModel = "qwen2.5-coder:14b-instruct"

// settingsFileName is the persisted user-settings file, stored under
// Config.DataDir alongside tracker.db.
const settingsFileName = "settings.json"

// Config holds resolved filesystem paths and settings for one invocation
// of the launcher.
type Config struct {
	Root         string // repo root
	ExercisesDir string // Root/exercises
	TestsDir     string // Root/tests (hidden tests, never mounted until submit)
	DataDir      string // Root/data
	DBPath       string // DataDir/tracker.db
	DockerImage  string
	TutorModel   string // Ollama model tag passed to the container as TUTOR_MODEL
}

// Settings holds user preferences persisted across invocations, e.g. the
// last Ollama model picked in the TUI's model picker.
type Settings struct {
	TutorModel string `json:"tutor_model"`
}

// SettingsPath returns the path to the persisted settings file.
func (c Config) SettingsPath() string {
	return filepath.Join(c.DataDir, settingsFileName)
}

// LoadSettings reads persisted settings from path. A missing file returns
// a zero-value Settings, not an error — there's simply nothing persisted
// yet on first run. A present-but-malformed file is a real error (fail
// loud rather than silently discarding whatever the user last picked).
func LoadSettings(path string) (Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Settings{}, nil
		}
		return Settings{}, fmt.Errorf("config: read settings: %w", err)
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return Settings{}, fmt.Errorf("config: parse settings %s: %w", path, err)
	}
	return s, nil
}

// SaveSettings persists s to path, creating the parent directory if it
// doesn't exist yet.
func SaveSettings(path string, s Settings) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("config: create settings dir: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("config: marshal settings: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("config: write settings: %w", err)
	}
	return nil
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

	cfg := Config{
		Root:         root,
		ExercisesDir: filepath.Join(root, "exercises"),
		TestsDir:     filepath.Join(root, "tests"),
		DataDir:      filepath.Join(root, "data"),
		DBPath:       filepath.Join(root, "data", "tracker.db"),
		DockerImage:  image,
	}

	settings, err := LoadSettings(cfg.SettingsPath())
	if err != nil {
		return Config{}, err
	}
	cfg.TutorModel = settings.TutorModel
	if cfg.TutorModel == "" {
		cfg.TutorModel = DefaultTutorModel
	}

	return cfg, nil
}

// ExercisePath returns the path to exercise <id>'s definition file.
func (c Config) ExercisePath(id string) string {
	return filepath.Join(c.ExercisesDir, id, "exercise.json")
}

// TestsPath returns the path to exercise <id>'s hidden test directory.
func (c Config) TestsPath(id string) string {
	return filepath.Join(c.TestsDir, id)
}
