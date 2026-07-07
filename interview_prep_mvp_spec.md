# Interview Prep Environment — MVP Spec (for Claude Code handoff)

## 1. Product summary
A single-user, containerized practice environment for FT SWE interview prep. One Docker image contains a full dev environment (Go, C++, Python toolchains + debuggers/sanitizers), a terminal-based editor (Neovim), and a local AI tutor (via Ollama on the host) — all arranged in a `tmux` layout inside the container. A lightweight tracker (SQLite) logs timed practice attempts across coding patterns, debugging, concurrency, and multi-file implementation exercises.

Full design context and rationale lives in `interview_prep_plan.md` (companion doc, already written) — this spec is the scoped-down MVP cut of that plan, meant to be buildable in a first pass.

## 2. Goal for the MVP
Get one full working loop end-to-end:
**pick an exercise → work in a real terminal environment with AI help available → get an objective pass/fail → log it.**

Everything else (spaced repetition scheduling, dashboards, AI-usage analytics, per-exercise tutor scope tuning) is explicitly deferred — see Section 5.

## 3. MVP scope

### 3.1 Container / image
- **One unified Docker image** (not three per-language images), containing:
  - Go toolchain + `delve` + built-in race detector
  - C++ compiler (clang) + `cmake` + `gdb`/`lldb` + sanitizers (TSan, ASan, UBSan) as compile flags in exercise build scripts
  - Python + `pdb`/`debugpy` + `pytest`
  - Neovim (bare config is fine for MVP — no need to replicate a full personal config yet)
  - `tmux`
  - A minimal chat CLI script that sends messages to a local Ollama endpoint and prints responses
- **Entrypoint**: on `docker run`, immediately launches a `tmux` session with 3 panes: editor (nvim), terminal (shell), tutor (chat CLI).
- **Host networking**: chat CLI points at `host.docker.internal:11434` (confirm this resolves correctly on the target machine — see Section 6, this is a known risk to de-risk early).
- **Two run modes**:
  - `--exercise <id>`: mounts the specified exercise's repo/branch, starts a timer, tears down on exit.
  - `--sandbox`: mounts a persistent volume, no timer, no grading, survives across runs until explicitly reset.

### 3.2 Exercise definition (simplified for MVP)
A single YAML or JSON file per exercise, minimal fields only:
```
id, title, category (pattern | debug | concurrency | implementation | ai-assisted)
language (go | cpp | python)
time_limit_min
tutor_mode (syntax-only | hints-first | full-assist)   # single field, no per-exercise scope_notes yet
repo_path                                               # host path/branch to mount
test_command                                            # single shell command that runs the hidden test suite and exits 0/1
```
Defer: `tutor_scope_notes`, `logic_stuck_help` vs `syntax_help_always_on` split, `repro_reliability` flags, partial-credit thresholds. MVP is binary pass/fail from `test_command`'s exit code.

### 3.3 Verification (MVP cut)
- Each exercise ships with a hidden test script/directory, mounted but not visible in the editor until after submission.
- MVP only needs `test_command` to exit 0 (pass) or non-zero (fail) — no need for granular `tests_passed/tests_total` reporting yet, no API/SQL-specific harnesses. A unit-test-based exercise (Go `testing`, pytest, Catch2/GoogleTest) is enough to prove the loop works.
- API-endpoint and SQL-query verification types are real features but not MVP-blocking — defer to v2.

### 3.4 AI tutor (MVP cut)
- Single global permission per exercise (`tutor_mode` field above), not the fuller Part 7/8 permission matrix.
- `syntax-only`: tutor CLI's system prompt restricts it to fixing syntax/typos, no logic changes.
- `hints-first`: tutor gives hints, escalates to a direct answer only if asked twice (simple conversational rule in the prompt, not a hard mechanism yet).
- `full-assist`: no restriction — used for the ai-assisted category.
- Defer: separate always-on syntax help independent of mode, AI-interaction logging (`interaction_mode`, `accepted_as_is`), scope-creep prevention beyond a system prompt instruction.

### 3.5 Tracking (MVP cut)
Single SQLite table is enough for v1:
```
attempts: id, exercise_id, category, language, date, time_spent_min, result (pass/fail from test_command), notes
```
Defer: separate tables per category (Part 2's 8 entities), spaced-repetition scheduling, pattern-coverage dashboard, weekly review view. A flat log you can query manually (`SELECT * FROM attempts WHERE result='fail'`) is sufficient to prove value before building automation on top.

### 3.6 Sandbox (MVP cut)
- One persistent volume for the unified image (not per-language, since it's one image now).
- Manual reset via a documented command (e.g. `docker volume rm` + recreate) — no scripted "reset to base" automation yet.

## 4. Out of scope for MVP (explicitly deferred)
- Multi-architecture image builds (`amd64`/`arm64`) — build for your own machine first.
- Per-exercise `tutor_scope_notes` and the full tutor permission matrix from Part 8.
- API-endpoint and SQL-query-diff test harnesses — unit tests only for v1.
- AI-interaction logging and analytics (`interaction_mode`, `accepted_as_is` trends).
- Spaced-repetition scheduling and automatic review queueing.
- Pattern-coverage dashboard, weekly review auto-summary.
- Concurrency-specific stress-test-loop verification (flaky-by-design bug handling) — treat concurrency exercises like any other exercise for v1, revisit once the loop works.
- Any custom TUI app (Option A from Part 6) — MVP uses plain tmux + nvim.

## 5. Milestones
- **M0 — Environment**: Dockerfile builds successfully; `docker run -it <image>` opens tmux with 3 working panes; chat CLI can reach host Ollama and get a response.
- **M1 — Exercise loop**: exercise definition file loads; repo mounts correctly; timer runs; `test_command` executes and reports pass/fail; result logs to SQLite.
- **M2 — Tutor modes**: `syntax-only`, `hints-first`, `full-assist` produce visibly different tutor behavior on the same prompt.
- **M3 — Sandbox mode**: persistent volume survives container restarts; manual reset works.
- **M4 — First real exercise set**: 3-5 exercises seeded (at least one per category) to validate the format holds up beyond hypothetical design.

## 6. Known risks to de-risk early
- `host.docker.internal` networking — confirm the container can actually reach the host's Ollama instance before building anything else on top of it.
- Sanitizer output readability inside a `tmux` pane — TSan/ASan output can be verbose; worth checking it's usable in a narrow terminal pane rather than assuming it'll be fine.
- Neovim default config inside the image — bare nvim without LSP/plugins may be uncomfortable to actually use; decide minimum viable nvim config before M0 is "done."

## 7. Reference
Full rationale, alternatives considered, and the complete (non-MVP) design space are in the companion planning document (`interview_prep_plan.md`, Parts 1–10).
