package tutor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newToolCallOllama simulates an Ollama /api/chat server whose first
// response includes a real tool_calls entry (matching the wire shape
// documented in github.com/eino-contrib/ollama's api.Message —
// message.tool_calls[].function.{name,arguments}) — sequencedOllama
// (tutor_test.go) can only ever return plain text content, so it can't
// exercise CheckToolCalling's "the model actually called the tool"
// path. The second response is plain content, ending the ReAct loop
// once the tool result is fed back.
func newToolCallOllama(t *testing.T, toolName string) *httptest.Server {
	t.Helper()
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.Header().Set("Content-Type", "application/json")
		if requests == 1 {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"message": map[string]any{
					"role":    "assistant",
					"content": "",
					"tool_calls": []map[string]any{
						{"function": map[string]any{"name": toolName, "arguments": map[string]any{"reason": "check"}}},
					},
				},
				"done": true,
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": map[string]string{"role": "assistant", "content": "pong received"},
			"done":    true,
		})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestCheckToolCalling_ReportsTrueWhenModelCallsTheTool(t *testing.T) {
	srv := newToolCallOllama(t, "ping")

	supported, err := CheckToolCalling(context.Background(), srv.URL, "test-model", "")
	if err != nil {
		t.Fatalf("CheckToolCalling: %v", err)
	}
	if !supported {
		t.Error("supported = false, want true — the mock model called the ping tool")
	}
}

func TestCheckToolCalling_ReportsFalseWhenModelOnlyRepliesWithText(t *testing.T) {
	mock := newSequencedOllama(t, "I would call the tool but here is text instead")

	supported, err := CheckToolCalling(context.Background(), mock.URL, "test-model", "")
	if err != nil {
		t.Fatalf("CheckToolCalling: %v", err)
	}
	if supported {
		t.Error("supported = true, want false — the mock model never populated tool_calls")
	}
}

func TestCheckToolCalling_UnreachableHostReturnsError(t *testing.T) {
	_, err := CheckToolCalling(context.Background(), "http://127.0.0.1:1", "test-model", "")
	if err == nil {
		t.Fatal("expected an error for an unreachable Ollama host")
	}
}

// newToolCallOpenAICompatible is newToolCallOllama's OpenAI-compatible
// counterpart (OpenRouter's actual wire format) — first response
// includes a real tool_calls entry, second is plain content, ending the
// ReAct loop. Also records the Authorization header so the apiKey
// threading itself is verified, not just that a tool call happened.
func newToolCallOpenAICompatible(t *testing.T, toolName string) (srv *httptest.Server, gotAuth *string) {
	t.Helper()
	requests := 0
	var auth string
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		auth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		if requests == 1 {
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id": "chatcmpl-test", "object": "chat.completion", "created": 1, "model": "test-model",
				"choices": []map[string]any{
					{"index": 0, "finish_reason": "tool_calls", "message": map[string]any{
						"role": "assistant",
						"tool_calls": []map[string]any{
							{"id": "call_1", "type": "function", "function": map[string]any{"name": toolName, "arguments": `{"reason":"check"}`}},
						},
					}},
				},
			})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": "chatcmpl-test2", "object": "chat.completion", "created": 1, "model": "test-model",
			"choices": []map[string]any{
				{"index": 0, "finish_reason": "stop", "message": map[string]string{"role": "assistant", "content": "pong received"}},
			},
		})
	}))
	t.Cleanup(srv.Close)
	return srv, &auth
}

func TestCheckToolCalling_RoutesOpenRouterPrefixedModelAndThreadsAPIKey(t *testing.T) {
	srv, gotAuth := newToolCallOpenAICompatible(t, "ping")

	origBaseURL := openRouterBaseURL
	openRouterBaseURL = srv.URL
	t.Cleanup(func() { openRouterBaseURL = origBaseURL })

	supported, err := CheckToolCalling(context.Background(), "unused-for-openrouter", OpenRouterModelPrefix+"some/model", "sk-test-key")
	if err != nil {
		t.Fatalf("CheckToolCalling: %v", err)
	}
	if !supported {
		t.Error("supported = false, want true — the mock model called the ping tool")
	}
	if *gotAuth != "Bearer sk-test-key" {
		t.Errorf("Authorization header = %q, want %q", *gotAuth, "Bearer sk-test-key")
	}
}
