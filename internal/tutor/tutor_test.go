package tutor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// tutorChatRequest/sequencedOllama are independent of chat_test.go's
// mockOllama (which tests tutor/chat.sh and is deleted once the bash
// implementation is cut over) — self-contained on purpose so this file
// doesn't need to change when that one goes away.
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
	// runComprehensionCheck's doc comment — followed by the second line.
	second := reqs[1].Messages
	if len(second) != 4 {
		t.Fatalf("second request: expected [system, user1, assistant1, user2] = 4 messages, got %d: %+v", len(second), second)
	}
	if second[1].Content != "real question" {
		t.Errorf("second request messages[1] = %q, want the real first question %q", second[1].Content, "real question")
	}
	if second[2].Content != "restated + questions" {
		t.Errorf("second request messages[2] = %q, want the check's reply", second[2].Content)
	}
	if second[3].Content != "follow up" {
		t.Errorf("second request messages[3] = %q, want %q", second[3].Content, "follow up")
	}
}
