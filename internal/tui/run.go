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
			// A launch failure (Docker down, image build broken, ...)
			// carries into the resumed screen via launchErr instead of
			// being printed here -- the very next line replaces this
			// program with a fresh alt-screen one, which would wipe a
			// stderr write before anyone could read it (issue #230).
			runErr := orchestrator.RunExercise(cfg, final.exerciseToRun, final.draftDirToUse)
			resume = appResume{stage: stageProblems, category: final.category, launchErr: runErr}
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
// docker handoff outside bubbletea (see appOutcome). When resume carries
// a launchErr, it's seeded into the fresh model's err field (see
// resumedAppModel) so the resumed screen renders it.
func RunApp(cfg config.Config, resume appResume) (appModel, error) {
	final, err := tea.NewProgram(resumedAppModel(cfg, resume), tea.WithAltScreen()).Run()
	if err != nil {
		return appModel{}, err
	}
	return final.(appModel), nil
}

// resumedAppModel is newAppModel plus RunApp's own launchErr seeding,
// split out so the seeding is unit-testable without going through a
// real bubbletea program. A launchErr never overwrites an error
// newAppModel's own construction already produced (e.g. the catalog
// failing to load) -- that one is more fundamental (it leaves the model
// on stageMain, unable to even reach the problem it resumed to) and
// takes precedence.
func resumedAppModel(cfg config.Config, resume appResume) appModel {
	m := newAppModel(cfg, resume)
	if m.err == nil {
		m.err = resume.launchErr
	}
	return m
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

// attemptsFor returns every attempt logged against exerciseID, newest
// first -- the Stats drill-down's data source (issue #252, see
// statsdetail.go).
func attemptsFor(cfg config.Config, exerciseID string) ([]tracker.Attempt, error) {
	tr, err := tracker.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	defer tr.Close()
	return tr.ListAttemptsFor(exerciseID)
}
