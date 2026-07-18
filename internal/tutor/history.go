package tutor

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/schema"
)

// historyBudgetChars bounds trimHistory's default budget for the
// conversation portion of m.history (the leading system/persona prompt
// is exempt -- see trimHistory's own doc comment). A var, not a const,
// so tests -- and a throwaway live-verification override -- can shrink
// it rather than needing a real long-running session to actually
// outgrow the real budget.
//
// Chars, not tokens: a real tokenizer would be a new dependency for
// marginal accuracy over this proxy, and this only needs to land in the
// right ballpark to keep a session responsive and clear of a provider's
// hard context-length rejection, not be exact -- ~4 chars/token is the
// commonly cited rule of thumb for English text against most tokenizers
// in this model class. 24000 chars is roughly 6000 tokens: comfortably
// inside even a constrained free-tier OpenRouter model's context window
// (some are as small as 4k-8k tokens total, shared with the system
// prompt, tool catalog, and the model's own reply) while still keeping
// a useful amount of real conversation around.
var historyBudgetChars = 24000

// trimHistory keeps history's leading system prompt (see
// tutorModel.history's own doc comment: "history holds only the system
// prompt plus clean (user, assistant) pairs") unconditionally and in
// full, plus as many of the most RECENT complete (user, assistant)
// pairs as fit within budgetChars, walking from the newest pair
// backward. Always keeps at least the single most recent pair, even if
// it alone exceeds budgetChars -- an oversized last exchange is still
// more useful as context than none at all, and dropping it would mean a
// turn's own just-sent message doesn't even survive into its own
// follow-up reply.
//
// Trimming at a pair boundary is only safe because every append to
// history is atomic and every pair is clean -- confirmed directly from
// the code, not assumed: the only two writes to m.history in the whole
// package are newTutorModel (seeds it with exactly the one system
// message) and Update's turnCompleteMsg case, whose success branch
// appends a schema.UserMessage and a schema.AssistantMessage together
// in one statement -- the failure branch above it returns before ever
// reaching that line, so a failed turn is never added at all. Splitting
// a pair (keeping a user message without its reply, or vice versa)
// would send the model half a conversation turn with no way to make
// sense of it.
func trimHistory(history []*schema.Message, budgetChars int) []*schema.Message {
	if len(history) == 0 {
		return history
	}

	prefixLen := 0
	if history[0].Role == schema.System {
		prefixLen = 1
	}
	pairs := history[prefixLen:]
	n := len(pairs) / 2

	keep := 0
	size := 0
	for i := n - 1; i >= 0; i-- {
		user, assistant := pairs[2*i], pairs[2*i+1]
		pairSize := len(user.Content) + len(assistant.Content)
		if keep > 0 && size+pairSize > budgetChars {
			break
		}
		size += pairSize
		keep++
	}

	if keep == n {
		// Nothing to trim -- return the original slice rather than a
		// freshly built copy of the exact same content.
		return history
	}

	trimmed := make([]*schema.Message, 0, prefixLen+2*keep)
	trimmed = append(trimmed, history[:prefixLen]...)
	trimmed = append(trimmed, pairs[len(pairs)-2*keep:]...)
	return trimmed
}

// classifyTurnError distinguishes WHY a turn's Generate call failed --
// a context-window overflow, a per-turn timeout (turnTimeout, model.go),
// or a cancellation -- so the note reads differently from a plain
// connectivity failure; none of those three actually failed to reach
// the model the way "could not reach <endpoint>" implies. Returns "" for
// everything else, telling the caller (Update's turnCompleteMsg case) to
// keep its own existing "could not reach <endpoint>: <err>" wording --
// classifyTurnError has no endpoint to build that from itself, and
// today's generic wording is still the right description for a genuine
// connectivity or provider-rejection failure.
//
// Every match below is a substring check on err.Error(), not
// errors.Is/errors.As: eino and its underlying provider clients don't
// reliably wrap sentinel errors with %w (see fallback.go's
// isEmptyChoicesErr for the same reasoning applied to a different
// provider failure) -- but a context.DeadlineExceeded/context.Canceled
// sentinel's own message text, and a provider's own context-length
// error message, both survive any amount of %v/%w wrapping regardless,
// so substring matching is the more robust choice here, not a lesser
// one.
func classifyTurnError(err error) string {
	if err == nil {
		return ""
	}
	msg := err.Error()
	switch {
	case containsAny(strings.ToLower(msg), "context length", "maximum context", "context_length_exceeded", "too many tokens"):
		return fmt.Sprintf("the conversation is too long for the model's context window -- try a shorter message, or start a fresh session: %v", err)
	case strings.Contains(msg, context.DeadlineExceeded.Error()):
		return fmt.Sprintf("turn timed out after %v with no reply -- please try again", turnTimeout)
	case strings.Contains(msg, context.Canceled.Error()):
		return fmt.Sprintf("turn was cancelled: %v", err)
	default:
		return ""
	}
}

// containsAny reports whether s contains any of substrs.
func containsAny(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
