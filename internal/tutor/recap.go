package tutor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/schema"
)

// recapInstruction drives SessionRecap's single-purpose model call.
const recapInstruction = "You are writing a 2-3 sentence recap of a practice session for the user's own study notes. From the tutor-conversation transcript and the submission outcome, say what was attempted, what the user actually struggled with or asked about, and how it ended. Address the user as 'you'. Plain sentences only -- no preamble, no headers, no bullet points."

// recapOutputCap bounds how much raw test/grade output rides along in
// the recap request -- the tail matters most (the verdict and the last
// failures), the head is usually compiler noise.
const recapOutputCap = 2000

// SessionRecap summarizes a just-submitted session from the workspace
// transcript (transcript.go's export; a session with no tutor
// conversation still gets a recap of the outcome) plus the submission
// result and output. One blocking call on cfg.Model; errors bubble so
// the caller degrades to a notice.
func SessionRecap(ctx context.Context, cfg Config, result, output string) (string, error) {
	transcript := "(the user didn't talk to the tutor this session)"
	if b, err := os.ReadFile(filepath.Join(cfg.WorkDir, "transcript.md")); err == nil && len(b) > 0 {
		transcript = string(b)
	}
	if len(output) > recapOutputCap {
		output = "...\n" + output[len(output)-recapOutputCap:]
	}

	cm, err := newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
	if err != nil {
		return "", fmt.Errorf("tutor: session recap: %w", err)
	}
	messages := []*schema.Message{
		schema.SystemMessage(recapInstruction),
		schema.UserMessage(fmt.Sprintf("Submission result: %s\n\nSubmission output:\n%s\n\nTutor conversation transcript:\n\n%s", result, output, transcript)),
	}
	reply, err := generateWithEmptyChoicesRetry(ctx, cm, messages)
	if err != nil {
		return "", fmt.Errorf("tutor: session recap: %w", err)
	}
	return strings.TrimSpace(reply.Content), nil
}
