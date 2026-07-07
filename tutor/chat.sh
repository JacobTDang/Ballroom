#!/usr/bin/env bash
set -euo pipefail

# Minimal chat CLI: sends messages to the host Ollama endpoint and prints
# responses. tutor_mode -> system prompt wiring is M2 work; for now this
# always runs with a plain default system prompt (full-assist-equivalent).
# See interview_prep_mvp_spec.md Section 3.4.

OLLAMA_HOST="${OLLAMA_HOST:-http://host.docker.internal:11434}"
MODEL="${TUTOR_MODEL:-qwen2.5-coder:1.5b}"
SYSTEM_PROMPT="${TUTOR_SYSTEM_PROMPT:-You are a concise coding interview tutor.}"

echo "tutor (${MODEL}) — connected to ${OLLAMA_HOST}. Ctrl-D to exit."

while IFS= read -r -p '> ' line; do
  [ -z "$line" ] && continue

  payload=$(jq -n \
    --arg model "$MODEL" \
    --arg system "$SYSTEM_PROMPT" \
    --arg prompt "$line" \
    '{model: $model, system: $system, prompt: $prompt, stream: false}')

  if ! response=$(curl -sf "${OLLAMA_HOST}/api/generate" -d "$payload"); then
    echo "tutor: could not reach ${OLLAMA_HOST}" >&2
    continue
  fi

  echo "$response" | jq -r '.response'
done

echo
