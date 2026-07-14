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
// Takes ctx from the caller (rather than context.Background() internally,
// as this did before newTutorModel started calling it as a blocking
// session-construction step) so real cancellation -- e.g. quitting
// during startup -- actually works; the TUI's own fire-and-forget
// tea.Cmd use just passes context.Background() itself, unaffected.
//
// Ports the same wiring cmd/tutor-spike used to first confirm this
// failure mode, as a small reusable check instead of a throwaway
// binary.
func CheckToolCalling(ctx context.Context, ollamaHost, model, apiKey string) (bool, error) {
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

// checkToolCallingForSession is CheckToolCalling, as a var so tests can
// substitute a fake that skips the real network round-trip
// newTutorModel's strategy detection (detectStrategy, below) would
// otherwise make — same indirection pattern internal/tui/app.go's
// checkToolCallingFn and cmd/ballroom's own var of the same name use for
// the identical reason. See testmain_test.go for this package's own
// test-suite default.
var checkToolCallingForSession = CheckToolCalling

// detectStrategy reports which toolCallingStrategy modelName needs,
// via checkToolCallingForSession. Fails open toward nativeToolCalling —
// this project's sole behavior before this strategy existed — on
// anything short of a definitive "no": a check that errors out (bad
// host, timeout) can't actually tell us the model lacks real tool
// calling, so silently downgrading a previously-working native session
// to the fallback loop on a transient check failure would be a worse
// outcome than just trying native and letting the turn itself fail
// visibly if the host really is unreachable.
func detectStrategy(ctx context.Context, ollamaHost, modelName, apiKey string) toolCallingStrategy {
	supported, err := checkToolCallingForSession(ctx, ollamaHost, modelName, apiKey)
	if err != nil || supported {
		return nativeToolCalling
	}
	return jsonFallbackToolCalling
}
