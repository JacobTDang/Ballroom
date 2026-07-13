package tutor

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type pingArgs struct {
	Reason string `json:"reason" jsonschema:"description=why you are pinging"`
}

// CheckToolCalling reports whether model actually supports the
// provider's structured tool_calls response field, by binding one
// trivial "ping" tool and issuing a directive prompt only a real tool
// call can satisfy. Some models (confirmed: qwen2.5-coder:7b — see
// config.DefaultTutorModel's doc comment) emit tool-call-shaped JSON as
// plain text content instead of populating tool_calls; this call
// distinguishes that failure mode directly rather than a caller
// inferring it from unrelated symptoms during a real tutor session.
//
// ollamaHost and apiKey are only used for their respective providers
// (see newChatModel) — pass "" for whichever doesn't apply to model.
//
// Ports the same wiring cmd/tutor-spike used to first confirm this
// failure mode, as a small reusable check instead of a throwaway
// binary.
func CheckToolCalling(ollamaHost, model, apiKey string) (bool, error) {
	ctx := context.Background()

	cm, err := newChatModel(ctx, model, ollamaHost, apiKey)
	if err != nil {
		return false, fmt.Errorf("tutor: check tool calling: %w", err)
	}

	called := false
	pingTool, err := utils.InferTool("ping", "Call this tool to respond.", func(_ context.Context, _ pingArgs) (string, error) {
		called = true
		return "pong", nil
	})
	if err != nil {
		return false, fmt.Errorf("tutor: check tool calling: %w", err)
	}

	agent, err := newAgent(ctx, cm, []tool.BaseTool{pingTool})
	if err != nil {
		return false, fmt.Errorf("tutor: check tool calling: %w", err)
	}

	_, err = agent.Generate(ctx, []*schema.Message{
		schema.UserMessage(`Call the "ping" tool right now with reason "check". Do not reply with text — call the tool.`),
	})
	if err != nil {
		return false, fmt.Errorf("tutor: check tool calling: %w", err)
	}

	return called, nil
}
