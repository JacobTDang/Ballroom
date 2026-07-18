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
	"sync"
	"syscall"
	"time"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/draft"
	"github.com/JacobTDang/Ballroom/internal/exercise"
)

const revealPollInterval = 200 * time.Millisecond

// draftSnapshotInterval is how often SnapshotLoop polls the workspace
// for solution.* changes to persist to data/.drafts/ -- cheap enough
// (a short-circuiting hash compare on a no-op tick) to run for the
// whole session without any noticeable cost.
const draftSnapshotInterval = 2 * time.Second

const sandboxVolume = "ballroom-sandbox"

// RunExercise starts a graded, timed session and blocks until the
// container exits. ex.RepoPath (the permanent exercise source) is never
// mounted directly — PrepareWorkspace copies it into a disposable temp
// dir that gets mounted as /workspace and deleted on exit, so nothing
// written during the session (edits, or hidden tests revealed on submit
// — see WaitAndReveal) can leak back into the source repo.
//
// The workspace is disposable but the user's in-progress code must not
// be: SnapshotLoop polls it into data/.drafts/<exercise-id>/ every
// draftSnapshotInterval while the session runs, and a final finalize
// call (see newSessionFinalizer) catches whatever was last saved on
// every exit path -- normal completion, submit, `ballroom return`, a
// killed container, a docker error (issue #221), or the host terminal
// closing (SIGINT/SIGTERM/SIGHUP -- issue #231). Go doesn't run
// deferred functions for an unhandled fatal signal, so without
// installSignalCleanup catching those three, finalize's deferred call
// below would simply never run and workspaceDir/controlDir would leak
// forever.
//
// draftDir is the caller's answer to "resume or start fresh": the
// draft directory to overlay onto the starter, or empty for a pristine
// start. The decision belongs to the caller (the TUI asks the user;
// see internal/tui/resumedraft.go) rather than to this function, so a
// saved draft is never silently resumed or silently discarded.
func RunExercise(cfg config.Config, ex exercise.Exercise, draftDir string) error {
	if err := EnsureImage(cfg); err != nil {
		return err
	}

	if err := os.MkdirAll(cfg.DataDir, 0o755); err != nil {
		return fmt.Errorf("orchestrator: create data dir: %w", err)
	}

	workspaceDir, cleanupWorkspace, err := PrepareWorkspace(ex.RepoPath, ex.VideoURL, draftDir)
	if err != nil {
		return err
	}

	controlDir, err := os.MkdirTemp("", controlDirPrefix)
	if err != nil {
		cleanupWorkspace()
		return fmt.Errorf("orchestrator: create control dir: %w", err)
	}

	startedAt := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// finalize is the session's one true cleanup -- the final draft
	// snapshot (catching whatever was last saved before the workspace it
	// lives in is deleted) plus removing the workspace and control dirs.
	// Idempotent, so it's safe for both the deferred call below (the
	// normal exit path) and the signal handler registered next to reach
	// it -- e.g. a SIGHUP arriving the same instant the container exits
	// on its own.
	finalize := newSessionFinalizer(cfg, ex.ID, workspaceDir, controlDir, cleanupWorkspace)
	defer finalize()

	// Closing the terminal SIGHUPs this process; Ctrl-C/a kill sends
	// SIGINT/SIGTERM. On signal, run the same finalize the normal exit
	// path does, then re-raise so the process still dies promptly (see
	// installSignalCleanup). stop deregisters this and stops its
	// goroutine on every return from here on, so a long session of many
	// RunExercise calls (the Run loop in internal/tui/run.go) never
	// accumulates one abandoned handler per past session.
	stopSignalCleanup := installSignalCleanup(finalize, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stopSignalCleanup()

	revealErr := make(chan error, 1)
	go func() {
		revealErr <- WaitAndReveal(ctx, controlDir, cfg.TestsPath(ex.ID), workspaceDir, revealPollInterval)
	}()

	// Independent of the submit watcher above (see WaitAndRevealReference's
	// own doc comment): the reference solution can be requested by `ballroom
	// reference` (M-g) or by a passing `ballroom submit`, in any order or
	// not at all, so it gets its own watcher rather than being folded into
	// revealErr's single-purpose loop.
	referenceErr := make(chan error, 1)
	go func() {
		referenceErr <- WaitAndRevealReference(ctx, controlDir, cfg.ReferencePath(ex.ID), cfg.TestsPath(ex.ID), workspaceDir, revealPollInterval)
	}()

	snapshotErr := make(chan error, 1)
	go func() {
		snapshotErr <- SnapshotLoop(ctx, cfg.DataDir, ex.ID, workspaceDir, draftSnapshotInterval)
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
	if err := <-referenceErr; err != nil && ctx.Err() == nil {
		fmt.Fprintf(os.Stderr, "orchestrator: reference watcher: %v\n", err)
	}
	if err := <-snapshotErr; err != nil {
		fmt.Fprintf(os.Stderr, "orchestrator: draft snapshot loop: %v\n", err)
	}

	if runErr != nil {
		return fmt.Errorf("orchestrator: docker run: %w", runErr)
	}
	return nil
}

// newSessionFinalizer returns the session's idempotent cleanup: one
// final draft.Snapshot, then removing the workspace and control dirs.
// Wrapped in sync.Once because two different paths can both reach the
// returned func -- RunExercise's own deferred call on a normal return,
// and installSignalCleanup's handler on SIGINT/SIGTERM/SIGHUP -- and
// only the first one to run should actually do the work.
func newSessionFinalizer(cfg config.Config, exerciseID, workspaceDir, controlDir string, cleanupWorkspace func()) func() {
	var once sync.Once
	return func() {
		once.Do(func() {
			if _, err := draft.Snapshot(cfg.DataDir, exerciseID, workspaceDir); err != nil {
				fmt.Fprintf(os.Stderr, "orchestrator: final draft snapshot: %v\n", err)
			}
			cleanupWorkspace()
			os.RemoveAll(controlDir)
		})
	}
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
