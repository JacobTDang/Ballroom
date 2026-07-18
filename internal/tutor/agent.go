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

// newChatModel builds a chat model — Ollama-backed by default, or
// OpenRouter (an OpenAI-compatible API) when modelName has the
// OpenRouterModelPrefix. Returns the interface type (not either
// concrete *ollama.ChatModel or *openai.ChatModel) since callers
// (newAgent, CheckToolCalling) only ever need model.ToolCallingChatModel
// — both concrete types satisfy it.
//
// Takes modelName/ollamaHost/apiKey directly rather than a whole Config
// so it can build a chat model for either role once orchestrator/worker
// routing exists (internal/tutor.Run — a session may need two chat
// models, only one of which corresponds to cfg.Model): ollamaHost and
// apiKey are shared across both roles regardless (one local Ollama
// server, one OpenRouter account key, confirmed live this session --
// "one key works for all models"), only the model name actually
// differs per role.
//
// Temperature matches tutor/chat.sh's previous options.temperature —
// kept low so the tutor's tone/behavior stays consistent across turns
// rather than drifting with a high-temperature sample. Set on BOTH
// provider paths: the OpenRouter path used to omit it, which left the
// tutor's sampling to whatever each hosted model's default was —
// meaningfully spikier than 0.2 on the free-tier models actually in
// use.
const tutorTemperature float32 = 0.2

func newChatModel(ctx context.Context, modelName, ollamaHost, apiKey string) (model.ToolCallingChatModel, error) {
	if strings.HasPrefix(modelName, OpenRouterModelPrefix) {
		temperature := tutorTemperature
		cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
			BaseURL:     openRouterBaseURL,
			APIKey:      apiKey,
			Model:       strings.TrimPrefix(modelName, OpenRouterModelPrefix),
			Timeout:     ollamaRequestTimeout,
			Temperature: &temperature,
		})
		if err != nil {
			return nil, fmt.Errorf("tutor: new chat model (openrouter): %w", err)
		}
		return cm, nil
	}

	cm, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: ollamaHost,
		Timeout: ollamaRequestTimeout,
		Model:   modelName,
		Options: &ollama.Options{Temperature: tutorTemperature},
	})
	if err != nil {
		return nil, fmt.Errorf("tutor: new chat model: %w", err)
	}
	return cm, nil
}

// reactMaxStep bounds how many graph steps (each model call and each
// tool execution is one step) a single agent.Generate can take before
// eino's react.Agent gives up with "[GraphRunError] exceeds max steps".
// Left unset, react.AgentConfig.MaxStep defaults to eino's own internal
// ~12 (node count + 10, per react.go's own comment) — only ~5-6
// tool-call rounds, which a real OpenRouter session (openai/gpt-oss-120b:free)
// hit live mid-conversation. Isolated re-tests of the identical failing
// scenario succeeded cleanly afterward, pointing to transient
// OpenRouter free-tier rate-limit pressure (this project's own
// concurrent testing, moments earlier) as the likely trigger rather
// than genuine model looping — but raising this is real, low-risk
// headroom against exactly this failure mode regardless of root cause,
// not a guess: verified via a mock requiring 8 tool-call rounds (16
// steps), comfortably past the old default but within this one.
const reactMaxStep = 30

// newAgent wires a ReAct agent (github.com/cloudwego/eino/flow/agent/react)
// around cm and tools. The agent handles the full call-model -> tool-calls
// -> execute -> feed-back -> loop cycle internally. Callers use Generate
// by default; Stream only when streamingEnabled says the model can take
// it (stream.go) — cmd/tutor-spike found eino's Generate path sends
// every internal model call with "stream": false, which sidesteps a
// known Ollama bug where streaming responses drop tool_calls for some
// model families, so Ollama-backed agents stay on Generate exclusively.
//
// StreamToolCallChecker only affects the Stream path — see
// windowedStreamToolCallChecker for the live narrate-then-call bug the
// default first-chunk checker caused.
func newAgent(ctx context.Context, cm model.ToolCallingChatModel, tools []tool.BaseTool) (*react.Agent, error) {
	a, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel:      cm,
		ToolsConfig:           compose.ToolsNodeConfig{Tools: tools},
		MaxStep:               reactMaxStep,
		StreamToolCallChecker: windowedStreamToolCallChecker,
	})
	if err != nil {
		return nil, fmt.Errorf("tutor: new agent: %w", err)
	}
	return a, nil
}

// NewChatModel builds a chat model for modelName exactly the way a
// real session does -- Ollama by default, OpenRouter for
// OpenRouterModelPrefix-prefixed names. Exported for cmd/tutor-eval,
// which used to construct a raw Ollama client directly and therefore
// couldn't evaluate the OpenRouter models real sessions actually run.
func NewChatModel(ctx context.Context, modelName, ollamaHost, apiKey string) (model.ToolCallingChatModel, error) {
	return newChatModel(ctx, modelName, ollamaHost, apiKey)
}
