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
#
# Each request also gets the current contents of the active solution file
# injected as ephemeral context (not stored in the persisted conversation
# history) so the tutor can see what the user has actually written in
# pane 0 (nvim). The file is re-read fresh every turn — not just once at
# session start — so edits between questions are reflected. See #22.
#
# hints-first and full-assist also run a one-time comprehension check
# (restate the problem, ask 1-2 clarifying questions) before the first
# real answer. This is enforced by chat.sh itself (see
# run_comprehension_check), not left to a prompt instruction the model
# might ignore: testing found that asking the model to "hold off on
# answering, ask questions first" in the system prompt was not reliably
# followed once the user's real, concrete question was sitting right
# there in the same request — a smaller model's pull to just answer it
# won out often enough that prompt wording alone couldn't guarantee the
# check happens. Instead, the first turn's real question is deliberately
# withheld from the model in an isolated request that only asks for the
# restatement + questions, with nothing concrete to answer instead; the
# user's actual question is picked up normally starting their next
# message. syntax-only is excluded — it's deliberately restricted to
# syntax feedback only and doesn't discuss the problem at all. See #23.

OLLAMA_HOST="${OLLAMA_HOST:-http://host.docker.internal:11434}"
MODEL="${TUTOR_MODEL:-qwen2.5-coder:7b}"
MODE="${PRACTICE_TUTOR_MODE:-full-assist}"
WORKDIR="${WORKDIR:-/workspace}"

# Cap on how many bytes of the solution file get sent per request, so a
# huge file can't blow up the request payload. Overridable for testing.
MAX_CONTEXT_BYTES="${TUTOR_FILE_CONTEXT_MAX_BYTES:-8000}"

# RPC socket for the editor window's nvim instance (set by entrypoint.sh
# alongside the `nvim --listen` invocation in the editor window; see issue
# #24). Empty when unset/unreachable just means highlighting is silently
# unavailable — the chat loop still works.
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
    DEFAULT_PROMPT="You are a coding interview tutor in hints-first mode. The first time the user asks about a particular stuck point, give ONLY a short nudge (one or two sentences) toward the right approach. Do NOT say the name of the algorithm, pattern, or data structure (for example, never say phrases like 'two pointer', 'two-pointer technique', 'sliding window', 'binary search', 'dynamic programming', or 'hash map') — describe the idea only in plain, generic terms (e.g. 'think about what you can track as you scan from both ends'). Do not give pseudocode or a step-by-step solution. Only give a direct, explicit, fully-worked answer — including names — once the user asks again about that same stuck point later in the conversation."
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

# Modes that run the one-time comprehension check (#23) before the first
# real answer. syntax-only never discusses the problem at all, so there's
# nothing to check comprehension of.
COMPREHENSION_CHECK_INSTRUCTION="Before helping, restate the problem in your own words in 1-2 sentences, then ask 1-2 short clarifying questions about the problem itself (constraints, edge cases, expected output). Do not answer, hint, or give code yet — only the restatement and questions."

# Whether mode runs the comprehension check — a plain function (not
# inlined in main's case statement) so it's directly testable.
wants_comprehension_check() {
  local mode="$1"
  [ "$mode" != "syntax-only" ]
}

# Finds the active solution file (same glob entrypoint.sh uses to open
# nvim) and prints its contents, truncated to MAX_CONTEXT_BYTES with a
# trailing marker if it was cut. Prints nothing (and always returns 0)
# when there's no solution file yet (sandbox mode's fresh `nvim .`) or
# the file can't be read — the tutor should degrade gracefully, not
# error, since a missing file is an expected state, not a bug.
build_file_context() {
  local file content size
  file=$(find "$WORKDIR" -maxdepth 1 -name 'solution.*' -type f 2>/dev/null | head -n1)
  if [ -z "$file" ]; then
    return 0
  fi

  if ! content=$(head -c "$MAX_CONTEXT_BYTES" "$file" 2>/dev/null); then
    return 0
  fi

  size=$(wc -c <"$file" 2>/dev/null | tr -d ' ')
  if [ -n "$size" ] && [ "$size" -gt "$MAX_CONTEXT_BYTES" ]; then
    content="${content}
...[truncated, file is ${size} bytes; showing first ${MAX_CONTEXT_BYTES}]"
  fi

  printf '%s' "$content"
}

# Issues one isolated request asking the model ONLY to restate the
# problem and ask clarifying questions (#23) — deliberately never
# including the user's real first message in this request, so there's
# nothing concrete for the model to answer instead of doing the check
# (see the header comment for why this is enforced in code rather than
# left to a prompt instruction). Appends the exchange to $messages
# (a global set by main) using the user's real message as the persisted
# turn, so conversation history still reads naturally; the user's actual
# question itself is picked up normally on their next input.
run_comprehension_check() {
  local user_first_message="$1"
  local check_messages check_payload check_response check_reply file_context

  check_messages=$(jq -n --arg system "$SYSTEM_PROMPT" '[{role: "system", content: $system}]')
  file_context=$(build_file_context)
  if [ -n "$file_context" ]; then
    check_messages=$(jq --arg content "$file_context" \
      '. + [{role: "system", content: ("Current contents of the solution file:\n\n" + $content)}]' \
      <<<"$check_messages")
  fi
  check_messages=$(jq --arg content "$COMPREHENSION_CHECK_INSTRUCTION" \
    '. + [{role: "user", content: $content}]' <<<"$check_messages")

  check_payload=$(jq -n --arg model "$MODEL" --argjson messages "$check_messages" \
    '{model: $model, messages: $messages, stream: false, options: {temperature: 0.2}}')

  if ! check_response=$(curl -sf "${OLLAMA_HOST}/api/chat" -d "$check_payload"); then
    echo "tutor: could not reach ${OLLAMA_HOST}" >&2
    return 1
  fi

  check_reply=$(echo "$check_response" | jq -r '.message.content')
  echo "$check_reply"

  messages=$(jq --arg content "$user_first_message" '. + [{role: "user", content: $content}]' <<<"$messages")
  messages=$(jq --arg content "$check_reply" '. + [{role: "assistant", content: $content}]' <<<"$messages")
  return 0
}

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
  # nothing else is special inside them. Substituting via a variable
  # holding a literal quote (rather than a \' literal in the pattern)
  # avoids a real bash quoting trap: inside the surrounding double
  # quotes, "${x//\'/\'\'}" doesn't produce a doubled quote at all — it
  # inserts a literal backslash before each resulting quote (bash
  # preserves the backslash rather than treating it as escaping the
  # quote), which corrupts the escaping and makes nvim's expression
  # parser reject ANY input containing a quote — apostrophes in
  # ordinary text ("it's", "doesn't"), not just malicious input.
  local q="'"
  local vim_file="${file//$q/$q$q}"
  local vim_note="${note//$q/$q$q}"
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

main() {
  echo "tutor (${MODEL}, mode=${MODE}) — connected to ${OLLAMA_HOST}. Ctrl-D to exit."

  messages=$(jq -n --arg system "$SYSTEM_PROMPT" '[{role: "system", content: $system}]')

  comprehension_check_pending=false
  if wants_comprehension_check "$MODE"; then
    comprehension_check_pending=true
  fi

  while IFS= read -r -p '> ' line; do
    [ -z "$line" ] && continue

    if [ "$comprehension_check_pending" = true ]; then
      comprehension_check_pending=false
      if run_comprehension_check "$line"; then
        continue
      fi
      # Couldn't reach Ollama for the check — fall through and handle
      # this message normally below rather than silently dropping it.
    fi

    # Build the outgoing request on top of persisted history, but inject
    # a fresh read of the solution file as ephemeral context. It's not
    # folded into $messages, so conversation history stays clean and
    # every turn re-reads the file instead of resending a stale copy.
    request_messages="$messages"
    file_context=$(build_file_context)
    if [ -n "$file_context" ]; then
      request_messages=$(jq --arg content "$file_context" \
        '. + [{role: "system", content: ("Current contents of the solution file (re-read fresh each turn, may have changed since earlier turns):\n\n" + $content)}]' \
        <<<"$request_messages")
    fi
    request_messages=$(jq --arg content "$line" '. + [{role: "user", content: $content}]' <<<"$request_messages")

    payload=$(jq -n --arg model "$MODEL" --argjson messages "$request_messages" \
      '{model: $model, messages: $messages, stream: false, options: {temperature: 0.2}}')

    if ! response=$(curl -sf "${OLLAMA_HOST}/api/chat" -d "$payload"); then
      echo "tutor: could not reach ${OLLAMA_HOST}" >&2
      continue
    fi

    reply=$(echo "$response" | jq -r '.message.content')
    display_reply=$(process_highlights "$reply")
    echo "$display_reply"

    # Persist only the real conversation turn — the file context above
    # is intentionally left out so it doesn't accumulate turn over turn.
    messages=$(jq --arg content "$line" '. + [{role: "user", content: $content}]' <<<"$messages")
    messages=$(jq --arg content "$reply" '. + [{role: "assistant", content: $content}]' <<<"$messages")
  done

  echo
}

# Guards the interactive loop so this script can be `source`d for testing
# (see chat_test.sh) without dropping into a blocking read loop.
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
