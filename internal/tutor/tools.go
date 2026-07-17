package tutor

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
// build_file_context injection, now an explicit tool call. Also
// advances the shared snapshot so a later read_solution_diff answers
// "since you last looked", whichever tool looked.
func newReadSolutionFileTool(cfg Config, snap *solutionSnapshot) (tool.InvokableTool, error) {
	fn := func(_ context.Context, _ noInput) (readFileOutput, error) {
		snap.swap(rawSolutionContent(cfg.WorkDir, cfg.MaxContextBytes))
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

// newReadSolutionDiffTool answers "what changed since you last looked"
// as a unified diff -- much cheaper context than re-reading the whole
// file, and it focuses the model on the edit instead of the file.
func newReadSolutionDiffTool(cfg Config, snap *solutionSnapshot) (tool.InvokableTool, error) {
	fn := func(_ context.Context, _ noInput) (readFileOutput, error) {
		current := rawSolutionContent(cfg.WorkDir, cfg.MaxContextBytes)
		previous := snap.swap(current)
		if current == "" {
			return readFileOutput{Content: "(no solution file exists yet)"}, nil
		}
		diff := diffUnified(previous, current)
		if diff == "" {
			return readFileOutput{Content: "(no changes since the last read)"}, nil
		}
		return readFileOutput{Content: diff}, nil
	}
	t, err := utils.InferTool("read_solution_diff", "Show what changed in the user's solution file since you last looked at it (via read_solution_file or this tool), as a unified diff. Prefer this over re-reading the whole file when the user says they changed or added something.", fn)
	if err != nil {
		return nil, fmt.Errorf("tutor: infer read_solution_diff tool: %w", err)
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

// newReadGradingRubricTool lets the model read a design session's
// grading rubric -- hidden until the user's M-q submit reveals it into
// the workspace (the same reveal mechanic that delivers hidden tests
// for coding exercises). Present in every session's tool set: for
// coding exercises there is simply never a rubric.md, so it stays
// inert, and keeping buildTools kind-agnostic means no plumbing of the
// exercise kind into the tutor.
func newReadGradingRubricTool(cfg Config) (tool.InvokableTool, error) {
	fn := func(_ context.Context, _ noInput) (readFileOutput, error) {
		content := readRubric(cfg.WorkDir)
		if content == "" {
			return readFileOutput{Content: "(no rubric available yet -- it is revealed after the user submits with M-q)"}, nil
		}
		return readFileOutput{Content: content}, nil
	}
	t, err := utils.InferTool("read_grading_rubric", "Read the design session's grading rubric, if the user has submitted and revealed it yet.", fn)
	if err != nil {
		return nil, fmt.Errorf("tutor: infer read_grading_rubric tool: %w", err)
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

// flexibleInt unmarshals from either a JSON number or a JSON string
// containing digits. Found via manual testing (M8 cutover): llama3.1:8b
// intermittently emits tool-call arguments with an integer field quoted
// as a string (e.g. "end":"4" instead of "end":4) — roughly 1 in 5-6
// highlight_lines calls in isolated repro testing. Go's encoding/json
// rejects that as a type mismatch by default, which without this would
// surface as a spurious tool error even though the call was otherwise
// well-formed. Confirmed via cmd/tutor-debug-repro (throwaway, not
// committed) that this fix alone takes the failure rate to 0/16 across
// both a fresh conversation and one with prior history — the earlier
// hypothesis that conversation history was the cause was wrong; this
// was the real root cause the whole time.
type flexibleInt int

func (fi *flexibleInt) UnmarshalJSON(data []byte) error {
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*fi = flexibleInt(n)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("flexibleInt: not a number or numeric string: %s", data)
	}
	parsed, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("flexibleInt: %q is not a valid integer: %w", s, err)
	}
	*fi = flexibleInt(parsed)
	return nil
}

type highlightLinesInput struct {
	File  string      `json:"file" jsonschema:"description=filename to highlight (e.g. solution.go) -- basename is enough, matches what read_solution_file/read_cursor_position return"`
	Start flexibleInt `json:"start" jsonschema:"description=1-indexed start line (inclusive) -- use the exact line numbers shown by read_solution_file, do not count lines yourself"`
	End   flexibleInt `json:"end" jsonschema:"description=1-indexed end line (inclusive) -- same as start to highlight a single line; use the exact line numbers shown by read_solution_file, do not count lines yourself"`
	Note  string      `json:"note" jsonschema:"description=a short note to attach to the highlighted lines, shown in the user's editor"`
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
		expr := highlightExpr(in.File, int(in.Start), int(in.End), in.Note)
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
	// The diff tool's session-start baseline: taken here, once, so the
	// first read_solution_diff shows everything changed since launch.
	snap := &solutionSnapshot{last: rawSolutionContent(cfg.WorkDir, cfg.MaxContextBytes)}
	readSolution, err := newReadSolutionFileTool(cfg, snap)
	if err != nil {
		return nil, err
	}
	readSolutionDiff, err := newReadSolutionDiffTool(cfg, snap)
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
	readGradingRubric, err := newReadGradingRubricTool(cfg)
	if err != nil {
		return nil, err
	}

	// Wrap every tool so a failure (malformed model-generated arguments,
	// an out-of-bounds highlight range, an unreachable nvim socket that
	// somehow still errors, ...) becomes a string result fed back to the
	// model instead of aborting the whole turn — replaces
	// tutor/chat.sh's process_highlights bash-level fallthrough case.
	raw := []tool.BaseTool{readSolution, readSolutionDiff, readProblem, readTestOutput, highlightLines, readCursorPosition, readGradingRubric}
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
