#!/usr/bin/env bash
set -euo pipefail

# Launches the practice session as a single tmux window split into three
# panes: editor (nvim) full-width across the top, with tutor chat and a
# terminal below it side by side — tutor under the problem statement,
# terminal under the editor. Switch panes with M-1/M-2/M-3 or Ctrl-Tab
# (see tmux.conf).

SESSION="${SESSION_NAME:-practice}"
WORKDIR="${WORKDIR:-/workspace}"
TMUX_CONF="${TMUX_CONF:-/etc/practice/tmux.conf}"

# NVIM_SOCKET is resolved and exported here -- before the tmux server
# starts below -- so it lands in the server's captured global
# environment: every pane tmux spawns afterward (send-keys, split-window)
# inherits that environment, which is what lets both the M-h keybinding
# (tmux.conf's run-shell string, `${NVIM_SOCKET:-...}`) and this script's
# own later `nvim --listen` invocation agree on the same socket path
# instead of each hardcoding the default separately.
export NVIM_SOCKET="${NVIM_SOCKET:-/tmp/ballroom-nvim.sock}"

# Container uptime at session start, in seconds since this container's
# kernel booted (the Docker Desktop Linux VM, on a Mac host) -- read and
# exported here for the same env-propagation reason as NVIM_SOCKET
# above. internal/session's Submit reads this back (PRACTICE_START_UPTIME)
# to compute time_spent_min against uptime instead of wall clock: the VM
# is suspended along with the host laptop, so uptime elapsed reflects
# real working time across a sleep/lid-close, unlike PRACTICE_STARTED_AT's
# wall-clock RFC3339 timestamp, which keeps advancing through a suspend
# (issue #229; see docker/clock.sh for the matching fix to the visible
# countdown). `|| true` degrades to leaving the var unset rather than
# aborting the whole session under `set -e` in the -- essentially
# theoretical -- case /proc/uptime can't be read; every consumer already
# treats an absent/empty value as "fall back to wall clock".
read -r PRACTICE_START_UPTIME _ < /proc/uptime || true
export PRACTICE_START_UPTIME

# A session created detached (-d, no client ever attached) has no
# established window size on this image's tmux (3.4) until a client
# attaches — and `split-window -p` (percentage) errors with "size
# missing" in that state regardless, even when the session is later
# given an explicit size, so percentage splits are unusable here. Read
# the real terminal size directly instead: `docker run -it` already
# gives this script's own stdin a real pty matching the host terminal,
# so `stty size` reports it without needing to wait for `tmux attach`.
# Fall back to tmux's own default (80x24) if stdin isn't a terminal at
# all (e.g. a non -it invocation).
if REAL_SIZE=$(stty size 2>/dev/null); then
  read -r REAL_ROWS REAL_COLS <<<"$REAL_SIZE"
else
  REAL_COLS=80
  REAL_ROWS=24
fi
BOTTOM_LINES=$((REAL_ROWS / 4))
[ "$BOTTOM_LINES" -lt 5 ] && BOTTOM_LINES=5

tmux -f "$TMUX_CONF" new-session -d -s "$SESSION" -n MAIN -c "$WORKDIR" -x "$REAL_COLS" -y "$REAL_ROWS"

# Visible countdown clock: timed exercise sessions carry PRACTICE_STARTED_AT
# and PRACTICE_TIME_LIMIT_MIN (forwarded by the host's orchestrator; the
# tutor's own interview-clock note reads the same two), so the status
# bar can show the user the clock the model already sees. Appended to
# status-right at the far-right edge; refreshes on the conf's
# status-interval. Sandbox sessions set neither var and keep the plain
# status bar. The numeric guard matters under set -e: a bare
# `[ ... -gt 0 ]` on a non-numeric value would abort the whole
# entrypoint.
#
# The deadline itself is computed in uptime terms (PRACTICE_START_UPTIME,
# read above), not from PRACTICE_STARTED_AT -- clock.sh diffs it against
# a live /proc/uptime read on every status-bar refresh (issue #229).
# PRACTICE_STARTED_AT's presence is still the "is this a real timed
# session" guard, matching PRACTICE_TIME_LIMIT_MIN's own presence check
# in the exercise-session args (see internal/orchestrator's
# exerciseRunArgs, which sets both together).
if [ -n "${PRACTICE_STARTED_AT:-}" ] \
  && [ -n "${PRACTICE_START_UPTIME:-}" ] \
  && [[ "${PRACTICE_TIME_LIMIT_MIN:-}" =~ ^[0-9]+$ ]] \
  && [ "$PRACTICE_TIME_LIMIT_MIN" -gt 0 ]; then
  START_UPTIME_S="${PRACTICE_START_UPTIME%.*}"
  DEADLINE_UPTIME=$((START_UPTIME_S + PRACTICE_TIME_LIMIT_MIN * 60))
  tmux set -g status-right "$(tmux show -gv status-right)#[fg=#E0468C]·#[default]  #(/etc/practice/clock.sh $DEADLINE_UPTIME) "
fi

# pane 0 (top, full width): editor. Open directly into the solution file
# to implement, not the netrw directory listing — glob rather than
# hardcode an extension, since it varies by language (.go/.py/.cpp/.hpp).
# When the exercise ships a problem.md (statement + examples +
# constraints — see internal/verify's sibling authoring convention),
# open it as a read-only left split so it reads like NeetCode's own
# two-pane layout, with focus landing on the solution file for editing.
# Falls back to a single-file open (or `nvim .` with no solution file at
# all, e.g. sandbox mode) when there's no problem.md.
#
# --listen exposes an RPC socket so the tutor pane can drive
# highlights/notes in the running nvim instance (issue #24), at the
# well-known /tmp location NVIM_SOCKET was resolved to above (shared by
# every process in the container, regardless of WORKDIR). rm -f first: a
# stale socket file from a previous run in the same container would
# otherwise make nvim refuse to bind.
rm -f "$NVIM_SOCKET"
SOLUTION_FILE=$(find "$WORKDIR" -maxdepth 1 -name 'solution.*' -type f | head -n1)
# Prefer the plain-text render (problem.txt, written by the host's
# PrepareWorkspace) -- clean structured text with no markdown markers.
# problem.md is the fallback for workspaces prepared by an older host
# binary that didn't render it.
PROBLEM_FILE="$WORKDIR/problem.txt"
[ -f "$PROBLEM_FILE" ] || PROBLEM_FILE="$WORKDIR/problem.md"
if [ -n "$SOLUTION_FILE" ] && [ -f "$PROBLEM_FILE" ]; then
  # -c arguments are Vim ex-commands, not shell — no shell quoting inside
  # them (paths are already fully resolved by this point, so this is safe
  # even though it wouldn't handle a path containing spaces).
  tmux send-keys -t "${SESSION}:MAIN.0" "nvim --listen '$NVIM_SOCKET' -c \"vsplit $PROBLEM_FILE\" -c 'set readonly nomodifiable' -c 'wincmd l' '$SOLUTION_FILE'" C-m
elif [ -n "$SOLUTION_FILE" ]; then
  tmux send-keys -t "${SESSION}:MAIN.0" "nvim --listen '$NVIM_SOCKET' '$SOLUTION_FILE'" C-m
else
  tmux send-keys -t "${SESSION}:MAIN.0" "nvim --listen '$NVIM_SOCKET' ." C-m
fi

# Split off a bottom row (~25% of the real window height, computed
# above as BOTTOM_LINES — a fixed line count, not -p/--percentage,
# since percentage splits error on this image's tmux) below the editor,
# then split that row into tutor (left, under the problem statement) and
# terminal (right, under the editor). Pane indices are assigned in
# creation order — 0 editor, 1 tutor, 2 terminal — which tmux.conf's
# M-1/M-2/M-3 and M-q bindings target directly.
tmux split-window -v -l "$BOTTOM_LINES" -t "${SESSION}:MAIN.0" -c "$WORKDIR"
tmux split-window -h -t "${SESSION}:MAIN.1" -c "$WORKDIR"

# pane 1 (bottom-left): tutor chat. NVIM_SOCKET tells the tutor agent
# where to reach the editor pane's RPC server (see above).
tmux send-keys -t "${SESSION}:MAIN.1" "NVIM_SOCKET='$NVIM_SOCKET' /usr/local/bin/ballroom tutor" C-m

# pane 2 (bottom-right): terminal
tmux send-keys -t "${SESSION}:MAIN.2" "/bin/bash" C-m

tmux select-pane -t "${SESSION}:MAIN.0"
exec tmux attach -t "$SESSION"
