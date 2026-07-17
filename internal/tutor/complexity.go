package tutor

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"
)

// complexityInstruction drives CheckComplexity's single-purpose model
// call. Same inline-the-context design as gradingInstruction: the
// solution is small and the verdict is the whole point, so no tools.
const complexityInstruction = "You are checking a student's claimed time and space complexity for their passing solution. Judge the actual code, not the claim. Reply with a first line that is exactly 'AGREE' or 'DISAGREE', then at most three short lines: the correct time and space complexity and a one-sentence reason grounded in the code (name the loop, data structure, or recursion that determines it)."

// CheckComplexity asks the model whether the user's claimed complexity
// matches their actual solution -- the post-pass quiz cmd/ballroom
// wires into session.Submit. One blocking call on cfg.Model; every
// failure returns an error so the caller can degrade to a notice
// instead of blocking the submit.
func CheckComplexity(ctx context.Context, cfg Config, claim string) (string, error) {
	solution := buildFileContext(cfg.WorkDir, cfg.MaxContextBytes)
	if solution == "" {
		return "", fmt.Errorf("tutor: complexity check: no solution file in %s", cfg.WorkDir)
	}

	cm, err := newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
	if err != nil {
		return "", fmt.Errorf("tutor: complexity check: %w", err)
	}

	messages := []*schema.Message{
		schema.SystemMessage(complexityInstruction),
		schema.UserMessage("The solution:\n\n" + solution + "\n\nThe student's claimed complexity: " + claim),
	}
	reply, err := generateWithEmptyChoicesRetry(ctx, cm, messages)
	if err != nil {
		return "", fmt.Errorf("tutor: complexity check: %w", err)
	}
	return strings.TrimSpace(reply.Content), nil
}
