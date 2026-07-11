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
			if runComprehensionCheck(ctx, agent, cfg.OllamaHost, line, &history, stdout, stderr) {
				continue
			}
			// Couldn't reach Ollama for the check — fall through and
			// handle this message normally below rather than silently
			// dropping it.
		}

		requestMessages := append(append([]*schema.Message{}, history...), schema.UserMessage(line))
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

// runComprehensionCheck issues one isolated Generate call asking the
// agent ONLY to restate the problem and ask clarifying questions —
// deliberately never including userFirstMessage, so there's nothing
// concrete for the model to answer instead of doing the check (port of
// tutor/chat.sh's run_comprehension_check; see its header comment for
// why this is enforced here rather than left to a prompt instruction the
// model might ignore). The model can call read_problem_statement itself
// if it needs to see the problem text — unlike the bash version, nothing
// is manually stuffed into this request.
//
// On success, appends (userFirstMessage, reply) to *history — the
// caller's real message is the persisted turn, so history reads
// naturally even though it was never sent to the model in this request —
// and returns true. Returns false if Ollama couldn't be reached, so the
// caller falls through to handling userFirstMessage normally instead of
// dropping it.
func runComprehensionCheck(ctx context.Context, agent *react.Agent, ollamaHost, userFirstMessage string, history *[]*schema.Message, stdout, stderr io.Writer) bool {
	checkMessages := append(append([]*schema.Message{}, (*history)...), schema.UserMessage(comprehensionCheckInstruction))

	reply, err := agent.Generate(ctx, checkMessages)
	if err != nil {
		fmt.Fprintf(stderr, "tutor: could not reach %s\n", ollamaHost)
		return false
	}

	fmt.Fprintln(stdout, reply.Content)

	*history = append(*history, schema.UserMessage(userFirstMessage), schema.AssistantMessage(reply.Content, nil))
	return true
}
