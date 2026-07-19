// Command capture-svg turns a terminal capture into an SVG, so the
// README's screenshots can be regenerated instead of being one-off
// artifacts nobody can reproduce (which is exactly what they were
// before this existed).
//
// Reads ANSI text on stdin, writes SVG on stdout:
//
//	tmux capture-pane -e -p -t <target> | go run ./cmd/capture-svg 132 > docs/img/home.svg
//
// The width argument is the terminal's column count — every glyph is
// placed on that grid and stretched to exactly one cell, so the result
// lines up the way the terminal did rather than however the viewer's
// font would lay it out.
package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: capture-svg <columns>  (ANSI on stdin, SVG on stdout)")
		os.Exit(2)
	}
	cols, err := strconv.Atoi(os.Args[1])
	if err != nil || cols <= 0 {
		fmt.Fprintf(os.Stderr, "capture-svg: columns must be a positive number, got %q\n", os.Args[1])
		os.Exit(2)
	}
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "capture-svg: read stdin: %v\n", err)
		os.Exit(1)
	}
	if _, err := io.WriteString(os.Stdout, render(parse(string(input), cols), cols)); err != nil {
		fmt.Fprintf(os.Stderr, "capture-svg: write: %v\n", err)
		os.Exit(1)
	}
}
