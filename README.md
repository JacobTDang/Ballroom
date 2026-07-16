# Ballroom

A terminal app for technical-interview practice: coding exercises with
hidden tests, system-design and behavioral mock interviews with an LLM
interviewer, and an in-session AI tutor — all running locally in a
Docker practice environment.

<!-- screenshot: home menu -->
<!-- screenshot: a practice session (editor + tutor + terminal panes, clock in the status bar) -->

Every session is a tmux window with three panes: the problem statement
and your editor (nvim) on top, the tutor chat and a terminal below.
You write real code (or a real design doc), submit against hidden
tests (or a hidden grading rubric), and the result lands in a local
progress tracker.

## Requirements

- Docker (the practice environment is a container; the image builds
  itself on first launch)
- Go 1.25+ (to build the CLI)
- An LLM for the tutor — either of:
  - [Ollama](https://ollama.com) running locally (default model
    `llama3.1:8b`), or
  - an [OpenRouter](https://openrouter.ai) API key for hosted models

## Install & run

```sh
go build -o ballroom ./cmd/ballroom
./ballroom
```

The first launch runs boot checks (Docker daemon, practice image —
built automatically — Ollama, tutor model) and lands on the home menu:
**Practice** (pick a category and problem), **Daily** (today's pick,
one keypress), **Sandbox** (ungraded scratch environment), **Stats**
(progress, recent attempts, rubric weak spots), **Settings** (models).

### Configuring the tutor models

Through Settings in the TUI, or directly:

```sh
ballroom config set-model llama3.1:8b                    # local Ollama
ballroom config set-model openrouter:<model-slug>        # hosted
ballroom config set-key <openrouter-api-key>
ballroom config set-orchestrator-model <tag|none>        # optional routing model
ballroom config set-grader-model <tag|none>              # optional dedicated grader
```

The tutor needs a model with **real tool-calling support** (the model
picker probes and shows a "(no tools)" badge next to models that
can't). Models that only narrate tool calls as text get a degraded
JSON-fallback mode; models with native support get the full agent,
and OpenRouter models stream their replies progressively.
`TUTOR_STREAM=on|off` overrides streaming per invocation.

## The three tracks

**Coding** — DSA (the NeetCode 150, by topic), debugging,
concurrency, implementation, and OO-design exercises in Python, Go,
and C++. Hidden tests are mounted only when you submit: the visible
starter must fail them, your job is to make them pass.

**System design** — the
[system-design-primer](https://github.com/donnemartin/system-design-primer)
questions as guided **coach** sessions (the 4-step method, one step at
a time) and timed **interviewer** mocks (bare prompt, you scope it,
45 minutes). A hidden per-question rubric grades your `solution.md`
on submit. `docs/system-design-roadmap.md` is the full curriculum,
start there.

**Behavioral** — eight classic "tell me about a time…" questions as
**story-coach** sessions (build a STAR answer one section at a time)
and timed interviewer mocks. Graded against STAR rubrics: situation
specificity, stakes, ownership, evidence, reflection.

## The practice loop

1. Pick a problem (or press `2` for **Daily** — a date-stable pick
   among due and unsolved problems).
2. Work in the session; talk to the tutor as much as the mode allows
   (coding exercises choose syntax-only / hints-first / full-assist
   per exercise; design and behavioral sessions choose coach or
   interviewer per session).
3. Submit with `M-q`. Coding: hidden tests run. Design/behavioral:
   the rubric (and, for solved primer questions, a reference design)
   appears in your workspace and the grader model grades your
   solution dimension by dimension — pass/fail plus a summary, both
   recorded.
4. The picker resurfaces what needs attention: **mock due** (coach
   pass whose interviewer mock is untouched), **review due** (a
   failure ≥3 days old, or anything solved but untouched for 30
   days) — due problems float to the top.
5. **Stats → Rubric weak spots** ranks the rubric dimensions you keep
   losing points on; when one keeps showing up, that's your next
   study block.

## Inside a session

| Key | Action |
|---|---|
| `M-1` / `M-2` / `M-3` | jump to editor / tutor / terminal pane |
| `Ctrl-Tab` | cycle panes |
| `M-q` | submit (pre-types the command; Enter confirms) |
| `M-h` | toggle tutor highlights in the editor |
| `Ctrl-D` (tutor pane) | exit the tutor chat |

Timed sessions show a countdown clock at the right of the status bar
(gold under 10 minutes, red past time). The tutor pane scrolls with
PgUp/PgDn or the mouse wheel, renders code as editor cards, and shows
the tools the model calls in real time. Estimation help for design
sessions: `less ~/back-of-envelope.md` in the terminal pane.

**Voice input**: macOS built-in dictation works directly into the
tutor pane — press the dictation shortcut (default: `fn` twice) with
the tutor pane focused and speak. Nothing to configure; the container
never needs microphone access.

## Development

```sh
go test -race ./internal/...   # unit/integration tests (mocked providers)
go vet ./cmd/... ./internal/...
go run ./cmd/verify-exercises  # structural check of every exercise
go run ./cmd/tutor-eval        # live model-behavior eval (needs Ollama; slow)
```

Exercises and hidden tests live on the host (`exercises/`, `tests/`)
and need no image rebuild to change. Changes to `cmd/`, `internal/`,
or `docker/` are picked up by the boot screen's image build on the
next launch, or rebuild manually:

```sh
docker build -f docker/Dockerfile -t ballroom-practice .
```
