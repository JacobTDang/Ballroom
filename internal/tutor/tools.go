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

// buildTools assembles every tool the tutor agent has access to. Grows as
// later milestones add read_test_output, highlight_lines, and
// read_cursor_position.
func buildTools(cfg Config) ([]tool.BaseTool, error) {
	readSolution, err := newReadSolutionFileTool(cfg)
	if err != nil {
		return nil, err
	}
	readProblem, err := newReadProblemStatementTool(cfg)
	if err != nil {
		return nil, err
	}
	return []tool.BaseTool{readSolution, readProblem}, nil
}
