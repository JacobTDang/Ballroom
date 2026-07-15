package orchestrator

import (
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/config"
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
