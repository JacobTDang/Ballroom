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
// opens a chain of popups — category, then problem, then language
// (itself a loop — you can work through several exercises before going
// back) — Sandbox launches directly, and Stats shows progress. Each
// returns to the menu when done, until the user quits from there.
//
// The practice image itself is ensured (built if missing, stale builds
// cleaned up) inside orchestrator.RunExercise/RunSandbox, not here — that
// way `ballroom practice <id>` and `ballroom sandbox` get the same
// behavior whether or not they go through this TUI.
func Run(cfg config.Config) error {
	cfg, proceed, err := RunBoot(cfg)
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
		case menuModelPicker:
			cfg, runErr = runModelPicker(cfg)
		}
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "ballroom: %v\n", runErr)
		}
	}
}

// runModelPicker shows the model popup and, if the user picks or types a
// model, persists it to settings.json and returns an updated cfg so every
// subsequent orchestrator.RunExercise/RunSandbox call in this same
// invocation immediately uses it too — not just future launches of the
// TUI. Backing out of the popup returns cfg unchanged.
func runModelPicker(cfg config.Config) (config.Config, error) {
	model, ok, err := RunModelPicker(ollamaHost, cfg.TutorModel)
	if err != nil {
		return cfg, err
	}
	if !ok {
		return cfg, nil
	}

	cfg.TutorModel = model
	if err := config.SaveSettings(cfg.SettingsPath(), config.Settings{TutorModel: model}); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// runPracticeLoop shows the category popup and, for whichever category is
// picked, drills into the problem popup until the user backs out to the
// main menu.
func runPracticeLoop(cfg config.Config) error {
	for {
		statuses, err := catalog.List(cfg)
		if err != nil {
			return err
		}
		problems := catalog.GroupByProblem(statuses)

		category, ok, err := RunCategoryPicker(problems)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		if err := runProblemLoop(cfg, problems, category); err != nil {
			return err
		}
	}
}

// runProblemLoop shows the problem popup for one category and, each time
// an exercise finishes, reopens it with refreshed status until the user
// backs out to the category popup. Selecting a problem doesn't launch it
// directly — a language popup asks which variant to practice first;
// backing out of that popup returns to the problem list rather than the
// category popup.
func runProblemLoop(cfg config.Config, problems []catalog.ProblemStatus, category string) error {
	for {
		problem, ok, err := RunProblemPicker(problems, category)
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

		statuses, err := catalog.List(cfg)
		if err != nil {
			return err
		}
		problems = catalog.GroupByProblem(statuses)
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
