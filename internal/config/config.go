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
// picked a model). Must support Ollama's structured tool_calls response
// field, since the tutor agent (internal/tutor) is built around real
// tool calling — confirmed via cmd/tutor-spike and cmd/tutor-eval that
// the previous default, qwen2.5-coder:7b, does NOT: it emits
// tool-call-shaped JSON as plain text content instead of a real
// structured call, so switched to llama3.1:8b, which does.
const DefaultTutorModel = "llama3.1:8b"

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
//
// NOT re-verified against the tool-calling requirement above (see
// DefaultTutorModel) — this was confirmed back when the tutor was a
// single plain chat call with no tools. Treat it as unverified for the
// current agent until checked with cmd/tutor-eval.
const DeepSeekCoderV2LiteModel = "deepseek-coder-v2:16b-lite-instruct-q4_K_M"

// Qwen25Coder14BModel is qwen2.5-coder:14b-instruct (9.0GB, q4_K_M
// quantization, 32K context) — DO NOT use as the tutor model. Confirmed
// via cmd/tutor-eval and a raw /api/chat repro (both live against a
// real Ollama 0.31.1) that despite being a larger, "-instruct"-tuned
// variant, it has the same real-tool-calling failure as qwen2.5-coder:7b
// (see DefaultTutorModel): CheckToolCalling reports false, and the full
// eval suite scored 0/8 on every real tool-calling scenario. Root cause
// traced to the wire level, not theorized — Ollama's own chat template
// for this model instructs it to wrap a tool call in
// <tool_call>...</tool_call> tags, which Ollama's server-side parser
// requires to populate the response's structured tool_calls field, but
// raw /api/chat calls (4/4 in a row) show the model consistently
// emitting a correctly-named, correctly-argued tool-call JSON body
// *without* those wrapper tags — so Ollama never recognizes it as a
// tool call at all, and it leaks through as plain message content
// instead. This is a property of this model's current Ollama packaging/
// template, not something fixable in this codebase's own prompts or
// code. This const is kept only so the tag stays typo-proof for anyone
// re-testing it after an Ollama/model update — re-run cmd/tutor-eval
// before trusting it again, don't assume this comment is still current.
const Qwen25Coder14BModel = "qwen2.5-coder:14b-instruct"

// Qwen2514BModel is qwen2.5:14b-instruct — the general-purpose sibling
// of Qwen25Coder14BModel, NOT the coder-tuned variant. Confirmed via
// CheckToolCalling and a full cmd/tutor-eval run (both live against a
// real Ollama 0.31.1) to not share the coder variant's tool-calling
// failure: 140/168 (83%) overall, including 7-8/8 on almost every real
// tool-calling scenario. Real, currently-open weakness found by that
// same run: hints-first mode leaks the forbidden technique name on the
// first ask meaningfully more often than this project's history
// documents for DefaultTutorModel (5/8, 1/8, 6/8 on the three
// hints-first checks) — treat that specific mode as less reliable on
// this model until prompts.go's hints-first instruction is
// strengthened and re-verified for it specifically.
const Qwen2514BModel = "qwen2.5:14b-instruct"

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
	TutorModel   string // Ollama tag, or an "openrouter:"-prefixed model slug (see internal/tutor.OpenRouterModelPrefix), passed to the container as TUTOR_MODEL
	// OrchestratorModel, when non-empty, enables per-turn routing (see
	// internal/tutor.Run): this model decides whether a turn needs
	// TutorModel's coding-specialist attention or can be answered
	// directly. Empty (the default) means no routing at all -- TutorModel
	// handles every turn by itself, identical to this project's
	// single-model behavior before routing existed.
	OrchestratorModel string
	// GraderModel, when non-empty, is the model design-submit grading
	// runs on (see tutor.GradeDesign) -- a role separate from worker and
	// orchestrator so grading judgment quality can be chosen
	// independently. Empty means grade with TutorModel.
	GraderModel string
	// OpenRouterAPIKey authenticates OpenRouter requests when TutorModel
	// or OrchestratorModel is openrouter:-prefixed; unused otherwise.
	// Resolved in Load: the persisted settings.json value if present,
	// else the OPENROUTER_API_KEY env var, else empty (not an error at
	// Load time -- only matters if an openrouter: model is actually
	// used). One key authenticates every model on the account, so this
	// stays a single shared field even with two model roles.
	OpenRouterAPIKey string
	// DefaultLanguage, when non-empty (one of "python"/"go"/"cpp"),
	// skips the language-picker step for problems that have a variant
	// in that language; empty means ask every time (the original
	// behavior). Design/behavioral problems ride coach/interviewer in
	// the language slot, so they never match and always ask.
	DefaultLanguage string
	// VoiceModel selects the whisper model `ballroom voice` transcribes
	// with -- one of tiny/base/small/medium (see cmd/ballroom's
	// voicemodel.go). Empty means base, what existing installs have.
	VoiceModel string
	// DisableTutorNotes removes the tutor's highlight_lines tool, so
	// sessions run without editor highlights/notes at all — a durable
	// alternative to the in-session M-h rendering toggle.
	DisableTutorNotes bool
	// TutorModeOverride, when non-empty (one of "syntax-only"/
	// "hints-first"/"full-assist" -- the exercise.TutorMode* coding
	// vocabulary, duplicated here as raw literals rather than importing
	// internal/exercise, matching DefaultLanguage's own precedent of
	// keeping this package leaf-level), overrides a coding exercise's
	// own tutor_mode for every session until changed again; empty means
	// use whatever each exercise.json declares (the original behavior).
	// Deliberately exercise.Kind-gated at the one place that applies it
	// (internal/orchestrator.exerciseRunArgs) rather than here: a design
	// or behavioral exercise's TutorMode is its session persona
	// (interviewer, story coach, ...), not a coding assistance level, so
	// it must never be overridden by this setting.
	TutorModeOverride string
}

// Settings holds user preferences persisted across invocations, e.g. the
// last model picked in the TUI's model picker.
type Settings struct {
	TutorModel string `json:"tutor_model"`
	// OrchestratorModel mirrors Config.OrchestratorModel -- see its doc
	// comment. No env-var fallback (unlike OpenRouterAPIKey): this isn't
	// a secret, and routing being on/off should be an explicit, durable
	// choice rather than something that flips based on ambient env vars.
	OrchestratorModel string `json:"orchestrator_model"`
	// GraderModel mirrors Config.GraderModel -- the design-submit
	// grading role. Same no-env-fallback rationale as OrchestratorModel.
	GraderModel string `json:"grader_model"`
	// OpenRouterAPIKey is saved here so the TUI's model picker only ever
	// needs to ask for it once (see internal/tui/app.go's key-entry
	// stage) instead of requiring OPENROUTER_API_KEY to be exported in
	// the shell every session.
	OpenRouterAPIKey string `json:"openrouter_api_key"`
	// DefaultLanguage mirrors Config.DefaultLanguage ("" = ask every
	// time). Validated on Load, not here — LoadSettings stays a dumb
	// reader so a bad value still round-trips visibly.
	DefaultLanguage string `json:"default_language,omitempty"`
	// VoiceModel mirrors Config.VoiceModel.
	VoiceModel string `json:"voice_model,omitempty"`
	// DisableTutorNotes mirrors Config.DisableTutorNotes (issue #25's
	// durable variant: notes off at the source, not just hidden).
	DisableTutorNotes bool `json:"disable_tutor_notes,omitempty"`
	// TutorModeOverride mirrors Config.TutorModeOverride ("" = each
	// exercise's own tutor_mode). Validated on Load, not here — same
	// dumb-reader rationale as DefaultLanguage.
	TutorModeOverride string `json:"tutor_mode_override,omitempty"`
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
//
// An explicit PRACTICE_ROOT is trusted as-is, no questions asked — the
// documented escape hatch, and how this codebase's own tests sandbox
// Load into a throwaway temp dir. Without it, the working directory
// must actually look like a ballroom checkout (see ResolveRoot), else
// this falls back to the last checkout root cached from a previous run.
// Without that fallback, `ballroom` launched from PATH outside the
// checkout (its usual location once installed via `go install`) would
// silently resolve ExercisesDir/TestsDir/DataDir under whatever
// directory it happened to be launched from (e.g. $HOME) — an empty
// picker and "couldn't load your progress", not a working boot.
// ToSettings projects a Config back into the persistable Settings
// struct. It exists so the many callers that save settings don't each
// hand-copy the fields -- a copy that dropped one silently wiped it,
// which is exactly how a set-model command once erased a saved API key.
// Every setter goes through here now, so a new field is added in one
// place.
func (c Config) ToSettings() Settings {
	return Settings{
		TutorModel:        c.TutorModel,
		OrchestratorModel: c.OrchestratorModel,
		GraderModel:       c.GraderModel,
		OpenRouterAPIKey:  c.OpenRouterAPIKey,
		DefaultLanguage:   c.DefaultLanguage,
		VoiceModel:        c.VoiceModel,
		DisableTutorNotes: c.DisableTutorNotes,
		TutorModeOverride: c.TutorModeOverride,
	}
}

func Load() (Config, error) {
	var root string
	if envRoot := os.Getenv("PRACTICE_ROOT"); envRoot != "" {
		resolved, err := filepath.EvalSymlinks(envRoot)
		if err != nil {
			return Config{}, fmt.Errorf("config: resolve root %s: %w", envRoot, err)
		}
		root = resolved
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return Config{}, fmt.Errorf("config: getwd: %w", err)
		}
		resolved, err := filepath.EvalSymlinks(wd)
		if err != nil {
			return Config{}, fmt.Errorf("config: resolve root %s: %w", wd, err)
		}
		root, err = ResolveRoot(resolved)
		if err != nil {
			return Config{}, err
		}
	}

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
	cfg.OrchestratorModel = settings.OrchestratorModel
	cfg.GraderModel = settings.GraderModel
	cfg.OpenRouterAPIKey = settings.OpenRouterAPIKey
	if cfg.OpenRouterAPIKey == "" {
		cfg.OpenRouterAPIKey = os.Getenv("OPENROUTER_API_KEY")
	}
	// Fail loud on an unsupported default_language (hand-edited file):
	// silently treating it as "ask every time" would leave the user
	// wondering why their default is never honored.
	switch settings.DefaultLanguage {
	case "", "python", "go", "cpp":
		cfg.DefaultLanguage = settings.DefaultLanguage
	default:
		return Config{}, fmt.Errorf("config: default_language %q in %s: want python, go, cpp, or empty", settings.DefaultLanguage, cfg.SettingsPath())
	}
	cfg.VoiceModel = settings.VoiceModel
	cfg.DisableTutorNotes = settings.DisableTutorNotes
	// Same fail-loud contract as default_language above: a hand-edited,
	// unsupported value must not be silently treated as "exercise
	// default". The vocabulary (syntax-only/hints-first/full-assist) is
	// exercise.TutorMode*, duplicated as literals rather than imported —
	// see TutorModeOverride's doc comment.
	switch settings.TutorModeOverride {
	case "", "syntax-only", "hints-first", "full-assist":
		cfg.TutorModeOverride = settings.TutorModeOverride
	default:
		return Config{}, fmt.Errorf("config: tutor_mode_override %q in %s: want syntax-only, hints-first, full-assist, or empty", settings.TutorModeOverride, cfg.SettingsPath())
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

// ReferencePath returns the path to exercise <id>'s reference solution
// directory -- a sibling of its repo/, checked into git but (like
// TestsPath's hidden tests) never mounted into a session's workspace
// until requested (see orchestrator.WaitAndRevealReference). Same
// ".reference" directory internal/verify already reads to confirm a
// known-correct solution passes each exercise's hidden tests.
func (c Config) ReferencePath(id string) string {
	return filepath.Join(c.ExercisesDir, id, ".reference")
}
