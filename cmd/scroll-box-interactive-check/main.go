// Command scroll-box-interactive-check exercises the exact per-turn
// cycle tutor.go's now-reverted anchored-input-box integration used:
// show a prompt inside a bordered box at the bottom of the pane, read a
// real line from stdin, jump the cursor back into the scroll region
// above the box, echo the line there. This is the gap in
// cmd/scroll-box-check (which never reads real input at all) that let
// the box integration ship without ever being tested against the one
// thing that broke it live in a real practice session: typing and
// submitting a message.
//
// Run this inside the same kind of split-pane tmux layout the real
// practice session uses (docker/tmux.conf), not just a bare
// single-pane tmux window — the failure showed up specifically in that
// multi-pane layout, and cmd/scroll-box-check's own local-tmux pass was
// only ever tested single-pane. For example:
//
//	tmux -f docker/tmux.conf new-session \; split-window -h \; select-pane -t 0
//
// then, inside one of the panes:
//
//	go run ./cmd/scroll-box-interactive-check
//
// Type a few lines and press enter after each. The box (borders +
// prompt) should stay fixed at the bottom the whole time; each line you
// type plus its echo should scroll normally above it. If the box
// disappears, garbles, or the echoed text overlaps instead of stacking
// cleanly, that confirms the failure reproduces outside Docker, and
// exactly what you did right before it broke (typed a line? switched
// panes? waited a few seconds?) is the next clue.
package main

import (
	"fmt"
	"os"

	"github.com/JacobTDang/Ballroom/internal/tutor"
)

func main() {
	if err := tutor.ScrollBoxInteractiveLiveCheck(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "scroll-box-interactive-check:", err)
		os.Exit(1)
	}
	fmt.Println("done — did the box stay anchored and undamaged the whole time?")
}
