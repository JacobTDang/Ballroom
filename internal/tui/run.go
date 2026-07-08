package tui

import (
	"fmt"
	"os"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
)

// Run shows the boot screen once, then loops the picker: each selection
// runs the chosen exercise or sandbox (blocking, full terminal control,
// same as the CLI's direct commands), and the picker reopens afterward
// with refreshed status until the user quits.
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
		statuses, err := catalog.List(cfg)
		if err != nil {
			return err
		}

		sel, ok, err := RunPicker(statuses)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		var runErr error
		if sel.Sandbox {
			runErr = orchestrator.RunSandbox(cfg)
		} else {
			runErr = orchestrator.RunExercise(cfg, sel.Exercise.Exercise)
		}
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "ballroom: %v\n", runErr)
		}
	}
}
