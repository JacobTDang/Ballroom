package tutor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// scrollBoxHeight is how many terminal rows the anchored input box
// reserves at the bottom of the pane: top border, one content row, and
// bottom border.
const scrollBoxHeight = 3

// terminalSize shells out to `stty size`, the same technique
// docker/entrypoint.sh already uses to query the real terminal size
// (its own comment explains why: the pty already reflects the real host
// terminal size before `tmux attach` even happens) — reused here rather
// than adding a new dependency (e.g. golang.org/x/term) for the same
// question.
func terminalSize() (rows, cols int, err error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("tutor: stty size: %w", err)
	}
	return parseSttySize(string(out))
}

// parseSttySize parses `stty size`'s "ROWS COLS" output — split out as
// a pure function so it's testable without a real terminal.
func parseSttySize(output string) (rows, cols int, err error) {
	parts := strings.Fields(output)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("tutor: unexpected stty size output %q", output)
	}
	rows, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("tutor: stty size: parse rows: %w", err)
	}
	cols, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("tutor: stty size: parse cols: %w", err)
	}
	return rows, cols, nil
}

// boxTopLine/boxMiddleLine/boxBottomLine build the three rows of a
// cols-wide bordered box — pure string construction, testable without a
// terminal. Each is exactly cols runes: two corner/border characters
// plus (cols-2) of fill.
func boxTopLine(cols int) string {
	return "╭" + strings.Repeat("─", cols-2) + "╮"
}

func boxMiddleLine(cols int) string {
	return "│" + strings.Repeat(" ", cols-2) + "│"
}

func boxBottomLine(cols int) string {
	return "╰" + strings.Repeat("─", cols-2) + "╯"
}

// inputBox manages the anchored-input-box technique: a DECSTBM scroll
// region confines normal scrolling to the rows above a reserved bottom
// area, with a bordered box drawn in that reserved area. Rows outside
// the region are untouched by scrolling within it, which is what lets
// the box "stay" at the bottom of the pane while conversation text
// scrolls normally above it — the same mechanism ScrollBoxLiveCheck
// verifies standalone and tutor.go's real turn loop uses live.
type inputBox struct {
	w io.Writer
	// regionBottom is the last row of the confined scroll region — the
	// box occupies the scrollBoxHeight rows directly below it
	// (regionBottom+1: top border, regionBottom+2: content/prompt row,
	// regionBottom+3: bottom border, i.e. the terminal's last row).
	regionBottom int
	cols         int
}

// newInputBoxAt is the testable core of newInputBox, taking an
// already-known rows/cols instead of querying a real terminal.
func newInputBoxAt(w io.Writer, rows, cols int) (*inputBox, error) {
	regionBottom := rows - scrollBoxHeight
	if regionBottom < 1 {
		return nil, fmt.Errorf("tutor: terminal too short (%d rows) for a %d-row input box", rows, scrollBoxHeight)
	}
	b := &inputBox{w: w, regionBottom: regionBottom, cols: cols}
	b.setup()
	return b, nil
}

// newInputBox queries the real terminal size (terminalSize) and sets up
// an anchored input box for it. Returns an error whenever that fails —
// in particular, whenever stdin isn't a real terminal (tests,
// cmd/tutor-eval, or any environment stty size can't read), which
// callers must treat as "fall back to the plain prompt," not fatal: a
// tutor session needs to keep working exactly as it did before this
// feature existed in that case.
func newInputBox(w io.Writer) (*inputBox, error) {
	rows, cols, err := terminalSize()
	if err != nil {
		return nil, err
	}
	return newInputBoxAt(w, rows, cols)
}

// setup clears the screen, confines the scroll region, and draws the
// box's borders once, leaving the cursor at the top of the now-confined
// scroll region (Home) so whatever the caller prints next — tutor.go's
// startup banner, ScrollBoxLiveCheck's filler lines — lands there rather
// than wherever drawBorders' own absolute positioning last left the
// cursor (inside the box itself). Errors from these writes aren't
// checked individually — "push bytes to the terminal" is treated as
// fire-and-forget rather than something a caller could meaningfully
// recover from mid-render.
//
// The initial \033[2J (clear entire screen) is load-bearing, found via
// a real tmux repro after this broke live: a pane always has *something*
// on screen already (at minimum the shell prompt that was there right
// before this program started) before setup ever runs, and \033[H only
// moves the cursor — it doesn't erase anything. Printing new content
// starting at Home overwrites the old content only up to the new
// content's own length; anything the old content had beyond that (e.g.
// a longer shell prompt line) is left dangling, unerased, and reads as
// garbled/overlapping text mixed with the new output.
func (b *inputBox) setup() {
	io.WriteString(b.w, "\033[2J")
	fmt.Fprintf(b.w, "\033[1;%dr", b.regionBottom)
	b.drawBorders()
	io.WriteString(b.w, "\033[H")
}

func (b *inputBox) drawBorders() {
	fmt.Fprintf(b.w, "\033[%d;1H\033[2K%s", b.regionBottom+1, boxTopLine(b.cols))
	fmt.Fprintf(b.w, "\033[%d;1H\033[2K%s", b.regionBottom+2, boxMiddleLine(b.cols))
	fmt.Fprintf(b.w, "\033[%d;1H\033[2K%s", b.regionBottom+3, boxBottomLine(b.cols))
}

// showPrompt positions the cursor at the box's content row and prints
// the "> " prompt, ready for the terminal's own cooked-mode line
// editing to echo whatever gets typed next inside the box.
func (b *inputBox) showPrompt() {
	fmt.Fprintf(b.w, "\033[%d;1H\033[2K> ", b.regionBottom+2)
}

// returnToScroll repositions the cursor to the bottom row of the
// confined scroll region and clears it, ready for the caller to print.
// Call this after a line is submitted and before printing anything
// else: cooked-mode Enter handling leaves the cursor sitting on the
// box's own bottom-border row (outside the region), which is not where
// conversation text should print.
//
// The clear is load-bearing, found via a real live session after this
// broke: regionBottom is an absolute jump, not something reached by a
// genuine scroll, so it carries no guarantee of being blank — a scroll
// region blanks a row only when content actually scrolls past the
// bottom margin, and returnToScroll can land here well before (or after
// a short reply that never reached this far down) that's happened.
// Landing on a row that still holds real, longer leftover content from
// several turns back and printing a shorter new line over it just
// leaves the old tail dangling — the exact bug showPrompt already
// avoided by clearing before it prints, which this now matches.
//
// Also clears the box's own content row (regionBottom+2). Nothing else
// touches that row between a line being submitted and the next
// showPrompt() call at the start of the following turn, so without
// this it keeps showing cooked-mode's echo of the just-submitted line
// — a stale copy sitting directly below the real echo this function's
// caller is about to print into the scroll region. Found live: it
// self-corrects on the next prompt, but in between it reads as a
// duplicate, confusing enough to look like a bug.
func (b *inputBox) returnToScroll() {
	fmt.Fprintf(b.w, "\033[%d;1H\033[2K", b.regionBottom+2)
	fmt.Fprintf(b.w, "\033[%d;1H\033[2K", b.regionBottom)
}

// reconfigureAt is the testable core of reconfigure, taking an
// already-known rows/cols instead of querying a real terminal. A no-op
// when the computed region is unchanged (e.g. a spurious/duplicate
// SIGWINCH) or when the new size is too short for the box — in the
// latter case there's nothing sensible to draw, so this leaves the box
// as it was rather than corrupt it further.
//
// Any row-count change (oldRegionBottom != regionBottom) clears the
// whole screen before resetting the region and redrawing. Two real bugs
// found live drove this:
//
//  1. The startup race: docker/entrypoint.sh starts the tutor process
//     via tmux send-keys before tmux attach ever runs, and a detached
//     tmux session's pane size isn't necessarily final until a client
//     attaches (entrypoint.sh's own comment). A real repro showed the
//     pane reporting a transitional size for ~150ms before a genuine
//     SIGWINCH corrected it, so newInputBox's initial setup can draw
//     the box for the wrong size.
//  2. A genuine later, user-initiated resize mid-conversation: the
//     first fix here tried to surgically clear only the box's *old* 3
//     rows (computed from the old regionBottom) before drawing the new
//     one elsewhere, to avoid wiping real conversation history. A live
//     repro showed this doesn't work: tmux visibly reflows/shifts
//     already-printed rows when a pane's row *count* changes (confirmed
//     by capturing the pane and comparing exact row numbers before and
//     after — content that was written at absolute row 25 showed up at
//     row 30 after growing the pane), so clearing rows computed from
//     the *old* size can miss entirely, landing on the wrong physical
//     rows post-reflow and leaving the old box's stale content row
//     (almost always still holding a real "> <message>" at the moment
//     this runs — see showPrompt/Run's doc comments) looking like a
//     duplicate of the real echo.
//
// A pure column-width change (rows unchanged) does NOT clear the whole
// screen: the box's rows don't move, and drawBorders already clears
// each one (\033[2K) before repainting at the new width, which a real
// repro confirmed is reflow-safe (no row-shift risk when the row count
// itself doesn't change). This is a real trade-off, not a free win: any
// resize that changes row count loses currently-visible conversation
// history (still in the terminal's own native scrollback, just not on
// screen) in exchange for guaranteed-correct rendering, rather than
// attempting a "clever" partial redraw that's been shown twice now not
// to survive tmux's actual reflow behavior.
func (b *inputBox) reconfigureAt(rows, cols int) {
	regionBottom := rows - scrollBoxHeight
	if regionBottom < 1 {
		return
	}
	if regionBottom == b.regionBottom && cols == b.cols {
		return
	}
	rowsChanged := regionBottom != b.regionBottom
	b.regionBottom = regionBottom
	b.cols = cols
	if rowsChanged {
		io.WriteString(b.w, "\033[2J")
	}
	fmt.Fprintf(b.w, "\033[1;%dr", b.regionBottom)
	b.drawBorders()
}

// reconfigure re-queries the real terminal size and, if it changed,
// resets the scroll region and redraws the box at the new bottom — a
// resized terminal otherwise leaves the box operating on stale
// dimensions captured once at startup, reproducing the same class of
// misplaced-box/garbled-content bugs a stale size causes elsewhere.
// Unlike setup, this deliberately does not clear the screen: already-
// printed conversation content should stay visible across a resize.
// Some cosmetic overlap right at the moment of a shrink (the new,
// smaller region's bottom rows landing on old content that used to be
// safely above the old, larger box) is possible and self-corrects as
// new content prints — an accepted tradeoff against a full
// history-redraw model this package doesn't keep (see watchResize).
// Errors from terminalSize are swallowed the same way: nothing sensible
// to do without a size, so this leaves the box as it was.
func (b *inputBox) reconfigure() {
	rows, cols, err := terminalSize()
	if err != nil {
		return
	}
	b.reconfigureAt(rows, cols)
}

// watchResize subscribes to SIGWINCH (the terminal-resize signal),
// returning the channel it arrives on and a stop function to call (via
// defer) when the session ends. Deliberately just a signal subscription,
// not something that calls box.reconfigure() itself from a background
// goroutine: inputBox has no locking of its own, and tutor.go's Run does
// its own unsynchronized writes to the same underlying writer, so a
// concurrent redraw here could interleave with those and corrupt output
// exactly the way this package's other bugs already have. The caller is
// expected to drain this channel at a safe point in its own
// single-threaded loop (see Run) and call box.reconfigure() there
// instead.
func watchResize() (pendingResize chan os.Signal, stop func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)
	return sigCh, func() { signal.Stop(sigCh) }
}

// close resets the scroll region to the full screen. Always call this
// on the way out (via defer) so an ended or interrupted session doesn't
// leave the outer shell's own scrolling permanently confined.
func (b *inputBox) close() {
	io.WriteString(b.w, "\033[r")
}

const (
	scrollBoxCheckFillerLines = 60
	scrollBoxCheckLineDelay   = 150 * time.Millisecond
)

// ScrollBoxLiveCheck is a standalone diagnostic for the inputBox
// technique before it's wired into the real tutor pane — draws a
// static bordered box in the bottom scrollBoxHeight rows, then prints
// filler "conversation" lines that scroll normally in the region above
// it. The box must never move or get overwritten by that scrolling; if
// it does, something about this technique doesn't survive this
// project's specific nested terminal chain (tmux inside Docker inside
// Kitty) and needs to be looked at before any of this reaches tutor.go.
// See cmd/scroll-box-check.
//
// Like the Kitty graphics protocol work, this could not be verified in
// the environment it was built in — no live terminal here to confirm
// DECSTBM actually behaves as documented through this project's
// specific chain. Restores the full-screen scroll region on exit
// (including on Ctrl-C) so an interrupted run doesn't leave the outer
// shell's scrolling permanently confined.
func ScrollBoxLiveCheck(w io.Writer) error {
	box, err := newInputBox(w)
	if err != nil {
		return err
	}
	defer box.close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	defer signal.Stop(stop)

	for i := 0; i < scrollBoxCheckFillerLines; i++ {
		select {
		case <-stop:
			return nil
		default:
		}
		if _, err := fmt.Fprintf(w, "filler conversation line %d\n", i); err != nil {
			return fmt.Errorf("tutor: scroll box live check: %w", err)
		}
		time.Sleep(scrollBoxCheckLineDelay)
	}
	box.showPrompt()
	return nil
}

const scrollBoxInteractiveGreeting = "type a line and press enter -- it should scroll into view above the box, which should never move or break. Ctrl-D to exit."

// runScrollBoxInteractive drives the interactive loop against an
// already-constructed box — split out from ScrollBoxInteractiveLiveCheck
// so the loop logic is testable with a scripted io.Reader instead of a
// real terminal.
func runScrollBoxInteractive(r io.Reader, w io.Writer, box *inputBox) {
	fmt.Fprintln(w, scrollBoxInteractiveGreeting)
	scanner := bufio.NewScanner(r)
	for {
		box.showPrompt()
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		box.returnToScroll()
		fmt.Fprintf(w, "> %s\n(you typed: %s)\n", line, line)
	}
}

// ScrollBoxInteractiveLiveCheck exercises the exact per-turn cycle
// tutor.go's now-reverted box integration used: show the prompt inside
// the box, read a real line from stdin, jump the cursor back into the
// scroll region, echo it there. ScrollBoxLiveCheck (above) never reads
// real input at all — a gap that let the box integration ship without
// ever being tested against the one thing that broke it live: typing
// and submitting a message in a real tmux pane. See
// cmd/scroll-box-interactive-check; run this inside the same kind of
// split-pane tmux layout the real practice session uses
// (docker/tmux.conf), not just a bare single-pane tmux window.
func ScrollBoxInteractiveLiveCheck(r io.Reader, w io.Writer) error {
	box, err := newInputBox(w)
	if err != nil {
		return err
	}
	defer box.close()
	runScrollBoxInteractive(r, w, box)
	return nil
}
