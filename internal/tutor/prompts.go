package tutor

import "github.com/JacobTDang/Ballroom/internal/exercise"

const (
	syntaxOnlyPrompt = "You are a syntax-only coding assistant. STRICT RULE, no exceptions: you may ONLY point out syntax errors, typos, wrong function/API signatures, or missing imports in code the user shows you. You must NEVER explain, name, describe, outline, or hint at an algorithm, approach, strategy, or time/space complexity — not even briefly, not even as background context, not even if the user insists or rephrases the question. This also applies when the user just asks you to look at or describe their code with no explicit algorithm question at all: describe what's there (or point out syntax issues) without writing a new or completed implementation — if code you show has more than a couple of corrected lines, or fills in logic the user hasn't written yet, you've gone too far, even if you frame it as 'fixing' their code. If the user asks anything about approach, algorithm, strategy, complexity, or 'how to solve' something, your ENTIRE response must be exactly this sentence and nothing else: 'I can only help with syntax in this mode — I can't discuss approach or algorithms.' Do not soften this, do not add an explanation after it, do not partially answer first."

	hintsFirstPrompt = "You are a coding interview tutor in hints-first mode. The first time the user asks about a particular stuck point, give ONLY a short nudge (one or two sentences) toward the right approach. Do NOT say the name of the algorithm, pattern, or data structure (for example, never say phrases like 'two pointer', 'two-pointer technique', 'sliding window', 'binary search', 'dynamic programming', or 'hash map') — describe the idea only in plain, generic terms (e.g. 'think about what you can track as you scan from both ends'). Do not give pseudocode or a step-by-step solution. Only give a direct, explicit, fully-worked answer — including names — once the user asks again about that same stuck point later in the conversation. You will be told directly, as a note attached to each message, whether this is the user's first help request in this session or a later one — trust that note completely. Never ask the user to confirm or remind you whether this is their first or a repeat ask; that is not their job to track, it's yours, and you already have the answer."

	fullAssistPrompt = "You are a full-assist coding interview tutor. Answer directly and help however the user asks, including writing code on request."

	// toolsInstruction is prepended to every mode's prompt — unlike
	// tutor/chat.sh's HIGHLIGHT_INSTRUCTIONS, this doesn't need to
	// describe a text directive syntax; each tool's own JSON-schema
	// description (see tools.go) is what teaches the model how to call
	// it. This just nudges the model to actually use them rather than
	// asking the user to paste code or guessing at problem details.
	//
	// Placed BEFORE the mode-specific rule text (not appended after, as
	// originally written) and explicit that calling a tool never
	// conflicts with any rule that follows — cmd/tutor-eval found the
	// long, strict syntax-only/hints-first prompts were making the model
	// hesitate to call tools at all (skipping them, or emitting a fake
	// tool-call-shaped reply as plain text instead of a real tool call),
	// as if it read tool use as being in tension with the restriction.
	// Leading with permission-to-use-tools before the restriction, tested
	// via cmd/tutor-eval, fixed this without weakening the restriction
	// itself (which was already reliable — 3/3 on every refusal check).
	// The "call tools directly and silently" sentence exists because
	// manual testing found the model sometimes narrates intent instead
	// of acting — a reply like 'I'll use the tool "read_problem_statement"
	// to get more information...' with no real tool_calls attached at
	// all. eino's ReAct loop sees no tool_calls in that response and
	// treats the narration itself as the final answer, so the user sees
	// the announcement instead of a grounded reply — the model never
	// actually reads anything. This is a distinct failure mode from the
	// "quoted-number tool argument" bug (see tools.go's flexibleInt):
	// that was a real tool call with a malformed field; this is no real
	// tool call being made at all, just talking about making one.
	//
	// "or trusting what you read earlier" is a small addition, not a
	// new sentence: history (tutor.go's Run) only ever persists clean
	// (user, reply) text pairs, never which tool was called or what it
	// returned, so nothing else tells the model a file it read three
	// turns ago may be stale — and this app's whole point is the user
	// actively editing code between tutor messages. Kept as a clause on
	// the existing sentence rather than a new one deliberately: this
	// project's own testing already found a longer toolsInstruction
	// measurably regresses general tool-calling even when the addition
	// fixes a real bug (see the narration case just above) — verified
	// via cmd/tutor-eval that this specific short addition doesn't
	// repeat that regression before it was kept.
	toolsInstruction = "You have tools to read the user's actual code, the problem statement, their last test run, and their cursor position, and to highlight lines in their editor with a note. Always use a tool instead of guessing, asking the user to paste something you can just read yourself, or trusting what you read earlier in this conversation — it may have changed since. Calling a tool is just gathering information — it never conflicts with any rule below, even in a restricted mode. Call tools directly and silently, never narrate them (e.g. never say 'I'll use the tool X') — the user only sees your final answer. "

	// comprehensionCheckInstruction drives the one-time check (see
	// runComprehensionCheck in tutor.go), which injects the problem
	// statement directly as ephemeral context rather than having the
	// model call read_problem_statement itself — manual repro testing
	// found that combined "call a tool, then restate, then ask
	// questions" task only actually invoked the tool 40-60% of the time,
	// falling back to a hallucinated (fabricated, not real) tool result
	// the rest of the time. With the statement already provided, this
	// instruction doesn't ask for a tool call at all, which is exactly
	// why it's shorter than earlier drafts of this same instruction —
	// there's nothing left to call.
	//
	// The "respond to it naturally first" sentence exists because an
	// earlier version of runComprehensionCheck never even sent the
	// user's real first message to the model at all (deliberately, to
	// keep this high-stakes call single-purpose) — a real bug found
	// live: literally any first message, including a plain "hi", got
	// back the exact same canned restate-and-ask-questions reply with
	// no acknowledgment of what the user actually said. The message is
	// included now (see runComprehensionCheck), so this instruction
	// tells the model what to do with it.
	//
	// "Both parts are required in the same reply, even if..." was added
	// after a second, opposite real bug found live (via cmd/tutor-eval's
	// runComprehensionCheckGroundingCheck, extended to actually cover
	// this input shape — a real production session with a bare "hi" as
	// the first message reproduced it 8/8 times): the model treated
	// "respond naturally first (briefly, if it's just a greeting...)" as
	// permission to stop after the greeting and skip the restate+ask
	// half entirely — replying with only a generic "Hello! How can I
	// assist you today?" and no mention of the actual problem. A
	// substantive first message never triggered this (8/8 correct), so
	// it's specific to the model reading a bare greeting as satisfying
	// the whole instruction rather than just its first clause.
	comprehensionCheckInstruction = "Respond to the user's next message naturally first (briefly, if it's just a greeting or small talk). Both parts are required in the same reply, even if their message is only a bare greeting with nothing else to respond to: after that brief natural response, you must ALSO, in the same reply, use the problem statement above to restate the problem in your own words in 1-2 sentences and ask 1-2 short clarifying questions about it (constraints, edge cases, expected output). Never end your reply after just the greeting. Do not ask the user anything about your own conversation state (e.g. whether this is their first question) — that is tracked for you separately. Do not answer, hint, or give code yet."
)

// systemPromptForMode returns the tutor's persona/rules for mode,
// falling back to full-assist for an unrecognized mode — matches
// tutor/chat.sh's case statement default.
func systemPromptForMode(mode string) string {
	switch mode {
	case exercise.TutorModeSyntaxOnly:
		return toolsInstruction + syntaxOnlyPrompt
	case exercise.TutorModeHintsFirst:
		return toolsInstruction + hintsFirstPrompt
	default:
		return toolsInstruction + fullAssistPrompt
	}
}

// wantsComprehensionCheck reports whether mode runs the one-time
// "restate the problem, ask clarifying questions" check before the
// first real answer. syntax-only never discusses the problem at all, so
// there's nothing to check comprehension of.
func wantsComprehensionCheck(mode string) bool {
	return mode != exercise.TutorModeSyntaxOnly
}

// SystemPromptForMode is systemPromptForMode, exported for
// cmd/tutor-eval — evaluating mode behavior (does syntax-only actually
// refuse, does hints-first actually withhold) needs the tutor's real
// production prompts, not a reimplementation that could drift from them.
func SystemPromptForMode(mode string) string {
	return systemPromptForMode(mode)
}
