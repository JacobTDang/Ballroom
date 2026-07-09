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
# Falls back to `nvim .` when there's no solution file (e.g. sandbox mode).
tmux select-pane -t "${SESSION}:main.0" -T "EDITOR"
SOLUTION_FILE=$(find "$WORKDIR" -maxdepth 1 -name 'solution.*' -type f | head -n1)
if [ -n "$SOLUTION_FILE" ]; then
  tmux send-keys -t "${SESSION}:main.0" "nvim '$SOLUTION_FILE'" C-m
else
  tmux send-keys -t "${SESSION}:main.0" "nvim ." C-m
fi

# pane 1: terminal (top right)
tmux split-window -h -t "${SESSION}:main.0" -c "$WORKDIR"
tmux select-pane -t "${SESSION}:main.1" -T "TERMINAL"
tmux send-keys -t "${SESSION}:main.1" "/bin/bash" C-m

# pane 2: tutor chat (bottom right)
tmux split-window -v -t "${SESSION}:main.1" -c "$WORKDIR"
tmux select-pane -t "${SESSION}:main.2" -T "TUTOR CHAT"
tmux send-keys -t "${SESSION}:main.2" "/usr/local/bin/tutor-chat.sh" C-m

tmux select-pane -t "${SESSION}:main.0"
exec tmux attach -t "$SESSION"
