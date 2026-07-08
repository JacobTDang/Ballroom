package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/session"
)

func main() {
	var err error
	if len(os.Args) < 2 {
		err = homeCmd()
	} else {
		switch os.Args[1] {
		case "home":
			err = homeCmd()
		case "run":
			err = runCmd(os.Args[2:])
		case "submit":
			err = submitCmd()
		default:
			usage()
			os.Exit(1)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "practice: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: practice [home] | practice run --exercise <id> | practice run --sandbox | practice submit")
}

// homeCmd shows the exercise catalog with practice status, prompts for a
// choice, launches it, and loops back until the user quits — the "home
// base" you return to between sessions rather than having to remember
// exercise ids.
func homeCmd() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		statuses, err := catalog.List(cfg)
		if err != nil {
			return err
		}
		sandboxChoice := len(statuses) + 1

		fmt.Println()
		fmt.Println("  Ballroom — Interview Prep")
		fmt.Println()
		fmt.Print(catalog.FormatTable(statuses))
		fmt.Printf("  %-3d %s\n", sandboxChoice, "sandbox — free practice, no grading")
		fmt.Println()
		fmt.Println("  " + catalog.FormatSummary(statuses))
		fmt.Println()
		fmt.Print("Type a number to practice, or 'q' to quit: ")

		if !scanner.Scan() {
			fmt.Println()
			return nil
		}
		choice := strings.TrimSpace(scanner.Text())
		if choice == "q" || choice == "quit" {
			return nil
		}

		n, convErr := strconv.Atoi(choice)
		if convErr != nil || n < 1 || n > sandboxChoice {
			fmt.Println("invalid choice")
			continue
		}

		var runErr error
		if n == sandboxChoice {
			runErr = orchestrator.RunSandbox(cfg)
		} else {
			runErr = orchestrator.RunExercise(cfg, statuses[n-1].Exercise)
		}
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "practice: %v\n", runErr)
		}
	}
}

func runCmd(args []string) error {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	exerciseID := fs.String("exercise", "", "exercise id to run (mutually exclusive with --sandbox)")
	sandbox := fs.Bool("sandbox", false, "run a persistent, untimed sandbox session")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if (*exerciseID == "") == !*sandbox {
		return fmt.Errorf("exactly one of --exercise <id> or --sandbox is required")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if *sandbox {
		return orchestrator.RunSandbox(cfg)
	}

	ex, err := exercise.Load(cfg.ExercisePath(*exerciseID))
	if err != nil {
		return err
	}
	return orchestrator.RunExercise(cfg, ex)
}

func submitCmd() error {
	startedAtRaw := os.Getenv("PRACTICE_STARTED_AT")
	if startedAtRaw == "" {
		return fmt.Errorf("submit: not running inside a graded exercise session (are you in --sandbox mode?)")
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
