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
SYSTEM_PROMPT="${TUTOR_SYSTEM_PROMPT:-$DEFAULT_PROMPT}"

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
  echo "$reply"

  messages=$(jq --arg content "$reply" '. + [{role: "assistant", content: $content}]' <<<"$messages")
done

echo
