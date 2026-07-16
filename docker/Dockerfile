# Unified practice environment image (Go + C++ + Python toolchains, nvim, tmux).
# See interview_prep_mvp_spec.md Section 3.1.
#
# Build from the repo root (paths below assume that build context):
#   docker build -f docker/Dockerfile -t ballroom-practice .

# Builds the ballroom launcher CLI itself (used inside the container for
# `ballroom submit` — `ballroom practice`/`ballroom sandbox` are invoked
# from the host).
FROM golang:1.25-bookworm AS ballroom-builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -o /ballroom ./cmd/ballroom

FROM ubuntu:24.04

# Identifies images built by this Dockerfile so the launcher can find and
# clean up its own stale/dangling images after a rebuild without ever
# touching unrelated images from other projects (see
# internal/orchestrator's ballroomImageLabel).
LABEL com.ballroom.app=practice-image

ENV DEBIAN_FRONTEND=noninteractive \
    GOBIN=/usr/local/bin \
    PATH=/usr/local/go/bin:/opt/pyenv/bin:/usr/local/bin:${PATH}

RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates curl git jq \
      build-essential clang cmake gdb lldb libclang-rt-18-dev llvm-18 \
      python3 python3-venv python3-pip \
      neovim tmux \
      golang-go \
    && rm -rf /var/lib/apt/lists/* \
    && ln -sf /usr/bin/llvm-symbolizer-18 /usr/local/bin/llvm-symbolizer

# delve (Go debugger); race detector ships with the toolchain (-race flag), no extra install.
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Python venv (Ubuntu 24.04 is PEP 668 externally-managed, so no system-wide pip installs).
# pdb/faulthandler are stdlib; only debugpy + pytest need installing.
RUN python3 -m venv /opt/pyenv \
    && /opt/pyenv/bin/pip install --no-cache-dir pytest debugpy

RUN useradd -m -s /bin/bash practice \
    && mkdir -p /workspace \
    && chown -R practice:practice /workspace /home/practice

COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY docker/clock.sh /etc/practice/clock.sh
COPY docker/tmux.conf /etc/practice/tmux.conf
COPY docker/nvim/ /home/practice/.config/nvim/
# Estimation reference for system-design sessions -- `less ~/back-of-envelope.md`
COPY docs/back-of-envelope.md /home/practice/back-of-envelope.md
COPY --from=ballroom-builder /ballroom /usr/local/bin/ballroom
RUN chmod +x /usr/local/bin/entrypoint.sh /etc/practice/clock.sh /usr/local/bin/ballroom \
    && chown -R practice:practice /home/practice/.config

USER practice
WORKDIR /workspace

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
