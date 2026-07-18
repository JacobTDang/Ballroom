package orchestrator

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// installSignalCleanup runs cleanup the first time one of sigs arrives
// at this process, then re-raises the same signal with its default
// disposition so the process still terminates promptly -- signal.Notify
// below is what stops the Go runtime from doing that on its own the
// instant the signal is delivered.
//
// Returns a stop func that must be called on every return path of the
// caller (a plain `defer stop()` right after this call): it deregisters
// the signal channel and stops the goroutine below, so a long-running
// process that calls this once per session (RunExercise, from the Run
// loop in internal/tui/run.go) never accumulates one abandoned
// goroutine and registration per past session. stop is safe to call
// more than once, matching context.CancelFunc's own contract.
func installSignalCleanup(cleanup func(), sigs ...os.Signal) (stop func()) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sigs...)
	done := make(chan struct{})

	go func() {
		select {
		case sig := <-ch:
			cleanup()
			signal.Stop(ch)
			// Restore the default disposition for this signal and send
			// it to ourselves again, so the process actually dies (and
			// with the right signal-terminated exit status) instead of
			// silently surviving a signal we just finished "handling".
			signal.Reset(sig)
			if s, ok := sig.(syscall.Signal); ok {
				if err := syscall.Kill(os.Getpid(), s); err != nil {
					fmt.Fprintf(os.Stderr, "orchestrator: re-raise %v: %v\n", sig, err)
				}
			}
		case <-done:
		}
	}()

	var stopOnce sync.Once
	return func() {
		stopOnce.Do(func() {
			signal.Stop(ch)
			close(done)
		})
	}
}
