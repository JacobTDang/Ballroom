// Command tutor-spike is a throwaway diagnostic (Milestone 0 of the tutor
// rewrite plan) that answers one question before any production code gets
// built on top of the assumption: does qwen2.5-coder:7b actually call
// tools reliably through Ollama + eino's ReAct agent? eino's own official
// examples use a model fine-tuned for tool use (llama3-groq-tool-use),
// not a generic coder model, so this is unverified until checked here.
//
// Not wired into any CI job and never invoked by tests — run manually
// against a real local Ollama:
//
//	go run ./cmd/tutor-spike
//
// Requires Ollama running locally with qwen2.5-coder:7b pulled.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// loggingTransport wraps the default HTTP transport to print each request
// body Ollama actually receives — specifically to confirm eino sends
// "stream":false on every call in the ReAct loop (Agent.Generate is
// documented/sourced to always use blocking Invoke internally, never
// Stream, which is what sidesteps Ollama's known streaming-drops-
// tool_calls bug for Qwen models — this is the empirical check for that).
type loggingTransport struct{}

func (loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if dump, err := httputil.DumpRequestOut(req, true); err == nil {
		fmt.Printf("\n--- outgoing request ---\n%s\n------------------------\n", dump)
	}
	return http.DefaultTransport.RoundTrip(req)
}

// weatherInput/weatherOutput are a deliberately trivial tool — the point
// of this spike is only to observe whether the model calls a tool at
// all, not to test anything about Ballroom's real tools yet.
type weatherInput struct {
	City string `json:"city" jsonschema:"description=the city to look up the weather for"`
}

type weatherOutput struct {
	Forecast string `json:"forecast"`
}

func getWeather(_ context.Context, in weatherInput) (weatherOutput, error) {
	fmt.Printf(">>> TOOL CALLED: get_weather(city=%q)\n", in.City)
	return weatherOutput{Forecast: fmt.Sprintf("sunny and 72F in %s", in.City)}, nil
}

func main() {
	ctx := context.Background()

	httpClient := &http.Client{Transport: loggingTransport{}}

	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL:    "http://localhost:11434",
		Model:      "llama3.1:8b",
		HTTPClient: httpClient,
	})
	if err != nil {
		log.Fatalf("new chat model: %v", err)
	}

	weatherTool, err := utils.InferTool("get_weather", "Get the current weather forecast for a city", getWeather)
	if err != nil {
		log.Fatalf("infer tool: %v", err)
	}

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig:      compose.ToolsNodeConfig{Tools: []tool.BaseTool{weatherTool}},
	})
	if err != nil {
		log.Fatalf("new agent: %v", err)
	}

	prompts := []string{
		"What's the weather like in San Francisco right now?",
		"I'm planning a trip to Tokyo, can you check the forecast for me?",
		"Tell me about the weather in Berlin.",
		"weather in Chicago?",
		"Hi, how are you?", // control case: should NOT call the tool
	}

	for i, p := range prompts {
		fmt.Printf("\n========== PROMPT %d: %q ==========\n", i+1, p)
		msg, err := agent.Generate(ctx, []*schema.Message{schema.UserMessage(p)})
		if err != nil {
			fmt.Printf("!!! ERROR: %v\n", err)
			continue
		}
		fmt.Printf("<<< FINAL REPLY: %s\n", msg.Content)
	}
}
