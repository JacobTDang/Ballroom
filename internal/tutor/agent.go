package tutor

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ollama"
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

// newChatModel builds the Ollama-backed chat model the agent calls.
// Temperature matches tutor/chat.sh's previous options.temperature — kept
// low so the tutor's tone/behavior stays consistent across turns rather
// than drifting with a high-temperature sample.
func newChatModel(ctx context.Context, cfg Config) (*ollama.ChatModel, error) {
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
