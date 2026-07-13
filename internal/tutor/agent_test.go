package tutor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
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
	cm, err := newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
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
	cm, err := newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
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

// TestNewAgent_SurvivesManyToolCallRoundsWithinRaisedMaxStep is a
// regression test for a real bug found live: a real OpenRouter session
// (openai/gpt-oss-120b:free) failed mid-conversation with eino's own
// "[GraphRunError] exceeds max steps" — react.AgentConfig.MaxStep
// wasn't set, defaulting to eino's internal ~12 (node count + 10, per
// react.go's own comment), which only covers ~5-6 tool-call rounds
// (each round is 2 graph steps: one model-call node, one tool-exec
// node). Isolated re-tests of the identical failing scenario succeeded
// cleanly (twice), pointing to transient OpenRouter free-tier rate-limit
// pressure from concurrent testing as the likely trigger rather than
// genuine model looping — but raising MaxStep is a real, low-risk
// hardening regardless of root cause: it's cheap headroom against
// exactly this failure mode, verified here with a mock that requires 8
// tool-call rounds (16 steps), comfortably past the old ~12-step default
// but within the raised limit.
func TestNewAgent_SurvivesManyToolCallRoundsWithinRaisedMaxStep(t *testing.T) {
	const rounds = 8
	calls := 0
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Set("Content-Type", "application/json")
		if calls <= rounds {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"message": map[string]any{
					"role": "assistant", "content": "",
					"tool_calls": []map[string]any{
						{"function": map[string]any{"name": "ping", "arguments": map[string]any{"reason": "check"}}},
					},
				},
				"done": true,
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"role": "assistant", "content": "done after many rounds"},
			"done":    true,
		})
	}))
	t.Cleanup(mock.Close)

	ctx := context.Background()
	cfg := Config{OllamaHost: mock.URL, Model: "test-model"}
	cm, err := newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}
	pingTool, err := utils.InferTool("ping", "Call this tool.", func(_ context.Context, _ struct {
		Reason string `json:"reason"`
	}) (string, error) {
		return "pong", nil
	})
	if err != nil {
		t.Fatalf("InferTool: %v", err)
	}
	agent, err := newAgent(ctx, cm, []tool.BaseTool{pingTool})
	if err != nil {
		t.Fatalf("newAgent: %v", err)
	}

	reply, err := agent.Generate(ctx, []*schema.Message{schema.UserMessage("go")})
	if err != nil {
		t.Fatalf("agent.Generate: %v (want it to survive %d tool-call rounds without exceeding max steps)", err, rounds)
	}
	if reply.Content != "done after many rounds" {
		t.Errorf("reply.Content = %q, want %q", reply.Content, "done after many rounds")
	}
}

// TestNewChatModel_RoutesOpenRouterPrefixedModelToOpenAICompatibleClient
// proves newChatModel's branch: an OpenRouterModelPrefix-prefixed
// Config.Model reaches a mock standing in for OpenRouter's OpenAI-
// compatible /chat/completions endpoint (openRouterBaseURL is a var
// specifically so this test can redirect it) rather than trying to
// treat the prefix as a literal Ollama tag, and the prefix itself is
// stripped before the model name reaches the request body — OpenRouter
// wouldn't recognize "openrouter:anthropic/claude-3.5-sonnet" as a
// model slug, only "anthropic/claude-3.5-sonnet".
func TestNewChatModel_RoutesOpenRouterPrefixedModelToOpenAICompatibleClient(t *testing.T) {
	var gotModel string
	var gotAuth string
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		var req struct {
			Model string `json:"model"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		gotModel = req.Model
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "chatcmpl-test", "object": "chat.completion", "created": 1,
			"model": req.Model,
			"choices": []map[string]any{
				{"index": 0, "message": map[string]string{"role": "assistant", "content": "hello from openrouter mock"}, "finish_reason": "stop"},
			},
		})
	}))
	t.Cleanup(mock.Close)

	origBaseURL := openRouterBaseURL
	openRouterBaseURL = mock.URL
	t.Cleanup(func() { openRouterBaseURL = origBaseURL })

	ctx := context.Background()
	cfg := Config{Model: OpenRouterModelPrefix + "anthropic/claude-3.5-sonnet", APIKey: "sk-test-key"}
	cm, err := newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
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
	if reply.Content != "hello from openrouter mock" {
		t.Errorf("reply.Content = %q, want %q", reply.Content, "hello from openrouter mock")
	}
	if gotModel != "anthropic/claude-3.5-sonnet" {
		t.Errorf("request model = %q, want the OpenRouterModelPrefix stripped: %q", gotModel, "anthropic/claude-3.5-sonnet")
	}
	if gotAuth != "Bearer sk-test-key" {
		t.Errorf("Authorization header = %q, want %q", gotAuth, "Bearer sk-test-key")
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
	cm, err := newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
	if err != nil {
		t.Fatalf("newChatModel: %v", err)
	}
	agent, err := newAgent(ctx, cm, nil)
	if err != nil {
		t.Fatalf("newAgent: %v", err)
	}

	reply, err := GenerateWithLeakRetry(ctx, agent, []*schema.Message{schema.UserMessage("hi")})
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
