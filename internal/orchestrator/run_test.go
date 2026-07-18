package orchestrator

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/draft"
	"github.com/JacobTDang/Ballroom/internal/exercise"
)

func containsFlag(args []string, want string) bool {
	for _, a := range args {
		if a == want {
			return true
		}
	}
	return false
}

func TestExerciseRunArgs_IncludesTutorModelEnvFromConfig(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Category: "pattern", Language: "go", TutorMode: "hint", TestCommand: "go test ./..."}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "-e") || !containsFlag(args, "TUTOR_MODEL=llama3:8b") {
		t.Errorf("expected TUTOR_MODEL=llama3:8b to be passed as an -e flag, got %v", args)
	}
}

func TestExerciseRunArgs_DifferentModelProducesDifferentFlag(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "qwen2.5-coder:7b"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Category: "pattern", Language: "go", TutorMode: "hint", TestCommand: "go test ./..."}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "TUTOR_MODEL=qwen2.5-coder:7b") {
		t.Errorf("expected TUTOR_MODEL=qwen2.5-coder:7b, got %v", args)
	}
	if containsFlag(args, "TUTOR_MODEL=llama3:8b") {
		t.Errorf("did not expect a stale model flag, got %v", args)
	}
}

func TestExerciseRunArgs_StillIncludesExistingPracticeEnvVars(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Category: "pattern", Language: "go", TutorMode: "hint", TestCommand: "go test ./..."}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "PRACTICE_EXERCISE_ID=two-pointers-01-go") {
		t.Errorf("expected PRACTICE_EXERCISE_ID to still be set, got %v", args)
	}
	if !containsFlag(args, "ballroom-practice") {
		t.Errorf("expected the docker image to still be the final arg, got %v", args)
	}
}

func TestExerciseRunArgs_ForwardsKind(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b"}
	ex := exercise.Exercise{ID: "url-shortener-01-coach", Kind: exercise.KindDesign, Category: "system-design", Language: "coach", TutorMode: "design-coach"}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "PRACTICE_KIND=design") {
		t.Errorf("expected PRACTICE_KIND=design so the in-container submit takes the design path, got %v", args)
	}
}

// TestExerciseRunArgs_TutorModeOverride* cover issue #255's per-session
// tutor-mode override: cfg.TutorModeOverride must win over a coding
// exercise's own TutorMode, but a design/behavioral exercise's TutorMode
// IS its session persona (interviewer, story coach, ...), not a coding
// assistance level, so it must never be overridden.

func TestExerciseRunArgs_TutorModeOverride_AppliesToCodingExercise(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b", TutorModeOverride: "syntax-only"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Kind: exercise.KindCoding, Category: "pattern", Language: "go", TutorMode: "full-assist", TestCommand: "go test ./..."}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "PRACTICE_TUTOR_MODE=syntax-only") {
		t.Errorf("expected the override to win over the exercise's own full-assist TutorMode, got %v", args)
	}
}

func TestExerciseRunArgs_TutorModeOverride_EmptyUsesExerciseDefault(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Kind: exercise.KindCoding, Category: "pattern", Language: "go", TutorMode: "hints-first", TestCommand: "go test ./..."}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "PRACTICE_TUTOR_MODE=hints-first") {
		t.Errorf("expected the exercise's own TutorMode with no override set, got %v", args)
	}
}

func TestExerciseRunArgs_TutorModeOverride_NeverAppliesToDesignExercise(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b", TutorModeOverride: "syntax-only"}
	ex := exercise.Exercise{ID: "url-shortener-01-coach", Kind: exercise.KindDesign, Category: "system-design", Language: "coach", TutorMode: "design-coach"}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "PRACTICE_TUTOR_MODE=design-coach") {
		t.Errorf("expected a design exercise's persona (TutorMode) to stay design-coach, unaffected by the coding-only override, got %v", args)
	}
	if containsFlag(args, "PRACTICE_TUTOR_MODE=syntax-only") {
		t.Errorf("did not expect the coding-only override to leak into a design session, got %v", args)
	}
}

func TestExerciseRunArgs_TutorModeOverride_NeverAppliesToBehavioralExercise(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b", TutorModeOverride: "full-assist"}
	ex := exercise.Exercise{ID: "disagreement-01-coach", Kind: exercise.KindDesign, Category: "behavioral", Language: "coach", TutorMode: "story-coach"}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "PRACTICE_TUTOR_MODE=story-coach") {
		t.Errorf("expected a behavioral exercise's persona (TutorMode) to stay story-coach, unaffected by the coding-only override, got %v", args)
	}
}

func TestExerciseRunArgs_ForwardsGraderModel(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b", GraderModel: "openrouter:tencent/hy3:free"}
	ex := exercise.Exercise{ID: "url-shortener-01-coach", Kind: exercise.KindDesign, Category: "system-design", Language: "coach", TutorMode: "design-coach"}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "TUTOR_GRADER_MODEL=openrouter:tencent/hy3:free") {
		t.Errorf("expected TUTOR_GRADER_MODEL forwarded for design grading, got %v", args)
	}
}

func TestExerciseRunArgs_ForwardsTimeLimit(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b"}
	ex := exercise.Exercise{ID: "url-shortener-01-interviewer", Kind: exercise.KindDesign, Category: "system-design", Language: "interviewer", TutorMode: "interviewer", TimeLimitMin: 45}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "PRACTICE_TIME_LIMIT_MIN=45") {
		t.Errorf("expected PRACTICE_TIME_LIMIT_MIN forwarded for the interview clock, got %v", args)
	}
}

func TestRunArgs_ForwardTutorStreamOverrideFromHostEnv(t *testing.T) {
	t.Setenv("TUTOR_STREAM", "off")
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "openrouter:some/model"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Category: "pattern", Language: "go", TutorMode: "hint", TestCommand: "go test ./..."}

	if args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now()); !containsFlag(args, "TUTOR_STREAM=off") {
		t.Errorf("expected the host's TUTOR_STREAM override forwarded into the exercise container, got %v", args)
	}
	if args := sandboxRunArgs(cfg); !containsFlag(args, "TUTOR_STREAM=off") {
		t.Errorf("expected the host's TUTOR_STREAM override forwarded into the sandbox container, got %v", args)
	}
}

func TestSandboxRunArgs_MarksTheSandboxSessionContext(t *testing.T) {
	cfg := config.Config{DockerImage: "ballroom-practice", TutorModel: "llama3:8b"}

	args := sandboxRunArgs(cfg)

	if !containsFlag(args, "PRACTICE_SANDBOX=1") {
		t.Errorf("expected PRACTICE_SANDBOX=1 so `ballroom return` works inside a sandbox, got %v", args)
	}
}

func TestSandboxRunArgs_IncludesTutorModelEnvFromConfig(t *testing.T) {
	cfg := config.Config{DockerImage: "ballroom-practice", TutorModel: "llama3:8b"}

	args := sandboxRunArgs(cfg)

	if !containsFlag(args, "TUTOR_MODEL=llama3:8b") {
		t.Errorf("expected TUTOR_MODEL=llama3:8b to be passed as an -e flag, got %v", args)
	}
	if !containsFlag(args, "ballroom-practice") {
		t.Errorf("expected the docker image to still be the final arg, got %v", args)
	}
}

func TestExerciseRunArgs_ForwardsOpenRouterAPIKeyFromConfig(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "openrouter:anthropic/claude-3.5-sonnet", OpenRouterAPIKey: "sk-test-key"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Category: "pattern", Language: "go", TutorMode: "hint", TestCommand: "go test ./..."}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "OPENROUTER_API_KEY=sk-test-key") {
		t.Errorf("expected OPENROUTER_API_KEY=sk-test-key to be passed as an -e flag, got %v", args)
	}
}

func TestSandboxRunArgs_ForwardsOpenRouterAPIKeyFromConfig(t *testing.T) {
	cfg := config.Config{DockerImage: "ballroom-practice", TutorModel: "openrouter:anthropic/claude-3.5-sonnet", OpenRouterAPIKey: "sk-test-key"}

	args := sandboxRunArgs(cfg)

	if !containsFlag(args, "OPENROUTER_API_KEY=sk-test-key") {
		t.Errorf("expected OPENROUTER_API_KEY=sk-test-key to be passed as an -e flag, got %v", args)
	}
}

func TestExerciseRunArgs_ForwardsOrchestratorModelFromConfig(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b", OrchestratorModel: "nemotron"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Category: "pattern", Language: "go", TutorMode: "hint", TestCommand: "go test ./..."}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "TUTOR_ORCHESTRATOR_MODEL=nemotron") {
		t.Errorf("expected TUTOR_ORCHESTRATOR_MODEL=nemotron to be passed as an -e flag, got %v", args)
	}
}

func TestExerciseRunArgs_ForwardsEmptyOrchestratorModelWhenRoutingIsOff(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Category: "pattern", Language: "go", TutorMode: "hint", TestCommand: "go test ./..."}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())

	if !containsFlag(args, "TUTOR_ORCHESTRATOR_MODEL=") {
		t.Errorf("expected TUTOR_ORCHESTRATOR_MODEL= (empty, routing off) to still be passed as an -e flag, got %v", args)
	}
}

func TestSandboxRunArgs_ForwardsOrchestratorModelFromConfig(t *testing.T) {
	cfg := config.Config{DockerImage: "ballroom-practice", TutorModel: "llama3:8b", OrchestratorModel: "nemotron"}

	args := sandboxRunArgs(cfg)

	if !containsFlag(args, "TUTOR_ORCHESTRATOR_MODEL=nemotron") {
		t.Errorf("expected TUTOR_ORCHESTRATOR_MODEL=nemotron to be passed as an -e flag, got %v", args)
	}
}

func TestExerciseRunArgs_ForwardsVideoURL(t *testing.T) {
	cfg := config.Config{DataDir: "/data", DockerImage: "ballroom-practice", TutorModel: "llama3:8b"}
	ex := exercise.Exercise{ID: "two-pointers-01-go", Category: "pattern", Language: "go", TutorMode: "hint", TestCommand: "go test ./...", VideoURL: "https://youtu.be/abc123"}

	args := exerciseRunArgs(cfg, ex, "/control", "/workspace", time.Now())
	if !containsFlag(args, "PRACTICE_VIDEO_URL=https://youtu.be/abc123") {
		t.Errorf("expected the video url forwarded, got %v", args)
	}
}

func TestExerciseRunArgs_ForwardsTutorNotesToggle(t *testing.T) {
	cfg := config.Config{DockerImage: "img", DisableTutorNotes: true}
	args := exerciseRunArgs(cfg, exercise.Exercise{}, "/c", "/w", time.Now())
	if !containsFlag(args, "PRACTICE_TUTOR_NOTES=off") {
		t.Errorf("args = %v, want PRACTICE_TUTOR_NOTES=off forwarded when notes are disabled", args)
	}
	cfg.DisableTutorNotes = false
	args = exerciseRunArgs(cfg, exercise.Exercise{}, "/c", "/w", time.Now())
	if !containsFlag(args, "PRACTICE_TUTOR_NOTES=") {
		t.Errorf("args = %v, want the empty always-forwarded PRACTICE_TUTOR_NOTES", args)
	}
}

func TestSandboxRunArgs_ForwardsTutorNotesToggle(t *testing.T) {
	args := sandboxRunArgs(config.Config{DockerImage: "img", DisableTutorNotes: true})
	if !containsFlag(args, "PRACTICE_TUTOR_NOTES=off") {
		t.Errorf("args = %v, want PRACTICE_TUTOR_NOTES=off in sandbox args too", args)
	}
}

// --- issue #231: idempotent session cleanup ---

func TestNewSessionFinalizer_SnapshotsAndRemovesBothDirsOnFirstCall(t *testing.T) {
	dataDir := t.TempDir()
	workspaceDir := t.TempDir()
	controlDir := t.TempDir()
	exerciseID := "finalizer-test"

	if err := os.WriteFile(filepath.Join(workspaceDir, "solution.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("seed workspace: %v", err)
	}

	cleanupCalls := 0
	cleanupWorkspace := func() { cleanupCalls++; os.RemoveAll(workspaceDir) }

	finalize := newSessionFinalizer(config.Config{DataDir: dataDir}, exerciseID, workspaceDir, controlDir, cleanupWorkspace)
	finalize()

	if cleanupCalls != 1 {
		t.Errorf("expected cleanupWorkspace called exactly once, got %d", cleanupCalls)
	}
	if _, err := os.Stat(controlDir); !os.IsNotExist(err) {
		t.Errorf("expected controlDir removed, stat err = %v", err)
	}
	draftFile := filepath.Join(draft.Dir(dataDir, exerciseID), "solution.go")
	if got, err := os.ReadFile(draftFile); err != nil {
		t.Errorf("expected a final draft snapshot written to %s: %v", draftFile, err)
	} else if string(got) != "package main" {
		t.Errorf("draft content = %q, want %q", got, "package main")
	}
}

// TestNewSessionFinalizer_SecondCallIsANoOp is issue #231's explicit
// idempotency requirement: RunExercise's normal defer and its signal
// handler can both reach the same finalizer (e.g. a signal arriving the
// same moment the container exits on its own), and only the first call
// should do anything.
func TestNewSessionFinalizer_SecondCallIsANoOp(t *testing.T) {
	dataDir := t.TempDir()
	workspaceDir := t.TempDir()
	controlDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(workspaceDir, "solution.py"), []byte("x = 1"), 0o644); err != nil {
		t.Fatalf("seed workspace: %v", err)
	}

	cleanupCalls := 0
	cleanupWorkspace := func() { cleanupCalls++ }

	finalize := newSessionFinalizer(config.Config{DataDir: dataDir}, "finalizer-test-2", workspaceDir, controlDir, cleanupWorkspace)

	finalize()
	finalize() // must not error, panic, or re-run the work

	if cleanupCalls != 1 {
		t.Errorf("expected cleanupWorkspace called exactly once across two finalize() calls, got %d", cleanupCalls)
	}
}
