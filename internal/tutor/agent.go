package tutor

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

// ollamaRequestTimeout bounds each individual HTTP request the chat
// model makes to Ollama — a react.Agent turn can make several of these
// sequentially (model call, execute tools, model call again, ...), so
// this is a per-request bound, not a whole-turn one. Without it, a
// stalled Ollama server would hang a request (and by extension the
// whole synchronous tutor turn loop) indefinitely, with no way to
// recover short of killing the process. Generous — llama3.1:8b on CPU
// genuinely can take tens of seconds for a longer generation — this is
// meant to catch a truly stuck request, not a merely slow one. A var,
// not a const, so tests can shrink it rather than waiting out the real
// duration.
var ollamaRequestTimeout = 120 * time.Second

// OpenRouterModelPrefix marks a Config.Model value as an OpenRouter
// model slug (e.g. "openrouter:anthropic/claude-3.5-sonnet") rather
// than a local Ollama tag — a prefix on the existing model string
// rather than a separate provider field, so every place that already
// takes a model string (the TUI picker, settings.json, tutor-eval)
// keeps working unchanged; newChatModel is the only place that needs
// to branch on it.
const OpenRouterModelPrefix = "openrouter:"

// openRouterBaseURL is a var, not a const, so tests can point it at a
// mock server — same reason ollamaRequestTimeout above is a var rather
// than hardcoded.
var openRouterBaseURL = "https://openrouter.ai/api/v1"

// newChatModel builds the chat model the agent calls — Ollama-backed by
// default, or OpenRouter (an OpenAI-compatible API) when cfg.Model has
// the OpenRouterModelPrefix. Returns the interface type (not either
// concrete *ollama.ChatModel or *openai.ChatModel) since callers
// (newAgent, and eventually CheckToolCalling) only ever need
// model.ToolCallingChatModel — both concrete types satisfy it.
//
// Temperature matches tutor/chat.sh's previous options.temperature —
// kept low so the tutor's tone/behavior stays consistent across turns
// rather than drifting with a high-temperature sample. Only set on the
// Ollama path for now; openai.ChatModelConfig's equivalent (Temperature
// *float32) can be wired the same way if OpenRouter's default proves
// too high-variance in practice.
func newChatModel(ctx context.Context, cfg Config) (model.ToolCallingChatModel, error) {
	if strings.HasPrefix(cfg.Model, OpenRouterModelPrefix) {
		cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
			BaseURL: openRouterBaseURL,
			APIKey:  cfg.APIKey,
			Model:   strings.TrimPrefix(cfg.Model, OpenRouterModelPrefix),
			Timeout: ollamaRequestTimeout,
		})
		if err != nil {
			return nil, fmt.Errorf("tutor: new chat model (openrouter): %w", err)
		}
		return cm, nil
	}

	cm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: cfg.OllamaHost,
		Timeout: ollamaRequestTimeout,
		Model:   cfg.Model,
		Options: &ollama.Options{Temperature: 0.2},
	})
	if err != nil {
		return nil, fmt.Errorf("tutor: new chat model: %w", err)
	}
	return cm, nil
}

// newAgent wires a ReAct agent (github.com/cloudwego/eino/flow/agent/react)
// around cm and tools. The agent handles the full call-model -> tool-calls
// -> execute -> feed-back -> loop cycle internally; callers only ever call
// Generate, never Stream — see cmd/tutor-spike's finding that eino's
// Generate path sends every internal model call with "stream": false,
// which sidesteps a known Ollama bug where streaming responses drop
// tool_calls for some model families.
func newAgent(ctx context.Context, cm model.ToolCallingChatModel, tools []tool.BaseTool) (*react.Agent, error) {
	a, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: cm,
		ToolsConfig:      compose.ToolsNodeConfig{Tools: tools},
	})
	if err != nil {
		return nil, fmt.Errorf("tutor: new agent: %w", err)
	}
	return a, nil
}
