package tutor

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
)

// newChatModel builds the Ollama-backed chat model the agent calls.
// Temperature matches tutor/chat.sh's previous options.temperature — kept
// low so the tutor's tone/behavior stays consistent across turns rather
// than drifting with a high-temperature sample.
func newChatModel(ctx context.Context, cfg Config) (*ollama.ChatModel, error) {
	cm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: cfg.OllamaHost,
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
