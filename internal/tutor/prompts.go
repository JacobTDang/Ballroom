package tutor

import (
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/cloudwego/eino/schema"
)

const (
	syntaxOnlyPrompt = "You are a syntax-only coding assistant. STRICT RULE, no exceptions: you may ONLY point out syntax errors, typos, wrong function/API signatures, or missing imports in code the user shows you. You must NEVER explain, name, describe, outline, or hint at an algorithm, approach, strategy, or time/space complexity — not even briefly, not even as background context, not even if the user insists or rephrases the question. This also applies when the user just asks you to look at or describe their code with no explicit algorithm question at all: describe what's there (or point out syntax issues) without writing a new or completed implementation — if code you show has more than a couple of corrected lines, or fills in logic the user hasn't written yet, you've gone too far, even if you frame it as 'fixing' their code. If the user asks anything about approach, algorithm, strategy, complexity, or 'how to solve' something, your ENTIRE response must be exactly this sentence and nothing else: 'I can only help with syntax in this mode — I can't discuss approach or algorithms.' Do not soften this, do not add an explanation after it, do not partially answer first."

	// hintsFirstPrompt: 2026-07-12, tried three prompt-wording variants
	// against qwen2.5:14b-instruct via cmd/tutor-eval (real sample size,
	// not 2-3 runs) after switching this project's tutor model away from
	// qwen2.5-coder:14b-instruct (see config.Qwen25Coder14BModel — that
	// model can't do real tool calling at all) and finding this new
	// model leaks the forbidden technique name on a first ask more often
	// than this project's history documents for DefaultTutorModel.
	//
	// Baseline (this exact wording): 5/8, 1/8, 6/8 on the three
	// hints-first cmd/tutor-eval checks.
	//  1. Reworded to a "STRICT RULE" framing (matching syntaxOnlyPrompt's
	//     proven style) plus a wider forbidden-word list plus a
	//     self-check sentence: 0/8 on the main leak check (worse), and a
	//     THIRD, previously-unrelated check in the same mode ("can still
	//     call read_solution_file") dropped from 6/8 to 0/8 — the "a
	//     longer/stricter instruction measurably hurts tool-calling even
	//     while fixing something else" pattern this project has hit
	//     before (see toolsInstruction's own doc comment).
	//  2. Kept this wording, only widened the forbidden-word list to
	//     match cmd/tutor-eval's own forbiddenTechniqueTerms (13 words,
	//     not the original 6) — the tool-calling regression partly
	//     recovered (4/8) but the leak checks got WORSE, not better
	//     (1/8, 0/8): a failing reply used "complements" despite
	//     "complement" being explicitly banned in both the prompt and
	//     the eval's own list — the model reached for a synonym-adjacent
	//     form of banned vocabulary rather than avoiding the concept
	//     entirely.
	// Three attempts, no improvement (net regression each time) — this
	// reads as a real behavioral limit of this model for this specific
	// constraint, not a wording problem still waiting on the right
	// phrasing. Reverted to the original wording rather than keep
	// guessing at another variant. If this needs fixing, the next things
	// worth trying are a semantic check instead of a word-list ban, or a
	// different model for hints-first specifically — not another prompt
	// reword.
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
	// comprehensionCheckMessages/startTurn in model.go), which injects the problem
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
	// earlier version of the comprehension check never even sent the
	// user's real first message to the model at all (deliberately, to
	// keep this high-stakes call single-purpose) — a real bug found
	// live: literally any first message, including a plain "hi", got
	// back the exact same canned restate-and-ask-questions reply with
	// no acknowledgment of what the user actually said. The message is
	// included now (see comprehensionCheckMessages in model.go), so this
	// instruction tells the model what to do with it.
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

	// routingInstruction drives decideHandoff (tutor.go) — a single
	// yes/no classification call on the orchestrator's raw chat model
	// (no tools, no react.Agent), asking only whether the coding
	// specialist should handle this turn. Deliberately biased toward
	// YES on anything unclear: a wrong "No" silently leaves a real code
	// question with the weaker/cheaper model with no way for the user
	// to notice, while a wrong "Yes" just costs one unnecessary
	// specialist call — asymmetric costs, so the instruction (and
	// decideHandoff's own error/parse-failure fallback) both lean
	// toward the safer side.
	routingInstruction = "You are deciding whether a coding-interview tutor question needs a coding specialist's attention. Reply with exactly one word: NO if the message is a greeting, small talk, or a general clarifying question that doesn't require reading or reasoning about code. Reply YES for anything about the problem's approach, algorithm, code review, debugging, or hints. When unsure, reply YES."
)

// toolCallingStrategy selects which tool-calling instruction variant a
// session's model needs.
type toolCallingStrategy string

const nativeToolCalling toolCallingStrategy = "native"

// jsonFallbackToolCalling is for a model toolcheck.go's CheckToolCalling
// confirms doesn't populate a real tool_calls field -- it's told (via
// jsonFallbackInstruction below) to emit a specific JSON shape as its
// entire reply instead, which fallback.go's runFallbackToolLoop parses
// and executes by hand.
const jsonFallbackToolCalling toolCallingStrategy = "json_fallback"

// jsonFallbackInstruction is the fixed protocol text for
// jsonFallbackToolCalling -- the real tool catalog (names/descriptions/
// JSON schemas) is appended separately per-session, by prependToolsPrompt
// below via renderToolCatalog (fallback.go), since the actual tool set
// is only known at runtime (buildTools(cfg)), not something a package
// const can describe.
const jsonFallbackInstruction = `You do not have real tool-calling in this ` +
	`conversation. To use a tool, your ENTIRE reply must be ONLY a single ` +
	`JSON object of this exact shape and nothing else -- no explanation, ` +
	`no code fence, no text before or after it: ` +
	`{"name": "<tool name>", "arguments": {<its arguments>}}. Call only ` +
	`one tool at a time, then stop and wait -- you will be given the ` +
	`result as your next message before you continue. Once you have ` +
	`everything you need, reply normally in plain text with your real ` +
	`answer and no JSON at all -- that reply is what the user sees.`

// toolsInstructions maps a toolCallingStrategy to the instruction text
// that teaches the model how tools work at all. A dictionary instead of
// a bare constant so a second strategy can be added as a data entry
// instead of restructuring systemPromptForMode/prependToolsPrompt again.
var toolsInstructions = map[toolCallingStrategy]string{
	nativeToolCalling:       toolsInstruction,
	jsonFallbackToolCalling: jsonFallbackInstruction,
}

// modePrompts maps each tutor mode to its persona/rule text (excluding
// toolsInstruction, selected separately above). A plain data lookup
// instead of a switch, so adding/removing/swapping a mode's prompt is a
// one-line map edit -- same pattern internal/exercise's validTutorModes
// already uses for mode validation.
var modePrompts = map[string]string{
	exercise.TutorModeSyntaxOnly: syntaxOnlyPrompt,
	exercise.TutorModeHintsFirst: hintsFirstPrompt,
	exercise.TutorModeFullAssist: fullAssistPrompt,
}

// personaPromptForMode returns mode's persona/rule text alone (no tools
// instruction), falling back to full-assist for an unrecognized mode —
// matches tutor/chat.sh's case statement default. This, not
// systemPromptForMode, is what newTutorModel seeds tutorModel.history
// with: the tools instruction is strategy-dependent and (once routing
// is enabled) can differ turn to turn depending which role answers, so
// it's prepended fresh per call by prependToolsPrompt instead of being
// baked once into the session's fixed system message.
func personaPromptForMode(mode string) string {
	prompt, ok := modePrompts[mode]
	if !ok {
		prompt = fullAssistPrompt
	}
	return prompt
}

// systemPromptForMode returns mode's persona/rules composed with the
// native tool-calling strategy's instruction. Used by cmd/tutor-eval,
// which only ever evaluates the native path -- real sessions use
// personaPromptForMode + prependToolsPrompt instead, since a session's
// actual strategy isn't known until newTutorModel runs CheckToolCalling.
func systemPromptForMode(mode string) string {
	return toolsInstructions[nativeToolCalling] + personaPromptForMode(mode)
}

// prependToolsPrompt prepends an ephemeral, strategy-specific tools
// system message ahead of messages -- never persisted to
// tutorModel.history, the same ephemeral-context pattern turnMessages'
// hint-count note already uses. toolCatalogText (renderToolCatalog,
// fallback.go) is only appended for jsonFallbackToolCalling: a native
// model already has real tool schemas bound via the provider's API, so
// describing them again in prose would be redundant (and, per
// toolsInstruction's own doc comment, this codebase has already found
// longer prompts measurably regress native tool-calling reliability).
// Does not mutate messages -- returns a new slice.
func prependToolsPrompt(strategy toolCallingStrategy, toolCatalogText string, messages []*schema.Message) []*schema.Message {
	instruction := toolsInstructions[strategy]
	if strategy == jsonFallbackToolCalling {
		instruction += "\n\n" + toolCatalogText
	}
	out := make([]*schema.Message, 0, len(messages)+1)
	out = append(out, schema.SystemMessage(instruction))
	out = append(out, messages...)
	return out
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
