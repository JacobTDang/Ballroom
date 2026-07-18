// Command tutor-preview renders the tutor pane on canned fixture
// content with no network, agents, or docker session — the visual
// harness for pane-restyling work. Run it in a tmux pane and
// capture-pane the result; tutor.RunPreview's doc comment lists the
// keys that reach each visual state.
//
//	go run ./cmd/tutor-preview
package main

import (
	"fmt"
	"os"

	"github.com/JacobTDang/Ballroom/internal/tutor"
)

func main() {
	if err := tutor.RunPreview(); err != nil {
		fmt.Fprintf(os.Stderr, "tutor-preview: %v\n", err)
		os.Exit(1)
	}
}
