// Package orchestrator drives the exercise/sandbox lifecycle: mounting the
// exercise repo into the unified Docker image, starting/stopping the timer,
// and revealing the hidden test suite only on submit.
package orchestrator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
)

const revealPollInterval = 200 * time.Millisecond

const sandboxVolume = "practice-sandbox"

// RunExercise starts a graded, timed session and blocks until the
// container exits. ex.RepoPath (the permanent exercise source) is never
// mounted directly — PrepareWorkspace copies it into a disposable temp
// dir that gets mounted as /workspace and deleted on exit, so nothing
// written during the session (edits, or hidden tests revealed on submit
// — see WaitAndReveal) can leak back into the source repo.
func RunExercise(cfg config.Config, ex exercise.Exercise) error {
	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return fmt.Errorf("orchestrator: create data dir: %w", err)
	}

	workspaceDir, cleanupWorkspace, err := PrepareWorkspace(ex.RepoPath)
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

	args := []string{
		"run", "-it", "--rm",
		"-v", workspaceDir + ":/workspace",
		"-v", cfg.DataDir + ":/data",
		"-v", controlDir + ":/control",
		"-e", "PRACTICE_CONTROL_DIR=/control",
		"-e", "PRACTICE_WORKSPACE_DIR=/workspace",
		"-e", "PRACTICE_TEST_COMMAND=" + ex.TestCommand,
		"-e", "PRACTICE_EXERCISE_ID=" + ex.ID,
		"-e", "PRACTICE_CATEGORY=" + ex.Category,
		"-e", "PRACTICE_LANGUAGE=" + ex.Language,
		"-e", "PRACTICE_TUTOR_MODE=" + ex.TutorMode,
		"-e", "PRACTICE_STARTED_AT=" + startedAt.Format(time.RFC3339),
		"-e", "PRACTICE_DB_PATH=/data/tracker.db",
		cfg.DockerImage,
	}

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
// an ungraded, untimed session. See interview_prep_mvp_spec.md Section 3.6.
//
// Reset is manual by design (no scripted "reset to base" for MVP): wipe
// the volume with `docker volume rm practice-sandbox` — the next
// `practice run --sandbox` recreates it empty.
func RunSandbox(cfg config.Config) error {
	args := []string{
		"run", "-it", "--rm",
		"-v", sandboxVolume + ":/workspace",
		cfg.DockerImage,
	}

	cmd := exec.Command("docker", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("orchestrator: docker run --sandbox: %w", err)
	}
	return nil
}
