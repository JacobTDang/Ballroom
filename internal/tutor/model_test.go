package tutor

import (
	"context"
	"io"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

	total := got.viewport.Height + got.textarea.Height() + textareaBorderRows
	if total > 30 {
		t.Errorf("viewport(%d) + textarea(%d) + border(%d) = %d, want <= terminal height 30", got.viewport.Height, got.textarea.Height(), textareaBorderRows, total)
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
	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

// --- Stage 2: real turn-loop logic, re-homed from the old Run()-level
// tests in tutor_test.go to test tutorModel.Update() directly -- same
// proven pattern internal/tui/app_test.go already uses for appModel.
// submitAndRun types line into m's textarea, presses Enter, and
// synchronously executes the returned tea.Cmd (a plain func() tea.Msg --
// no real tea.Program needed) feeding its result back into Update(),
// mirroring internal/tui's own async-message test pattern.

func submitAndRun(t *testing.T, m tutorModel, line string) tutorModel {
	t.Helper()
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

func TestNewTutorModel_NoRoutingBannerAndHistorySeededWithSystemPrompt(t *testing.T) {
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	if len(m.history) != 1 || m.history[0].Role != schema.System {
		t.Fatalf("history = %+v, want exactly one System message seeded", m.history)
	}
	if !strings.Contains(m.banner, cfg.Model) || strings.Contains(m.banner, "orchestrator=") {
		t.Errorf("banner = %q, want it to name the model and not mention routing", m.banner)
	}
}

func TestNewTutorModel_RoutingEnabledBannerNamesBothModels(t *testing.T) {
	mock := newSequencedOllama(t, "reply")
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	if !strings.Contains(m.banner, "worker-model") || !strings.Contains(m.banner, "orchestrator-model") {
		t.Errorf("banner = %q, want it to name both models", m.banner)
	}
}

func TestNewTutorModel_ComprehensionCheckPendingMatchesMode(t *testing.T) {
	mock := newSequencedOllama(t, "reply")

	syntaxCfg := testConfig(mock.URL)
	syntaxCfg.Mode = exercise.TutorModeSyntaxOnly
	m, err := newTutorModel(context.Background(), syntaxCfg, io.Discard)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	if m.comprehensionCheckPending {
		t.Error("syntax-only mode must never want the comprehension check")
	}

	fullCfg := testConfig(mock.URL)
	fullCfg.Mode = exercise.TutorModeFullAssist
	m, err = newTutorModel(context.Background(), fullCfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

func TestTutorModel_ComprehensionCheckAlwaysUsesOrchestratorWhenRoutingEnabled(t *testing.T) {
	mock := newSequencedOllama(t, "restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Model = "worker-model"
	cfg.OrchestratorModel = "orchestrator-model"
	cfg.Mode = exercise.TutorModeFullAssist

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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
	if !strings.Contains(m.View(), "read_solution_file") {
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

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
