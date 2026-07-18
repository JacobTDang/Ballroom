package orchestrator

import (
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

// testSignal is what these tests fire at the process to exercise
// installSignalCleanup -- SIGWINCH rather than SIGINT/SIGTERM/SIGHUP
// (what RunExercise actually registers) because SIGWINCH's default,
// unhandled disposition is Ignore. That makes it safe to send for real,
// including through installSignalCleanup's own re-raise-with-default-
// disposition step, without any risk of it tearing down the test
// process itself if a step here doesn't behave as expected.
const testSignal = syscall.SIGWINCH

func TestInstallSignalCleanup_RunsCleanupOnSignal(t *testing.T) {
	done := make(chan struct{}, 1)
	stop := installSignalCleanup(func() { done <- struct{}{} }, testSignal)
	defer stop()

	if err := syscall.Kill(os.Getpid(), testSignal); err != nil {
		t.Fatalf("self-signal: %v", err)
	}

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("cleanup was not called within 2s of the signal")
	}
}

func TestInstallSignalCleanup_DeregisteredOnNormalReturn(t *testing.T) {
	var calls int32
	stop := installSignalCleanup(func() { atomic.AddInt32(&calls, 1) }, testSignal)

	stop() // simulate RunExercise's normal return

	if err := syscall.Kill(os.Getpid(), testSignal); err != nil {
		t.Fatalf("self-signal: %v", err)
	}
	// Give a (incorrectly) still-registered handler a chance to fire --
	// there's nothing to block on for an absence, so this is a bounded
	// wait rather than a signal of its own.
	time.Sleep(200 * time.Millisecond)

	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Errorf("cleanup ran %d time(s) after stop() -- the signal handler was not deregistered", got)
	}
}

func TestInstallSignalCleanup_StopIsSafeToCallTwice(t *testing.T) {
	stop := installSignalCleanup(func() {}, testSignal)
	stop()
	stop() // must not panic (e.g. a double close on an internal channel)
}
