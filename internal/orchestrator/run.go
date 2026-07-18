// Package orchestrator drives the exercise/sandbox lifecycle: mounting the
// exercise repo into the unified Docker image, starting/stopping the timer,
// and revealing the hidden test suite only on submit.
package orchestrator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
)

const revealPollInterval = 200 * time.Millisecond

const sandboxVolume = "ballroom-sandbox"

// RunExercise starts a graded, timed session and blocks until the
// container exits. ex.RepoPath (the permanent exercise source) is never
// mounted directly — PrepareWorkspace copies it into a disposable temp
// dir that gets mounted as /workspace and deleted on exit, so nothing
// written during the session (edits, or hidden tests revealed on submit
// — see WaitAndReveal) can leak back into the source repo.
func RunExercise(cfg config.Config, ex exercise.Exercise) error {
	if err := EnsureImage(cfg); err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return fmt.Errorf("orchestrator: create data dir: %w", err)
	}

	workspaceDir, cleanupWorkspace, err := PrepareWorkspace(ex.RepoPath, ex.VideoURL)
	if err != nil {
		return err
	}
	defer cleanupWorkspace()

	controlDir, err := os.MkdirTemp("", "practice-control-")
	if err != nil {
		return fmt.Errorf("orchestrator: create control dir: %w", err)
	}
	defer os.RemoveAll(controlDir)

	startedAt := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	revealErr := make(chan error, 1)
	go func() {
		revealErr <- WaitAndReveal(ctx, controlDir, cfg.TestsPath(ex.ID), workspaceDir, revealPollInterval)
	}()

	args := exerciseRunArgs(cfg, ex, controlDir, workspaceDir, startedAt)

	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	runErr := cmd.Run()

	cancel()
	if err := <-revealErr; err != nil && ctx.Err() == nil {
		// Only surface reveal-side errors that happened for a reason other
		// than us cancelling the context on normal container exit.
		fmt.Fprintf(os.Stderr, "orchestrator: reveal watcher: %v\n", err)
	}

	if runErr != nil {
		return fmt.Errorf("orchestrator: docker run: %w", runErr)
	}
	return nil
}

// RunSandbox mounts a persistent volume (survives across runs) and starts
// an ungraded, untimed session — scratch space for trying something out
// without it landing in the tracker.
//
// Reset is manual by design (no scripted "reset to base" for MVP): wipe
// the volume with `docker volume rm ballroom-sandbox` — the next
// `ballroom sandbox` recreates it empty.
func RunSandbox(cfg config.Config) error {
	if err := EnsureImage(cfg); err != nil {
		return err
	}

	args := sandboxRunArgs(cfg)

	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("orchestrator: docker run --sandbox: %w", err)
	}
	return nil
}

// exerciseRunArgs builds the `docker run` argument list for a graded
// exercise session. Pulled out as a pure function (no docker/filesystem
// side effects) so the flags it produces — notably TUTOR_MODEL, sourced
// from cfg.TutorModel rather than a hardcoded const — are unit-testable
// without shelling out to docker.
func exerciseRunArgs(cfg config.Config, ex exercise.Exercise, controlDir, workspaceDir string, startedAt time.Time) []string {
	return []string{
		"run", "-it", "--rm",
		// Labeled so host-side tooling (ballroom voice) can find the
		// running session container to docker-exec into.
		"--label", sessionContainerLabel + "=1",
		"-v", workspaceDir + ":/workspace",
		"-v", cfg.DataDir + ":/data",
		"-v", controlDir + ":/control",
		"-e", "PRACTICE_CONTROL_DIR=/control",
		"-e", "PRACTICE_WORKSPACE_DIR=/workspace",
		"-e", "PRACTICE_TEST_COMMAND=" + ex.TestCommand,
		"-e", "PRACTICE_EXERCISE_ID=" + ex.ID,
		"-e", "PRACTICE_CATEGORY=" + ex.Category,
		"-e", "PRACTICE_LANGUAGE=" + ex.Language,
		"-e", "PRACTICE_KIND=" + ex.Kind,
		"-e", "PRACTICE_TUTOR_MODE=" + ex.TutorMode,
		"-e", "PRACTICE_TIME_LIMIT_MIN=" + strconv.Itoa(ex.TimeLimitMin),
		"-e", "PRACTICE_STARTED_AT=" + startedAt.Format(time.RFC3339),
		"-e", "PRACTICE_DB_PATH=/data/tracker.db",
		"-e", "TUTOR_MODEL=" + cfg.TutorModel,
		// Always forwarded, empty or not -- the tutor process only reads
		// it when TutorModel is tutor.OpenRouterModelPrefix-prefixed
		// (see cmd/ballroom/main.go's tutorCmd), same as TUTOR_MODEL
		// being forwarded even when it's just the default.
		"-e", "OPENROUTER_API_KEY=" + cfg.OpenRouterAPIKey,
		// Always forwarded, empty or not -- an empty value means routing
		// is off (see cmd/ballroom/main.go's tutorCmd and
		// internal/tutor.Config.OrchestratorModel), same rationale as
		// OPENROUTER_API_KEY above.
		"-e", "TUTOR_ORCHESTRATOR_MODEL=" + cfg.OrchestratorModel,
		// Same always-forwarded contract -- empty means design grading
		// uses the worker model (cmd/ballroom's graderModelFromEnv).
		"-e", "TUTOR_GRADER_MODEL=" + cfg.GraderModel,
		// exercise.json never enters the container, so the submit-time
		// solution-video line rides an env var; empty means none.
		"-e", "PRACTICE_VIDEO_URL=" + ex.VideoURL,
		// Host env, not a persisted setting: TUTOR_STREAM=on|off is the
		// per-invocation escape hatch over internal/tutor's
		// streamingEnabled default (stream OpenRouter replies, never
		// Ollama's). Forwarded because the tutor process that reads it
		// runs inside the container; empty means "use the default".
		"-e", "TUTOR_STREAM=" + os.Getenv("TUTOR_STREAM"),
		// Always forwarded, empty or "off" -- the Settings toggle that
		// removes the tutor's editor highlight/note tool at the source
		// (cmd/ballroom's tutorCmd reads it into
		// tutor.Config.DisableEditorNotes).
		"-e", "PRACTICE_TUTOR_NOTES=" + tutorNotesEnvValue(cfg),
		cfg.DockerImage,
	}
}

// tutorNotesEnvValue renders cfg.DisableTutorNotes for the container
// env: "off" disables the tool, empty means the default (enabled) --
// an env var can't carry a Go bool, and "off" reads clearer in `docker
// inspect` output than "false"/"0".
func tutorNotesEnvValue(cfg config.Config) string {
	if cfg.DisableTutorNotes {
		return "off"
	}
	return ""
}

// sandboxRunArgs builds the `docker run` argument list for an ungraded
// sandbox session. Same rationale as exerciseRunArgs for being a pure
// function.
func sandboxRunArgs(cfg config.Config) []string {
	return []string{
		"run", "-it", "--rm",
		// Same session label as exerciseRunArgs -- ballroom voice works
		// in sandbox sessions too.
		"--label", sessionContainerLabel + "=1",
		"-v", sandboxVolume + ":/workspace",
		// Sandbox sessions set none of the exercise PRACTICE_* vars, so
		// they carry their own marker -- isSessionContext
		// (cmd/ballroom/main.go) needs it for `ballroom return` to work
		// inside a sandbox.
		"-e", "PRACTICE_SANDBOX=1",
		"-e", "TUTOR_MODEL=" + cfg.TutorModel,
		"-e", "OPENROUTER_API_KEY=" + cfg.OpenRouterAPIKey,
		"-e", "TUTOR_ORCHESTRATOR_MODEL=" + cfg.OrchestratorModel,
		// See exerciseRunArgs -- same host-env streaming override.
		"-e", "TUTOR_STREAM=" + os.Getenv("TUTOR_STREAM"),
		// See exerciseRunArgs -- same notes-toggle forwarding.
		"-e", "PRACTICE_TUTOR_NOTES=" + tutorNotesEnvValue(cfg),
		cfg.DockerImage,
	}
}

// SessionContainerLabel marks running practice containers so host-side
// tooling (ballroom voice) can discover them via
// `docker ps --filter label=...` -- containers are otherwise anonymous
// (--rm, no --name).
const sessionContainerLabel = "com.ballroom.session"

// SessionContainerFilter is the docker ps filter that finds running
// practice session containers.
const SessionContainerFilter = "label=" + sessionContainerLabel
