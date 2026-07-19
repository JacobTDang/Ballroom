package tutor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/cloudwego/eino/schema"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

func TestEstimatedTextareaRows_ShortLineIsOneRow(t *testing.T) {
	if got := estimatedTextareaRows("hello", 80); got != 1 {
		t.Errorf("estimatedTextareaRows(%q, 80) = %d, want 1", "hello", got)
	}
}

func TestEstimatedTextareaRows_LongSingleLineWrapsAcrossMultipleRows(t *testing.T) {
	// A real bug found live: LineCount() (explicit newlines only) never
	// grows for a single long line that visually WRAPS across the box's
	// width -- exactly the scenario that used to corrupt the old
	// hand-rolled ANSI box. Width-aware estimation from the raw value is
	// what makes growth track a wrapped single line too, not just
	// explicit paragraph breaks.
	line := strings.Repeat("x", 25) // 25 runes, width 10 -> ceil(25/10) = 3 rows
	if got := estimatedTextareaRows(line, 10); got != 3 {
		t.Errorf("estimatedTextareaRows(25 x's, width 10) = %d, want 3", got)
	}
}

func TestEstimatedTextareaRows_MultipleLogicalLinesSumTheirWrappedRows(t *testing.T) {
	value := strings.Repeat("x", 15) + "\n" + strings.Repeat("y", 5) // width 10: ceil(15/10)=2, ceil(5/10)=1 -> 3
	if got := estimatedTextareaRows(value, 10); got != 3 {
		t.Errorf("estimatedTextareaRows(...) = %d, want 3", got)
	}
}

func TestEstimatedTextareaRows_EmptyValueIsOneRow(t *testing.T) {
	if got := estimatedTextareaRows("", 80); got != 1 {
		t.Errorf("estimatedTextareaRows(\"\", 80) = %d, want 1", got)
	}
}

func TestEstimatedTextareaRows_ZeroOrNegativeWidthTreatedAsOne(t *testing.T) {
	// Guards the division below from ever seeing a non-positive width --
	// a real risk during startup/resize edge cases before a real size is
	// known yet.
	if got := estimatedTextareaRows("hello", 0); got != 5 {
		t.Errorf("estimatedTextareaRows(hello, 0) = %d, want 5 (treated as width 1)", got)
	}
}

func TestNewTutorModel_TextareaIsFocusedAndEnterDoesNotInsertANewline(t *testing.T) {
	m := newTutorLayoutOnly()
	if !m.textarea.Focused() {
		t.Error("expected the textarea to start focused")
	}
	if m.textarea.KeyMap.InsertNewline.Enabled() {
		t.Error("expected InsertNewline disabled -- Enter must submit, not insert a newline")
	}
}

func TestTutorModel_WindowSizeMsg_SetsViewportAndTextareaWidth(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	got := newM.(tutorModel)

	if got.width != 100 || got.height != 30 {
		t.Errorf("width,height = %d,%d, want 100,30", got.width, got.height)
	}
	// The thinking aurora is a background (aurora.go), not a frame --
	// it reserves no cells, so content sizes against the full terminal
	// width.
	if got.viewport.Width != 100 {
		t.Errorf("viewport.Width = %d, want 100", got.viewport.Width)
	}
	if got.textarea.Width() <= 0 || got.textarea.Width() >= 100 {
		t.Errorf("textarea.Width() = %d, want less than 100 (room left for the border)", got.textarea.Width())
	}
}

func TestTutorModel_WindowSizeMsg_ViewportAndTextareaHeightsSumWithinTerminal(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	got := newM.(tutorModel)

	border := textareaBoxStyle.GetVerticalFrameSize()
	total := got.viewport.Height + got.textarea.Height() + border + statusBarHeight
	if total > 30 {
		t.Errorf("viewport(%d) + textarea(%d) + border(%d) + statusbar(%d) = %d, want <= terminal height 30", got.viewport.Height, got.textarea.Height(), border, statusBarHeight, total)
	}
}

func TestTutorModel_TypingALongLineGrowsTextareaHeight(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 30})
	m = newM.(tutorModel)
	before := m.textarea.Height()

	long := strings.Repeat("a very long message that should wrap across several rows ", 5)
	for _, r := range long {
		newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = newM.(tutorModel)
	}

	if m.textarea.Height() <= before {
		t.Errorf("textarea.Height() = %d after a long line, want it to have grown past the starting %d", m.textarea.Height(), before)
	}
}

func TestTutorModel_TextareaHeightNeverExceedsHalfTheTerminal(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 20})
	m = newM.(tutorModel)

	// A single absurdly long line -- growth must be capped, not
	// unbounded, or a pathological message could starve the viewport
	// entirely. The cap scales with the real terminal height (not a
	// hardcoded row count), matching the "dynamic, not hardcoded" ask.
	huge := strings.Repeat("x", 2000)
	m.textarea.SetValue(huge)
	newM, _ = m.Update(tea.WindowSizeMsg{Width: 40, Height: 20})
	m = newM.(tutorModel)

	if m.textarea.Height() > 10 { // height/2
		t.Errorf("textarea.Height() = %d, want capped at half the terminal height (10)", m.textarea.Height())
	}
	if m.viewport.Height < minViewportRows {
		t.Errorf("viewport.Height = %d, want it to keep at least minViewportRows (%d) even with a huge textarea", m.viewport.Height, minViewportRows)
	}
}

func TestTutorModel_EnterOnEmptyTextareaDoesNothing(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := newM.(tutorModel)
	if got.textarea.Value() != "" {
		t.Errorf("textarea.Value() = %q, want still empty", got.textarea.Value())
	}
	if cmd != nil {
		t.Error("expected no command from submitting an empty message")
	}
}

func TestTutorModel_EnterOnNonEmptyTextareaResetsItAndEchoesIntoViewport(t *testing.T) {
	// Checks submit()'s SYNCHRONOUS behavior specifically -- the textarea
	// reset and the immediate echo, both of which happen before the
	// turn's tea.Cmd is ever run -- decoupled from whether the async
	// turn itself succeeds (see the Stage 2/3 tests below for that).
	// Needs a real newTutorModel (not newTutorLayoutOnly): submit() now
	// unconditionally starts a real turn, which would nil-pointer-panic
	// against newTutorLayoutOnly's unset agents.
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly
	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m.textarea.SetValue("hello there")

	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := newM.(tutorModel)

	if got.textarea.Value() != "" {
		t.Errorf("textarea.Value() = %q, want reset to empty after submit", got.textarea.Value())
	}
	if !strings.Contains(got.viewport.View(), "hello there") {
		t.Errorf("viewport view %q, want it to contain the submitted message", got.viewport.View())
	}
}

func TestTutorModel_View_RendersBothViewportAndTextarea(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	out := m.View()
	if out == "" {
		t.Fatal("View() returned empty output")
	}
}

// TestTutorModel_InputIsFullRoundedBox pins the input's opencode-style
// full rounded frame — an explicit user choice (2026-07-17 restyle)
// superseding the earlier top-rule-only design. The original
// "sidebar" complaint that killed the first full box was about its
// bright teal left edge; this frame is the dim structural paneRule,
// which the user accepted knowingly (see textareaBoxStyle's doc
// comment in styles.go).
func TestTutorModel_InputIsFullRoundedBox(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	out := m.View()
	for _, corner := range []string{"╭", "╮", "╰", "╯"} {
		if !strings.Contains(out, corner) {
			t.Errorf("View() missing %q, want the input framed as a full rounded box", corner)
		}
	}
}

// TestTutorModel_PageUpScrollsViewportNotTextarea is a regression test
// for a real bug found live: the user had no way to scroll up and see
// earlier conversation history at all -- Update() forwarded every key
// to the textarea unconditionally, and refreshViewport's GotoBottom
// meant the viewport was always pinned to the latest content with no
// way to move it. PageUp/PageDown are dedicated to viewport scrolling
// specifically because they're never used for normal text editing
// (unlike arrow keys, which move the cursor within a multi-line draft,
// or bubbles/viewport's own default single-letter bindings like "j"/"k",
// which would otherwise swallow normal typing).
func TestTutorModel_PageUpScrollsViewportNotTextarea(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	m = newM.(tutorModel)

	blocks := make([]displayBlock, 0, 40)
	for i := 0; i < 40; i++ {
		blocks = append(blocks, displayBlock{kind: blockNote, raw: fmt.Sprintf("line %d", i)})
	}
	m.displayBlocks = blocks
	m.refreshViewport()
	if !m.viewport.AtBottom() {
		t.Fatal("setup: expected viewport to start at the bottom (refreshViewport's GotoBottom)")
	}

	m.textarea.SetValue("draft")
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	got := newM.(tutorModel)

	if got.viewport.AtBottom() {
		t.Error("expected PageUp to scroll the viewport up, but it's still at the bottom")
	}
	if got.textarea.Value() != "draft" {
		t.Errorf("textarea.Value() = %q, want the draft untouched by PageUp", got.textarea.Value())
	}
}

func TestTutorModel_PageDownScrollsViewportBackToBottom(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	m = newM.(tutorModel)

	blocks := make([]displayBlock, 0, 40)
	for i := 0; i < 40; i++ {
		blocks = append(blocks, displayBlock{kind: blockNote, raw: fmt.Sprintf("line %d", i)})
	}
	m.displayBlocks = blocks
	m.refreshViewport()

	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	m = newM.(tutorModel)
	if m.viewport.AtBottom() {
		t.Fatal("setup: expected PageUp to have scrolled away from the bottom")
	}

	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	got := newM.(tutorModel)
	if !got.viewport.AtBottom() {
		t.Error("expected PageDown to scroll back down to the bottom")
	}
}

func TestTutorModel_MouseWheelScrollsTheViewport(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	m = newM.(tutorModel)

	blocks := make([]displayBlock, 0, 40)
	for i := 0; i < 40; i++ {
		blocks = append(blocks, displayBlock{kind: blockNote, raw: fmt.Sprintf("line %d", i)})
	}
	m.displayBlocks = blocks
	m.refreshViewport()

	newM, _ = m.Update(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelUp})
	got := newM.(tutorModel)

	if got.viewport.AtBottom() {
		t.Error("expected a mouse wheel-up event to scroll the viewport up")
	}
}

// TestTutorModel_NormalLetterKeysStillGoToTheTextarea guards against a
// too-broad fix for the PageUp/PageDown scrolling above: bubbles/viewport's
// own default key bindings include plain letters ("j"/"k"/"h"/"l"/"f"/"b"/
// "d"/"u") for vim-style scrolling -- forwarding all keys to the viewport
// instead of being selective about it would silently eat normal typing.
func TestTutorModel_NormalLetterKeysStillGoToTheTextarea(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	m = newM.(tutorModel)

	for _, r := range "hello jkl" {
		newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		m = newM.(tutorModel)
	}

	if m.textarea.Value() != "hello jkl" {
		t.Errorf("textarea.Value() = %q, want %q -- letters that double as viewport scroll bindings must still type normally", m.textarea.Value(), "hello jkl")
	}
}

// TestTutorModel_LongReplyWrapsInsteadOfBeingCutOff is a regression test
// for a real bug found live: bubbles/viewport does not wrap long lines
// -- it only splits on explicit "\n", and a line wider than the
// viewport gets hard-truncated at the frame edge with the rest silently
// discarded, not shown on a wrapped row below. A real assistant reply is
// normally one long unbroken line (no embedded newlines), so every
// reply longer than the pane was getting cut off mid-sentence instead of
// wrapping. refreshViewport now pre-wraps content before SetContent
// (see its doc comment) specifically to fix this.
func TestTutorModel_LongReplyWrapsInsteadOfBeingCutOff(t *testing.T) {
	long := strings.Repeat("supercalifragilistic ", 20) // one long unbroken line, no "\n"
	mock := newSequencedOllama(t, long)
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "give me a long reply")

	view := got.viewport.View()
	if !strings.Contains(view, "supercalifragilistic") {
		t.Fatalf("viewport view %q, want it to contain the reply", view)
	}
	for _, line := range strings.Split(view, "\n") {
		if w := lipgloss.Width(line); w > 40 {
			t.Errorf("view line %q is %d cells wide, want it wrapped within the 40-wide viewport, not overflowing", line, w)
		}
	}
	// Every repetition of the word must still be visible somewhere in
	// the view (each on its own wrapped row, since one repetition plus a
	// trailing space already exceeds the 38-wide content area), not
	// silently dropped by a hard cut after the first screen-width.
	if n := strings.Count(view, "supercalifragilistic"); n != 20 {
		t.Errorf("view contains %d occurrences of the repeated word, want all 20 -- the tail of a long reply must not be truncated", n)
	}
}

// --- Stage 2: real turn-loop logic, re-homed from the old Run()-level
// tests in tutor_test.go to test tutorModel.Update() directly -- same
// proven pattern internal/tui/app_test.go already uses for appModel.
// submitAndRun types line into m's textarea, presses Enter, and
// synchronously executes the returned tea.Cmd (a plain func() tea.Msg --
// no real tea.Program needed) feeding its result back into Update(),
// mirroring internal/tui's own async-message test pattern.

func submitAndRun(t *testing.T, m tutorModel, line string) tutorModel {
	t.Helper()
	// Detection is async (from Init) in a real session and always lands
	// before the user can type -- deliver it here the same way if this
	// model hasn't resolved its strategies yet, so the submit below
	// starts a real turn instead of queueing.
	if !m.strategiesDetected {
		m = detectAndApply(t, m)
	}
	m.textarea.SetValue(line)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)
	if cmd == nil {
		t.Fatal("submit produced no command -- expected a turn to start")
	}
	// Drains however many activityEventMsgs precede the final
	// turnCompleteMsg (zero for a turn with no tool calls -- the
	// activity channel closes immediately and the first cmd() call
	// already yields the result; one or more for a real tool-calling
	// turn), each re-arming the wait the same way the real bubbletea
	// runtime would.
	for i := 0; i < 100; i++ {
		msg := cmd()
		newM, cmd = m.Update(msg)
		m = newM.(tutorModel)
		if _, ok := msg.(turnCompleteMsg); ok {
			return m
		}
		if cmd == nil {
			t.Fatal("submitAndRun: turn ended without a turnCompleteMsg")
		}
	}
	t.Fatal("submitAndRun: too many iterations, possible infinite activity-event loop")
	return m
}

func TestNewTutorModel_NoRoutingHeaderAndHistorySeededWithSystemPrompt(t *testing.T) {
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	if len(m.history) != 1 || m.history[0].Role != schema.System {
		t.Fatalf("history = %+v, want exactly one System message seeded", m.history)
	}
	status := m.statusLeftText()
	if !strings.Contains(status, cfg.Model) || strings.Contains(status, "+") {
		t.Errorf("status bar left = %q, want it to name the model and not mention a second model", status)
	}
	if bar := m.statusBarView(); strings.Contains(bar, "orchestrator") {
		t.Errorf("status bar = %q, want no routing mention without an orchestrator", bar)
	}
}

// TestTutorModel_StatusBarPinnedBelowInputWithEndpointAndExitHint pins
// the status bar's contract: it renders as the last line of every
// frame (not as a transcript entry that scrolls away), shows the
// session's identity and endpoint, and keeps the exit hint the old
// header carried.
// TestNewTutorModel_WritesTutorStateDotfileAtStartup covers the initial,
// pre-any-turn write: session/submit.go and session/reference.go read
// this file to attach hints_used/tutor_mode/model to an attempt, and
// must see real (zero, not missing) values even for a submit or
// reference reveal that happens before the user ever asks the tutor
// anything.
func TestNewTutorModel_WritesTutorStateDotfileAtStartup(t *testing.T) {
	mock := newSequencedOllama(t, "hi")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly
	cfg.Model = "llama3.1:8b"
	cfg.WorkDir = t.TempDir()

	if _, err := newTutorModel(context.Background(), cfg); err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(cfg.WorkDir, tutorStateFile))
	if err != nil {
		t.Fatalf("read tutor state dotfile: %v", err)
	}
	var got tutorState
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal tutor state: %v", err)
	}
	if got.HintsUsed != 0 {
		t.Errorf("HintsUsed = %d, want 0 before any turn", got.HintsUsed)
	}
	if got.TutorMode != exercise.TutorModeSyntaxOnly {
		t.Errorf("TutorMode = %q, want %q", got.TutorMode, exercise.TutorModeSyntaxOnly)
	}
	if got.Model != "llama3.1:8b" {
		t.Errorf("Model = %q, want %q", got.Model, "llama3.1:8b")
	}
}

// TestTutorModel_TurnCompleteUpdatesTutorStateDotfileWithHintsUsed covers
// the live-updating half: the count session/submit.go and
// session/reference.go pick up must track m.helpRequestCount as it
// actually grows, not just the startup snapshot.
func TestTutorModel_TurnCompleteUpdatesTutorStateDotfileWithHintsUsed(t *testing.T) {
	mock := newSequencedOllama(t, "reply one", "reply two")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly // skips the comprehension check
	cfg.WorkDir = t.TempDir()

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	m = submitAndRun(t, m, "first question")

	readState := func() tutorState {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(cfg.WorkDir, tutorStateFile))
		if err != nil {
			t.Fatalf("read tutor state dotfile: %v", err)
		}
		var s tutorState
		if err := json.Unmarshal(data, &s); err != nil {
			t.Fatalf("unmarshal tutor state: %v", err)
		}
		return s
	}

	if got := readState(); got.HintsUsed != 1 {
		t.Errorf("HintsUsed after one turn = %d, want 1", got.HintsUsed)
	}

	m = submitAndRun(t, m, "second question")
	if got := readState(); got.HintsUsed != 2 {
		t.Errorf("HintsUsed after two turns = %d, want 2", got.HintsUsed)
	}
}

// TestWriteTutorState_EmptyWorkDirIsANoOp guards against a real
// footgun: most of this package's tests build a Config via testConfig,
// which leaves WorkDir empty, and newTutorModel now writes on every
// construction -- without this guard, every one of those tests would
// litter a stray dotfile into the process's current working directory
// (the package source tree under `go test`) instead of a real
// exercise workspace.
func TestWriteTutorState_EmptyWorkDirIsANoOp(t *testing.T) {
	writeTutorState("", tutorState{HintsUsed: 3, TutorMode: "hints-first", Model: "x"})
	if _, err := os.Stat(tutorStateFile); err == nil {
		os.Remove(tutorStateFile)
		t.Fatal("writeTutorState with an empty workDir wrote into the current directory")
	}
}

func TestTutorModel_StatusBarPinnedBelowInputWithEndpointAndExitHint(t *testing.T) {
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 24})
	m = newM.(tutorModel)

	view := m.View()
	lines := strings.Split(view, "\n")
	plainLast := stripAnsiTest(lines[len(lines)-1])
	if !strings.Contains(plainLast, cfg.Model) || !strings.Contains(plainLast, letterspacedUpper(cfg.Mode)) {
		t.Errorf("View's last line = %q, want the model and mode pill pinned there", plainLast)
	}
	if !strings.Contains(plainLast, mock.URL) {
		t.Errorf("View's last line = %q, want the endpoint (%s) shown on the right", plainLast, mock.URL)
	}
	if !strings.Contains(plainLast, "ctrl+d") {
		t.Errorf("View's last line = %q, want the ctrl+d exit hint kept from the old header", plainLast)
	}
	if len(m.displayBlocks) != 0 {
		t.Errorf("displayBlocks = %+v, want the transcript to start empty -- no banner block", m.displayBlocks)
	}
}

// withFakeCheckToolCallingForSession overrides checkToolCallingForSession
// for the duration of one test (TestMain already defaults it to "always
// native"), restoring it on cleanup -- same save/restore pattern
// internal/tui/app_test.go's fakeCheckToolCalling uses.
func withFakeCheckToolCallingForSession(t *testing.T, fn func(ctx context.Context, ollamaHost, model, apiKey string) (bool, error)) {
	t.Helper()
	orig := checkToolCallingForSession
	checkToolCallingForSession = fn
	t.Cleanup(func() { checkToolCallingForSession = orig })
}

// detectAndApply synchronously runs m's strategy-detection command and
// feeds its message back through Update -- what the real bubbletea
// runtime does asynchronously after Init(). Tests that need a model
// with strategies already resolved call this right after newTutorModel
// (detection no longer happens during construction -- see
// TestNewTutorModel_DoesNotProbeDuringConstruction for why).
func detectAndApply(t *testing.T, m tutorModel) tutorModel {
	t.Helper()
	msg := m.detectStrategies()()
	if _, ok := msg.(strategiesDetectedMsg); !ok {
		t.Fatalf("detectStrategies produced %T, want strategiesDetectedMsg", msg)
	}
	newM, _ := m.Update(msg)
	return newM.(tutorModel)
}

// TestNewTutorModel_DoesNotProbeDuringConstruction locks in the fix for
// a real bug found live: strategy detection used to run (and block on
// wg.Wait) inside newTutorModel, BEFORE Run() ever started the
// tea.Program -- so with slow free-tier models the whole tutor pane,
// chat box included, stayed completely blank for the duration of up to
// two live LLM round-trips. Construction must make zero probe calls;
// detection belongs to Init()'s async command.
func TestNewTutorModel_DoesNotProbeDuringConstruction(t *testing.T) {
	var probes atomic.Int32
	withFakeCheckToolCallingForSession(t, func(context.Context, string, string, string) (bool, error) {
		probes.Add(1)
		return false, nil
	})
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	if n := probes.Load(); n != 0 {
		t.Fatalf("newTutorModel made %d strategy probes, want 0 -- detection must not block construction", n)
	}
	m = detectAndApply(t, m)
	if n := probes.Load(); n != 1 {
		t.Errorf("detectStrategies made %d probes for a worker-only session, want 1", n)
	}
	if m.workerStrategy != jsonFallbackToolCalling {
		t.Errorf("workerStrategy = %v after detection, want jsonFallbackToolCalling", m.workerStrategy)
	}
}

// Init must actually carry the detection command -- if it didn't, a
// real session would paint fine but any message submitted before
// detection would sit queued forever, waiting on a strategiesDetectedMsg
// that never comes.
func TestTutorModel_InitIncludesStrategyDetection(t *testing.T) {
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	batch, ok := m.Init()().(tea.BatchMsg)
	if !ok {
		t.Fatalf("Init() produced %T, want tea.BatchMsg", m.Init()())
	}
	for _, cmd := range batch {
		if _, ok := cmd().(strategiesDetectedMsg); ok {
			return
		}
	}
	t.Error("no command in Init()'s batch produced a strategiesDetectedMsg")
}

func TestTutorModel_SubmitBeforeDetectionQueuesTurnUntilStrategiesArrive(t *testing.T) {
	mock := newSequencedOllama(t, "queued reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	// Submit while detection is still pending -- the message must echo
	// and the thinking state must engage, but no turn can start yet.
	m.textarea.SetValue("early question")
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)
	if cmd != nil {
		t.Fatal("submit before detection returned a command, want the turn queued until strategies arrive")
	}
	if !m.turnInFlight {
		t.Error("turnInFlight = false after a queued submit, want true -- the thinking indicator must engage immediately")
	}
	if !strings.Contains(ansi.Strip(m.View()), "│ early question") {
		t.Error("queued submit's echo missing from the view")
	}

	// Detection lands -> the queued turn must start and complete.
	newM, cmd = m.Update(m.detectStrategies()())
	m = newM.(tutorModel)
	if cmd == nil {
		t.Fatal("strategiesDetectedMsg with a queued submit produced no command, want the turn to start")
	}
	for i := 0; i < 100; i++ {
		msg := cmd()
		newM, cmd = m.Update(msg)
		m = newM.(tutorModel)
		if _, ok := msg.(turnCompleteMsg); ok {
			if !strings.Contains(ansi.Strip(m.View()), "queued reply") {
				t.Error("queued turn's reply missing from the view after completion")
			}
			return
		}
		if cmd == nil {
			t.Fatal("queued turn ended without a turnCompleteMsg")
		}
	}
	t.Fatal("queued turn never completed")
}

func TestNewTutorModel_DetectsFallbackStrategyWhenCheckReportsUnsupported(t *testing.T) {
	withFakeCheckToolCallingForSession(t, func(context.Context, string, string, string) (bool, error) {
		return false, nil
	})
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	m = detectAndApply(t, m)
	if m.workerStrategy != jsonFallbackToolCalling {
		t.Errorf("workerStrategy = %v, want jsonFallbackToolCalling", m.workerStrategy)
	}
}

func TestNewTutorModel_DetectsNativeStrategyWhenCheckReportsSupported(t *testing.T) {
	withFakeCheckToolCallingForSession(t, func(context.Context, string, string, string) (bool, error) {
		return true, nil
	})
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	m = detectAndApply(t, m)
	if m.workerStrategy != nativeToolCalling {
		t.Errorf("workerStrategy = %v, want nativeToolCalling", m.workerStrategy)
	}
}

// TestNewTutorModel_DetectsFallbackStrategyOnCheckError locks in a
// real, live-reproduced case (not a hypothetical): OpenRouter's free
// meta-llama/llama-3.2-3b-instruct:free returns a 404 "No endpoints
// found that support tool use" for ANY request that binds a tools
// parameter at all -- including CheckToolCalling's own probe. Defaulting
// to nativeToolCalling on that error (an earlier version of this
// function did) breaks the session completely: the real worker/
// orchestrator agent ALSO binds tools via WithTools, so every single
// real turn 404s identically -- confirmed live via cmd/ballroom against
// the real OpenRouter API before this test existed. runFallbackToolLoop
// never binds tools via the API at all (it teaches the model about
// tools through the text prompt instead), so it's immune to this
// specific failure -- meaning an error here should default to
// fallback, not native, since fallback is the strategy actually capable
// of working when the check itself couldn't even complete.
func TestNewTutorModel_DetectsFallbackStrategyOnCheckError(t *testing.T) {
	withFakeCheckToolCallingForSession(t, func(context.Context, string, string, string) (bool, error) {
		return false, fmt.Errorf("simulated check failure")
	})
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	m = detectAndApply(t, m)
	if m.workerStrategy != jsonFallbackToolCalling {
		t.Errorf("workerStrategy = %v, want jsonFallbackToolCalling (fail toward the strategy that can actually work when the check itself fails)", m.workerStrategy)
	}
}

func TestNewTutorModel_DetectsOrchestratorStrategyIndependentlyFromWorker(t *testing.T) {
	withFakeCheckToolCallingForSession(t, func(_ context.Context, _, model, _ string) (bool, error) {
		return model == "orchestrator-model", nil
	})
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	m = detectAndApply(t, m)
	if m.workerStrategy != jsonFallbackToolCalling {
		t.Errorf("workerStrategy = %v, want jsonFallbackToolCalling", m.workerStrategy)
	}
	if m.orchestratorStrategy != nativeToolCalling {
		t.Errorf("orchestratorStrategy = %v, want nativeToolCalling", m.orchestratorStrategy)
	}
}

// TestTutorModel_ComprehensionCheckAndRoutedTurnUseOrchestratorFallbackStrategy
// is the gap a naive per-turn-only fix would silently miss: worker is
// native, orchestrator is jsonFallbackToolCalling, and BOTH the
// comprehension check (first message) and a routed-to real turn
// (second message) must each honor the orchestrator's own strategy, not
// the worker's. Proven by scripting the orchestrator's real-turn reply
// as a tool-call-shaped {"name": ..., "arguments": ...} blob: dispatched
// to the native path, that text would only ever be caught by
// leakedToolCallPattern and produce the honest "couldn't get a grounded
// answer" fallback after a failed retry -- it would never actually
// invoke the tool. Dispatched correctly to runFallbackToolLoop, it's
// parsed as a real call, actually executes read_solution_file, and the
// conversation reaches the real final answer scripted after it.
func TestTutorModel_ComprehensionCheckAndRoutedTurnUseOrchestratorFallbackStrategy(t *testing.T) {
	withFakeCheckToolCallingForSession(t, func(_ context.Context, _, model, _ string) (bool, error) {
		return model != "orchestrator-model", nil
	})
	mock := newSequencedOllama(t,
		"restated problem, clarifying questions", // comprehension check (orchestrator)
		"NO",                                     // routing decision -> orchestrator handles the real turn
		`{"name": "read_solution_file", "arguments": {}}`,  // orchestrator's real-turn reply: a tool call
		"orchestrator final answer after reading the file", // final answer once the tool result is fed back
	)
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"
	cfg.Mode = exercise.TutorModeFullAssist
	cfg.WorkDir = t.TempDir()

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	m = detectAndApply(t, m)
	if m.workerStrategy != nativeToolCalling {
		t.Fatalf("workerStrategy = %v, want nativeToolCalling", m.workerStrategy)
	}
	if m.orchestratorStrategy != jsonFallbackToolCalling {
		t.Fatalf("orchestratorStrategy = %v, want jsonFallbackToolCalling", m.orchestratorStrategy)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "hi, first message")
	if !strings.Contains(got.viewport.View(), "restated problem") {
		t.Fatalf("viewport after comprehension check = %q, want the comprehension check reply", got.viewport.View())
	}

	got = submitAndRun(t, got, "please look at my code")
	view := got.viewport.View()
	if !strings.Contains(view, "orchestrator final answer after reading the file") {
		t.Errorf("viewport = %q, want the real final answer -- the orchestrator's tool call should have been executed via the fallback loop, not treated as a native leak", view)
	}
	if strings.Contains(view, leakedToolCallFallbackReply) {
		t.Error("viewport contains the native leak-retry's honest-fallback message -- the routed turn was dispatched through generateWithLeakRetry instead of runFallbackToolLoop")
	}

	reqs := mock.allRequests()
	if len(reqs) != 4 {
		t.Fatalf("got %d requests, want 4 (comprehension check, routing decision, tool-call round, final-answer round): %+v", len(reqs), reqs)
	}
	foundToolResult := false
	for _, m := range reqs[3].Messages {
		if strings.Contains(m.Content, "Tool result:") {
			foundToolResult = true
		}
	}
	if !foundToolResult {
		t.Errorf("final request's messages = %+v, want one containing the real tool result", reqs[3].Messages)
	}
}

func TestNewTutorModel_ComprehensionCheckPendingMatchesMode(t *testing.T) {
	mock := newSequencedOllama(t, "reply")

	syntaxCfg := testConfig(mock.URL)
	syntaxCfg.Mode = exercise.TutorModeSyntaxOnly
	m, err := newTutorModel(context.Background(), syntaxCfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	if m.comprehensionCheckPending {
		t.Error("syntax-only mode must never want the comprehension check")
	}

	fullCfg := testConfig(mock.URL)
	fullCfg.Mode = exercise.TutorModeFullAssist
	m, err = newTutorModel(context.Background(), fullCfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	if !m.comprehensionCheckPending {
		t.Error("full-assist mode must want the comprehension check on the first message")
	}
}

func TestTutorModel_SubmitAppendsUserMessageImmediatelyAndReplyAfterCompletion(t *testing.T) {
	mock := newSequencedOllama(t, "the answer is 42")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "what does my code look like?")

	if !strings.Contains(got.viewport.View(), "the answer is 42") {
		t.Errorf("viewport view %q, want the assistant's reply", got.viewport.View())
	}
	if got.turnInFlight {
		t.Error("turnInFlight = true after the result arrived, want false")
	}
	if len(got.history) != 3 { // system + user + assistant
		t.Errorf("history has %d messages, want 3 (system, user, assistant)", len(got.history))
	}
}

func TestTutorModel_TurnFailureShowsFallbackAndDoesNotPersistToHistory(t *testing.T) {
	cfg := testConfig("http://127.0.0.1:1") // refuses immediately
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "hello")

	if !strings.Contains(got.viewport.View(), turnFailedFallbackReply) {
		t.Errorf("viewport view %q, want the honest fallback message", got.viewport.View())
	}
	if len(got.history) != 1 {
		t.Errorf("history has %d messages, want just the seeded system prompt (a failed turn is never persisted)", len(got.history))
	}
}

func TestTutorModel_NoRoutingWhenOrchestratorModelEmpty(t *testing.T) {
	mock := newSequencedOllama(t, "the only reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	submitAndRun(t, m, "hi")

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected exactly 1 request (no routing decision), got %d: %+v", len(reqs), reqs)
	}
	if reqs[0].Model != "test-model" {
		t.Errorf("request model = %q, want %q", reqs[0].Model, "test-model")
	}
}

func TestTutorModel_RoutesToOrchestratorWhenDecisionIsNo(t *testing.T) {
	mock := newSequencedOllama(t, "NO", "orchestrator answered")
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "hi")

	if !strings.Contains(got.viewport.View(), "orchestrator answered") {
		t.Errorf("viewport view %q, want the orchestrator's reply", got.viewport.View())
	}
	reqs := mock.allRequests()
	if len(reqs) != 2 {
		t.Fatalf("expected exactly 2 requests (routing decision + orchestrator answer), got %d: %+v", len(reqs), reqs)
	}
	for i, req := range reqs {
		if req.Model != "orchestrator-model" {
			t.Errorf("request[%d].Model = %q, want %q -- worker must never be touched when the decision is No", i, req.Model, "orchestrator-model")
		}
	}
}

func TestTutorModel_RoutesToWorkerWhenDecisionIsYes(t *testing.T) {
	mock := newSequencedOllama(t, "YES", "worker answered")
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "hi")

	if !strings.Contains(got.viewport.View(), "worker answered") {
		t.Errorf("viewport view %q, want the worker's reply", got.viewport.View())
	}
	reqs := mock.allRequests()
	if len(reqs) != 2 {
		t.Fatalf("expected exactly 2 requests (routing decision + worker answer), got %d: %+v", len(reqs), reqs)
	}
}

// TestTutorModel_RoutingDecisionFailureShowsWarningButTurnStillSucceeds
// is a regression test for a real bug found live: this warning used to
// go to a raw fmt.Fprintf(m.stderr, ...) call, which corrupted the
// real terminal (see activityErrorNote's doc comment in activity.go).
// It's rendered into the viewport instead now. decideHandoff already
// defaults to handoff (true) on its own request failure, so the turn
// still completes normally via the worker -- this only checks that the
// warning explaining *why* it defaulted is now visible in the chat.
func TestTutorModel_RoutingDecisionFailureShowsWarningButTurnStillSucceeds(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req tutorChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Model == "orchestrator-model" {
			// The routing decision itself always goes to the
			// orchestrator (see decideHandoff) -- fail only that
			// request, so the worker's own answer still succeeds.
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"orchestrator unreachable"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"role": "assistant", "content": "worker answered"},
			"done":    true,
		})
	}))
	defer mock.Close()

	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "hi")

	view := got.viewport.View()
	if !strings.Contains(view, "routing decision failed") {
		t.Errorf("viewport view %q, want the routing-failure warning visible", view)
	}
	if !strings.Contains(view, "worker answered") {
		t.Errorf("viewport view %q, want the turn to still succeed via the defaulted handoff", view)
	}
}

func TestTutorModel_ComprehensionCheckAlwaysUsesOrchestratorWhenRoutingEnabled(t *testing.T) {
	mock := newSequencedOllama(t, "restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"
	cfg.Mode = exercise.TutorModeFullAssist

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	submitAndRun(t, m, "hi")

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected exactly 1 request (the comprehension check, no routing decision), got %d: %+v", len(reqs), reqs)
	}
	if reqs[0].Model != "orchestrator-model" {
		t.Errorf("request[0].Model = %q, want %q -- the comprehension check always uses the orchestrator", reqs[0].Model, "orchestrator-model")
	}
}

func TestTutorModel_ComprehensionCheckClearsPendingRegardlessOfOutcome(t *testing.T) {
	mock := newSequencedOllama(t, "restated + questions", "second reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "hi")
	if got.comprehensionCheckPending {
		t.Error("comprehensionCheckPending must be false after the first message, win or lose")
	}

	got2 := submitAndRun(t, got, "a second message")
	reqs := mock.allRequests()
	if len(reqs) != 2 {
		t.Fatalf("expected exactly 2 requests total (one check, one normal turn), got %d: %+v", len(reqs), reqs)
	}
	if !strings.Contains(got2.viewport.View(), "second reply") {
		t.Errorf("viewport view %q, want the second turn's real reply", got2.viewport.View())
	}
}

func TestTutorModel_RetriesWhenReplyLeaksFakeToolCallJSON(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "read_solution_file", "parameters": {}}`, "your code looks fine")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "what does my code look like?")

	view := got.viewport.View()
	if strings.Contains(view, `{"name"`) {
		t.Errorf("viewport view still contains leaked tool-call JSON: %q", view)
	}
	if !strings.Contains(view, "your code looks fine") {
		t.Errorf("viewport view %q, want the retry's clean reply", view)
	}
}

// --- Stage 3: live tool-call activity display, channel-based instead of
// scrollbox.go's raw ANSI writes.

func TestActivityEventMsg_UpdatesActiveCallsAndRearmsWait(t *testing.T) {
	activityCh := make(chan []activityCall, 1)
	doneCh := make(chan turnCompleteMsg, 1)
	calls := []activityCall{{name: "read_solution_file", status: "running"}}

	m := newTutorLayoutOnly()
	newM, cmd := m.Update(activityEventMsg{calls: calls, activityCh: activityCh, doneCh: doneCh})
	got := newM.(tutorModel)

	if len(got.activeCalls) != 1 || got.activeCalls[0].name != "read_solution_file" {
		t.Errorf("activeCalls = %+v, want the pushed snapshot", got.activeCalls)
	}
	if cmd == nil {
		t.Fatal("expected activityEventMsg to re-arm the wait")
	}

	// Re-armed wait must target the SAME channels -- closing activityCh
	// and pushing a result onto doneCh must be what the re-armed cmd
	// picks up next.
	close(activityCh)
	doneCh <- turnCompleteMsg{userMessage: "done"}
	msg := cmd()
	if _, ok := msg.(turnCompleteMsg); !ok {
		t.Errorf("re-armed wait produced %T, want turnCompleteMsg once activityCh closed", msg)
	}
}

func TestPulseTickMsg_IncrementsPhaseAndAlwaysRearms(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, cmd := m.Update(pulseTickMsg{})
	got := newM.(tutorModel)

	if got.pulsePhase != 1 {
		t.Errorf("pulsePhase = %d, want 1", got.pulsePhase)
	}
	if cmd == nil {
		t.Error("expected pulseTickMsg to always re-arm, even when idle -- cheap to run continuously, avoids needing to start/stop it per turn")
	}
}

func TestTutorModel_View_ShowsActivityRegionOnlyWhileTurnInFlight(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m.activeCalls = []activityCall{{name: "read_solution_file", status: "running"}}

	m.turnInFlight = false
	if strings.Contains(m.View(), "read_solution_file") {
		t.Error("expected no activity region when no turn is in flight")
	}

	m.turnInFlight = true
	// Strip styling first: the thinking aurora paints a background
	// escape before every glyph (see overlayAurora), so the call name
	// is present but never as one contiguous raw substring.
	if !strings.Contains(ansi.Strip(m.View()), "read_solution_file") {
		t.Error("expected the activity region to show the active call while a turn is in flight")
	}
}

func TestTutorModel_TurnCompleteClearsActiveCallsAndTurnInFlight(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnInFlight = true
	m.activeCalls = []activityCall{{name: "read_solution_file", status: "running"}}

	newM, _ := m.Update(turnCompleteMsg{reply: schema.AssistantMessage("done", nil), userMessage: "x"})
	got := newM.(tutorModel)

	if got.turnInFlight {
		t.Error("turnInFlight = true after turnCompleteMsg, want false")
	}
	if len(got.activeCalls) != 0 {
		t.Errorf("activeCalls = %+v, want cleared", got.activeCalls)
	}
}

func TestTutorModel_SubmitEchoCarriesUserAccentBar(t *testing.T) {
	m := newTutorLayoutOnly()
	m.textarea.SetValue("what about sharding?")
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	got := newM.(tutorModel)

	display := renderedBlocks(got)
	if !strings.Contains(ansi.Strip(display), "│ what about sharding?") {
		t.Errorf("echo missing the │ accent bar:\n%s", ansi.Strip(display))
	}
	if !strings.Contains(display, ansiFg(panePink)) {
		t.Error("the accent bar should carry the pink foreground escape")
	}
}

// TestRenderUserBlock_EveryWrappedLineCarriesTheBar pins the block's
// two invariants: the bar marks every wrapped line (not just the
// first), and every emitted line is already within the given width —
// which is what keeps refreshViewport's outer word-wrap from ever
// re-breaking a line and orphaning its continuation without the bar.
func TestRenderUserBlock_EveryWrappedLineCarriesTheBar(t *testing.T) {
	raw := "this is a long user question that will definitely wrap across several rows at a narrow width"
	got := renderUserBlock(raw, 40)

	lines := strings.Split(got, "\n")
	if len(lines) < 2 {
		t.Fatalf("renderUserBlock at width 40 produced %d line(s), want the text wrapped across several", len(lines))
	}
	for i, line := range lines {
		plain := ansi.Strip(line)
		if !strings.HasPrefix(plain, "│ ") {
			t.Errorf("line %d = %q, want every wrapped line to start with the │ bar", i, plain)
		}
		if w := lipgloss.Width(line); w > 40 {
			t.Errorf("line %d is %d cells wide, want within the 40-cell width", i, w)
		}
	}
	joined := strings.ReplaceAll(ansi.Strip(got), "│ ", "")
	for _, word := range strings.Fields(raw) {
		if !strings.Contains(joined, word) {
			t.Errorf("wrapped block lost the word %q from the raw text", word)
		}
	}
}

func TestRenderUserBlock_WidthZeroPrefixesWithoutWrapping(t *testing.T) {
	got := renderUserBlock("early question", 0)
	if plain := ansi.Strip(got); plain != "│ early question" {
		t.Errorf("renderUserBlock(width 0) = %q, want the bare barred line with no wrapping or padding", plain)
	}
}

// TestTutorModel_ResizeReWrapsUserAccentBars mirrors the editor cards'
// resize test: a user block re-flows at the new width and the bar
// still marks every line.
func TestTutorModel_ResizeReWrapsUserAccentBars(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m.displayBlocks = []displayBlock{{kind: blockUser, raw: "a question long enough to stay on one row when wide but wrap once narrow"}}
	m.refreshViewport()
	if wide := strings.Split(renderedBlocks(m), "\n"); len(wide) != 1 {
		t.Fatalf("user block spans %d lines at width 80, want 1 (the fixture must only wrap when narrow)", len(wide))
	}

	newM, _ = m.Update(tea.WindowSizeMsg{Width: 40, Height: 24})
	m = newM.(tutorModel)
	narrow := renderedBlocks(m)
	narrowLines := strings.Split(narrow, "\n")
	if len(narrowLines) <= 1 {
		t.Fatalf("user block did not re-wrap at width 40:\n%s", ansi.Strip(narrow))
	}
	for i, line := range narrowLines {
		if plain := ansi.Strip(line); !strings.HasPrefix(plain, "│ ") {
			t.Errorf("post-resize line %d = %q, want the bar on every wrapped line", i, plain)
		}
	}
}

// renderedBlocks renders m's transcript blocks the same way
// refreshViewport does, for tests asserting display styling without
// caring about viewport frame details.
func renderedBlocks(m tutorModel) string {
	w := m.viewport.Width - m.viewport.Style.GetHorizontalFrameSize()
	rendered := make([]string, 0, len(m.displayBlocks))
	for _, b := range m.displayBlocks {
		rendered = append(rendered, renderBlock(b, w))
	}
	return strings.Join(rendered, "\n\n")
}

func TestTutorModel_TurnCompleteStylesReplyForDisplayButKeepsHistoryRaw(t *testing.T) {
	m := newTutorLayoutOnly()
	m.turnInFlight = true
	raw := "use a **hash set** and call `add()`"

	newM, _ := m.Update(turnCompleteMsg{reply: schema.AssistantMessage(raw, nil), userMessage: "how?"})
	got := newM.(tutorModel)

	display := renderedBlocks(got)
	if strings.Contains(display, "**") || strings.Contains(display, "`") {
		t.Errorf("rendered transcript still carries raw markdown markers:\n%s", display)
	}
	if !strings.Contains(display, "\x1b[1m") {
		t.Error("rendered transcript has no bold escape -- reply wasn't styled for display")
	}
	last := got.history[len(got.history)-1]
	if last.Content != raw {
		t.Errorf("history got %q, want the raw unstyled reply %q -- escape codes must never enter model context", last.Content, raw)
	}
}

func TestTutorModel_SubmitShowsLiveToolCallActivity(t *testing.T) {
	// newToolCallOllama (toolcheck_test.go) simulates a real tool_calls
	// response for its first request, then a plain-text reply for its
	// second -- driving a real read_solution_file call through the real
	// eino callback -> channel -> activityEventMsg pipeline, not a
	// synthetic message.
	mock := newToolCallOllama(t, "read_solution_file")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly
	cfg.WorkDir = t.TempDir()

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m = detectAndApply(t, m)

	m.textarea.SetValue("what does my code look like?")
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)
	if cmd == nil {
		t.Fatal("submit produced no command")
	}

	sawActivity := false
	for i := 0; i < 100; i++ {
		msg := cmd()
		if ev, ok := msg.(activityEventMsg); ok {
			for _, c := range ev.calls {
				if c.name == "read_solution_file" {
					sawActivity = true
				}
			}
		}
		newM, cmd = m.Update(msg)
		m = newM.(tutorModel)
		if _, ok := msg.(turnCompleteMsg); ok {
			break
		}
		if cmd == nil {
			t.Fatal("turn ended without a turnCompleteMsg")
		}
	}

	if !sawActivity {
		t.Error("expected a real activityEventMsg naming read_solution_file during the turn")
	}
	if m.turnInFlight {
		t.Error("turnInFlight = true after the turn completed, want false")
	}
	if len(m.activeCalls) != 0 {
		t.Errorf("activeCalls = %+v, want cleared once the turn completed", m.activeCalls)
	}
	if !strings.Contains(m.viewport.View(), "pong received") {
		t.Errorf("viewport view %q, want the final reply", m.viewport.View())
	}
}

// TestTutorModel_ToolNameLeftBehindInHistoryAfterTurnCompletes is a
// regression test for a real feature request: once a turn ends, the
// live activity region disappears entirely (turnInFlight clears), so
// there was previously no lasting trace in the conversation that a tool
// was ever called -- only the final reply remained visible. The tool
// name must now survive as a permanent part of displayLines, still
// visible after a later, unrelated turn has happened.
func TestTutorModel_ToolNameLeftBehindInHistoryAfterTurnCompletes(t *testing.T) {
	mock := newToolCallOllama(t, "read_solution_file")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly
	cfg.WorkDir = t.TempDir()

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "what does my code look like?")

	view := got.viewport.View()
	if !strings.Contains(view, "read_solution_file") {
		t.Fatalf("viewport view %q, want the tool name left behind after the turn completed", view)
	}
	if !strings.Contains(view, "pong received") {
		t.Fatalf("viewport view %q, want the final reply too", view)
	}

	// A later, unrelated turn (no tool calls) must not push the earlier
	// tool name out of the permanent record.
	mock2 := newSequencedOllama(t, "a plain reply")
	cfg2 := testConfig(mock2.URL)
	cfg2.Mode = exercise.TutorModeSyntaxOnly
	m2, err := newTutorModel(context.Background(), cfg2)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM2, _ := m2.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m2 = newM2.(tutorModel)
	m2.displayBlocks = got.displayBlocks
	m2.refreshViewport()

	got2 := submitAndRun(t, m2, "a follow-up question")
	view2 := got2.viewport.View()
	if !strings.Contains(view2, "read_solution_file") {
		t.Errorf("viewport view %q, want the earlier tool name still present after a later turn", view2)
	}
	if !strings.Contains(view2, "a plain reply") {
		t.Errorf("viewport view %q, want the later turn's own reply too", view2)
	}
}

// --- Stage 4: remaining regression coverage re-homed from tutor_test.go's
// old Run()-level suite -- scenarios not already exercised by the Stage
// 2/3 tests above.

func TestTutorModel_ErrorMessageIncludesRealUnderlyingDetail(t *testing.T) {
	// Ollama's own real error responses are JSON with an "error" field
	// (see eino-contrib/ollama/api's checkError), not a plain-text body --
	// matching that shape here so the client actually decodes the message
	// instead of failing on JSON unmarshal first. Regression test for a
	// real bug found live: a model picked without real tool-calling
	// support made Ollama reject every request with 400 "does not support
	// tools", but a generic "could not reach <host>" message swallowed
	// that detail entirely.
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"does not support tools"}`))
	}))
	defer mock.Close()

	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	// Not stderr -- a real bug found live: writing this detail directly
	// to a raw stderr stream corrupted the terminal, since a real
	// interactive session has stderr and stdout on the very same tty
	// (see activityErrorNote's doc comment in activity.go). It's
	// rendered into the viewport instead now, same safe pipeline as
	// everything else on screen.
	got := submitAndRun(t, m, "hello")

	view := got.viewport.View()
	if !strings.Contains(view, "could not reach") {
		t.Errorf("viewport view %q, want the generic message preserved", view)
	}
	if !strings.Contains(view, "does not support tools") {
		t.Errorf("viewport view %q, want the real underlying error detail included, not just the generic host message", view)
	}
}

func TestTutorModel_ComprehensionCheckErrorMessageIncludesRealUnderlyingDetail(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"does not support tools"}`))
	}))
	defer mock.Close()

	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist // wants the comprehension check

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "hi")

	view := got.viewport.View()
	if !strings.Contains(view, "could not reach") || !strings.Contains(view, "does not support tools") {
		t.Errorf("viewport view %q, want the real underlying error detail included", view)
	}
}

// TestTutorModel_OpenRouterModelShowsOpenRouterInBannerAndErrors is a
// regression test for a real bug found live (via a real OpenRouter
// session): the banner and error messages printed cfg.OllamaHost
// directly, which is meaningless -- empty, in practice -- for an
// OpenRouterModelPrefix-prefixed model.
func TestTutorModel_OpenRouterModelShowsOpenRouterInBannerAndErrors(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":{"message":"rate limited"}}`))
	}))
	defer mock.Close()

	origBaseURL := openRouterBaseURL
	openRouterBaseURL = mock.URL
	defer func() { openRouterBaseURL = origBaseURL }()

	cfg := testConfig("") // OllamaHost deliberately empty/unused for this path
	cfg.Model = OpenRouterModelPrefix + "some/model"
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	if endpoint := m.statusEndpointText(); !strings.Contains(endpoint, "OpenRouter") {
		t.Errorf("status bar endpoint = %q, want it to say \"OpenRouter\"", endpoint)
	}

	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	got := submitAndRun(t, m, "hello")

	view := got.viewport.View()
	if !strings.Contains(view, "could not reach OpenRouter:") {
		t.Errorf("viewport view %q, want \"could not reach OpenRouter:\", not the empty/meaningless OllamaHost", view)
	}
	if !strings.Contains(view, "rate limited") {
		t.Errorf("viewport view %q, want the real underlying error detail included too", view)
	}
}

func TestTutorModel_ComprehensionCheckIncludesUsersRealFirstMessage(t *testing.T) {
	// A real bug found live: an earlier version excluded the user's actual
	// first message from the comprehension-check request, so literally
	// any first message -- including a plain "hi" -- got back the exact
	// same canned restate-and-ask-questions reply with no acknowledgment
	// of what the user said.
	mock := newSequencedOllama(t, "hey! restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	greeting := "hi"
	got := submitAndRun(t, m, greeting)

	reqs := mock.allRequests()
	if len(reqs) == 0 {
		t.Fatal("expected at least 1 request (the comprehension check), got 0")
	}
	found := false
	for _, msg := range reqs[0].Messages {
		if msg.Content == greeting {
			found = true
		}
	}
	if !found {
		t.Errorf("comprehension check request never included the user's real first message %q: %+v", greeting, reqs[0].Messages)
	}
	if !strings.Contains(got.viewport.View(), "hey! restated problem + questions") {
		t.Errorf("viewport view %q, want the comprehension check's reply", got.viewport.View())
	}
}

func TestTutorModel_ComprehensionCheckRetriesWhenReplyLeaksFakeToolCallJSON(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "read_problem_statement", "parameters": {}}`, "clean restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "hi")

	view := got.viewport.View()
	if strings.Contains(view, `{"name"`) {
		t.Errorf("viewport view still contains leaked tool-call JSON: %q", view)
	}
	if !strings.Contains(view, "clean restated problem + questions") {
		t.Errorf("viewport view %q, want the retry's clean reply", view)
	}
	if n := len(mock.allRequests()); n != 2 {
		t.Errorf("requests = %d, want exactly 2 (original comprehension check + one retry)", n)
	}
}

func TestTutorModel_ComprehensionCheckInjectsProblemStatementDirectly(t *testing.T) {
	mock := newSequencedOllama(t, "restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist
	cfg.WorkDir = t.TempDir()

	want := "# Two Sum\n\nReturn indices of the two numbers that add up to target."
	if err := os.WriteFile(filepath.Join(cfg.WorkDir, "problem.md"), []byte(want), 0o644); err != nil {
		t.Fatalf("write problem.md: %v", err)
	}

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	submitAndRun(t, m, "what's the problem?")

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 request (the comprehension check), got %d", len(reqs))
	}
	found := false
	for _, msg := range reqs[0].Messages {
		if strings.Contains(msg.Content, want) {
			found = true
		}
	}
	if !found {
		t.Errorf("comprehension check request never included the injected problem statement %q: %+v", want, reqs[0].Messages)
	}
}

func TestTutorModel_ComprehensionCheckHistoryPersistsBothTurns(t *testing.T) {
	mock := newSequencedOllama(t, "restated + questions", "real answer")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeHintsFirst

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "real question")
	submitAndRun(t, got, "follow up")

	reqs := mock.allRequests()
	// Only 2 requests total: the check consumes the first input line's
	// turn entirely (no separate normal request is also sent for it),
	// then the second input line is one normal turn.
	if len(reqs) != 2 {
		t.Fatalf("expected 2 requests (check, then 1 real turn), got %d", len(reqs))
	}

	// Second request's history should include the check's exchange --
	// persisted using the real first question as the user turn -- followed
	// by the ephemeral hint-count note (hints-first mode, see
	// turnMessages) and the second line. Leading with an extra
	// tools-instruction system message (prependToolsPrompt, prompts.go)
	// ahead of the persona system message that used to be history[0]
	// alone -- the tools instruction is strategy-dependent now, so it's
	// prepended fresh per call instead of baked into history itself.
	second := reqs[1].Messages
	if len(second) != 6 {
		t.Fatalf("second request: expected [tools, persona, user1, assistant1, hint-note, user2] = 6 messages, got %d: %+v", len(second), second)
	}
	if second[2].Content != "real question" {
		t.Errorf("second request messages[2] = %q, want the real first question %q", second[2].Content, "real question")
	}
	if second[3].Content != "restated + questions" {
		t.Errorf("second request messages[3] = %q, want the check's reply", second[3].Content)
	}
	if second[4].Role != "system" || !strings.Contains(second[4].Content, "1st help request") {
		t.Errorf("second request messages[4] = %+v, want an ephemeral system note about the 1st help request", second[4])
	}
	if second[5].Content != "follow up" {
		t.Errorf("second request messages[5] = %q, want %q", second[5].Content, "follow up")
	}
}

func TestTutorModel_FallsBackToHonestMessageWhenRetryAlsoLeaks(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "read_solution_file", "parameters": {}}`, `{"name": "read_cursor_position", "parameters": {}}`)
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	// Wide enough that leakedToolCallFallbackReply renders on one
	// unbroken line -- the viewport word-wraps its content at the real
	// terminal width, so a narrower width would split this assertion's
	// target string across a line break.
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 200, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "where is my cursor?")

	view := got.viewport.View()
	if strings.Contains(view, `{"name"`) {
		t.Errorf("viewport view still contains leaked tool-call JSON: %q", view)
	}
	if !strings.Contains(view, leakedToolCallFallbackReply) {
		t.Errorf("viewport view %q, want the honest fallback message", view)
	}
}

func TestTutorModel_DoesNotRetryWhenReplyIsClean(t *testing.T) {
	mock := newSequencedOllama(t, "the answer is 42")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	submitAndRun(t, m, "question")

	if n := len(mock.allRequests()); n != 1 {
		t.Errorf("requests = %d, want exactly 1 (no retry for a clean reply)", n)
	}
}

func TestTutorModel_LeakedReplyNeverPersistedToHistory(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "read_solution_file", "parameters": {}}`, "your code looks fine", "second reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "what does my code look like?")
	submitAndRun(t, got, "another question")

	reqs := mock.allRequests()
	if len(reqs) != 3 {
		t.Fatalf("requests = %d, want exactly 3 (leaked original + retry + second turn)", len(reqs))
	}
	// The second turn's request carries history from the first turn --
	// confirm the leaked (never-shown) reply isn't in it, only the clean
	// retry reply.
	for _, msg := range reqs[2].Messages {
		if strings.Contains(msg.Content, `{"name"`) {
			t.Errorf("second turn's request carries leaked JSON in history: %+v", reqs[2].Messages)
		}
	}
}

func TestTutorModel_RetainsConversationHistoryAcrossTurns(t *testing.T) {
	mock := newSequencedOllama(t, "assistant-reply-1", "assistant-reply-2")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "first line")
	submitAndRun(t, got, "second line")

	reqs := mock.allRequests()
	if len(reqs) != 2 {
		t.Fatalf("expected 2 requests (one per input line), got %d", len(reqs))
	}
	// [tools, persona, user1] = 3: an extra leading tools-instruction
	// system message (prependToolsPrompt) ahead of the persona system
	// message that used to be the sole history[0] -- see
	// TestTutorModel_ComprehensionCheckHistoryPersistsBothTurns' own
	// comment for why.
	if len(reqs[0].Messages) != 3 {
		t.Errorf("first request: expected [tools, persona, user1] = 3 messages, got %d: %+v", len(reqs[0].Messages), reqs[0].Messages)
	}
	second := reqs[1].Messages
	if len(second) != 5 {
		t.Fatalf("second request: expected [tools, persona, user1, assistant1, user2] = 5 messages, got %d: %+v", len(second), second)
	}
	if second[2].Content != "first line" {
		t.Errorf("second request messages[2] (user1) = %q, want %q", second[2].Content, "first line")
	}
	if second[3].Role != "assistant" || second[3].Content != "assistant-reply-1" {
		t.Errorf("second request messages[3] (assistant1) = %+v, want role=assistant content=%q", second[3], "assistant-reply-1")
	}
	if second[4].Content != "second line" {
		t.Errorf("second request messages[4] (user2) = %q, want %q", second[4].Content, "second line")
	}
}

// --- Stage 5: per-turn cancellation (Esc/Ctrl-C) and a bounded per-turn
// timeout, issue #239 -- before this, the key handler covered exactly
// Ctrl-D (quit)/Enter (submit)/PgUp/PgDn (scroll); Ctrl-C reached the
// textarea unhandled, a turn had no bound of its own (only
// ollamaRequestTimeout's PER-REQUEST 120s, and a react.Agent turn can
// make several requests), and a cancelled or failed turn discarded the
// user's message outright.

func TestTutorModel_CtrlCWhileTurnInFlightCancelsAndDoesNotQuit(t *testing.T) {
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m = detectAndApply(t, m)

	m.textarea.SetValue("a question")
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)
	if !m.turnInFlight {
		t.Fatal("setup: expected turnInFlight after submit")
	}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	got := newM.(tutorModel)

	if got.turnInFlight {
		t.Error("turnInFlight = true after Ctrl-C mid-turn, want false -- Ctrl-C should cancel the turn")
	}
	if cmd != nil {
		if _, isQuit := cmd().(tea.QuitMsg); isQuit {
			t.Error("Ctrl-C mid-turn produced tea.Quit, want the turn cancelled instead of the program exiting")
		}
	}
}

func TestTutorModel_CtrlCWhileIdleQuits(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("Ctrl-C while idle produced no command, want tea.Quit")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Error("Ctrl-C while idle did not quit -- matching every other program in this codebase")
	}
}

func TestTutorModel_EscWhileTurnInFlightCancels(t *testing.T) {
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m = detectAndApply(t, m)

	m.textarea.SetValue("a question")
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)
	if !m.turnInFlight {
		t.Fatal("setup: expected turnInFlight after submit")
	}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	got := newM.(tutorModel)

	if got.turnInFlight {
		t.Error("turnInFlight = true after Esc mid-turn, want false")
	}
	if cmd != nil {
		t.Error("expected Esc cancellation to return no further command -- it must not wait on the turn's own goroutine")
	}
	if !strings.Contains(ansi.Strip(got.viewport.View()), "cancelled") {
		t.Errorf("viewport view %q, want a cancellation note", ansi.Strip(got.viewport.View()))
	}
}

func TestTutorModel_CancelledTurnRestoresTextIntoEmptyTextarea(t *testing.T) {
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m = detectAndApply(t, m)

	m.textarea.SetValue("please explain recursion")
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)
	if m.textarea.Value() != "" {
		t.Fatal("setup: expected submit to clear the textarea")
	}

	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	got := newM.(tutorModel)

	if got.textarea.Value() != "please explain recursion" {
		t.Errorf("textarea.Value() = %q after cancel, want the original message restored verbatim", got.textarea.Value())
	}
}

// TestTutorModel_CancelledTurnWithTypedAheadTextSurfacesOldTextInsteadOfClobbering
// covers the "they may have typed ahead" case the issue calls out
// explicitly: a cancel must never overwrite a fresh, unsent draft the
// user has already started typing while the old turn was in flight.
func TestTutorModel_CancelledTurnWithTypedAheadTextSurfacesOldTextInsteadOfClobbering(t *testing.T) {
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m = detectAndApply(t, m)

	m.textarea.SetValue("please explain recursion")
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)

	// Typed ahead while the turn was in flight -- a fresh, unsent draft
	// that a cancel must not clobber.
	m.textarea.SetValue("actually never mind, something else")

	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	got := newM.(tutorModel)

	if got.textarea.Value() != "actually never mind, something else" {
		t.Errorf("textarea.Value() = %q after cancel, want the typed-ahead draft left untouched", got.textarea.Value())
	}
	if !strings.Contains(ansi.Strip(got.viewport.View()), "please explain recursion") {
		t.Errorf("viewport view %q, want the cancelled message surfaced in the note since it couldn't go back in the textarea", ansi.Strip(got.viewport.View()))
	}
}

func TestTutorModel_FailedTurnRestoresText(t *testing.T) {
	cfg := testConfig("http://127.0.0.1:1") // refuses immediately
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "please explain recursion")

	if got.textarea.Value() != "please explain recursion" {
		t.Errorf("textarea.Value() = %q after a failed turn, want the original message restored instead of discarded", got.textarea.Value())
	}
}

// TestTutorModel_FailedTurnWithTypedAheadTextSurfacesOldText is the
// failure-path counterpart of the cancel-path test above -- the same
// "never clobber a fresh draft" contract applies regardless of why the
// turn didn't complete.
func TestTutorModel_FailedTurnWithTypedAheadTextSurfacesOldText(t *testing.T) {
	cfg := testConfig("http://127.0.0.1:1") // refuses immediately
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m = detectAndApply(t, m)

	m.textarea.SetValue("please explain recursion")
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(tutorModel)
	if cmd == nil {
		t.Fatal("submit produced no command")
	}
	m.textarea.SetValue("a fresh draft typed while waiting")

	for i := 0; i < 100; i++ {
		msg := cmd()
		newM, cmd = m.Update(msg)
		m = newM.(tutorModel)
		if _, ok := msg.(turnCompleteMsg); ok {
			break
		}
		if cmd == nil {
			t.Fatal("turn ended without a turnCompleteMsg")
		}
	}

	if m.textarea.Value() != "a fresh draft typed while waiting" {
		t.Errorf("textarea.Value() = %q, want the fresh draft left untouched", m.textarea.Value())
	}
	if !strings.Contains(ansi.Strip(m.viewport.View()), "please explain recursion") {
		t.Errorf("viewport view %q, want the failed message surfaced in the note", ansi.Strip(m.viewport.View()))
	}
}

// TestTutorModel_PerTurnTimeoutProducesTimeoutFlavoredNote proves
// turnTimeout (not just ollamaRequestTimeout's per-REQUEST bound) is
// actually wired to the whole turn, and that hitting it reads
// differently from a plain connectivity failure -- a timeout genuinely
// reached the model, it just never got a reply back in time.
func TestTutorModel_PerTurnTimeoutProducesTimeoutFlavoredNote(t *testing.T) {
	hang := make(chan struct{})
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-hang // never respond until the test cleans up
	}))
	t.Cleanup(func() {
		close(hang)
		mock.Close()
	})

	origTimeout := turnTimeout
	turnTimeout = 100 * time.Millisecond
	t.Cleanup(func() { turnTimeout = origTimeout })

	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "a question that will hang forever")

	view := ansi.Strip(got.viewport.View())
	if !strings.Contains(view, "timed out") {
		t.Errorf("viewport view %q, want a timeout-flavored note", view)
	}
	if strings.Contains(view, "could not reach") {
		t.Errorf("viewport view %q, want the timeout note, not the generic connectivity wording", view)
	}
}

// TestTutorModel_StaleTurnCompleteMsgAfterCancelIsIgnored locks in
// turnSeq's whole reason for existing: bubbletea's runtime starts a
// turn's waitForActivityEvent goroutine the moment submit returns it,
// independent of anything Update does afterward -- cancelling can't
// un-start it. That goroutine WILL eventually deliver the cancelled
// turn's own result (ctx cancellation aborts the request, but not
// instantly), tagged with the seq it was started under. Without the
// seq guard, that late delivery would silently flip turnInFlight back
// on (or off, stomping a NEW turn already running by then) and leak a
// stray note into the transcript.
func TestTutorModel_StaleTurnCompleteMsgAfterCancelIsIgnored(t *testing.T) {
	m := newTutorLayoutOnly()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	m.turnInFlight = true
	m.turnSeq = 5 // this generation is the one currently tracked

	// A result tagged with an earlier (already-retired, e.g. by a
	// cancel) generation must not touch state that belongs to 5.
	newM, cmd := m.Update(turnCompleteMsg{seq: 4, err: fmt.Errorf("stale failure"), userMessage: "old message"})
	got := newM.(tutorModel)

	if !got.turnInFlight {
		t.Error("a stale turnCompleteMsg (old seq) changed turnInFlight, want the current turn's state left untouched")
	}
	if cmd != nil {
		t.Error("a stale turnCompleteMsg re-armed a wait, want it dropped outright")
	}
	if strings.Contains(ansi.Strip(got.viewport.View()), "stale failure") {
		t.Error("a stale turnCompleteMsg's error note leaked into the transcript")
	}
	if len(got.displayBlocks) != 0 {
		t.Errorf("displayBlocks = %+v, want untouched by a stale message", got.displayBlocks)
	}
}

// --- Stage 6: bounded conversation history and distinguishable turn
// failures, issue #240 -- before this, m.history grew unbounded and was
// resent in full every turn, so a long session got slower and slower
// before eventually failing outright with a provider context-length
// error that surfaced as the exact same generic "could not reach
// <model>" wording as a real network problem. trimHistory/
// classifyTurnError themselves are unit-tested directly in
// history_test.go; these prove the wiring into the real turn path.

func TestTutorModel_HistoryIsTrimmedToBudgetAfterEachTurn(t *testing.T) {
	origBudget := historyBudgetChars
	historyBudgetChars = 60 // tiny -- room for roughly one short pair
	t.Cleanup(func() { historyBudgetChars = origBudget })

	mock := newSequencedOllama(t, "short reply one", "short reply two", "short reply three")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "first question")
	got = submitAndRun(t, got, "second question")
	got = submitAndRun(t, got, "third question")

	// Untrimmed, 3 turns would leave 1 (system) + 2*3 (pairs) = 7
	// messages -- a tiny budget must keep it well below that.
	if untrimmed := 1 + 2*3; len(got.history) >= untrimmed {
		t.Errorf("history has %d messages after 3 turns under a tiny budget, want it trimmed below the untrimmed %d", len(got.history), untrimmed)
	}
	if got.history[0].Role != schema.System {
		t.Errorf("history[0].Role = %v, want System -- the persona prompt must survive trimming", got.history[0].Role)
	}
	last := got.history[len(got.history)-1]
	if last.Content != "short reply three" {
		t.Errorf("history's last message = %q, want the most recent reply always kept regardless of budget", last.Content)
	}
}

// TestTutorModel_ContextOverflowFailureShowsDistinctNoteFromNetworkFailure
// drives a real (mocked) provider context-length rejection through the
// actual Update path, not just classifyTurnError in isolation -- proving
// the wiring, not just the classifier.
func TestTutorModel_ContextOverflowFailureShowsDistinctNoteFromNetworkFailure(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"this model's maximum context length is 4096 tokens"}`))
	}))
	defer mock.Close()

	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	got := submitAndRun(t, m, "hello")

	view := ansi.Strip(got.viewport.View())
	if !strings.Contains(view, "context window") {
		t.Errorf("viewport view %q, want a context-overflow note", view)
	}
	if strings.Contains(view, "could not reach") {
		t.Errorf("viewport view %q, want the overflow note distinct from the generic connectivity wording", view)
	}
}
