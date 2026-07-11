package tutor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

// noInput is used for tools that take no model-supplied arguments — they
// always operate on cfg (the exercise workspace already known from the
// session), not on anything the model chooses.
type noInput struct{}

type readFileOutput struct {
	Content string `json:"content" jsonschema:"description=the file's current contents, or a short note if it does not exist yet"`
}

// newReadSolutionFileTool lets the model read the exercise's active
// solution file on demand instead of having it stuffed into every
// request whether needed or not. Port of tutor/chat.sh's per-turn
// build_file_context injection, now an explicit tool call.
func newReadSolutionFileTool(cfg Config) (tool.InvokableTool, error) {
	fn := func(_ context.Context, _ noInput) (readFileOutput, error) {
		content := buildFileContext(cfg.WorkDir, cfg.MaxContextBytes)
		if content == "" {
			return readFileOutput{Content: "(no solution file exists yet)"}, nil
		}
		return readFileOutput{Content: content}, nil
	}
	t, err := utils.InferTool("read_solution_file", "Read the current contents of the exercise's solution file the user is editing right now.", fn)
	if err != nil {
		return nil, fmt.Errorf("tutor: infer read_solution_file tool: %w", err)
	}
	return t, nil
}

// newReadProblemStatementTool lets the model read the exercise's
// problem.md (statement, examples, constraints) on demand.
func newReadProblemStatementTool(cfg Config) (tool.InvokableTool, error) {
	fn := func(_ context.Context, _ noInput) (readFileOutput, error) {
		content := readProblemStatement(cfg.WorkDir)
		if content == "" {
			return readFileOutput{Content: "(no problem statement available)"}, nil
		}
		return readFileOutput{Content: content}, nil
	}
	t, err := utils.InferTool("read_problem_statement", "Read the exercise's problem statement: description, examples, and constraints.", fn)
	if err != nil {
		return nil, fmt.Errorf("tutor: infer read_problem_statement tool: %w", err)
	}
	return t, nil
}

type readTestOutputOutput struct {
	Available   bool   `json:"available" jsonschema:"description=whether a test result exists yet"`
	Result      string `json:"result,omitempty" jsonschema:"description=pass or fail"`
	Output      string `json:"output,omitempty" jsonschema:"description=the raw test command output"`
	TestCommand string `json:"test_command,omitempty" jsonschema:"description=the shell command that was run"`
}

// newReadTestOutputTool lets the model read the result of the user's
// most recent `ballroom submit` — see internal/session's
// writeLastTestResult, which produces the file this reads.
func newReadTestOutputTool(cfg Config) (tool.InvokableTool, error) {
	fn := func(_ context.Context, _ noInput) (readTestOutputOutput, error) {
		result, ok, err := readLastTestResult(cfg.WorkDir)
		if err != nil {
			return readTestOutputOutput{}, fmt.Errorf("read test output: %w", err)
		}
		if !ok {
			return readTestOutputOutput{Available: false}, nil
		}
		return readTestOutputOutput{
			Available:   true,
			Result:      result.Result,
			Output:      result.Output,
			TestCommand: result.TestCommand,
		}, nil
	}
	t, err := utils.InferTool("read_test_output", "Read the pass/fail result and raw output from the user's most recent test submission, if they've submitted yet.", fn)
	if err != nil {
		return nil, fmt.Errorf("tutor: infer read_test_output tool: %w", err)
	}
	return t, nil
}

type highlightLinesInput struct {
	File  string `json:"file" jsonschema:"description=filename to highlight (e.g. solution.go) -- basename is enough, matches what read_solution_file/read_cursor_position return"`
	Start int    `json:"start" jsonschema:"description=1-indexed start line (inclusive)"`
	End   int    `json:"end" jsonschema:"description=1-indexed end line (inclusive) -- same as start to highlight a single line"`
	Note  string `json:"note" jsonschema:"description=a short note to attach to the highlighted lines, shown in the user's editor"`
}

type highlightLinesOutput struct {
	Status string `json:"status" jsonschema:"description=ok on success, or a note that no editor is currently attached"`
}

// newHighlightLinesTool lets the model highlight a range of lines with a
// note directly in the user's editor pane (docker/nvim/lua/ballroom_highlight.lua),
// instead of the previous approach of asking the model to emit a magic
// text string in its reply for a regex to scrape out. file/note are
// model-controlled, so escaping (nvimrpc.go's escapeVimSingleQuoted) is
// load-bearing here — see nvimrpc_test.go's live-nvim injection tests.
func newHighlightLinesTool(cfg Config) (tool.InvokableTool, error) {
	fn := func(ctx context.Context, in highlightLinesInput) (highlightLinesOutput, error) {
		expr := highlightExpr(in.File, in.Start, in.End, in.Note)
		out, err := remoteExpr(ctx, cfg.NvimSocket, expr)
		if err != nil {
			return highlightLinesOutput{}, err
		}
		if out == "" {
			return highlightLinesOutput{Status: "no editor is currently attached, highlight not shown"}, nil
		}
		if strings.HasPrefix(out, "ballroom_highlight error:") {
			return highlightLinesOutput{}, fmt.Errorf("%s", out)
		}
		return highlightLinesOutput{Status: out}, nil
	}
	t, err := utils.InferTool("highlight_lines", "Highlight a range of lines in the user's solution file with a short note, visible directly in their editor.", fn)
	if err != nil {
		return nil, fmt.Errorf("tutor: infer highlight_lines tool: %w", err)
	}
	return t, nil
}

type cursorPositionOutput struct {
	Available  bool   `json:"available" jsonschema:"description=whether an editor is currently attached"`
	File       string `json:"file,omitempty" jsonschema:"description=basename of the file currently focused"`
	Line       int    `json:"line,omitempty" jsonschema:"description=1-indexed cursor line"`
	Col        int    `json:"col,omitempty" jsonschema:"description=1-indexed cursor column"`
	TotalLines int    `json:"total_lines,omitempty" jsonschema:"description=total lines in the focused file"`
}

// newReadCursorPositionTool lets the model see roughly where the user is
// currently looking/working in the editor, rather than only ever seeing
// the whole file dumped every turn. The underlying expression
// (nvimrpc.go's cursorPositionExpr) is static, not model-controlled, so
// unlike highlight_lines there's no injection surface here.
func newReadCursorPositionTool(cfg Config) (tool.InvokableTool, error) {
	fn := func(ctx context.Context, _ noInput) (cursorPositionOutput, error) {
		out, err := remoteExpr(ctx, cfg.NvimSocket, cursorPositionExpr())
		if err != nil {
			return cursorPositionOutput{}, err
		}
		if out == "" {
			return cursorPositionOutput{Available: false}, nil
		}
		var pos struct {
			File       string `json:"file"`
			Line       int    `json:"line"`
			Col        int    `json:"col"`
			TotalLines int    `json:"total_lines"`
		}
		if err := json.Unmarshal([]byte(out), &pos); err != nil {
			return cursorPositionOutput{}, fmt.Errorf("tutor: parse cursor position: %w", err)
		}
		return cursorPositionOutput{
			Available:  true,
			File:       pos.File,
			Line:       pos.Line,
			Col:        pos.Col,
			TotalLines: pos.TotalLines,
		}, nil
	}
	t, err := utils.InferTool("read_cursor_position", "See where the user's cursor currently is in the editor: filename, line, and column.", fn)
	if err != nil {
		return nil, fmt.Errorf("tutor: infer read_cursor_position tool: %w", err)
	}
	return t, nil
}

// buildTools assembles every tool the tutor agent has access to.
func buildTools(cfg Config) ([]tool.BaseTool, error) {
	readSolution, err := newReadSolutionFileTool(cfg)
	if err != nil {
		return nil, err
	}
	readProblem, err := newReadProblemStatementTool(cfg)
	if err != nil {
		return nil, err
	}
	readTestOutput, err := newReadTestOutputTool(cfg)
	if err != nil {
		return nil, err
	}
	highlightLines, err := newHighlightLinesTool(cfg)
	if err != nil {
		return nil, err
	}
	readCursorPosition, err := newReadCursorPositionTool(cfg)
	if err != nil {
		return nil, err
	}

	// Wrap every tool so a failure (malformed model-generated arguments,
	// an out-of-bounds highlight range, an unreachable nvim socket that
	// somehow still errors, ...) becomes a string result fed back to the
	// model instead of aborting the whole turn — replaces
	// tutor/chat.sh's process_highlights bash-level fallthrough case.
	raw := []tool.BaseTool{readSolution, readProblem, readTestOutput, highlightLines, readCursorPosition}
	wrapped := make([]tool.BaseTool, len(raw))
	for i, t := range raw {
		wrapped[i] = utils.WrapToolWithErrorHandler(t, toolErrorHandler)
	}
	return wrapped, nil
}

// toolErrorHandler is the shared ErrorHandler every tool is wrapped
// with in buildTools.
func toolErrorHandler(_ context.Context, err error) string {
	return fmt.Sprintf("tool error: %v", err)
}

// BuildTools is buildTools, exported for cmd/tutor-eval — evaluating
// whether the model actually calls tools correctly needs the tutor's
// real tool implementations, not stand-ins.
func BuildTools(cfg Config) ([]tool.BaseTool, error) {
	return buildTools(cfg)
}
