package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/session"
	"github.com/JacobTDang/Ballroom/internal/tui"
	"github.com/JacobTDang/Ballroom/internal/tutor"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		exitOnErr(homeCmd())
		return
	}

	switch args[0] {
	case "-h", "--help", "help":
		printUsage(os.Stdout)
	case "home":
		exitOnErr(homeCmd())
	case "practice":
		exitOnErr(practiceCmd(args[1:]))
	case "sandbox":
		exitOnErr(sandboxCmd())
	case "submit":
		exitOnErr(submitCmd())
	case "return":
		exitOnErr(returnCmd())
	case "tutor":
		exitOnErr(tutorCmd())
	default:
		fmt.Fprintf(os.Stderr, "ballroom: unknown command %q\n\n", args[0])
		printUsage(os.Stderr)
		os.Exit(1)
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ballroom: %v\n", err)
		os.Exit(1)
	}
}

func printUsage(w *os.File) {
	fmt.Fprint(w, `Ballroom — interview practice CLI

Usage:
  ballroom                    Open the homepage (pick an exercise interactively)
  ballroom home                Same as above
  ballroom practice <id>       Practice a specific exercise by id
  ballroom sandbox             Free practice, no grading, persists across sessions
  ballroom submit              Submit your solution (run this inside an active session)
  ballroom tutor               Start the tutor chat (run this inside an active session)
  ballroom return              Return to the host homepage (run this inside an active session)
  ballroom help | -h | --help  Show this help

Examples:
  ballroom
  ballroom practice two-pointers-01-go
  ballroom sandbox

Reset the sandbox volume:
  docker volume rm ballroom-sandbox
`)
}

// homeCmd shows the full-screen boot check + exercise picker (see
// internal/tui) — the "home base" you return to between sessions rather
// than having to remember exercise ids. The ballroom binary baked into
// the practice image (docker/Dockerfile) means this same code path can
// run either on the host or inside an active session's container; the
// boot screen's preflight checks (CheckDocker, CheckOllama against
// localhost:11434, ...) assume host networking and there's no Docker
// client inside the container, so booting a nested instance there
// doesn't fail cleanly. If we're inside a session, return to the host
// homepage instead of attempting that nested boot.
func homeCmd() error {
	if isSessionContext() {
		return returnToHost()
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return tui.Run(cfg)
}

// isSessionContext reports whether this process is running inside an
// active practice session's container, as opposed to on the host. It
// checks all three session-scoped env vars orchestrator.RunExercise sets
// via `docker run -e` (PRACTICE_WORKSPACE_DIR, PRACTICE_CONTROL_DIR,
// PRACTICE_STARTED_AT) rather than any single one, so a host shell that
// happens to have one of these set for unrelated reasons isn't
// misdetected as a session.
func isSessionContext() bool {
	return os.Getenv("PRACTICE_WORKSPACE_DIR") != "" &&
		os.Getenv("PRACTICE_CONTROL_DIR") != "" &&
		os.Getenv("PRACTICE_STARTED_AT") != ""
}

// returnCmd is `ballroom return`, run from a session's TERMINAL pane to
// get back to the host homepage. Unlike homeCmd, it's only meaningful
// inside a session, so it errors instead of silently falling back to
// tui.Run when there's nothing to return from.
func returnCmd() error {
	if !isSessionContext() {
		return fmt.Errorf("return: not running inside an active practice session (did you mean `ballroom home`?)")
	}
	return returnToHost()
}

// returnToHost ends the practice session so control lands back on the
// host homepage. The container can't reach out and control the host's
// `docker run -it --rm` (orchestrator.RunExercise) directly — no Docker
// client is installed inside it — but that `docker run` is blocking on
// the container's PID 1, which docker/entrypoint.sh sets to `tmux
// attach` after starting the session's tmux server. Killing that server
// tears down every window, which ends the attach client, which exits the
// container, which is what makes `docker run -it --rm` on the host
// return. RunExercise returning is what lets practiceCmd continue on to
// homeCmd and open the homepage picker.
func returnToHost() error {
	if err := exec.Command("tmux", "kill-server").Run(); err != nil {
		return fmt.Errorf("return: tmux kill-server: %w", err)
	}
	return nil
}

func practiceCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ballroom practice <exercise-id>")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := runExercise(cfg, args[0]); err != nil {
		return err
	}
	// The session container just exited (see returnToHost); land back on
	// the host homepage rather than dropping to a bare shell prompt.
	return homeCmd()
}

func runExercise(cfg config.Config, id string) error {
	ex, err := exercise.Load(cfg.ExercisePath(id))
	if err != nil {
		return fmt.Errorf("unknown exercise %q — run `ballroom help` for usage: %w", id, err)
	}
	return orchestrator.RunExercise(cfg, ex)
}

func sandboxCmd() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return orchestrator.RunSandbox(cfg)
}

func submitCmd() error {
	startedAtRaw := os.Getenv("PRACTICE_STARTED_AT")
	if startedAtRaw == "" {
		return fmt.Errorf("submit: not running inside a graded exercise session (did you mean to run `ballroom sandbox`?)")
	}
	startedAt, err := time.Parse(time.RFC3339, startedAtRaw)
	if err != nil {
		return fmt.Errorf("submit: parse PRACTICE_STARTED_AT: %w", err)
	}

	cfg := session.Config{
		ControlDir:    os.Getenv("PRACTICE_CONTROL_DIR"),
		WorkspaceDir:  os.Getenv("PRACTICE_WORKSPACE_DIR"),
		TestCommand:   os.Getenv("PRACTICE_TEST_COMMAND"),
		ExerciseID:    os.Getenv("PRACTICE_EXERCISE_ID"),
		Category:      os.Getenv("PRACTICE_CATEGORY"),
		Language:      os.Getenv("PRACTICE_LANGUAGE"),
		StartedAt:     startedAt,
		DBPath:        os.Getenv("PRACTICE_DB_PATH"),
		PollInterval:  200 * time.Millisecond,
		RevealTimeout: 30 * time.Second,
	}

	attempt, err := session.Submit(cfg, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	fmt.Printf("logged attempt #%d\n", attempt.ID)
	return nil
}

// tutorCmd is `ballroom tutor`, launched in the tutor pane by
// docker/entrypoint.sh (env vars below match what it sets — see
// NVIM_SOCKET/OLLAMA_HOST/TUTOR_MODEL/PRACTICE_TUTOR_MODE there, plus
// WORKDIR which every pane shares). Defaults mirror tutor/chat.sh's own
// fallbacks so a standalone run (e.g. local testing outside a real
// session) behaves the same way.
func tutorCmd() error {
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://host.docker.internal:11434"
	}
	model := os.Getenv("TUTOR_MODEL")
	if model == "" {
		model = config.DefaultTutorModel
	}
	mode := os.Getenv("PRACTICE_TUTOR_MODE")
	if mode == "" {
		mode = exercise.TutorModeFullAssist
	}
	workDir := os.Getenv("WORKDIR")
	if workDir == "" {
		workDir = "/workspace"
	}
	maxContextBytes := 8000
	if raw := os.Getenv("TUTOR_FILE_CONTEXT_MAX_BYTES"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			maxContextBytes = n
		}
	}

	cfg := tutor.Config{
		OllamaHost:      ollamaHost,
		Model:           model,
		Mode:            mode,
		WorkDir:         workDir,
		NvimSocket:      os.Getenv("NVIM_SOCKET"),
		MaxContextBytes: maxContextBytes,
	}
	return tutor.Run(context.Background(), cfg, os.Stdin, os.Stdout, os.Stderr)
}
