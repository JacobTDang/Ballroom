#!/usr/bin/env bash
set -euo pipefail

# Minimal chat CLI: sends messages to the host Ollama endpoint
# (host.docker.internal:11434) and prints responses. System prompt varies by
# tutor_mode (syntax-only | hints-first | full-assist), per
# interview_prep_mvp_spec.md Section 3.4.
#
# TODO(M0): pin default Ollama model.
# TODO(M2): wire tutor_mode -> system prompt variants.

echo "tutor chat: not yet implemented" >&2
exit 1
