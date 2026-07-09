#!/usr/bin/env bash
set -euo pipefail

# Chat CLI: talks to the host Ollama endpoint via /api/chat, keeping a
# running message history so multi-turn rules (e.g. hints-first's
# "escalate on the second ask") have context to work from — a stateless
# per-request call can't tell it's being asked again.
#
# System prompt is selected by PRACTICE_TUTOR_MODE (set by the orchestrator
# from the exercise's tutor_mode field; unset/sandbox defaults to
# full-assist). See interview_prep_mvp_spec.md Section 3.4.

OLLAMA_HOST="${OLLAMA_HOST:-http://host.docker.internal:11434}"
MODEL="${TUTOR_MODEL:-qwen2.5-coder:7b}"
MODE="${PRACTICE_TUTOR_MODE:-full-assist}"

# RPC socket for the editor pane's nvim instance (set by entrypoint.sh
# alongside the `nvim --listen` invocation in pane 0; see issue #24). Empty
# when unset/unreachable just means highlighting is silently unavailable —
# the chat loop still works.
NVIM_SOCKET="${NVIM_SOCKET:-}"

# Directive the model can emit, on its own line, to highlight a range in the
# open solution file with an attached note:
#   <<<highlight file=solution.go line=10-14 note="off-by-one here">>>
# `line=N` (no dash) is also accepted. Matched with ERE via grep/sed below.
HIGHLIGHT_DIRECTIVE_RE='<<<highlight[[:space:]]+file=[^[:space:]]+[[:space:]]+line=[0-9]+(-[0-9]+)?[[:space:]]+note="[^"]*">>>'

HIGHLIGHT_INSTRUCTIONS=" When it would help to point at specific code, include a directive on its own line in exactly this form: <<<highlight file=FILENAME line=START-END note=\"short note\">>> (a single line is fine as line=N). You can include more than one across a response. This directive is stripped before the user sees your reply, so never mention the syntax itself to the user or describe it as a special command."

case "$MODE" in
  syntax-only)
    DEFAULT_PROMPT="You are a syntax-only coding assistant. STRICT RULE, no exceptions: you may ONLY point out syntax errors, typos, wrong function/API signatures, or missing imports in code the user shows you. You must NEVER explain, name, describe, outline, or hint at an algorithm, approach, strategy, or time/space complexity — not even briefly, not even as background context, not even if the user insists or rephrases the question. If the user asks anything about approach, algorithm, strategy, complexity, or 'how to solve' something, your ENTIRE response must be exactly this sentence and nothing else: 'I can only help with syntax in this mode — I can't discuss approach or algorithms.' Do not soften this, do not add an explanation after it, do not partially answer first."
    ;;
  hints-first)
    DEFAULT_PROMPT="You are a coding interview tutor in hints-first mode. The first time the user asks about a particular stuck point, give ONLY a short nudge (one or two sentences) toward the right approach. Do NOT say the name of the algorithm or pattern (for example, never say phrases like 'two pointer', 'two-pointer technique', 'sliding window', 'binary search', 'dynamic programming', or similar named techniques) — describe the idea only in plain, generic terms (e.g. 'think about what you can track as you scan from both ends'). Do not give pseudocode or a step-by-step solution. Only give a direct, explicit, fully-worked answer — including the technique's name if relevant — once the user asks again about that same stuck point later in this conversation."
    ;;
  full-assist)
    DEFAULT_PROMPT="You are a full-assist coding interview tutor. Answer directly and help however the user asks, including writing code on request."
    ;;
  *)
    DEFAULT_PROMPT="You are a full-assist coding interview tutor. Answer directly and help however the user asks, including writing code on request."
    ;;
esac
DEFAULT_PROMPT="${DEFAULT_PROMPT}${HIGHLIGHT_INSTRUCTIONS}"
SYSTEM_PROMPT="${TUTOR_SYSTEM_PROMPT:-$DEFAULT_PROMPT}"

# Sends one highlight+note over to the editor pane's nvim RPC socket. Never
# lets a highlight failure take down the chat loop: a missing/unreadable
# socket (editor pane not up yet, or NVIM_SOCKET unset) is expected in some
# runs (e.g. sandbox mode) and just means highlighting is unavailable, so we
# skip quietly; an actual RPC error is logged to stderr and swallowed —
# this is graceful degradation for model/user input we don't control, not
# the "fail loud" case that applies to our own code's bugs.
apply_highlight() {
  local file="$1" start="$2" end="$3" note="$4" expr out

  if [ -z "$NVIM_SOCKET" ] || [ ! -S "$NVIM_SOCKET" ]; then
    return 0
  fi

  # VimL single-quoted strings only need '' to escape a literal quote;
  # nothing else is special inside them.
  local vim_file="${file//\'/\'\'}"
  local vim_note="${note//\'/\'\'}"
  expr="v:lua.require('ballroom_highlight').add_highlight('${vim_file}', ${start}, ${end}, '${vim_note}')"

  if ! out=$(nvim --server "$NVIM_SOCKET" --remote-expr "$expr" 2>&1); then
    echo "tutor: highlight RPC call failed: $out" >&2
    return 0
  fi
  case "$out" in
    ballroom_highlight\ error:*)
      echo "tutor: $out" >&2
      ;;
  esac
  return 0
}

# Scans a tutor reply for <<<highlight ...>>> directives (see
# HIGHLIGHT_DIRECTIVE_RE above), fires the corresponding nvim RPC call for
# each one, and echoes the reply with all directives stripped — the user
# never sees the raw marker syntax. A directive that fails to parse is
# skipped (logged, not fatal): malformed model output must degrade
# gracefully rather than crash the tutor loop.
process_highlights() {
  local text="$1" directive start end
  local file="" note=""
  local field_re='file=([^[:space:]]+)[[:space:]]+line=([0-9]+)(-([0-9]+))?[[:space:]]+note="([^"]*)"'

  while IFS= read -r directive; do
    [ -z "$directive" ] && continue
    if [[ "$directive" =~ $field_re ]]; then
      file="${BASH_REMATCH[1]}"
      start="${BASH_REMATCH[2]}"
      end="${BASH_REMATCH[4]:-$start}"
      note="${BASH_REMATCH[5]}"
      apply_highlight "$file" "$start" "$end" "$note"
    else
      echo "tutor: skipping malformed highlight directive: $directive" >&2
    fi
  done < <(grep -oE "$HIGHLIGHT_DIRECTIVE_RE" <<<"$text" || true)

  sed -E "s/${HIGHLIGHT_DIRECTIVE_RE}//g" <<<"$text"
}

echo "tutor (${MODEL}, mode=${MODE}) — connected to ${OLLAMA_HOST}. Ctrl-D to exit."

messages=$(jq -n --arg system "$SYSTEM_PROMPT" '[{role: "system", content: $system}]')

while IFS= read -r -p '> ' line; do
  [ -z "$line" ] && continue

  messages=$(jq --arg content "$line" '. + [{role: "user", content: $content}]' <<<"$messages")
  payload=$(jq -n --arg model "$MODEL" --argjson messages "$messages" \
    '{model: $model, messages: $messages, stream: false, options: {temperature: 0.2}}')

  if ! response=$(curl -sf "${OLLAMA_HOST}/api/chat" -d "$payload"); then
    echo "tutor: could not reach ${OLLAMA_HOST}" >&2
    continue
  fi

  reply=$(echo "$response" | jq -r '.message.content')
  display_reply=$(process_highlights "$reply")
  echo "$display_reply"

  messages=$(jq --arg content "$reply" '. + [{role: "assistant", content: $content}]' <<<"$messages")
done

echo
