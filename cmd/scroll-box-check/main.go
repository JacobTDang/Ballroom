// Command scroll-box-check is a standalone diagnostic for the anchored-
// input-box technique (internal/tutor/scrollbox.go) — draws a static
// bordered box in the bottom 3 rows of the terminal via a DECSTBM
// scroll region, then prints filler "conversation" lines that scroll
// normally in the region above it. The box must never move or get
// overwritten by that scrolling.
//
// Already tried live and reverted: wired directly into tutor.go's real
// Run() loop once (skipping this check), and broke in the actual
// practice session's tutor pane — garbled, overlapping text, and the
// box itself never rendered at all — despite passing every isolated
// test, including a real pty capture outside tmux. Almost certainly a
// tmux-specific incompatibility (tmux's own pane redraw model doesn't
// know to preserve an application's DECSTBM boundary when it repaints
// a pane), not something wrong with the DECSTBM sequences themselves.
// tutor.go no longer uses this technique.
//
// If revisiting this: run outside Docker/tmux entirely first, in a
// real Kitty window:
//
//	go run ./cmd/scroll-box-check
//
// If the box stays put while filler lines scroll above it, that only
// confirms the technique works standalone — it did last time too, and
// still failed once actually running inside tmux. Test inside local
// tmux (tmux -f docker/tmux.conf new-session) next, and don't wire it
// into tutor.go again until that specifically is confirmed working.
package main

import (
	"fmt"
	"os"

	"github.com/JacobTDang/Ballroom/internal/tutor"
)

func main() {
	if err := tutor.ScrollBoxLiveCheck(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "scroll-box-check:", err)
		os.Exit(1)
	}
	fmt.Println("done — did the box stay anchored at the bottom the whole time?")
}
