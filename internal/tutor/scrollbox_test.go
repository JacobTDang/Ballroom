package tutor

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestParseSttySize_ParsesRowsAndCols(t *testing.T) {
	rows, cols, err := parseSttySize("24 80\n")
	if err != nil {
		t.Fatalf("parseSttySize: %v", err)
	}
	if rows != 24 || cols != 80 {
		t.Errorf("rows=%d cols=%d, want 24 80", rows, cols)
	}
}

func TestParseSttySize_MalformedOutputReturnsError(t *testing.T) {
	if _, _, err := parseSttySize("not a size"); err == nil {
		t.Error("expected an error for malformed stty size output")
	}
	if _, _, err := parseSttySize(""); err == nil {
		t.Error("expected an error for empty stty size output")
	}
	if _, _, err := parseSttySize("24 eighty\n"); err == nil {
		t.Error("expected an error for a non-numeric column count")
	}
}

func TestBoxTopLine_SpansRequestedWidthWithCorners(t *testing.T) {
	line := boxTopLine(10)
	if got := len([]rune(line)); got != 10 {
		t.Errorf("len(boxTopLine(10)) = %d, want 10", got)
	}
	if line[:len("╭")] != "╭" {
		t.Errorf("boxTopLine = %q, want it to start with ╭", line)
	}
}

func TestBoxMiddleLine_SpansRequestedWidthWithSidesAndBlankInterior(t *testing.T) {
	line := boxMiddleLine(10)
	runes := []rune(line)
	if len(runes) != 10 {
		t.Errorf("len(boxMiddleLine(10)) = %d, want 10", len(runes))
	}
	if string(runes[0]) != "│" || string(runes[len(runes)-1]) != "│" {
		t.Errorf("boxMiddleLine = %q, want │ at both ends", line)
	}
}

func TestBoxBottomLine_SpansRequestedWidthWithCorners(t *testing.T) {
	line := boxBottomLine(10)
	runes := []rune(line)
	if len(runes) != 10 {
		t.Errorf("len(boxBottomLine(10)) = %d, want 10", len(runes))
	}
	if string(runes[0]) != "╰" || string(runes[len(runes)-1]) != "╯" {
		t.Errorf("boxBottomLine = %q, want ╰...╯", line)
	}
}

func TestNewInputBoxAt_ErrorsWhenTerminalTooShortForTheBox(t *testing.T) {
	var buf bytes.Buffer
	if _, err := newInputBoxAt(&buf, scrollBoxHeight, 80); err == nil {
		t.Error("expected an error when rows leaves no room for a scroll region above the box")
	}
}

func TestNewInputBoxAt_SetsScrollRegionAndDrawsBorders(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\033[1;16r") {
		t.Errorf("output %q does not confine the scroll region to rows 1-16 (24 - scrollBoxHeight(3) - activityHeight(5))", out)
	}
	if !strings.Contains(out, boxTopLine(80)) {
		t.Errorf("output %q does not draw the box's top border", out)
	}
	if !strings.Contains(out, boxBottomLine(80)) {
		t.Errorf("output %q does not draw the box's bottom border", out)
	}
	if box.regionBottom != 16 {
		t.Errorf("regionBottom = %d, want 16", box.regionBottom)
	}
	// drawBorders' own absolute cursor positioning leaves the cursor
	// inside the box (wherever it last wrote); setup must reposition
	// back to Home afterward so whatever the caller prints next lands
	// at the top of the scroll region, not inside the box. A real pty
	// capture caught this failing when Home was sent before, not after,
	// drawBorders.
	if !strings.HasSuffix(out, "\033[H") {
		t.Errorf("output %q does not end with \\033[H (cursor left inside the box instead of at Home)", out)
	}
}

func TestInputBox_ShowPromptPositionsAtContentRowAndPrintsPrompt(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	box.showPrompt()

	out := buf.String()
	if !strings.Contains(out, "\033[23;1H") {
		t.Errorf("showPrompt output %q does not position at the box's content row", out)
	}
	if !strings.HasSuffix(out, "> ") {
		t.Errorf("showPrompt output %q does not end with the \"> \" prompt", out)
	}
}

func TestInputBox_ReturnToScrollPositionsAtScrollRegionBottom(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	box.returnToScroll()

	// Clears the row too (\033[2K), not just positions the cursor --
	// regionBottom is an absolute jump, not something reached by a
	// genuine scroll, so it carries no guarantee of being blank; a real
	// live session found this the hard way (see returnToScroll's doc
	// comment). Also clears the box's own content row (23 = regionBottom
	// + activityHeight + 2) first, so a stale cooked-mode echo of the
	// submitted line doesn't sit there looking like a duplicate until the
	// next prompt.
	if got := buf.String(); got != "\033[23;1H\033[2K\033[16;1H\033[2K" {
		t.Errorf("returnToScroll = %q, want \\033[23;1H\\033[2K\\033[16;1H\\033[2K", got)
	}
}

func TestInputBox_CloseResetsScrollRegionToFullScreen(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	box.close()

	if got := buf.String(); got != "\033[r" {
		t.Errorf("close = %q, want \\033[r", got)
	}
}

func TestRunScrollBoxInteractive_ShowsPromptReturnsToScrollAndEchoesEachLine(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	runScrollBoxInteractive(strings.NewReader("hello\nworld\n"), &buf, box)

	out := buf.String()
	if !strings.Contains(out, "\033[23;1H\033[2K> ") {
		t.Errorf("output %q missing a showPrompt call at the box's content row", out)
	}
	if !strings.Contains(out, "\033[16;1H") {
		t.Errorf("output %q missing a returnToScroll call to the scroll region's bottom row", out)
	}
	if !strings.Contains(out, "hello") || !strings.Contains(out, "world") {
		t.Errorf("output %q missing echoed input lines", out)
	}
	// showPrompt must run one more time than there are lines -- once
	// per real turn, plus once more before the read that hits EOF and
	// breaks the loop (matches tutor.go's own Run loop structure).
	// Counts the row-23-specific sequence, not the bare "\033[2K> "
	// substring -- returnToScroll (row 21) now also clears its line
	// (\033[2K) immediately before the caller's own "> "-prefixed echo
	// print, which would otherwise accidentally match the same substring.
	if got := strings.Count(out, "\033[23;1H\033[2K> "); got != 3 {
		t.Errorf("showPrompt ran %d times, want 3 (2 lines + 1 final EOF prompt)", got)
	}
}

func TestInputBox_ReconfigureAt_RowCountChangeClearsScreen(t *testing.T) {
	// Any row-count change clears the whole screen before redrawing.
	// Two real bugs drove this (see reconfigureAt's doc comment for the
	// full account of both):
	//
	//  1. The entrypoint.sh startup race, where the box's initial setup
	//     can run against a not-yet-final pane size.
	//  2. A genuine later, user-initiated resize mid-conversation. A
	//     first attempt tried to surgically clear only the box's *old*
	//     3 rows instead of the whole screen, to preserve visible
	//     history -- a live repro showed tmux visibly reflows/shifts
	//     already-printed rows when a pane's row count changes, so
	//     clearing rows computed from the *old* size can miss entirely
	//     and land on the wrong physical rows post-reflow. A full clear
	//     is the only approach confirmed live to survive that reflow.
	//
	// This test doesn't distinguish "before any real turn" from "mid
	// conversation" -- both clear identically now; there's no separate
	// state to gate on.
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	box.returnToScroll() // a real turn already happened -- still clears
	buf.Reset()

	box.reconfigureAt(40, 120) // grew taller and wider

	if box.regionBottom != 32 {
		t.Errorf("regionBottom = %d, want 32 (40 - scrollBoxHeight(3) - activityHeight(5))", box.regionBottom)
	}
	if box.cols != 120 {
		t.Errorf("cols = %d, want 120", box.cols)
	}
	out := buf.String()
	if !strings.Contains(out, "\033[2J") {
		t.Errorf("output %q does not clear the screen despite the row count changing", out)
	}
	if !strings.Contains(out, "\033[1;32r") {
		t.Errorf("output %q does not reset the scroll region to the new bounds", out)
	}
	if !strings.Contains(out, boxTopLine(120)) {
		t.Errorf("output %q does not redraw the box's top border at the new width", out)
	}
}

func TestInputBox_ReconfigureAt_EndsWithCursorAtScrollRegionBottomNotInsideTheBox(t *testing.T) {
	// drawBorders' own absolute positioning leaves the cursor at the
	// bottom border row -- reconfigureAt must reposition afterward the
	// same way setup() already does, or a caller that reconfigures and
	// then prints directly (tutor.go's drainResize, called right before
	// printing a reply) ends up printing into the box instead of the
	// scroll region. A real bug found live via a mid-generation resize:
	// the reply printed mixed in with the box's bottom border characters.
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	box.reconfigureAt(40, 120)

	want := "\033[32;1H\033[2K"
	if got := buf.String(); !strings.HasSuffix(got, want) {
		t.Errorf("reconfigureAt output %q does not end with %q (cursor left inside the box instead of at the scroll region's bottom row)", got, want)
	}
}

func TestInputBox_ReconfigureAt_ColsOnlyChangeRedrawsSameRowsWithoutClearingScreen(t *testing.T) {
	// A pure width change (rows unchanged) does NOT clear the whole
	// screen: the box stays at the same physical rows (no reflow risk,
	// confirmed live -- reflow was only observed when the row count
	// itself changed), and drawBorders already clears each of those
	// rows (\033[2K) before repainting them at the new width.
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	box.returnToScroll()
	buf.Reset()

	box.reconfigureAt(24, 120) // same rows, wider cols

	if box.regionBottom != 16 {
		t.Errorf("regionBottom = %d, want unchanged at 16", box.regionBottom)
	}
	out := buf.String()
	if strings.Contains(out, "\033[2J") {
		t.Errorf("output %q clears the screen -- a cols-only change must not, since the box's rows don't move and drawBorders already handles repainting them safely", out)
	}
	if !strings.Contains(out, boxTopLine(120)) {
		t.Errorf("output %q does not redraw the box's top border at the new width", out)
	}
	// Each row (22, 23, 24) should be cleared exactly once, by
	// drawBorders -- confirms there's no redundant extra clearing on
	// top of it for this case.
	for _, row := range []int{22, 23, 24} {
		want := fmt.Sprintf("\033[%d;1H\033[2K", row)
		if got := strings.Count(out, want); got != 1 {
			t.Errorf("row %d cleared %d times, want exactly 1 (drawBorders only, no redundant old-row clear for a cols-only change)", row, got)
		}
	}
}

func TestInputBox_ReconfigureAt_NoOpWhenSizeUnchanged(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	box.reconfigureAt(24, 80) // same size that newInputBoxAt was already given

	if got := buf.String(); got != "" {
		t.Errorf("reconfigureAt on an unchanged size wrote %q, want no output", got)
	}
}

func TestTruncateLine_ShortStringUnchanged(t *testing.T) {
	if got := truncateLine("hello", 10); got != "hello" {
		t.Errorf("truncateLine(%q, 10) = %q, want unchanged", "hello", got)
	}
}

func TestTruncateLine_LongStringTruncatedWithEllipsis(t *testing.T) {
	// ASCII "..." only -- a real bug found live: the Unicode ellipsis
	// (…) and every other symbol this package originally used (⟳ → ✓ ✗)
	// rendered as unrecognizable fallback glyphs (tofu, looking like
	// stray underscores) in a real user's terminal font. Everything this
	// package writes must be plain ASCII plus the one glyph confirmed to
	// render everywhere: ● (see activityThinkingStatus/formatActivityLine).
	got := truncateLine("this is a much longer string than the limit allows", 10)
	if runes := []rune(got); len(runes) != 10 {
		t.Errorf("truncateLine(...) = %q (len %d), want exactly 10 runes", got, len(runes))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("truncateLine(...) = %q, want it to end with \"...\"", got)
	}
}

func TestTruncateLine_MaxOfZeroOrLessReturnsEmpty(t *testing.T) {
	if got := truncateLine("anything", 0); got != "" {
		t.Errorf("truncateLine(_, 0) = %q, want empty", got)
	}
}

func TestTruncateLine_MaxSmallerThanEllipsisTruncatesTheEllipsisItself(t *testing.T) {
	if got := truncateLine("anything long enough to truncate", 2); got != ".." {
		t.Errorf("truncateLine(_, 2) = %q, want \"..\"", got)
	}
}

// --- inputBox.showActivity / clearActivity ---

func TestInputBox_ShowActivity_WritesStatusAndToolLinesAtReservedRows(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	box.showActivity("Thinking...", []string{"-> read_solution_file", "-> read_problem_statement"})

	out := buf.String()
	// Status row: regionBottom(16) + 1 = 17.
	if !strings.Contains(out, "\033[17;1H\033[2KThinking...") {
		t.Errorf("output %q missing the status line at row 17", out)
	}
	// Tool-call rows: regionBottom + 2..5 = 18-21 -- only the first two
	// carry real content, the remaining two (no calls yet) must still be
	// positioned and cleared so a shorter feed doesn't leave stale lines
	// from a longer one.
	if !strings.Contains(out, "\033[18;1H\033[2K-> read_solution_file") {
		t.Errorf("output %q missing the first tool line at row 18", out)
	}
	if !strings.Contains(out, "\033[19;1H\033[2K-> read_problem_statement") {
		t.Errorf("output %q missing the second tool line at row 19", out)
	}
	if !strings.Contains(out, "\033[20;1H\033[2K") {
		t.Errorf("output %q missing the (empty) third tool-line row 20 being cleared", out)
	}
	if !strings.Contains(out, "\033[21;1H\033[2K") {
		t.Errorf("output %q missing the (empty) fourth tool-line row 21 being cleared", out)
	}
}

func TestInputBox_ShowActivity_TruncatesLinesToBoxWidth(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 10) // narrow box
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	box.showActivity("a status line far longer than ten columns", nil)

	// Every line written (after each "\033[R;1H\033[2K" position+clear
	// pair) must be truncated to the box's current width, not just
	// passed through -- a real risk given tool names/results are
	// arbitrary-length strings, and an untruncated line could wrap onto
	// (and corrupt) the box's own reserved rows below it.
	out := buf.String()
	for _, part := range strings.Split(out, "\033[2K")[1:] {
		line := strings.SplitN(part, "\033[", 2)[0]
		if runes := []rune(line); len(runes) > 10 {
			t.Errorf("line %q is %d runes, want <= 10 (box width)", line, len(runes))
		}
	}
}

func TestInputBox_ShowActivity_SingleBufferedWrite(t *testing.T) {
	// One Write call for the whole frame, not one per row -- avoids
	// partial-frame flicker/interleaving risk, especially since this can
	// be called far more often (once per tool-call event, potentially
	// from concurrent goroutines) than drawBorders (once per
	// setup/resize).
	var buf countingWriter
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.writes = 0

	box.showActivity("status", []string{"one", "two"})

	if buf.writes != 1 {
		t.Errorf("showActivity made %d Write calls, want exactly 1", buf.writes)
	}
}

func TestInputBox_ClearActivity_BlanksEveryActivityRow(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	box.showActivity("Thinking...", []string{"-> read_solution_file"})
	buf.Reset()

	box.clearActivity()

	out := buf.String()
	for _, row := range []int{17, 18, 19, 20, 21} {
		want := fmt.Sprintf("\033[%d;1H\033[2K", row)
		if !strings.Contains(out, want) {
			t.Errorf("output %q missing clear of activity row %d", out, row)
		}
	}
	if strings.Contains(out, "Thinking") || strings.Contains(out, "read_solution_file") {
		t.Errorf("output %q still contains prior activity content", out)
	}
}

func TestInputBox_ShowActivity_ConcurrentCallsDoNotRace(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}

	// Simulates eino's own concurrent tool-call execution -- multiple
	// callback goroutines can call showActivity at the same time within
	// one turn. Run with -race to actually catch a data race, not just
	// rely on this not panicking.
	done := make(chan struct{})
	for i := 0; i < 8; i++ {
		go func(n int) {
			box.showActivity(fmt.Sprintf("status %d", n), []string{"line"})
			done <- struct{}{}
		}(i)
	}
	for i := 0; i < 8; i++ {
		<-done
	}
}

// countingWriter counts how many times Write is called, without caring
// about the bytes themselves -- used to assert a caller batches into one
// Write rather than many small ones.
type countingWriter struct {
	writes int
}

func (c *countingWriter) Write(p []byte) (int, error) {
	c.writes++
	return len(p), nil
}

func TestInputBox_ReconfigureAt_NoOpWhenTerminalTooShort(t *testing.T) {
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	box.reconfigureAt(scrollBoxHeight, 80) // shrunk to leave no room for a scroll region above the box

	if got := buf.String(); got != "" {
		t.Errorf("reconfigureAt on a too-short size wrote %q, want no output", got)
	}
	if box.regionBottom != 16 {
		t.Errorf("regionBottom changed to %d despite the too-short size being rejected, want it left at 16", box.regionBottom)
	}
}

func TestInputBox_ReconfigureAt_NoOpWhenTerminalTooShortForActivityRegionEvenIfTallEnoughForJustTheBox(t *testing.T) {
	// scrollBoxHeight(3) + activityHeight(5) = 8 rows minimum; a terminal
	// tall enough for the box alone (e.g. rows=6) but not tall enough once
	// the activity region is also reserved must still be rejected, not
	// silently drop the activity region.
	var buf bytes.Buffer
	box, err := newInputBoxAt(&buf, 24, 80)
	if err != nil {
		t.Fatalf("newInputBoxAt: %v", err)
	}
	buf.Reset()

	box.reconfigureAt(scrollBoxHeight+activityHeight, 80) // regionBottom would be exactly 0

	if got := buf.String(); got != "" {
		t.Errorf("reconfigureAt on a too-short size wrote %q, want no output", got)
	}
	if box.regionBottom != 16 {
		t.Errorf("regionBottom changed to %d despite the too-short size being rejected, want it left at 16", box.regionBottom)
	}
}
