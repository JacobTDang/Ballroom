# Unified practice environment image (Go + C++ + Python toolchains, nvim, tmux).
# See interview_prep_mvp_spec.md Section 3.1. Not yet buildable — M0 work.

FROM ubuntu:24.04

# TODO(M0): Go toolchain + delve
# TODO(M0): clang + cmake + gdb/lldb + sanitizers (TSan/ASan/UBSan)
# TODO(M0): Python + pdb/debugpy + pytest
# TODO(M0): neovim (minimum viable config, see spec Section 6 risk)
# TODO(M0): tmux
# TODO(M0): tutor chat CLI (see tutor/chat.sh) + its runtime deps
# TODO(M0): non-root user, /workspace convention

COPY entrypoint.sh /usr/local/bin/entrypoint.sh
COPY tmux.conf /etc/practice/tmux.conf
RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
