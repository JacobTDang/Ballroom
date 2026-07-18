package tutor

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cloudwego/eino/schema"
)

// --- trimHistory ---

func TestTrimHistory_KeepsLastPairEvenWhenItAloneExceedsBudget(t *testing.T) {
	history := []*schema.Message{
		schema.SystemMessage("persona"),
		schema.UserMessage(strings.Repeat("u", 500)),
		schema.AssistantMessage(strings.Repeat("a", 500), nil),
	}
	got := trimHistory(history, 10)

	if len(got) != 3 {
		t.Fatalf("trimHistory(budget=10) kept %d messages, want 3 (system + the one mandatory pair): %+v", len(got), got)
	}
	if got[0].Role != schema.System {
		t.Errorf("got[0].Role = %v, want System", got[0].Role)
	}
	if got[1].Content != history[1].Content || got[2].Content != history[2].Content {
		t.Errorf("got[1:] = %+v, want the original last pair kept intact even though it alone exceeds the budget", got[1:])
	}
}

// TestTrimHistory_NeverSplitsAPair uses a budget that lands mid-way
// through some pair's own size -- a trim that walked message-by-message
// instead of pair-by-pair would be tempted to keep, say, just an
// assistant reply without the user message that prompted it.
func TestTrimHistory_NeverSplitsAPair(t *testing.T) {
	history := []*schema.Message{schema.SystemMessage("persona")}
	for i := 0; i < 5; i++ {
		history = append(history,
			schema.UserMessage(fmt.Sprintf("user message number %d, with some padding", i)),
			schema.AssistantMessage(fmt.Sprintf("assistant reply number %d", i), nil),
		)
	}
	last := history[len(history)-2:]
	budget := len(last[0].Content) + len(last[1].Content) + 5 // room for the last pair plus a few spare chars, not a whole extra pair
	got := trimHistory(history, budget)

	if got[0].Role != schema.System {
		t.Fatalf("got[0].Role = %v, want System", got[0].Role)
	}
	rest := got[1:]
	if len(rest)%2 != 0 {
		t.Fatalf("trimHistory left %d post-system messages (odd count) -- a pair got split: %+v", len(rest), rest)
	}
	for i := 0; i < len(rest); i += 2 {
		if rest[i].Role != schema.User {
			t.Errorf("message at offset %d has role %v, want User (pair start)", i, rest[i].Role)
		}
		if rest[i+1].Role != schema.Assistant {
			t.Errorf("message at offset %d has role %v, want Assistant (pair end)", i+1, rest[i+1].Role)
		}
	}
}

func TestTrimHistory_UnderBudgetIsANoOp(t *testing.T) {
	history := []*schema.Message{
		schema.SystemMessage("persona"),
		schema.UserMessage("u1"), schema.AssistantMessage("a1", nil),
		schema.UserMessage("u2"), schema.AssistantMessage("a2", nil),
	}
	got := trimHistory(history, 1_000_000)

	if len(got) != len(history) {
		t.Fatalf("trimHistory trimmed under budget: got %d messages, want all %d kept", len(got), len(history))
	}
	for i := range history {
		if got[i].Role != history[i].Role || got[i].Content != history[i].Content {
			t.Errorf("message %d changed under a no-op trim: got %+v, want %+v", i, got[i], history[i])
		}
	}
}

func TestTrimHistory_EmptyHistoryIsUnchanged(t *testing.T) {
	var history []*schema.Message
	if got := trimHistory(history, 10); len(got) != 0 {
		t.Errorf("trimHistory(nil) = %+v, want empty", got)
	}
}

// TestTrimHistory_SystemPromptSurvivesAggressiveTrim pins the structural
// guarantee the whole design leans on: the persona prompt lives at
// history[0] (see tutorModel.history's own doc comment), so even a trim
// so aggressive it can't fully afford the single mandatory last pair
// must still never drop it -- losing it mid-session would silently
// change the tutor's whole persona/mode instructions partway through a
// conversation.
func TestTrimHistory_SystemPromptSurvivesAggressiveTrim(t *testing.T) {
	const persona = "you are a hints-first coding tutor..."
	history := []*schema.Message{schema.SystemMessage(persona)}
	for i := 0; i < 30; i++ {
		history = append(history,
			schema.UserMessage(strings.Repeat("q", 800)),
			schema.AssistantMessage(strings.Repeat("r", 800), nil),
		)
	}

	got := trimHistory(history, 1) // as aggressive as it gets

	if len(got) == 0 {
		t.Fatal("trimHistory returned an empty slice, want at least the system prompt")
	}
	if got[0].Role != schema.System || got[0].Content != persona {
		t.Fatalf("got[0] = %+v, want the original system prompt preserved verbatim", got[0])
	}
}

// --- classifyTurnError ---

func TestClassifyTurnError_ContextOverflowVariants(t *testing.T) {
	cases := []string{
		"provider error: context length exceeded",
		"400 Bad Request: this model's maximum context length is 4096 tokens",
		"invalid_request_error: context_length_exceeded",
		"Error: too many tokens in the request",
	}
	for _, raw := range cases {
		got := classifyTurnError(fmt.Errorf("%s", raw))
		if !strings.Contains(got, "context window") {
			t.Errorf("classifyTurnError(%q) = %q, want a context-window note", raw, got)
		}
		if strings.Contains(got, "could not reach") {
			t.Errorf("classifyTurnError(%q) = %q, want it distinct from the generic connectivity wording", raw, got)
		}
	}
}

func TestClassifyTurnError_DeadlineExceededGetsTimeoutNote(t *testing.T) {
	err := fmt.Errorf("Post \"http://x\": %w", context.DeadlineExceeded)
	got := classifyTurnError(err)

	if !strings.Contains(got, "timed out") {
		t.Errorf("classifyTurnError(deadline exceeded) = %q, want a timeout note", got)
	}
	if strings.Contains(got, "could not reach") {
		t.Errorf("classifyTurnError(deadline exceeded) = %q, want it distinct from the generic connectivity wording", got)
	}
	if strings.Contains(got, "context window") {
		t.Errorf("classifyTurnError(deadline exceeded) = %q, want it distinct from the context-overflow note", got)
	}
}

func TestClassifyTurnError_CanceledGetsCancelNote(t *testing.T) {
	err := fmt.Errorf("Post \"http://x\": %w", context.Canceled)
	got := classifyTurnError(err)

	if !strings.Contains(got, "cancel") {
		t.Errorf("classifyTurnError(canceled) = %q, want a cancellation note", got)
	}
	if strings.Contains(got, "could not reach") {
		t.Errorf("classifyTurnError(canceled) = %q, want it distinct from the generic connectivity wording", got)
	}
	if strings.Contains(got, "timed out") {
		t.Errorf("classifyTurnError(canceled) = %q, want it distinct from the timeout note", got)
	}
}

// TestClassifyTurnError_UnknownReturnsEmptySoCallerKeepsExistingWording
// is the "keep today's wording for the unrecognized case" requirement:
// classifyTurnError has no endpoint to build "could not reach <endpoint>"
// from itself, so it signals "no special classification" with an empty
// string and leaves that wording to the caller (Update's turnCompleteMsg
// case). TestTutorModel_ErrorMessageIncludesRealUnderlyingDetail
// (model_test.go) is the end-to-end confirmation that the caller's
// fallback still produces exactly that wording.
func TestClassifyTurnError_UnknownReturnsEmptySoCallerKeepsExistingWording(t *testing.T) {
	got := classifyTurnError(fmt.Errorf("connection refused"))
	if got != "" {
		t.Errorf("classifyTurnError(unrecognized) = %q, want empty", got)
	}
}

func TestClassifyTurnError_NilErrReturnsEmpty(t *testing.T) {
	if got := classifyTurnError(nil); got != "" {
		t.Errorf("classifyTurnError(nil) = %q, want empty", got)
	}
}
