package tutor

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	agentopt "github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// fallbackToolCall is one parsed tool-call attempt from a
// jsonFallbackToolCalling model's reply.
type fallbackToolCall struct {
	Name      string
	Arguments json.RawMessage
}

// fallbackToolCallHint pre-filters before a full parse: cheap way to
// decide "this reply is at least attempting a tool call" before paying
// for a brace-balanced scan. Accepts both "arguments" (the shape
// jsonFallbackInstruction asks for) and "parameters" (the shape real
// leaked native tool calls have actually used in this codebase -- see
// leakedToolCallPattern's doc comment and agent_test.go's own fixture)
// as aliases.
var fallbackToolCallHint = regexp.MustCompile(`\{"name"\s*:\s*"[a-zA-Z_][a-zA-Z0-9_]*"\s*,\s*"(arguments|parameters)"\s*:`)

// rawFallbackCall is the wire shape parseFallbackToolCall decodes into
// before folding the arguments/parameters alias down to a single field.
type rawFallbackCall struct {
	Name       string          `json:"name"`
	Arguments  json.RawMessage `json:"arguments"`
	Parameters json.RawMessage `json:"parameters"`
}

// parseFallbackToolCall reports whether content is a tool-call attempt
// and, if so, parses it. Tries the whole trimmed (and code-fence-
// stripped, if the entire content is one) reply as a call first -- the
// well-behaved case jsonFallbackInstruction asks for. Only if that fails
// AND fallbackToolCallHint matches somewhere does it fall back to a
// brace-balanced scan for the first {...} object, tolerating a preface
// the model wasn't supposed to add.
//
// matched=false (err always nil) means content is a final answer, not a
// call attempt at all -- the common case once the model is done calling
// tools. matched=true, err!=nil means it looks like an attempted call
// but doesn't parse -- the caller should treat this as a failed attempt
// (corrective retry), not silently fall through to treating content as
// the final answer.
//
// Known limitation, accepted rather than chased: a final answer that
// happens to describe or quote a tool-call-shaped JSON example (e.g.
// explaining what a call "would look like") can be misread as a real
// attempt. jsonFallbackInstruction tells the model its entire reply must
// be ONLY the JSON when calling a tool specifically to keep this
// vanishingly rare in practice, matching this codebase's own established
// lesson (see toolsInstruction's doc comment) that chasing every edge
// case with more parsing complexity tends to regress reliability more
// than it fixes.
func parseFallbackToolCall(content string) (fallbackToolCall, bool, error) {
	unfenced := stripJSONFence(strings.TrimSpace(content))
	if call, err := decodeFallbackCall(unfenced); err == nil {
		return call, true, nil
	}

	if !fallbackToolCallHint.MatchString(content) {
		return fallbackToolCall{}, false, nil
	}

	obj, ok := extractFirstJSONObject(content)
	if !ok {
		return fallbackToolCall{}, true, fmt.Errorf("tutor: reply looks like a tool call attempt but contains no valid JSON object")
	}
	call, err := decodeFallbackCall(string(obj))
	if err != nil {
		return fallbackToolCall{}, true, err
	}
	return call, true, nil
}

// stripJSONFence removes a ``` or ```json fence wrapping the entire
// (already-trimmed) string s, if present -- a model asked to reply with
// ONLY JSON sometimes wraps it in a code fence anyway. Returns s
// unchanged if it isn't fenced.
func stripJSONFence(s string) string {
	if !strings.HasPrefix(s, "```") {
		return s
	}
	body := strings.TrimPrefix(s, "```")
	body = strings.TrimPrefix(body, "json")
	body = strings.TrimPrefix(body, "\n")
	body = strings.TrimSuffix(strings.TrimSpace(body), "```")
	return strings.TrimSpace(body)
}

// decodeFallbackCall parses s as a rawFallbackCall and folds its
// arguments/parameters alias down to a single field, defaulting to an
// empty object for a tool that takes no arguments. Returns an error if s
// isn't valid JSON or has no "name".
func decodeFallbackCall(s string) (fallbackToolCall, error) {
	var raw rawFallbackCall
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		return fallbackToolCall{}, fmt.Errorf("tutor: parse fallback tool call: %w", err)
	}
	if raw.Name == "" {
		return fallbackToolCall{}, fmt.Errorf("tutor: fallback tool call has no \"name\"")
	}
	args := raw.Arguments
	if len(args) == 0 {
		args = raw.Parameters
	}
	if len(args) == 0 {
		args = json.RawMessage("{}")
	}
	return fallbackToolCall{Name: raw.Name, Arguments: args}, nil
}

// extractFirstJSONObject brace-scans s for the first top-level {...}
// span, honoring string literals and escaped quotes so a "}" or "{"
// inside a quoted value (e.g. a tool argument's own text) doesn't
// desync the scan. Returns ok=false if s has no balanced object.
func extractFirstJSONObject(s string) (json.RawMessage, bool) {
	start := strings.IndexByte(s, '{')
	if start == -1 {
		return nil, false
	}

	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if inString {
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == '"':
				inString = false
			}
			continue
		}
		switch c {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return json.RawMessage(s[start : i+1]), true
			}
		}
	}
	return nil, false
}

// fallbackRoundCap mirrors reactMaxStep (agent.go): react.Agent's
// MaxStep=30 gives ~14 tool rounds + 1 final call = 15 model calls (each
// round = model call + tool exec = 2 steps). fallbackRoundCap counts
// model calls directly, so 15 gives the fallback path the same real
// budget, not just a superficially similar-looking number.
const fallbackRoundCap = 15

// fallbackLoopExhaustedReply is shown when the model never settles on a
// final answer within fallbackRoundCap rounds -- the raw last tool-call
// JSON (or whatever it was mid-attempt) must never reach the user,
// matching generateWithLeakRetry's own fallback-reply guarantee.
const fallbackLoopExhaustedReply = "Sorry, I wasn't able to work through that just now -- could you try asking again?"

// findTool returns the tool named name from tools if it's both present
// and invokable, or nil otherwise. Every tool buildTools returns is
// InvokableTool (confirmed: utils.WrapToolWithErrorHandler's wrapper
// still implements it for an Invokable-only input), so a nil return in
// practice only ever means "no tool with that name exists".
func findTool(ctx context.Context, tools []tool.BaseTool, name string) tool.InvokableTool {
	for _, t := range tools {
		info, err := t.Info(ctx)
		if err != nil || info == nil || info.Name != name {
			continue
		}
		if it, ok := t.(tool.InvokableTool); ok {
			return it
		}
	}
	return nil
}

// safeInvokeTool wraps t.InvokableRun with its own recover(). eino's
// compose graph runner recovers tool panics internally
// (compose/tool_node.go, compose/graph_manager.go) -- calling
// InvokableRun by hand here bypasses that safety net entirely, and an
// unrecovered panic in startTurn's goroutine would crash the whole
// process, not just the turn.
func safeInvokeTool(ctx context.Context, t tool.InvokableTool, argsJSON string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("tutor: tool panicked: %v", r)
		}
	}()
	return t.InvokableRun(ctx, argsJSON)
}

// pushActivity publishes feed's current state onto activityCh without
// blocking -- same non-blocking-send-or-drop semantics as
// buildActivityChannelOption's own push closure (model.go), extracted
// here so runFallbackToolLoop uses identical behavior instead of a
// second copy of the same select/default.
func pushActivity(feed *activityFeed, activityCh chan<- []activityCall) {
	select {
	case activityCh <- feed.currentCalls():
	default:
	}
}

// renderToolCatalog renders each tool's name/description/JSON schema as
// text, teaching a jsonFallbackToolCalling model what tools exist and
// how to call them -- computed once per session (newTutorModel), not
// per call, since buildTools(cfg)'s output is identical for both roles
// and doesn't change turn to turn.
func renderToolCatalog(ctx context.Context, tools []tool.BaseTool) (string, error) {
	var b strings.Builder
	b.WriteString("Available tools:\n")
	for _, t := range tools {
		info, err := t.Info(ctx)
		if err != nil {
			return "", fmt.Errorf("tutor: render tool catalog: %w", err)
		}
		fmt.Fprintf(&b, "\n- %s: %s\n", info.Name, info.Desc)

		params, err := info.ParamsOneOf.ToJSONSchema()
		if err != nil {
			return "", fmt.Errorf("tutor: render tool catalog: %s: %w", info.Name, err)
		}
		if params == nil {
			b.WriteString("  Takes no arguments.\n")
			continue
		}
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			return "", fmt.Errorf("tutor: render tool catalog: %s: marshal schema: %w", info.Name, err)
		}
		fmt.Fprintf(&b, "  Arguments JSON schema: %s\n", paramsJSON)
	}
	return b.String(), nil
}

// runFallbackToolLoop is generateWithLeakRetry's counterpart (tutor.go)
// for a role whose model CheckToolCalling found doesn't populate real
// tool_calls. Calls cm.Generate directly (bare -- no WithTools; the
// model only ever learns about tools from the text prompt the caller
// already prepended via prependToolsPrompt/renderToolCatalog), parses
// each reply, and either executes the named tool and loops with the
// result appended as an ephemeral schema.SystemMessage("Tool result:
// ...") -- not schema.RoleType("tool"), which requires a real
// ToolCallID an OpenAI-compatible backend validates against a preceding
// tool_calls entry that doesn't exist here -- or returns the first
// reply that doesn't parse as a call attempt at all as the final
// answer. Never returns raw tool-call JSON to the caller, matching
// generateWithLeakRetry's own guarantee.
//
// Guards identical-consecutive-call looping (comparing only against the
// immediately preceding round, not full history -- re-checking the same
// file after a few other calls in between is a legitimate pattern, not
// a bug) since react.Agent has no such guard of its own (confirmed via
// its source: MaxStep is its only bound).
func runFallbackToolLoop(ctx context.Context, cm model.ToolCallingChatModel, tools []tool.BaseTool, messages []*schema.Message, feed *activityFeed, activityCh chan<- []activityCall) (*schema.Message, error) {
	convo := append([]*schema.Message{}, messages...)

	var lastCall fallbackToolCall
	haveLastCall := false

	for round := 0; round < fallbackRoundCap; round++ {
		reply, err := cm.Generate(ctx, convo)
		if err != nil {
			return nil, err
		}

		call, matched, parseErr := parseFallbackToolCall(reply.Content)
		if !matched {
			return reply, nil
		}
		if parseErr != nil {
			convo = append(convo, schema.AssistantMessage(reply.Content, nil), schema.SystemMessage(
				fmt.Sprintf("Your last reply looked like a tool call but wasn't valid: %v. Reply again with ONLY a valid JSON object of the shape {\"name\": \"...\", \"arguments\": {...}}, or with your real answer if you don't need a tool.", parseErr),
			))
			continue
		}

		if haveLastCall && call.Name == lastCall.Name && string(call.Arguments) == string(lastCall.Arguments) {
			convo = append(convo, schema.AssistantMessage(reply.Content, nil), schema.SystemMessage(
				fmt.Sprintf("You already called %s with those exact arguments -- you have that result already. Use it, call a different tool, or give your real answer now.", call.Name),
			))
			continue
		}

		t := findTool(ctx, tools, call.Name)
		if t == nil {
			// Not recorded as lastCall: nothing was actually invoked, so
			// a repeat of this exact (unfound) call shouldn't be told
			// "you already have that result" -- it should just be told
			// again that the tool doesn't exist.
			convo = append(convo, schema.AssistantMessage(reply.Content, nil), schema.SystemMessage(
				fmt.Sprintf("There is no tool named %q. Check the tool catalog and try again, or give your real answer if you don't need a tool.", call.Name),
			))
			continue
		}
		lastCall, haveLastCall = call, true

		callID := fmt.Sprintf("fallback-%d", round)
		feed.started(callID, call.Name, truncateLine(string(call.Arguments), activityArgsPreviewMax))
		pushActivity(feed, activityCh)

		result, err := safeInvokeTool(ctx, t, string(call.Arguments))
		if err != nil {
			feed.failed(callID, truncateLine(err.Error(), activityResultPreviewMax))
			pushActivity(feed, activityCh)
			convo = append(convo, schema.AssistantMessage(reply.Content, nil), schema.SystemMessage(fmt.Sprintf("Tool error: %v", err)))
			continue
		}

		feed.finished(callID, truncateLine(result, activityResultPreviewMax))
		pushActivity(feed, activityCh)
		convo = append(convo, schema.AssistantMessage(reply.Content, nil), schema.SystemMessage("Tool result: "+result))
	}

	return schema.AssistantMessage(fallbackLoopExhaustedReply, nil), nil
}

// callRole dispatches one Generate call for a session role by its
// detected strategy -- through the real react.Agent
// (generateWithLeakRetry, tutor.go, unchanged) for nativeToolCalling, or
// runFallbackToolLoop directly against the role's raw chat model for
// jsonFallbackToolCalling. Called identically from startTurn's
// checkComprehension and real-turn branches (model.go), so a role's
// strategy is honored everywhere it can answer, not just the main turn.
func callRole(ctx context.Context, strategy toolCallingStrategy, agent *react.Agent, cm model.ToolCallingChatModel, tools []tool.BaseTool, messages []*schema.Message, feed *activityFeed, activityCh chan<- []activityCall, activityOpt agentopt.AgentOption) (*schema.Message, error) {
	if strategy == jsonFallbackToolCalling {
		return runFallbackToolLoop(ctx, cm, tools, messages, feed, activityCh)
	}
	return generateWithLeakRetry(ctx, agent, messages, activityOpt)
}
