package tutor

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	template "github.com/cloudwego/eino/utils/callbacks"
)

// activityThinkingStatus is the activity region's status line while a
// turn is in flight. Phase 1: static text. A later pass replaces this
// with a moving gradient-shimmer redraw (see the plan for this feature)
// — nothing else about the activity region changes for that, since the
// status line is just one more argument to showActivity.
const activityThinkingStatus = "⟳ Thinking…"

// activityArgsPreviewMax/activityResultPreviewMax cap how much of a raw
// tool-call argument/result string appears on one activity line — both
// are already subject to showActivity's own width truncation, but
// truncating here first keeps the preview a readable "at a glance"
// length even in a wide terminal, rather than filling the whole line
// with e.g. a long file's raw content.
const activityArgsPreviewMax = 40
const activityResultPreviewMax = 60

// activityCall is one tool invocation's current display state.
type activityCall struct {
	callID, name, status, detail string // status: "running" | "done" | "failed"
}

// activityFeed tracks the tool calls happening during one Generate call
// (a turn, or a comprehension check) and formats them into the lines
// showActivity displays. A fresh feed is used per call (see
// newActivityOption) — it is not a session-wide log.
type activityFeed struct {
	mu    sync.Mutex
	calls []activityCall
}

// started records a new running call, capping the list at
// activityToolLines by dropping the oldest entry — the same "most
// recent N" trade-off the activity region's fixed row count already
// requires. Returns the current formatted lines under the same lock, so
// a caller's redraw is never built from a state that's already stale by
// the time it reads it.
func (f *activityFeed) started(callID, name, argsPreview string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls = append(f.calls, activityCall{callID: callID, name: name, status: "running", detail: argsPreview})
	if len(f.calls) > activityToolLines {
		f.calls = f.calls[len(f.calls)-activityToolLines:]
	}
	return f.linesLocked()
}

// finished marks callID done with resultPreview. A callID that isn't
// found (already dropped by started's cap) is a no-op, not an error —
// eino always pairs a real OnStart with the matching OnEnd/OnError, but
// that OnStart's entry may have aged out of the capped list already.
func (f *activityFeed) finished(callID, resultPreview string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.updateLocked(callID, "done", resultPreview)
	return f.linesLocked()
}

// failed marks callID failed with errDetail. Same no-op-if-unknown
// behavior as finished.
func (f *activityFeed) failed(callID, errDetail string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.updateLocked(callID, "failed", errDetail)
	return f.linesLocked()
}

func (f *activityFeed) updateLocked(callID, status, detail string) {
	for i := range f.calls {
		if f.calls[i].callID == callID {
			f.calls[i].status = status
			f.calls[i].detail = detail
			return
		}
	}
}

func (f *activityFeed) linesLocked() []string {
	lines := make([]string, len(f.calls))
	for i, c := range f.calls {
		lines[i] = formatActivityLine(c)
	}
	return lines
}

// formatActivityLine renders one call's current state. "{}" (an empty
// JSON object -- what eino sends for a no-argument tool) is treated the
// same as no args at all, since showing "({})" on every no-arg tool call
// (most of this package's tools) would be noise, not information.
func formatActivityLine(c activityCall) string {
	switch c.status {
	case "done":
		if c.detail == "" {
			return fmt.Sprintf("✓ %s", c.name)
		}
		return fmt.Sprintf("✓ %s  %s", c.name, c.detail)
	case "failed":
		return fmt.Sprintf("✗ %s: %s", c.name, c.detail)
	default: // "running"
		if c.detail == "" || c.detail == "{}" {
			return fmt.Sprintf("→ %s", c.name)
		}
		return fmt.Sprintf("→ %s(%s)", c.name, c.detail)
	}
}

// newActivityOption builds the agent.AgentOption that wires a fresh
// activityFeed into box's activity display via real eino tool-call
// callbacks (react.BuildAgentCallback / utils/callbacks.ToolCallbackHandler
// — OnStart/OnEnd/OnError fire live, while Generate is still running, not
// after the fact). box == nil (no real terminal, e.g. cmd/tutor-eval)
// returns a harmless no-op option — there's no reserved region to draw
// into, matching how every other box-dependent feature in this package
// already degrades.
//
// Call fresh for each Generate call (a turn, or a comprehension check),
// never shared/reused across calls: each one's tool-call window starts
// empty, the same way each turn's own conversation starts fresh context.
//
// compose.GetToolCallID(ctx) is what correlates a specific call's OnEnd
// back to its own OnStart entry (not just by tool name) — eino sets this
// on ctx before invoking the tool, and that same ctx (returned by
// OnStart, per eino's own callback contract) is what OnEnd receives, so
// two concurrent calls to the same tool are tracked as separate entries
// rather than one clobbering the other's status.
func newActivityOption(box *inputBox) agent.AgentOption {
	if box == nil {
		return agent.WithComposeOptions()
	}

	feed := &activityFeed{}
	toolHandler := &template.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			argsPreview := ""
			if input != nil {
				argsPreview = truncateLine(input.ArgumentsInJSON, activityArgsPreviewMax)
			}
			lines := feed.started(compose.GetToolCallID(ctx), info.Name, argsPreview)
			box.showActivity(activityThinkingStatus, lines)
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
			resultPreview := ""
			if output != nil {
				resultPreview = truncateLine(output.Response, activityResultPreviewMax)
			}
			lines := feed.finished(compose.GetToolCallID(ctx), resultPreview)
			box.showActivity(activityThinkingStatus, lines)
			return ctx
		},
		OnError: func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
			lines := feed.failed(compose.GetToolCallID(ctx), truncateLine(err.Error(), activityResultPreviewMax))
			box.showActivity(activityThinkingStatus, lines)
			return ctx
		},
	}
	handler := react.BuildAgentCallback(nil, toolHandler)
	return agent.WithComposeOptions(compose.WithCallbacks(handler))
}
