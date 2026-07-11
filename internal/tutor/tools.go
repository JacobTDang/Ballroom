package tutor

import (
	"context"
	"fmt"

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

// buildTools assembles every tool the tutor agent has access to. Grows as
// later milestones add highlight_lines and read_cursor_position.
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
	return []tool.BaseTool{readSolution, readProblem, readTestOutput}, nil
}
