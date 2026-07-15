package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// Run shows the boot screen once, then runs the merged menu/practice/
// stats/model-picker program (see app.go) in a loop. Sandbox and
// finishing a language variant both hand the terminal to `docker run -it`
// (orchestrator.RunSandbox/RunExercise) — bubbletea can't render inside
// that external interactive process, so each one fully tears the program
// down and, once docker returns, a fresh one launches picking back up
// where it makes sense: stageProblems for the category just practiced,
// stageMain for sandbox.
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

	resume := appResume{}
	for {
		final, err := RunApp(cfg, resume)
		if err != nil {
			return err
		}
		if final.quit {
			return nil
		}
		cfg = final.cfg

		switch final.outcome {
		case outcomeRunExercise:
			if runErr := orchestrator.RunExercise(cfg, final.exerciseToRun); runErr != nil {
				fmt.Fprintf(os.Stderr, "ballroom: %v\n", runErr)
			}
			resume = appResume{stage: stageProblems, category: final.category}
		case outcomeRunSandbox:
			if runErr := orchestrator.RunSandbox(cfg); runErr != nil {
				fmt.Fprintf(os.Stderr, "ballroom: %v\n", runErr)
			}
			resume = appResume{}
		default:
			resume = appResume{}
		}
	}
}

// RunApp shows the merged menu/practice/stats/model-picker program and
// blocks until it exits — either the user quit, or an outcome needs a
// docker handoff outside bubbletea (see appOutcome).
func RunApp(cfg config.Config, resume appResume) (appModel, error) {
	final, err := tea.NewProgram(newAppModel(cfg, resume), tea.WithAltScreen()).Run()
	if err != nil {
		return appModel{}, err
	}
	return final.(appModel), nil
}

// recentAttempts returns up to n attempts, newest first.
// allAttempts returns the full attempt log -- Stats aggregates rubric
// weak spots across every graded design attempt, not just the recent
// window recentAttempts trims to.
func allAttempts(cfg config.Config) ([]tracker.Attempt, error) {
	tr, err := tracker.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	defer tr.Close()
	return tr.ListAttempts()
}

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
