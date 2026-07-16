package tutor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	agentopt "github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// streamingEnabled reports whether a role backed by modelName should
// stream its final replies (agent.Stream) instead of blocking on
// Generate. Default: on only for OpenRouter models -- Ollama drops
// tool_calls entirely when streaming (cmd/tutor-spike's original
// finding, the reason this package used Generate exclusively until
// now), so Ollama stays on the blocking path. TUTOR_STREAM=on|off
// overrides the default either way; any other value is ignored rather
// than guessed at.
//
// Read per call (not cached at startup) so tests can flip it with
// t.Setenv; a real session's environment never changes mid-run.
func streamingEnabled(modelName string) bool {
	switch os.Getenv("TUTOR_STREAM") {
	case "on":
		return true
	case "off":
		return false
	}
	return strings.HasPrefix(modelName, OpenRouterModelPrefix)
}

// streamHoldBackRunes is how long a "{"- or "<"-prefixed reply must
// grow before it's shown: leakedToolCallPattern needs `{"name": "x` --
// 11 runes -- to fire, and the <tool_call tag needs 10, so by 16
// accumulated runes a leak has definitively either matched (suppressed
// permanently below) or can't be one, and legitimate content that
// merely starts with a brace (a quoted JSON snippet) gets displayed.
const streamHoldBackRunes = 16

// heldBackFromDisplay reports whether accumulated streamed text is
// still too short to distinguish a leaked tool call from legitimate
// content that happens to start like one.
func heldBackFromDisplay(text string) bool {
	if !strings.HasPrefix(text, "{") && !strings.HasPrefix(text, "<") {
		return false
	}
	return utf8.RuneCountInString(text) < streamHoldBackRunes
}

// errEmptyStream marks a stream that ended cleanly without yielding a
// single chunk -- the streaming shape of the empty-choices failure
// generateWithEmptyChoicesRetry (fallback.go) retries on the blocking
// path, and retried with the same backoff by streamWithLeakGuard.
var errEmptyStream = errors.New("tutor: model streamed an empty reply")

// streamToolCallDecisionRunes is how much narration
// windowedStreamToolCallChecker reads past before concluding a streamed
// model reply really is a final answer and not a preamble to a tool
// call. Cost of a bigger window: the caller's first painted chunk is
// delayed until the checker decides, so this trades display latency
// (one sentence's worth, a few chunks) against catching the
// narrate-then-call pattern. The live case that set it: ~58 runes of
// preamble before the tool call.
const streamToolCallDecisionRunes = 200

// windowedStreamToolCallChecker replaces eino's default
// firstChunkStreamToolCallChecker, which decides "no tool call" at the
// very first non-empty content chunk -- a real bug found live: a model
// that narrates before calling ("I'll read your solution file to see
// what you're working on." ... tool_calls) had its tool call silently
// dropped, turning a grounded answer into pure narration. This keeps
// reading until it sees a tool call, hits EOF, or has read
// streamToolCallDecisionRunes of content; a call hidden past the window
// is still recovered -- see streamWithLeakGuard's unexecuted-tool-calls
// safety net.
func windowedStreamToolCallChecker(_ context.Context, sr *schema.StreamReader[*schema.Message]) (bool, error) {
	defer sr.Close()
	runes := 0
	for {
		msg, err := sr.Recv()
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if len(msg.ToolCalls) > 0 {
			return true, nil
		}
		runes += utf8.RuneCountInString(msg.Content)
		if runes > streamToolCallDecisionRunes {
			return false, nil
		}
	}
}

// streamOnce drives one agent.Stream call to completion: every chunk is
// accumulated for the final reassembled reply, and onText receives the
// full accumulated text each time it grows displayably. Two display
// gates, both about never painting tool-call JSON:
//
//   - hold-back: text starting with "{"/"<" stays hidden until
//     streamHoldBackRunes runes have accumulated (see heldBackFromDisplay).
//   - permanent suppression: the moment looksLikeLeakedToolCall matches
//     the accumulated text, display stops for good -- the rest of the
//     stream is still drained (the caller needs the full leaked content
//     to build its corrective retry), it just never paints.
//
// eino's react graph consumes tool-call rounds internally on the Stream
// path exactly as on Generate (verified in v0.9.12's react.go: the
// stream branch checker peeks each model reply's first chunks and
// routes tool calls back into the graph) -- the StreamReader returned
// here yields only the final answer's chunks, and activity callbacks
// keep firing unchanged.
func streamOnce(ctx context.Context, agent *react.Agent, messages []*schema.Message, onText func(string), opts ...agentopt.AgentOption) (*schema.Message, error) {
	sr, err := agent.Stream(ctx, messages, opts...)
	if err != nil {
		return nil, err
	}
	defer sr.Close()

	var chunks []*schema.Message
	var accum strings.Builder
	lastShown := ""
	suppressed := false
	for {
		chunk, err := sr.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
		accum.WriteString(chunk.Content)

		if suppressed || onText == nil {
			continue
		}
		text := accum.String()
		if looksLikeLeakedToolCall(text) {
			suppressed = true
			continue
		}
		if heldBackFromDisplay(text) || text == lastShown {
			continue
		}
		lastShown = text
		onText(text)
	}

	if len(chunks) == 0 {
		return nil, errEmptyStream
	}
	full, err := schema.ConcatMessages(chunks)
	if err != nil {
		return nil, fmt.Errorf("tutor: concat streamed reply: %w", err)
	}
	return full, nil
}

// streamWithLeakGuard is generateWithLeakRetry's streaming counterpart:
// the same two protections that wrap every blocking Generate call --
// the leaked-tool-call retry and the empty-choices backoff -- applied
// around agent.Stream, plus progressive display via onText. The final
// returned reply is the reassembled full message, so callers persist
// exactly what a blocking call would have returned.
//
// The leak retry itself is deliberately non-streaming (retryLeakedReply's
// plain Generate): a model that just leaked is the worst candidate for
// optimistic progressive display, and the blocking path's guarantees
// are already proven.
func streamWithLeakGuard(ctx context.Context, agent *react.Agent, messages []*schema.Message, onText func(string), opts ...agentopt.AgentOption) (*schema.Message, error) {
	var lastErr error
	for attempt := 0; attempt <= emptyChoicesMaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * emptyChoicesRetryBackoff):
			}
		}
		reply, err := streamOnce(ctx, agent, messages, onText, opts...)
		if err == nil {
			if len(reply.ToolCalls) > 0 {
				// The reassembled reply carries tool calls that never
				// executed: the model hid its call past the checker's
				// decision window, so the graph routed the stream out as a
				// final answer. Whatever painted so far is ungrounded
				// narration -- discard it and redo the whole turn on the
				// blocking path, which inspects the complete reply and
				// executes tools regardless of narration.
				return generateWithLeakRetry(ctx, agent, messages, opts...)
			}
			if !looksLikeLeakedToolCall(reply.Content) {
				return reply, nil
			}
			return retryLeakedReply(ctx, agent, messages, reply.Content, opts...)
		}
		// Only the two empty-reply shapes are retried -- a clean EOF
		// with zero chunks, or the provider's own "empty choices" error
		// -- both the free-tier rate-limit signature. Anything else
		// propagates unchanged; retrying unknown failures would mask
		// real bugs.
		if !errors.Is(err, errEmptyStream) && !isEmptyChoicesErr(err) {
			return nil, err
		}
		lastErr = err
	}
	return nil, fmt.Errorf("%w (the model streamed an empty response %d times in a row -- free-tier models rate-limit quickly, wait a few seconds and try again)", lastErr, emptyChoicesMaxRetries+1)
}

// pushLatestStreamText publishes text on ch with replace-latest
// semantics: ch is buffered (capacity 1), and when the UI hasn't
// consumed the previous snapshot yet, that stale snapshot is dropped in
// favor of this newer one -- each send is the full accumulated text, so
// intermediate snapshots are strictly redundant. Never blocks the
// streaming goroutine on a slow UI.
func pushLatestStreamText(ch chan string, text string) {
	for {
		select {
		case ch <- text:
			return
		default:
			select {
			case <-ch:
			default:
			}
		}
	}
}
