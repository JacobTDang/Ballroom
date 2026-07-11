package tutor

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
)

// TestAgent_RoundTripsThroughMockOllama proves the eino wiring itself
// (newChatModel -> newAgent -> Generate) actually reaches an Ollama
// /api/chat endpoint and returns its reply — no tools yet, that's
// covered once the real tools exist (see tools_test.go). Reuses the
// mockOllama helper from chat_test.go (same package).
func TestAgent_RoundTripsThroughMockOllama(t *testing.T) {
	mock := newMockOllama(t, "hello from the mock")
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
