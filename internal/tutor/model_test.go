package tutor

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

	var stderr strings.Builder
	m, err := newTutorModel(context.Background(), cfg, &stderr)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	submitAndRun(t, m, "hello")

	got := stderr.String()
	if !strings.Contains(got, "could not reach") {
		t.Errorf("stderr = %q, want the generic message preserved", got)
	}
	if !strings.Contains(got, "does not support tools") {
		t.Errorf("stderr = %q, want the real underlying error detail included, not just the generic host message", got)
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

	var stderr strings.Builder
	m, err := newTutorModel(context.Background(), cfg, &stderr)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)

	submitAndRun(t, m, "hi")

	got := stderr.String()
	if !strings.Contains(got, "could not reach") || !strings.Contains(got, "does not support tools") {
		t.Errorf("stderr = %q, want the real underlying error detail included", got)
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

	var stderr strings.Builder
	m, err := newTutorModel(context.Background(), cfg, &stderr)
	if err != nil {
		t.Fatalf("newTutorModel: %v", err)
	}
	if !strings.Contains(m.banner, "connected to OpenRouter") {
		t.Errorf("banner = %q, want it to say \"connected to OpenRouter\"", m.banner)
	}

	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = newM.(tutorModel)
	submitAndRun(t, m, "hello")

	errOut := stderr.String()
	if !strings.Contains(errOut, "could not reach OpenRouter:") {
		t.Errorf("stderr = %q, want \"could not reach OpenRouter:\", not the empty/meaningless OllamaHost", errOut)
	}
	if !strings.Contains(errOut, "rate limited") {
		t.Errorf("stderr = %q, want the real underlying error detail included too", errOut)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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
	// turnMessages) and the second line.
	second := reqs[1].Messages
	if len(second) != 5 {
		t.Fatalf("second request: expected [system, user1, assistant1, hint-note, user2] = 5 messages, got %d: %+v", len(second), second)
	}
	if second[1].Content != "real question" {
		t.Errorf("second request messages[1] = %q, want the real first question %q", second[1].Content, "real question")
	}
	if second[2].Content != "restated + questions" {
		t.Errorf("second request messages[2] = %q, want the check's reply", second[2].Content)
	}
	if second[3].Role != "system" || !strings.Contains(second[3].Content, "1st help request") {
		t.Errorf("second request messages[3] = %+v, want an ephemeral system note about the 1st help request", second[3])
	}
	if second[4].Content != "follow up" {
		t.Errorf("second request messages[4] = %q, want %q", second[4].Content, "follow up")
	}
}

func TestTutorModel_FallsBackToHonestMessageWhenRetryAlsoLeaks(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "read_solution_file", "parameters": {}}`, `{"name": "read_cursor_position", "parameters": {}}`)
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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

	m, err := newTutorModel(context.Background(), cfg, io.Discard)
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
	if len(reqs[0].Messages) != 2 {
		t.Errorf("first request: expected [system, user1] = 2 messages, got %d: %+v", len(reqs[0].Messages), reqs[0].Messages)
	}
	second := reqs[1].Messages
	if len(second) != 4 {
		t.Fatalf("second request: expected [system, user1, assistant1, user2] = 4 messages, got %d: %+v", len(second), second)
	}
	if second[1].Content != "first line" {
		t.Errorf("second request messages[1] (user1) = %q, want %q", second[1].Content, "first line")
	}
	if second[2].Role != "assistant" || second[2].Content != "assistant-reply-1" {
		t.Errorf("second request messages[2] (assistant1) = %+v, want role=assistant content=%q", second[2], "assistant-reply-1")
	}
	if second[3].Content != "second line" {
		t.Errorf("second request messages[3] (user2) = %q, want %q", second[3].Content, "second line")
	}
}
