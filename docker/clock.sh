#!/usr/bin/env bash
# Prints one tmux status-bar segment: time remaining until a deadline,
# as MM:SS in the session palette's teal, switching to gold under 10
# minutes and a bold red "TIME UP" once the deadline passes. Invoked by
# tmux's own #() substitution in status-right (see entrypoint.sh, which
# computes the deadline from PRACTICE_STARTED_AT + PRACTICE_TIME_LIMIT_MIN),
# refreshing on the conf's status-interval.
#
# The optional second argument substitutes "now" so the color
# transitions are testable without waiting out a real session:
#   clock.sh <deadline-epoch> [now-epoch]
set -euo pipefail

DEADLINE="${1:?usage: clock.sh <deadline-epoch-seconds> [now-epoch-seconds]}"
NOW="${2:-$(date +%s)}"

REMAINING=$((DEADLINE - NOW))
if [ "$REMAINING" -le 0 ]; then
  printf '#[bold,fg=#F03C3C]TIME UP#[default]'
  exit 0
fi

COLOR='#2FA6A6'
if [ "$REMAINING" -lt 600 ]; then
  COLOR='#E8A93C'
fi
printf '#[fg=%s]%02d:%02d#[default]' "$COLOR" $((REMAINING / 60)) $((REMAINING % 60))
