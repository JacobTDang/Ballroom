package tutor

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"
)

// gradingInstruction drives GradeDesign's single-purpose model call.
// The rubric and the candidate's design are inlined as context rather
// than fetched via tools -- this codebase measured combined
// "call a tool, then reason" requests completing the tool call only
// 40-60% of the time (see comprehensionCheckInstruction's doc comment),
// and a grading verdict is too load-bearing for that failure rate.
// The rigid VERDICT first-line contract exists so the caller can parse
// pass/fail mechanically and fail loud when the model doesn't comply.
// Category-neutral on purpose: the same grading path serves
// system-design mocks and behavioral STAR stories -- the rubric
// carries everything category-specific.
const gradingInstruction = "You are grading a mock interview answer against a rubric. Judge only what is written in the answer -- do not fill gaps with charitable assumptions. Your reply MUST start with a first line that is exactly 'VERDICT: pass' or 'VERDICT: fail' (a passing answer is adequate or better on every rubric dimension), followed by a short assessment of each rubric dimension with specific evidence from the answer."

// GradeDesign grades a design session's solution.md against the
// revealed rubric.md with one model call, returning the parsed verdict
// (tracker.ResultPass/"pass" or "fail") and the model's per-dimension
// summary. Every failure -- missing files, transport errors, a reply
// without a parseable verdict -- returns an error so the caller can
// fall back to explicit self-assessment instead of recording a guess.
func GradeDesign(ctx context.Context, cfg Config) (verdict, summary string, err error) {
	rubric := readRubric(cfg.WorkDir)
	if rubric == "" {
		return "", "", fmt.Errorf("tutor: grade design: no rubric.md in %s (has the submit reveal run?)", cfg.WorkDir)
	}
	solution := buildFileContext(cfg.WorkDir, cfg.MaxContextBytes)
	if solution == "" {
		return "", "", fmt.Errorf("tutor: grade design: no solution file in %s", cfg.WorkDir)
	}

	cm, err := newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
	if err != nil {
		return "", "", fmt.Errorf("tutor: grade design: %w", err)
	}

	messages := []*schema.Message{
		schema.SystemMessage(gradingInstruction),
		schema.UserMessage("Rubric:\n\n" + rubric + "\n\nThe candidate's design:\n\n" + solution),
	}
	reply, err := generateWithEmptyChoicesRetry(ctx, cm, messages)
	if err != nil {
		return "", "", fmt.Errorf("tutor: grade design: %w", err)
	}

	verdict, err = parseVerdict(reply.Content)
	if err != nil {
		return "", "", err
	}
	return verdict, strings.TrimSpace(reply.Content), nil
}

// parseVerdict finds the VERDICT line in a grading reply. Tolerant of
// case and of prose before the line (small models rarely obey format
// instructions perfectly), but never defaults: no verdict, no result.
func parseVerdict(content string) (string, error) {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.ToLower(strings.TrimSpace(line))
		if !strings.HasPrefix(trimmed, "verdict:") {
			continue
		}
		switch strings.TrimSpace(strings.TrimPrefix(trimmed, "verdict:")) {
		case "pass":
			return "pass", nil
		case "fail":
			return "fail", nil
		}
	}
	return "", fmt.Errorf("tutor: grade design: reply has no parseable VERDICT line: %.120q", content)
}
