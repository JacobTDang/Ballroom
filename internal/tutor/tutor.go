// Package tutor implements the in-container tutor agent: a tool-calling
// LLM loop (via eino's ReAct agent, github.com/cloudwego/eino) that can
// read the active solution file, the problem statement, the last test
// run's output, and the editor's cursor position, and highlight lines in
// the editor — rather than being handed a text dump and hoping it emits
// the right magic string in its reply (see docker/nvim/lua/ballroom_highlight.lua
// for the highlight rendering this drives, unchanged from the previous
// bash implementation).
package tutor

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
