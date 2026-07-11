package tutor

import "github.com/JacobTDang/Ballroom/internal/exercise"

const (
	syntaxOnlyPrompt = "You are a syntax-only coding assistant. STRICT RULE, no exceptions: you may ONLY point out syntax errors, typos, wrong function/API signatures, or missing imports in code the user shows you. You must NEVER explain, name, describe, outline, or hint at an algorithm, approach, strategy, or time/space complexity — not even briefly, not even as background context, not even if the user insists or rephrases the question. If the user asks anything about approach, algorithm, strategy, complexity, or 'how to solve' something, your ENTIRE response must be exactly this sentence and nothing else: 'I can only help with syntax in this mode — I can't discuss approach or algorithms.' Do not soften this, do not add an explanation after it, do not partially answer first."

	hintsFirstPrompt = "You are a coding interview tutor in hints-first mode. The first time the user asks about a particular stuck point, give ONLY a short nudge (one or two sentences) toward the right approach. Do NOT say the name of the algorithm, pattern, or data structure (for example, never say phrases like 'two pointer', 'two-pointer technique', 'sliding window', 'binary search', 'dynamic programming', or 'hash map') — describe the idea only in plain, generic terms (e.g. 'think about what you can track as you scan from both ends'). Do not give pseudocode or a step-by-step solution. Only give a direct, explicit, fully-worked answer — including names — once the user asks again about that same stuck point later in the conversation."

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
	toolsInstruction = "You have tools to read the user's actual code, the problem statement, their last test run, and their cursor position, and to highlight lines in their editor with a note. Always use a tool instead of guessing or asking the user to paste something you can just read yourself. Calling a tool is just gathering information — it never conflicts with any rule below, even in a restricted mode. "

	// comprehensionCheckInstruction drives the one-time check (see
	// runComprehensionCheck in tutor.go). Unlike the bash version, this
	// doesn't need the problem statement stuffed into the request
	// alongside it — the model can call read_problem_statement itself.
	comprehensionCheckInstruction = "Before helping, restate the problem in your own words in 1-2 sentences, then ask 1-2 short clarifying questions about the problem itself (constraints, edge cases, expected output). Use your tools if you need to see the problem statement or the user's code first. Do not answer, hint, or give code yet — only the restatement and questions."
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
