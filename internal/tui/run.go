package tui

import (
	"fmt"
	"os"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

const recentAttemptsLimit = 10

// Run shows the boot screen once, then loops the main menu: Practice
// opens the pattern tree (itself a loop — you can work through several
// exercises before going back), Sandbox launches directly, and Stats
// shows progress. Each returns to the menu when done, until the user
// quits from there.
//
// The practice image itself is ensured (built if missing, stale builds
// cleaned up) inside orchestrator.RunExercise/RunSandbox, not here — that
// way `ballroom practice <id>` and `ballroom sandbox` get the same
// behavior whether or not they go through this TUI.
func Run(cfg config.Config) error {
	proceed, err := RunBoot(cfg)
	if err != nil {
		return err
	}
	if !proceed {
		return nil
	}

	for {
		choice, ok, err := RunMenu()
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		var runErr error
		switch choice {
		case menuPractice:
			runErr = runPracticeLoop(cfg)
		case menuSandbox:
			runErr = orchestrator.RunSandbox(cfg)
		case menuStats:
			runErr = runStats(cfg)
		}
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "ballroom: %v\n", runErr)
		}
	}
}

// runPracticeLoop shows the pattern tree and, each time an exercise
// finishes, reopens it with refreshed status until the user backs out to
// the main menu. Selecting a problem doesn't launch it directly — a
// language popup asks which variant to practice first; backing out of
// that popup returns to the tree rather than the main menu.
func runPracticeLoop(cfg config.Config) error {
	for {
		statuses, err := catalog.List(cfg)
		if err != nil {
			return err
		}

		problem, ok, err := RunTree(statuses)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		variant, ok, err := RunLangPicker(problem)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}

		if runErr := orchestrator.RunExercise(cfg, variant.Exercise); runErr != nil {
			fmt.Fprintf(os.Stderr, "ballroom: %v\n", runErr)
		}
	}
}

func runStats(cfg config.Config) error {
	statuses, err := catalog.List(cfg)
	if err != nil {
		return err
	}
	recent, err := recentAttempts(cfg, recentAttemptsLimit)
	if err != nil {
		return err
	}
	return RunStats(statuses, recent)
}

// recentAttempts returns up to n attempts, newest first.
func recentAttempts(cfg config.Config, n int) ([]tracker.Attempt, error) {
	tr, err := tracker.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	defer tr.Close()

	all, err := tr.ListAttempts()
	if err != nil {
		return nil, err
	}

	if len(all) > n {
		all = all[len(all)-n:]
	}
	recent := make([]tracker.Attempt, len(all))
	for i, a := range all {
		recent[len(all)-1-i] = a
	}
	return recent, nil
}
