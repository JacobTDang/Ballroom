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

	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// Config describes one tutor invocation. All paths are as seen from
// inside the practice container.
type Config struct {
	// OllamaHost is the base URL of the Ollama server (e.g.
	// http://host.docker.internal:11434).
	OllamaHost string
	// Model is the Ollama model tag to use. Must support Ollama's
	// structured tool_calls response field — confirmed via
	// cmd/tutor-spike that qwen2.5-coder:7b does not (it emits
	// tool-call-shaped JSON as plain text content instead), while
	// llama3.1:8b does.
	Model string
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
func Run(ctx context.Context, cfg Config, stdin io.Reader, stdout, stderr io.Writer) error {
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

	scanner := bufio.NewScanner(stdin)
	for {
		fmt.Fprint(stdout, "> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "" {
			continue
		}

		if comprehensionCheckPending {
			comprehensionCheckPending = false
			if runComprehensionCheck(ctx, agent, cfg.OllamaHost, cfg.WorkDir, line, &history, stdout, stderr) {
				continue
			}
			// Couldn't reach Ollama for the check — fall through and
			// handle this message normally below rather than silently
			// dropping it.
		}

		helpRequestCount++
		requestMessages := append(append([]*schema.Message{}, history...), turnMessages(cfg.Mode, helpRequestCount, line)...)
		reply, err := agent.Generate(ctx, requestMessages)
		if err != nil {
			fmt.Fprintf(stderr, "tutor: could not reach %s\n", cfg.OllamaHost)
			continue
		}

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

// runComprehensionCheck issues one isolated Generate call asking the
// agent ONLY to restate the problem and ask clarifying questions —
// deliberately never including userFirstMessage, so there's nothing
// concrete for the model to answer instead of doing the check (port of
// tutor/chat.sh's run_comprehension_check; see its header comment for
// why this is enforced here rather than left to a prompt instruction the
// model might ignore).
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
// On success, appends (userFirstMessage, reply) to *history — the
// caller's real message is the persisted turn, so history reads
// naturally even though it was never sent to the model in this request —
// and returns true. Returns false if Ollama couldn't be reached, so the
// caller falls through to handling userFirstMessage normally instead of
// dropping it.
func runComprehensionCheck(ctx context.Context, agent *react.Agent, ollamaHost, workDir, userFirstMessage string, history *[]*schema.Message, stdout, stderr io.Writer) bool {
	checkMessages := append([]*schema.Message{}, (*history)...)
	if problem := readProblemStatement(workDir); problem != "" {
		checkMessages = append(checkMessages, schema.SystemMessage("The exercise's problem statement:\n\n"+problem))
	}
	checkMessages = append(checkMessages, schema.UserMessage(comprehensionCheckInstruction))

	reply, err := agent.Generate(ctx, checkMessages)
	if err != nil {
		fmt.Fprintf(stderr, "tutor: could not reach %s\n", ollamaHost)
		return false
	}

	fmt.Fprintln(stdout, reply.Content)

	*history = append(*history, schema.UserMessage(userFirstMessage), schema.AssistantMessage(reply.Content, nil))
	return true
}
