package tutor

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino/schema"
)

// TestAgent_RoundTripsThroughMockOllama proves the eino wiring itself
// (newChatModel -> newAgent -> Generate) actually reaches an Ollama
// /api/chat endpoint and returns its reply — no tools yet, that's
// covered once the real tools exist (see tools_test.go). Reuses
// sequencedOllama from tutor_test.go (same package).
func TestAgent_RoundTripsThroughMockOllama(t *testing.T) {
	mock := newSequencedOllama(t, "hello from the mock")
	ctx := context.Background()

	cfg := Config{OllamaHost: mock.URL, Model: "test-model"}
	cm, err := newChatModel(ctx, cfg)
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}

	agent, err := newAgent(ctx, cm, nil)
	if err != nil {
		t.Fatalf("newAgent: %v", err)
	}

	reply, err := agent.Generate(ctx, []*schema.Message{schema.UserMessage("hi")})
	if err != nil {
		t.Fatalf("agent.Generate: %v", err)
	}
	if reply.Content != "hello from the mock" {
		t.Errorf("reply.Content = %q, want %q", reply.Content, "hello from the mock")
	}

	reqs := mock.allRequests()
	if len(reqs) != 1 {
		t.Fatalf("expected 1 request to the mock, got %d", len(reqs))
	}
	if reqs[0].Stream {
		t.Error("request had stream=true — Generate must always send stream:false (see agent.go's newAgent doc comment)")
	}
}

// TestNewChatModel_TimesOutIfOllamaHangs simulates a stalled Ollama
// server (e.g. an overloaded or wedged instance never responding) —
// without ollamaRequestTimeout, this would hang the HTTP request (and
// by extension the whole synchronous tutor turn loop) indefinitely.
func TestNewChatModel_TimesOutIfOllamaHangs(t *testing.T) {
	hang := make(chan struct{})
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-hang // never respond until the test cleans up
	}))
	// One cleanup, not two separately t.Cleanup-registered ones: Close
	// blocks until in-flight requests finish, so hang must be closed
	// first, in the same func, not left to (LIFO, easy to get backwards)
	// cleanup-ordering to get right.
	t.Cleanup(func() {
		close(hang)
		mock.Close()
	})

	origTimeout := ollamaRequestTimeout
	ollamaRequestTimeout = 200 * time.Millisecond
	t.Cleanup(func() { ollamaRequestTimeout = origTimeout })

	ctx := context.Background()
	cfg := Config{OllamaHost: mock.URL, Model: "test-model"}
	cm, err := newChatModel(ctx, cfg)
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}
	agent, err := newAgent(ctx, cm, nil)
	if err != nil {
		t.Fatalf("newAgent: %v", err)
	}

	start := time.Now()
	_, err = agent.Generate(ctx, []*schema.Message{schema.UserMessage("hi")})
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error when Ollama hangs past ollamaRequestTimeout")
	}
	if elapsed > 5*time.Second {
		t.Errorf("agent.Generate took %v to return, want it to time out quickly (~%v)", elapsed, ollamaRequestTimeout)
	}
}

// TestGenerateWithLeakRetry_ExportedWrapperProtectsAgainstLeaks proves
// the exported wrapper cmd/tutor-eval calls actually gets the same
// leak-detection retry production Run() gets, not a bare agent.Generate
// call — a real gap found live: cmd/tutor-eval calling agent.Generate
// directly showed a hints-first scenario failing ~25% of the time on
// exactly this leaked-JSON pattern, a failure mode that can't reach a
// real user (Run always retries/falls back) but was muddying the eval's
// signal by reporting it as a scenario failure anyway.
func TestGenerateWithLeakRetry_ExportedWrapperProtectsAgainstLeaks(t *testing.T) {
	mock := newSequencedOllama(t, `{"name": "suggest_hash_map", "parameters": {}}`, "clean answer")
	ctx := context.Background()

	cfg := Config{OllamaHost: mock.URL, Model: "test-model"}
	cm, err := newChatModel(ctx, cfg)
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}
	agent, err := newAgent(ctx, cm, nil)
	if err != nil {
		t.Fatalf("newAgent: %v", err)
	}

	reply, err := GenerateWithLeakRetry(ctx, agent, []*schema.Message{schema.UserMessage("hi")}, io.Discard)
	if err != nil {
		t.Fatalf("GenerateWithLeakRetry: %v", err)
	}
	if strings.Contains(reply.Content, `{"name"`) {
		t.Errorf("reply.Content = %q, still contains leaked tool-call JSON", reply.Content)
	}
	if reply.Content != "clean answer" {
		t.Errorf("reply.Content = %q, want the retry's clean reply %q", reply.Content, "clean answer")
	}
	if n := len(mock.allRequests()); n != 2 {
		t.Errorf("requests = %d, want exactly 2 (original + one retry)", n)
	}
}
