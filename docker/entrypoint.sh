#!/usr/bin/env bash
set -euo pipefail

# Launches the 3-pane tmux session: editor (nvim) | terminal (shell) /
# tutor (chat CLI). See interview_prep_mvp_spec.md Section 3.1.

SESSION="${SESSION_NAME:-practice}"
WORKDIR="${WORKDIR:-/workspace}"
TMUX_CONF="${TMUX_CONF:-/etc/practice/tmux.conf}"

tmux -f "$TMUX_CONF" new-session -d -s "$SESSION" -n main -c "$WORKDIR"

# pane 0: editor (left, full height)
tmux send-keys -t "${SESSION}:main.0" "nvim ." C-m

# pane 1: terminal (top right)
tmux split-window -h -t "${SESSION}:main.0" -c "$WORKDIR"
tmux send-keys -t "${SESSION}:main.1" "/bin/bash" C-m

# pane 2: tutor chat (bottom right)
tmux split-window -v -t "${SESSION}:main.1" -c "$WORKDIR"
tmux send-keys -t "${SESSION}:main.2" "/usr/local/bin/tutor-chat.sh" C-m

tmux select-pane -t "${SESSION}:main.0"
exec tmux attach -t "$SESSION"
