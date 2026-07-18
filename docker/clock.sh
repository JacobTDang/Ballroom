#!/usr/bin/env bash
# Prints one tmux status-bar segment: time remaining until a deadline,
# as MM:SS in the session palette's teal, switching to gold under 10
# minutes and a bold red "TIME UP" once the deadline passes. Invoked by
# tmux's own #() substitution in status-right (see entrypoint.sh, which
# computes the deadline from PRACTICE_START_UPTIME + PRACTICE_TIME_LIMIT_MIN),
# refreshing on the conf's status-interval.
#
# Both times are container uptime (/proc/uptime's first field, seconds
# since this container's kernel booted), not wall clock -- the Docker
# Desktop Linux VM a session runs in is suspended along with the host
# laptop, so uptime does not advance across a lid-close the way
# `date +%s` does, and a lunch break no longer resumes to a false
# TIME UP (issue #229). /proc/uptime's fields carry two decimals (e.g.
# "12345.67"); truncated to whole seconds below -- bash arithmetic has
# no float support, and a countdown clock has no use for sub-second
# precision anyway.
#
# The optional second argument substitutes "now" so the color
# transitions are testable without waiting out a real session:
#   clock.sh <deadline-uptime-seconds> [now-uptime-seconds]
set -euo pipefail

DEADLINE="${1:?usage: clock.sh <deadline-uptime-seconds> [now-uptime-seconds]}"
if [ -n "${2:-}" ]; then
  NOW="$2"
else
  read -r NOW _ < /proc/uptime
fi

DEADLINE_S="${DEADLINE%.*}"
NOW_S="${NOW%.*}"
REMAINING=$((DEADLINE_S - NOW_S))
if [ "$REMAINING" -le 0 ]; then
  printf '#[bold,fg=#F03C3C]TIME UP#[default]'
  exit 0
fi

COLOR='#2FA6A6'
if [ "$REMAINING" -lt 600 ]; then
  COLOR='#E8A93C'
fi
printf '#[fg=%s]%02d:%02d#[default]' "$COLOR" $((REMAINING / 60)) $((REMAINING % 60))
