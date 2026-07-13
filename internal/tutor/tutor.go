// Package tutor implements the in-container tutor agent: a tool-calling
// LLM loop (via eino's ReAct agent, github.com/cloudwego/eino) that can
// read the active solution file, the problem statement, the last test
// run's output, and the editor's cursor position, and highlight lines in
// the editor — rather than being handed a text dump and hoping it emits
// the right magic string in its reply (see docker/nvim/lua/ballroom_highlight.lua
// for the highlight rendering this drives, unchanged from the previous
// bash implementation).
package tutor

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cloudwego/eino/components/model"
	agentopt "github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// Config describes one tutor invocation. All paths are as seen from
// inside the practice container.
type Config struct {
	// OllamaHost is the base URL of the Ollama server (e.g.
	// http://host.docker.internal:11434). Unused when Model is
	// OpenRouterModelPrefix-prefixed.
	OllamaHost string
	// Model is the Ollama model tag to use, or an
	// OpenRouterModelPrefix-prefixed OpenRouter model slug (see
	// agent.go's newChatModel). Must support the provider's structured
	// tool_calls response field — confirmed via cmd/tutor-spike that
	// qwen2.5-coder:7b does not (it emits tool-call-shaped JSON as plain
	// text content instead), while llama3.1:8b does.
	Model string
	// OrchestratorModel, when non-empty, enables per-turn routing: this
	// model decides (see decideHandoff) whether a turn needs Model's
	// coding-specialist attention or can be answered directly, and
	// always handles the comprehension check. Empty (the default) means
	// no routing at all — Model handles every turn by itself, identical
	// to this project's behavior before routing existed.
	OrchestratorModel string
	// APIKey authenticates OpenRouter requests when Model or
	// OrchestratorModel is OpenRouterModelPrefix-prefixed; unused
	// otherwise. One key authenticates every model on an OpenRouter
	// account, so this stays a single shared field even with two roles.
	APIKey string
	// Mode is the tutor_mode (syntax-only / hints-first / full-assist)
	// selecting the system prompt and whether the comprehension check
	// runs.
	Mode string
	// WorkDir is the exercise workspace directory, where the active
	// solution.*, problem.md, and (after a submit) the last test result
	// file are read from.
	WorkDir string
	// NvimSocket is the path to the editor pane's nvim --listen socket
	// (see docker/entrypoint.sh). Empty means highlighting/cursor-position
	// are unavailable; tools degrade gracefully rather than failing.
	NvimSocket string
	// MaxContextBytes caps how much of the solution file gets sent to the
	// model per read_solution_file call.
	MaxContextBytes int
}

// providerEndpoint returns a human-readable description of where a
// request for model actually goes, for display in the startup banner
// and "could not reach" error messages — ollamaHost (cfg.OllamaHost) is
// meaningless (empty, in practice) once model is
// OpenRouterModelPrefix-prefixed, since that path never touches it (see
// newChatModel, agent.go). A real bug found live: an OpenRouter session
// showed "connected to ." and "could not reach :" verbatim before this
// existed.
func providerEndpoint(model, ollamaHost string) string {
	if strings.HasPrefix(model, OpenRouterModelPrefix) {
		return "OpenRouter"
	}
	return ollamaHost
}

// decideHandoff asks orchestratorCM (a raw chat model, not a
// react.Agent — no tools needed for a classification call, and
// avoiding the agent graph means this can't hit the MaxStep issue
// react.Agent turns can) whether userMessage needs Model's coding-
// specialist attention. See prompts.go's routingInstruction for the
// exact instruction and the reasoning behind biasing toward handoff on
// anything unclear.
//
// Defaults to (true, err) on a request failure — same asymmetric-cost
// reasoning as an ambiguous reply: silently leaving a real code
// question with the orchestrator on a routing bug is a much worse
// failure than one unnecessary specialist call.
func decideHandoff(ctx context.Context, orchestratorCM model.ToolCallingChatModel, userMessage string) (bool, error) {
	reply, err := orchestratorCM.Generate(ctx, []*schema.Message{
		schema.SystemMessage(routingInstruction),
		schema.UserMessage(userMessage),
	})
	if err != nil {
		return true, err
	}
	return !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(reply.Content)), "NO"), nil
}

// Run drives one tutor session as a bubbletea program (model.go's
// tutorModel) — a scrolling conversation viewport above a
// dynamically-growing multi-line input box, replacing an earlier
// hand-rolled raw-ANSI anchored-box implementation
// (internal/tutor/scrollbox.go, deleted) after it repeatedly broke live:
// a typed message that wrapped past the box's single reserved content
// row landed on the terminal's actual last row, triggering the outer
// terminal's own native scroll — something a manually DECSTBM-confined
// region has no way to detect or recover from. bubbletea's renderer
// (already proven bug-free in this codebase's other TUI, internal/tui)
// owns all cursor positioning instead, so this class of bug can't
// recur — bubbles/textarea's real per-keystroke height recompute is
// what makes the input box's growth genuinely dynamic (see model.go's
// recomputeLayout/estimatedTextareaRows), not a fixed row count.
//
// tea.WithInput/tea.WithOutput let this run against any io.Reader/
// io.Writer, not just a real terminal — cmd/tutor-eval's grounding
// check and this package's own tests rely on that to drive a real
// session against a fake stdin.
func Run(ctx context.Context, cfg Config, stdin io.Reader, stdout, stderr io.Writer) error {
	m, err := newTutorModel(ctx, cfg, stderr)
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(m, tea.WithAltScreen(), tea.WithInput(stdin), tea.WithOutput(stdout)).Run()
	return err
}

// leakedToolCallPattern matches a reply that describes a tool call as
// literal JSON text in Content instead of the model making a real eino
// tool_calls invocation — e.g. `{"name": "read_solution_file", "parameters": {}}`
// showing up as part of the assistant's visible reply. prompts.go's
// toolsInstruction already tells the model never to do this, but a
// real-sample-size repro (12 sessions x 4 turns each against
// llama3.1:8b) found the leak rate climbs as a conversation goes on
// regardless: 0/12 on turn 1, up to 5/12 by turn 4. In every observed
// case the tool name and arguments were well-formed and real (never a
// hallucinated tool) — the model correctly decided what to call, it
// just wrote that decision out as text instead of actually calling it.
// Chasing this further with prompt wording alone risks repeating a
// known regression (see toolsInstruction's doc comment: a longer
// instruction fixed a related narration case but measurably hurt
// tool-calling on unrelated tools) — this is handled as model output
// the caller must detect and recover from instead.
var leakedToolCallPattern = regexp.MustCompile(`\{"name"\s*:\s*"[a-zA-Z_]+"`)

func looksLikeLeakedToolCall(content string) bool {
	return leakedToolCallPattern.MatchString(content)
}

// leakedToolCallRetryNote is appended as an ephemeral system message
// (never persisted to history, same pattern as turnMessages' hint-count
// note) when generateWithLeakRetry needs a second attempt.
const leakedToolCallRetryNote = "Your last reply described calling a tool by writing out JSON like {\"name\": ...} instead of actually calling it. Try again: call the tool for real this time. Your reply must not contain any JSON or description of what tool you're using — only your real answer, written after you actually have the tool's result."

// leakedToolCallFallbackReply is shown when even the retry leaks (or
// the retry itself can't reach Ollama) — the user must never be shown
// raw tool-call JSON, so this is an honest admission instead of a
// second garbled attempt.
const leakedToolCallFallbackReply = "Sorry, I wasn't able to get a grounded answer for that just now — could you try asking again?"

// turnFailedFallbackReply is shown in the chat (not just logged to
// stderr, which the user may never see once bubbletea's alt-screen has
// taken over the terminal) when a turn's Generate call fails outright — a bad
// host, a rejected request, an upstream rate limit, anything. A real bug
// found live: without this, a failed turn printed nothing to stdout at
// all and just silently moved on to the next prompt, so any transient
// failure looked exactly like the tutor being completely unresponsive
// rather than a one-off hiccup worth retrying.
const turnFailedFallbackReply = "Sorry, I couldn't reach the model just now — please try asking again."

// generateWithLeakRetry wraps agent.Generate with one retry: if the
// model leaks a fake tool-call JSON blob into its reply instead of
// making a real tool call (see leakedToolCallPattern), retries once
// with a corrective note, and falls back to an honest message rather
// than ever showing the user raw tool-call JSON. The leaked reply is
// never added to messages/history beyond this one retry attempt, so it
// can't bias later turns toward repeating the same pattern.
//
// opts is threaded into both Generate calls unchanged — in practice this
// is model.go's buildActivityChannelOption (see startTurn), so a retry's
// own tool calls stay visible in the same activity feed as the original
// attempt's, not a separate one. Variadic and additive: the only direct
// external caller (cmd/tutor-eval, via GenerateWithLeakRetry) passes
// none, unaffected.
func generateWithLeakRetry(ctx context.Context, agent *react.Agent, messages []*schema.Message, opts ...agentopt.AgentOption) (*schema.Message, error) {
	reply, err := agent.Generate(ctx, messages, opts...)
	if err != nil {
		return nil, err
	}
	if !looksLikeLeakedToolCall(reply.Content) {
		return reply, nil
	}

	retryMessages := append(append([]*schema.Message{}, messages...),
		schema.AssistantMessage(reply.Content, nil),
		schema.SystemMessage(leakedToolCallRetryNote),
	)
	retryReply, err := agent.Generate(ctx, retryMessages, opts...)
	if err == nil && !looksLikeLeakedToolCall(retryReply.Content) {
		return retryReply, nil
	}
	return schema.AssistantMessage(leakedToolCallFallbackReply, nil), nil
}

// GenerateWithLeakRetry is generateWithLeakRetry, exported for
// cmd/tutor-eval — evaluating tool-calling/mode behavior needs the
// tutor's real protected Generate path, not a bare agent.Generate call
// that can fail in ways a real session (which always goes through this)
// never would surface to a user. A real, newly-found gap: without this,
// cmd/tutor-eval's own runs showed a hints-first scenario failing ~25%
// of the time on a leaked fake tool-call JSON — a failure mode that
// can't actually reach a real user, since Run() always retries and
// falls back to an honest message instead, but the eval was reporting
// it as a raw scenario failure anyway.
func GenerateWithLeakRetry(ctx context.Context, agent *react.Agent, messages []*schema.Message) (*schema.Message, error) {
	return generateWithLeakRetry(ctx, agent, messages)
}

// turnMessages returns the messages appended to history for one real
// (non-comprehension-check) turn's outgoing request: an ephemeral
// hint-count note for hints-first mode, followed by the user's message.
// The note is never persisted into history (recomputed fresh from
// helpRequestCount each turn instead) — same ephemeral-context pattern
// already used elsewhere in this package.
//
// Stating the count directly, rather than leaving the model to infer
// "is this a first or repeat ask" from re-reading the conversation,
// fixes a real observed bug: the model, uncertain of its own state,
// asked the user to confirm whether this was their first question —
// exactly the kind of self-tracking a human tutor would never need to
// ask about out loud.
func turnMessages(mode string, helpRequestCount int, line string) []*schema.Message {
	if mode != exercise.TutorModeHintsFirst {
		return []*schema.Message{schema.UserMessage(line)}
	}

	var note string
	if helpRequestCount <= 1 {
		note = "(This is the user's 1st help request in this session. If they seem stuck on a specific point, give only a short nudge per your instructions — do not name the technique yet.)"
	} else {
		note = fmt.Sprintf("(This is help request #%d in this session. If the user is asking again about a point you already nudged them on, give a full, direct answer now, including names — use your own judgment, don't ask them to confirm.)", helpRequestCount)
	}
	return []*schema.Message{schema.SystemMessage(note), schema.UserMessage(line)}
}

// TurnMessages is turnMessages, exported for cmd/tutor-eval — a real
// gap found live: cmd/tutor-eval's runScenario built each turn as a
// bare schema.UserMessage(turn), never including the hint-count note
// hintsFirstPrompt (prompts.go) explicitly tells the model to trust.
// That made a hints-first scenario's own eval numbers (5-6/8) look far
// worse than real production reliability (15/16 in a direct real-Ollama
// repro through the actual Run() loop) — the eval wasn't testing the
// real code path.
func TurnMessages(mode string, helpRequestCount int, line string) []*schema.Message {
	return turnMessages(mode, helpRequestCount, line)
}
