// Package tutor tests tutor/chat.sh against a mock Ollama server. It
// deliberately never talks to a real model — CI has no Ollama running, and
// asserting on model output quality would be flaky and non-deterministic
// anyway. What's testable and worth gating on is our own plumbing: does the
// right system prompt go out for each tutor_mode, does conversation history
// actually accumulate (hints-first's "escalate on second ask" depends on
// it), and does a network failure get handled instead of crashing.
package tutor

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type mockOllama struct {
	*httptest.Server
	reply string

	mu       sync.Mutex
	requests []chatRequest
}

func newMockOllama(t *testing.T, reply string) *mockOllama {
	t.Helper()
	m := &mockOllama{reply: reply}
	m.Server = httptest.NewServer(http.HandlerFunc(m.handle))
	t.Cleanup(m.Server.Close)
	return m
}

func (m *mockOllama) handle(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	m.requests = append(m.requests, req)
	m.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": map[string]string{"role": "assistant", "content": m.reply},
		"done":    true,
	})
}

func (m *mockOllama) allRequests() []chatRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]chatRequest, len(m.requests))
	copy(out, m.requests)
	return out
}

func chatScriptPath(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	abs, err := filepath.Abs(filepath.Join(filepath.Dir(thisFile), "..", "..", "tutor", "chat.sh"))
	if err != nil {
		t.Fatalf("resolve chat.sh path: %v", err)
	}
	return abs
}

// runChat runs the real tutor/chat.sh against mock (or ollamaHost if set,
// for the unreachable-host case), feeding stdin and returning captured
// stdout/stderr. Does not fail the test on a non-zero exit — chat.sh is
// expected to exit cleanly even when a request fails, so callers assert on
// that themselves.
func runChat(t *testing.T, ollamaHost, mode, stdin string, extraEnv map[string]string) (stdout, stderr string, err error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", chatScriptPath(t))
	env := []string{
		"PATH=" + os.Getenv("PATH"),
		"OLLAMA_HOST=" + ollamaHost,
		"TUTOR_MODEL=test-model",
		"PRACTICE_TUTOR_MODE=" + mode,
	}
	for k, v := range extraEnv {
		env = append(env, k+"="+v)
	}
	cmd.Env = env
	cmd.Stdin = strings.NewReader(stdin)

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	return outBuf.String(), errBuf.String(), err
}

func TestChatSh_SystemPromptDiffersByMode(t *testing.T) {
	modes := []string{"syntax-only", "hints-first", "full-assist"}
	prompts := make(map[string]string)

	for _, mode := range modes {
		mock := newMockOllama(t, "ok")
		_, stderr, err := runChat(t, mock.URL, mode, "hello\n", nil)
		if err != nil {
			t.Fatalf("mode %s: chat.sh failed: %v\nstderr: %s", mode, err, stderr)
		}

		reqs := mock.allRequests()
		if len(reqs) != 1 {
			t.Fatalf("mode %s: expected 1 request, got %d", mode, len(reqs))
		}
		if len(reqs[0].Messages) == 0 || reqs[0].Messages[0].Role != "system" {
			t.Fatalf("mode %s: expected first message to have role=system, got %+v", mode, reqs[0].Messages)
		}
		prompts[mode] = reqs[0].Messages[0].Content
		if prompts[mode] == "" {
			t.Errorf("mode %s: system prompt is empty", mode)
		}
	}

	if prompts["syntax-only"] == prompts["hints-first"] {
		t.Error("syntax-only and hints-first got the same system prompt — modes must differ")
	}
	if prompts["syntax-only"] == prompts["full-assist"] {
		t.Error("syntax-only and full-assist got the same system prompt — modes must differ")
	}
	if prompts["hints-first"] == prompts["full-assist"] {
		t.Error("hints-first and full-assist got the same system prompt — modes must differ")
	}
}

func TestChatSh_RetainsConversationHistory(t *testing.T) {
	mock := newMockOllama(t, "assistant-reply")

	_, stderr, err := runChat(t, mock.URL, "full-assist", "first line\nsecond line\n", nil)
	if err != nil {
		t.Fatalf("chat.sh failed: %v\nstderr: %s", err, stderr)
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
	if second[2].Role != "assistant" || second[2].Content != "assistant-reply" {
		t.Errorf("second request messages[2] (assistant1) = %+v, want role=assistant content=%q", second[2], "assistant-reply")
	}
	if second[3].Content != "second line" {
		t.Errorf("second request messages[3] (user2) = %q, want %q", second[3].Content, "second line")
	}
}

func TestChatSh_PrintsAssistantReply(t *testing.T) {
	mock := newMockOllama(t, "the answer is 42")

	stdout, stderr, err := runChat(t, mock.URL, "full-assist", "question\n", nil)
	if err != nil {
		t.Fatalf("chat.sh failed: %v\nstderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "the answer is 42") {
		t.Errorf("stdout = %q, want it to contain the assistant reply", stdout)
	}
}

func TestChatSh_TutorSystemPromptOverrideRespected(t *testing.T) {
	mock := newMockOllama(t, "ok")

	_, stderr, err := runChat(t, mock.URL, "syntax-only", "hi\n", map[string]string{
		"TUTOR_SYSTEM_PROMPT": "custom override prompt",
	})
	if err != nil {
		t.Fatalf("chat.sh failed: %v\nstderr: %s", err, stderr)
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 || len(reqs[0].Messages) == 0 {
		t.Fatalf("expected 1 request with at least a system message, got %+v", reqs)
	}
	if got := reqs[0].Messages[0].Content; got != "custom override prompt" {
		t.Errorf("system prompt = %q, want override %q", got, "custom override prompt")
	}
}

func TestChatSh_HandlesUnreachableHostGracefully(t *testing.T) {
	// Port 1 is reserved/unlisted — connection is refused immediately
	// rather than hanging, so this stays fast without an unreachable IP.
	stdout, stderr, err := runChat(t, "http://127.0.0.1:1", "full-assist", "hello\n", nil)
	if err != nil {
		t.Fatalf("chat.sh should exit cleanly even when Ollama is unreachable, got error: %v\nstdout: %s\nstderr: %s", err, stdout, stderr)
	}
	if !strings.Contains(stderr, "could not reach") {
		t.Errorf("stderr = %q, want a message about being unable to reach the host", stderr)
	}
}
