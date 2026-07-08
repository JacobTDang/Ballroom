package main

import (
	"fmt"
	"os"
	"time"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/session"
	"github.com/JacobTDang/Ballroom/internal/tui"
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
// than having to remember exercise ids.
func homeCmd() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return tui.Run(cfg)
}

func practiceCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ballroom practice <exercise-id>")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return runExercise(cfg, args[0])
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
