package tutor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// ansiPattern/stripAnsi strip terminal escape sequences (cursor
// movement, color) out of captured stdout so tests can assert on the
// visible text content — shared by discoball_test.go, thinkingdisplay_test.go,
// and this file's own Run-level smoke test.
var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*[A-Za-z]`)

func stripAnsi(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

// tutorChatRequest/sequencedOllama are the package's Ollama mock server
// for Go-native tests — used here and by model_test.go/agent_test.go.
// Originally kept independent of chat_test.go's now-deleted mockOllama
// (which tested tutor/chat.sh, the bash implementation this package
// replaced) so this file wouldn't need to change when that one went away.
type tutorChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type tutorChatRequest struct {
	Model    string             `json:"model"`
	Messages []tutorChatMessage `json:"messages"`
	Stream   bool               `json:"stream"`
}

// sequencedOllama serves one reply per request in order, repeating the
// last reply if more requests arrive than replies were given — lets a
// test simulate a multi-turn conversation (e.g. the comprehension
// check's reply, then a real answer) without needing a real model.
type sequencedOllama struct {
	*httptest.Server
	replies []string

	mu       sync.Mutex
	requests []tutorChatRequest
}

func newSequencedOllama(t *testing.T, replies ...string) *sequencedOllama {
	t.Helper()
	m := &sequencedOllama{replies: replies}
	m.Server = httptest.NewServer(http.HandlerFunc(m.handle))
	t.Cleanup(m.Server.Close)
	return m
}

func (m *sequencedOllama) handle(w http.ResponseWriter, r *http.Request) {
	var req tutorChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	idx := len(m.requests)
	m.requests = append(m.requests, req)
	m.mu.Unlock()

	reply := m.replies[len(m.replies)-1]
	if idx < len(m.replies) {
		reply = m.replies[idx]
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": map[string]string{"role": "assistant", "content": reply},
		"done":    true,
	})
}

func (m *sequencedOllama) allRequests() []tutorChatRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]tutorChatRequest, len(m.requests))
	copy(out, m.requests)
	return out
}

func testConfig(ollamaHost string) Config {
	return Config{
		OllamaHost:      ollamaHost,
		Model:           "test-model",
		Mode:            exercise.TutorModeFullAssist,
		MaxContextBytes: 8000,
	}
}

// TestRun_ExitsCleanlyOnCtrlDWithNoInput is Run()'s one true
// program-level smoke test: it proves the thin wrapper actually builds a
// real tutorModel and drives a real tea.Program against a fake
// io.Reader/io.Writer end to end. It deliberately submits no message --
// a scripted input stream that both submits a message AND appends Ctrl-D
// races the turn's own async completion against the queued Ctrl-D (see
// RunOneTurn's doc comment), so real turn-loop behavior is covered by
// model_test.go's direct tutorModel.Update() tests instead, which have
// no such race.
//
// Asserts on the textarea's placeholder, not the viewport's banner text
// -- a real tty triggers bubbletea's own initial tea.WindowSizeMsg query,
// but a fake, non-tty io.Writer (as used here and by every other test in
// this package) never does, so the viewport (whose View() renders
// exactly Height rows) stays at its zero-value height and paints
// nothing; the textarea renders unconditionally regardless of that.
func TestRun_ExitsCleanlyOnCtrlDWithNoInput(t *testing.T) {
	mock := newSequencedOllama(t, "unused")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	var stdout strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("\x04"), &stdout); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if got := stripAnsi(stdout.String()); !strings.Contains(got, "Ask a question") {
		t.Errorf("stdout = %q, want the textarea's placeholder, proving a real frame was rendered", got)
	}
}

func TestLooksLikeLeakedToolCall_DetectsRawJSONShape(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    bool
	}{
		{"bare JSON", `{"name": "read_solution_file", "parameters": {}}`, true},
		{"prose then JSON", "To answer this, I need to check your code.\n\n" + `{"name": "read_cursor_position", "parameters": {}}`, true},
		{"hallucinated tool name still matches the shape", `{"name": "read_user_code", "parameters": {}}`, true},
		{"clean reply", "your code looks fine to me", false},
		{"empty", "", false},
	}
	for _, c := range cases {
		if got := looksLikeLeakedToolCall(c.content); got != c.want {
			t.Errorf("%s: looksLikeLeakedToolCall(%q) = %v, want %v", c.name, c.content, got, c.want)
		}
	}
}

func TestTurnMessages_NonHintsFirstModeHasNoNote(t *testing.T) {
	for _, mode := range []string{exercise.TutorModeSyntaxOnly, exercise.TutorModeFullAssist} {
		msgs := turnMessages(mode, 1, "hello")
		if len(msgs) != 1 {
			t.Errorf("mode %s: turnMessages returned %d messages, want 1 (just the user message)", mode, len(msgs))
			continue
		}
		if msgs[0].Content != "hello" {
			t.Errorf("mode %s: messages[0].Content = %q, want %q", mode, msgs[0].Content, "hello")
		}
	}
}

func TestTurnMessages_HintsFirstFirstRequestNotesFirstAsk(t *testing.T) {
	msgs := turnMessages(exercise.TutorModeHintsFirst, 1, "help")
	if len(msgs) != 2 {
		t.Fatalf("turnMessages returned %d messages, want 2 (note + user message)", len(msgs))
	}
	if msgs[0].Role != "system" || !strings.Contains(msgs[0].Content, "1st help request") {
		t.Errorf("messages[0] = %+v, want a system note mentioning the 1st help request", msgs[0])
	}
	if msgs[1].Content != "help" {
		t.Errorf("messages[1].Content = %q, want %q", msgs[1].Content, "help")
	}
}

func TestTurnMessages_HintsFirstRepeatRequestNotesCount(t *testing.T) {
	msgs := turnMessages(exercise.TutorModeHintsFirst, 3, "still stuck")
	if len(msgs) != 2 {
		t.Fatalf("turnMessages returned %d messages, want 2 (note + user message)", len(msgs))
	}
	if msgs[0].Role != "system" || !strings.Contains(msgs[0].Content, "#3") {
		t.Errorf("messages[0] = %+v, want a system note mentioning request #3", msgs[0])
	}
	if !strings.Contains(msgs[0].Content, "don't ask them to confirm") {
		t.Errorf("messages[0].Content = %q, want it to explicitly tell the model not to ask the user to confirm", msgs[0].Content)
	}
}

// --- decideHandoff (orchestrator/worker routing) ---

func TestDecideHandoff_YesRepliesMeanHandoff(t *testing.T) {
	mock := newSequencedOllama(t, "YES")
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", mock.URL, "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "how do I solve this with a hash map")
	if err != nil {
		t.Fatalf("decideHandoff: %v", err)
	}
	if !handoff {
		t.Error("handoff = false, want true for a YES reply")
	}
}

func TestDecideHandoff_NoRepliesMeanNoHandoff(t *testing.T) {
	mock := newSequencedOllama(t, "no")
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", mock.URL, "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "hi")
	if err != nil {
		t.Fatalf("decideHandoff: %v", err)
	}
	if handoff {
		t.Error("handoff = true, want false for a (lowercase) no reply")
	}
}

func TestDecideHandoff_NoWithTrailingTextStillMeansNoHandoff(t *testing.T) {
	mock := newSequencedOllama(t, "No, this is just a greeting.")
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", mock.URL, "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "hi there")
	if err != nil {
		t.Fatalf("decideHandoff: %v", err)
	}
	if handoff {
		t.Error("handoff = true, want false -- reply starts with No even with trailing explanation")
	}
}

func TestDecideHandoff_AmbiguousReplyDefaultsToHandoff(t *testing.T) {
	mock := newSequencedOllama(t, "I'm not sure, could go either way")
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", mock.URL, "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "something ambiguous")
	if err != nil {
		t.Fatalf("decideHandoff: %v", err)
	}
	if !handoff {
		t.Error("handoff = false, want true -- anything not clearly starting with No defaults to handoff (safer to over-use the specialist)")
	}
}

func TestDecideHandoff_RequestErrorDefaultsToHandoffAndReturnsError(t *testing.T) {
	ctx := context.Background()
	cm, err := newChatModel(ctx, "test-model", "http://127.0.0.1:1", "")
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	handoff, err := decideHandoff(ctx, cm, "hi")
	if err == nil {
		t.Fatal("expected an error for an unreachable host")
	}
	if !handoff {
		t.Error("handoff = false, want true -- a routing failure must default to handoff, not silently leave the turn unanswered")
	}
}

func TestLooksLikeLeakedToolCall_DetectsToolCallTag(t *testing.T) {
	// The XML-style tag dialect (observed live from
	// poolside/laguna-xs-2.1) must never reach the user from the native
	// path's final answers either.
	if !looksLikeLeakedToolCall("Let me look.<tool_call>read_solution_file</tool_call>") {
		t.Error("looksLikeLeakedToolCall missed a <tool_call> tag")
	}
	if looksLikeLeakedToolCall("use a set instead of a list here") {
		t.Error("looksLikeLeakedToolCall false positive on plain prose")
	}
}

func TestInterviewClockNote_InterviewerGetsElapsedOverLimit(t *testing.T) {
	started := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	now := started.Add(23 * time.Minute)
	note := interviewClockNote(exercise.TutorModeInterviewer, started, 45, now)
	if note == nil {
		t.Fatal("interviewClockNote = nil for an interviewer session with a live clock")
	}
	if !strings.Contains(note.Content, "23 of 45 minutes") {
		t.Errorf("note = %q, want elapsed-over-limit minutes", note.Content)
	}
}

func TestInterviewClockNote_TimeUpMessageWhenPastLimit(t *testing.T) {
	started := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	now := started.Add(50 * time.Minute)
	note := interviewClockNote(exercise.TutorModeInterviewer, started, 45, now)
	if note == nil {
		t.Fatal("interviewClockNote = nil past the limit, want a wrap-up note")
	}
	if !strings.Contains(note.Content, "45 of 45 minutes") || !strings.Contains(note.Content, "wrap up") {
		t.Errorf("note = %q, want clamped time and a wrap-up push", note.Content)
	}
}

func TestInterviewClockNote_OtherModesAndUnknownClockGetNothing(t *testing.T) {
	started := time.Date(2026, 7, 15, 10, 0, 0, 0, time.UTC)
	now := started.Add(10 * time.Minute)
	if note := interviewClockNote(exercise.TutorModeDesignCoach, started, 90, now); note != nil {
		t.Errorf("design-coach got a clock note %q, want nil -- coaching is untimed", note.Content)
	}
	if note := interviewClockNote(exercise.TutorModeHintsFirst, started, 25, now); note != nil {
		t.Errorf("hints-first got a clock note %q, want nil", note.Content)
	}
	if note := interviewClockNote(exercise.TutorModeInterviewer, time.Time{}, 45, now); note != nil {
		t.Errorf("zero StartedAt got a clock note %q, want nil (sandbox/tests have no clock)", note.Content)
	}
	if note := interviewClockNote(exercise.TutorModeInterviewer, started, 0, now); note != nil {
		t.Errorf("zero TimeLimitMin got a clock note %q, want nil", note.Content)
	}
}
