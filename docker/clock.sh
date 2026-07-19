#!/usr/bin/env bash
# Prints one tmux status-bar segment: time remaining until a deadline,
# bracketed like the host TUI's own meters (progressBar,
# internal/tui/homeboard.go) so the clock reads as the same kind of
# instrument as the rest of the app -- as MM:SS in the session palette's
# teal, switching to gold under 10 minutes and a bold red "TIME UP" once
# the deadline passes, with the brackets themselves always the same dim
# structural gray the host meters frame theirs in. Invoked by tmux's own
# #() substitution in status-right (see entrypoint.sh, which computes
# the deadline from PRACTICE_START_UPTIME + PRACTICE_TIME_LIMIT_MIN),
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

# BRACKET is palette.PaleGray -- the same dim structural gray
# internal/tui/homeboard.go's progressBar frames its own meters with
# (checkDimStyle), so this clock and the host's progress bars read as
# one family of instrument, not two different styles.
BRACKET='#D9D3C4'

DEADLINE_S="${DEADLINE%.*}"
NOW_S="${NOW%.*}"
REMAINING=$((DEADLINE_S - NOW_S))
if [ "$REMAINING" -le 0 ]; then
  printf '#[fg=%s][#[fg=#F03C3C]#[bold]TIME UP#[nobold]#[fg=%s]]#[default]' "$BRACKET" "$BRACKET"
  exit 0
fi

COLOR='#2FA6A6'
if [ "$REMAINING" -lt 600 ]; then
  COLOR='#E8A93C'
fi
printf '#[fg=%s][#[fg=%s]%02d:%02d#[fg=%s]]#[default]' "$BRACKET" "$COLOR" $((REMAINING / 60)) $((REMAINING % 60)) "$BRACKET"
