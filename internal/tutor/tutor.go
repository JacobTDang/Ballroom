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
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"

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
	// APIKey authenticates OpenRouter requests when Model is
	// OpenRouterModelPrefix-prefixed; unused otherwise.
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

// Run drives one tutor session: read a line from stdin, get the agent's
// reply, print it, repeat until stdin closes. Port of tutor/chat.sh's
// main() loop — but unlike that version, this never stuffs the solution
// file into every request; the model calls read_solution_file (and the
// other 4 tools) itself when it decides it needs to, which is the whole
// point of this rewrite over the old regex-directive approach.
//
// Where the real terminal supports it, input happens in a bordered box
// anchored at the bottom of the pane (internal/tutor/scrollbox.go), with
// conversation scrolling above it — the same technique Claude Code's
// own CLI uses. This was tried once before and reverted after breaking
// live (garbled, overlapping text) — root-caused via a real tmux repro
// (not just a raw pty, which never caught it) to inputBox.setup() never
// clearing the screen before drawing: a pane always has *something* on
// it already (at minimum the shell prompt from right before this
// program started), and a bare cursor-home doesn't erase that, so
// shorter new lines left old content dangling. Fixed in setup() itself
// (see its doc comment) and reverified through the real docker/tmux.conf,
// a split-pane layout, real typed input, pane switching, and a
// status-bar refresh boundary — not just a fresh empty terminal — before
// being wired back in here.
//
// newInputBox fails whenever stdin isn't a real terminal (tests,
// cmd/tutor-eval) or the terminal is too short; box is nil in that case
// and every use below is guarded so the session runs exactly as it did
// before this feature existed.
//
// The box's dimensions are captured once at construction — a resized
// terminal window would otherwise leave it operating on stale
// dimensions indefinitely, found live as a real bug (user report: "it
// happens when I readjust the size of the terminal"). watchResize
// subscribes to SIGWINCH; this loop drains that channel once per turn
// (a safe point between this goroutine's own writes) and calls
// box.reconfigure() there, rather than reconfiguring from a background
// goroutine, which could otherwise interleave with this loop's own
// unsynchronized stdout writes and corrupt output.
func Run(ctx context.Context, cfg Config, stdin io.Reader, stdout, stderr io.Writer) error {
	box, boxErr := newInputBox(stdout)
	var pendingResize chan os.Signal
	if boxErr == nil {
		defer box.close()
		var stopWatch func()
		pendingResize, stopWatch = watchResize()
		defer stopWatch()
	}

	fmt.Fprintf(stdout, "tutor (%s, mode=%s) — connected to %s. Ctrl-D to exit.\n", cfg.Model, cfg.Mode, cfg.OllamaHost)

	cm, err := newChatModel(ctx, cfg)
	if err != nil {
		return err
	}
	tools, err := buildTools(cfg)
	if err != nil {
		return err
	}
	agent, err := newAgent(ctx, cm, tools)
	if err != nil {
		return err
	}

	history := []*schema.Message{schema.SystemMessage(systemPromptForMode(cfg.Mode))}
	comprehensionCheckPending := wantsComprehensionCheck(cfg.Mode)
	helpRequestCount := 0

	// drainResize applies a pending resize signal (if any) to the box's
	// dimensions, non-blocking, a no-op when box is nil. Called at the
	// top of the loop (before showPrompt) AND again right before a
	// reply prints (both here and inside runComprehensionCheck) — a
	// real gap found live via user screenshot: agent.Generate can run
	// for many seconds against a real model, and a resize landing
	// during that wait was only ever drained on the *next* loop
	// iteration, so that turn's reply printed against the box's stale,
	// pre-resize dimensions — producing exactly the kind of garbled,
	// overlapping text a size mismatch between the confined scroll
	// region and the real terminal causes. The signal channel is
	// buffered (watchResize) so nothing is lost while waiting, only
	// delayed; this just drains it at more points instead of one.
	drainResize := func() {
		if box == nil {
			return
		}
		select {
		case <-pendingResize:
			box.reconfigure()
		default:
		}
	}

	scanner := bufio.NewScanner(stdin)
	for {
		if box != nil {
			drainResize()
			box.showPrompt()
		} else {
			fmt.Fprint(stdout, "> ")
		}
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		if box != nil {
			// The box's content row is about to be reused for the next
			// prompt, so nothing else preserves what was just typed —
			// echo it into the scrolling region as part of the
			// permanent conversation history. In the box == nil (plain
			// prompt) path, the terminal's own cooked-mode echo already
			// put "> line" in the real scrollback, so this must not
			// also run there or the line would print twice.
			box.returnToScroll()
			fmt.Fprintf(stdout, "> %s\n", line)
		}

		if comprehensionCheckPending {
			comprehensionCheckPending = false
			if runComprehensionCheck(ctx, agent, cfg.OllamaHost, cfg.WorkDir, line, &history, stdout, stderr, drainResize) {
				continue
			}
			// Couldn't reach Ollama for the check — fall through and
			// handle this message normally below rather than silently
			// dropping it.
		}

		helpRequestCount++
		requestMessages := append(append([]*schema.Message{}, history...), turnMessages(cfg.Mode, helpRequestCount, line)...)
		reply, err := generateWithLeakRetry(ctx, agent, requestMessages)
		if err != nil {
			// The real err is included, not just cfg.OllamaHost -- a
			// real bug found live: "could not reach" reads as a network
			// problem, but the actual cause is just as often a real API
			// rejection (e.g. Ollama returning 400 "does not support
			// tools" for a model that was picked without real
			// tool-calling support) that has nothing to do with
			// reachability at all. Showing only the generic message
			// sent a live debugging session chasing a nonexistent
			// Docker-networking problem instead of straight to the
			// actual cause.
			fmt.Fprintf(stderr, "tutor: could not reach %s: %v\n", cfg.OllamaHost, err)
			continue
		}

		drainResize()
		fmt.Fprintln(stdout, reply.Content)

		// Persist only the clean (user, assistant reply) pair — no
		// intermediate tool-call scaffolding — matching chat.sh's
		// history, which never accumulated the per-turn file context
		// either.
		history = append(history, schema.UserMessage(line), schema.AssistantMessage(reply.Content, nil))
	}

	fmt.Fprintln(stdout)
	return nil
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

// generateWithLeakRetry wraps agent.Generate with one retry: if the
// model leaks a fake tool-call JSON blob into its reply instead of
// making a real tool call (see leakedToolCallPattern), retries once
// with a corrective note, and falls back to an honest message rather
// than ever showing the user raw tool-call JSON. The leaked reply is
// never added to messages/history beyond this one retry attempt, so it
// can't bias later turns toward repeating the same pattern.
func generateWithLeakRetry(ctx context.Context, agent *react.Agent, messages []*schema.Message) (*schema.Message, error) {
	reply, err := agent.Generate(ctx, messages)
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
	retryReply, err := agent.Generate(ctx, retryMessages)
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

// runComprehensionCheck issues one isolated Generate call that responds
// to the user's real first message and asks the agent to restate the
// problem and ask clarifying questions (port of tutor/chat.sh's
// run_comprehension_check; see its header comment for why this is
// enforced here rather than left to a prompt instruction the model
// might ignore).
//
// The problem statement is injected directly as ephemeral context
// (readProblemStatement, not a tool call) rather than leaving the model
// to call read_problem_statement itself. Manual repro testing found the
// check's combined "call a tool, then restate, then ask questions" task
// only actually invoked the tool 40-60% of the time regardless of
// instruction wording/length — well below the ~100% reliability normal
// single-purpose tool-calling turns get (see cmd/tutor-eval) — and on
// failure the model would hallucinate a plausible-looking but entirely
// fabricated tool result instead. This is the one place in the package
// that still stuffs context directly rather than using a real tool
// call, deliberately: it's the single highest-stakes moment (every
// session's first exchange) to get reliably right, and with the content
// provided directly there's nothing left to call.
//
// userFirstMessage is included as a real user turn (an earlier version
// deliberately excluded it) — a real bug found live: excluding it meant
// literally any first message, including a plain "hi", got the exact
// same canned restate-and-ask-questions reply with no acknowledgment of
// what the user actually said. comprehensionCheckInstruction (prompts.go)
// now tells the model to respond to it. Routed through
// generateWithLeakRetry for the same reason every other turn is: this
// call can leak fake tool-call JSON exactly like any other
// (cmd/tutor-eval's grounding check already tested for this here, but
// nothing actually protected it before).
//
// On success, appends (userFirstMessage, reply) to *history and returns
// true. Returns false if Ollama couldn't be reached, so the caller
// falls through to handling userFirstMessage normally instead of
// dropping it.
//
// drainResize is called right before the reply prints, same as Run's
// own main-loop call — this check's own agent.Generate call is exactly
// as susceptible to the resize-during-generation gap drainResize fixes
// (see its doc comment in Run), and it happens to be the very first
// Generate call of a session, so it's a likely place for a user to
// resize their terminal while waiting.
func runComprehensionCheck(ctx context.Context, agent *react.Agent, ollamaHost, workDir, userFirstMessage string, history *[]*schema.Message, stdout, stderr io.Writer, drainResize func()) bool {
	checkMessages := append([]*schema.Message{}, (*history)...)
	if problem := readProblemStatement(workDir); problem != "" {
		checkMessages = append(checkMessages, schema.SystemMessage("The exercise's problem statement:\n\n"+problem))
	}
	checkMessages = append(checkMessages, schema.SystemMessage(comprehensionCheckInstruction), schema.UserMessage(userFirstMessage))

	reply, err := generateWithLeakRetry(ctx, agent, checkMessages)
	if err != nil {
		// See Run's identical fix for why the real err is included, not
		// just ollamaHost.
		fmt.Fprintf(stderr, "tutor: could not reach %s: %v\n", ollamaHost, err)
		return false
	}

	drainResize()
	fmt.Fprintln(stdout, reply.Content)

	*history = append(*history, schema.UserMessage(userFirstMessage), schema.AssistantMessage(reply.Content, nil))
	return true
}
