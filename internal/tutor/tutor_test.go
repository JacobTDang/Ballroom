package tutor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// tutorChatRequest/sequencedOllama are the package's Ollama mock server
// for Go-native tests — used here and by agent_test.go. Originally kept
// independent of chat_test.go's now-deleted mockOllama (which tested
// tutor/chat.sh, the bash implementation this package replaced) so this
// file wouldn't need to change when that one went away.
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

func TestRun_PrintsAssistantReply(t *testing.T) {
	mock := newSequencedOllama(t, "the answer is 42")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check for a single-turn test

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("question\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !strings.Contains(stdout.String(), "the answer is 42") {
		t.Errorf("stdout = %q, want it to contain the assistant reply", stdout.String())
	}
}

func TestRun_RetainsConversationHistory(t *testing.T) {
	mock := newSequencedOllama(t, "assistant-reply-1", "assistant-reply-2")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly // skip the comprehension check to isolate history behavior

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("first line\nsecond line\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

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

func TestRun_HandlesUnreachableHostGracefully(t *testing.T) {
	cfg := testConfig("http://127.0.0.1:1") // port 1 is reserved/unlisted, refuses immediately
	cfg.Mode = exercise.TutorModeSyntaxOnly

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("hello\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run should exit cleanly even when Ollama is unreachable, got error: %v", err)
	}
	if !strings.Contains(stderr.String(), "could not reach") {
		t.Errorf("stderr = %q, want a message about being unable to reach the host", stderr.String())
	}
}

func TestRun_ComprehensionCheckIsolatesRealQuestion(t *testing.T) {
	mock := newSequencedOllama(t, "restated problem + questions", "real answer")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist // wants the comprehension check

	secretQuestion := "what is the secret algorithm here"
	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader(secretQuestion+"\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) == 0 {
		t.Fatal("expected at least 1 request (the comprehension check), got 0")
	}

	checkReq := reqs[0]
	for _, m := range checkReq.Messages {
		if strings.Contains(m.Content, secretQuestion) {
			t.Errorf("comprehension check request contained the user's real question %q — isolation broken: %+v", secretQuestion, checkReq.Messages)
		}
	}

	if !strings.Contains(stdout.String(), "restated problem + questions") {
		t.Errorf("stdout = %q, want the comprehension check's reply printed", stdout.String())
	}
}

func TestRun_ComprehensionCheckInjectsProblemStatementDirectly(t *testing.T) {
	mock := newSequencedOllama(t, "restated problem + questions")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeFullAssist
	cfg.WorkDir = t.TempDir()

	want := "# Two Sum\n\nReturn indices of the two numbers that add up to target."
	if err := os.WriteFile(filepath.Join(cfg.WorkDir, "problem.md"), []byte(want), 0o644); err != nil {
		t.Fatalf("write problem.md: %v", err)
	}

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("what's the problem?\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 request (the comprehension check), got %d", len(reqs))
	}

	found := false
	for _, m := range reqs[0].Messages {
		if strings.Contains(m.Content, want) {
			found = true
		}
	}
	if !found {
		t.Errorf("comprehension check request never included the injected problem statement %q: %+v", want, reqs[0].Messages)
	}
}

func TestRun_ComprehensionCheckSkippedForSyntaxOnly(t *testing.T) {
	mock := newSequencedOllama(t, "direct reply")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeSyntaxOnly

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("what's wrong with my code\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected exactly 1 request (no comprehension check for syntax-only), got %d", len(reqs))
	}
	if !strings.Contains(stdout.String(), "direct reply") {
		t.Errorf("stdout = %q, want the direct reply printed with no check", stdout.String())
	}
}

func TestRun_ComprehensionCheckHistoryPersistsBothTurns(t *testing.T) {
	mock := newSequencedOllama(t, "restated + questions", "real answer")
	cfg := testConfig(mock.URL)
	cfg.Mode = exercise.TutorModeHintsFirst

	var stdout, stderr strings.Builder
	if err := Run(context.Background(), cfg, strings.NewReader("real question\nfollow up\n"), &stdout, &stderr); err != nil {
		t.Fatalf("Run: %v", err)
	}

	reqs := mock.allRequests()
	// Only 2 requests total: the check consumes the first input line's
	// turn entirely (no separate normal request is also sent for it),
	// then the second input line is one normal turn.
	if len(reqs) != 2 {
		t.Fatalf("expected 2 requests (check, then 1 real turn), got %d", len(reqs))
	}

	// Second request's history should include the check's exchange —
	// persisted using the real first question as the user turn, per
	// runComprehensionCheck's doc comment — followed by the ephemeral
	// hint-count note (hints-first mode, see turnMessages) and the
	// second line.
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
