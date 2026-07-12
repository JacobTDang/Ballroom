// Command kitty-image-check is a standalone diagnostic for the disco
// ball's Kitty graphics protocol rendering path
// (internal/tutor/kittyimage.go) — loops through all 16 real sprite
// frames in place, so the protocol (and, run inside tmux, the tmux
// allow-passthrough chain — see docker/tmux.conf) can be checked in
// isolation from a full practice session.
//
// This path could not be verified in the environment it was built in
// (no Kitty install, no display) — everything in kittyimage.go is a
// best-effort implementation of Kitty's documented protocol. Run this
// first, outside Docker/tmux entirely, in a real Kitty window:
//
//	go run ./cmd/kitty-image-check
//
// Runs until interrupted — Ctrl-C to stop (cleans up its uploaded
// images either way). If the ball animates smoothly, the protocol
// itself works. Only then is it worth testing through the full chain
// (inside a practice session's tutor pane, where TMUX is set and
// passthrough framing kicks in automatically).
package main

import (
	"fmt"
	"os"

	"github.com/JacobTDang/Ballroom/internal/tutor"
)

func main() {
	fmt.Println("looping — Ctrl-C to stop")
	if err := tutor.KittyLiveCheck(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "kitty-image-check:", err)
		os.Exit(1)
	}
}
