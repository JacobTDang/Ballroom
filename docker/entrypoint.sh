#!/usr/bin/env bash
set -euo pipefail

# Launches the 3-pane tmux session: editor (nvim) | terminal (shell) /
# tutor (chat CLI). See interview_prep_mvp_spec.md Section 3.1.

SESSION="${SESSION_NAME:-practice}"
WORKDIR="${WORKDIR:-/workspace}"
TMUX_CONF="${TMUX_CONF:-/etc/practice/tmux.conf}"

tmux -f "$TMUX_CONF" new-session -d -s "$SESSION" -n main -c "$WORKDIR"

# pane 0: editor (left, full height). Open directly into the solution
# file to implement, not the netrw directory listing — glob rather than
# hardcode an extension, since it varies by language (.go/.py/.cpp/.hpp).
# When the exercise ships a problem.md (statement + examples +
# constraints — see internal/verify's sibling authoring convention), open
# it as a read-only left split so it reads like NeetCode's own two-pane
# layout, with focus landing on the solution file for editing. Falls back
# to a single-file open (or `nvim .` with no solution file at all, e.g.
# sandbox mode) when there's no problem.md.
#
# --listen exposes an RPC socket so the tutor pane (pane 2) can drive
# highlights/notes in the running nvim instance (issue #24). The socket
# path is a well-known /tmp location (shared by every process in the
# container, regardless of WORKDIR) rather than something discovered via
# tmux env propagation — we're already constructing both panes' commands
# right here, so just pass it to both explicitly. rm -f first: a stale
# socket file from a previous run in the same container would otherwise
# make nvim refuse to bind.
NVIM_SOCKET="${NVIM_SOCKET:-/tmp/ballroom-nvim.sock}"
rm -f "$NVIM_SOCKET"
tmux select-pane -t "${SESSION}:main.0" -T "EDITOR"
SOLUTION_FILE=$(find "$WORKDIR" -maxdepth 1 -name 'solution.*' -type f | head -n1)
PROBLEM_FILE="$WORKDIR/problem.md"
if [ -n "$SOLUTION_FILE" ] && [ -f "$PROBLEM_FILE" ]; then
  # -c arguments are Vim ex-commands, not shell — no shell quoting inside
  # them (paths are already fully resolved by this point, so this is safe
  # even though it wouldn't handle a path containing spaces).
  tmux send-keys -t "${SESSION}:main.0" "nvim --listen '$NVIM_SOCKET' -c \"vsplit $PROBLEM_FILE\" -c 'set readonly nomodifiable' -c 'wincmd l' '$SOLUTION_FILE'" C-m
elif [ -n "$SOLUTION_FILE" ]; then
  tmux send-keys -t "${SESSION}:main.0" "nvim --listen '$NVIM_SOCKET' '$SOLUTION_FILE'" C-m
else
  tmux send-keys -t "${SESSION}:main.0" "nvim --listen '$NVIM_SOCKET' ." C-m
fi

# pane 1: terminal (top right)
tmux split-window -h -t "${SESSION}:main.0" -c "$WORKDIR"
tmux select-pane -t "${SESSION}:main.1" -T "TERMINAL"
tmux send-keys -t "${SESSION}:main.1" "/bin/bash" C-m

# pane 2: tutor chat (bottom right). NVIM_SOCKET tells chat.sh where to
# reach the editor pane's RPC server (see above).
tmux split-window -v -t "${SESSION}:main.1" -c "$WORKDIR"
tmux select-pane -t "${SESSION}:main.2" -T "TUTOR CHAT"
tmux send-keys -t "${SESSION}:main.2" "NVIM_SOCKET='$NVIM_SOCKET' /usr/local/bin/tutor-chat.sh" C-m

tmux select-pane -t "${SESSION}:main.0"
exec tmux attach -t "$SESSION"
