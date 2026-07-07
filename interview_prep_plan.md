# FT SWE Interview Prep: Plan + Practice Tool Spec

## Part 1: The Prep Plan

### Where you actually stand
You've got real backend depth (CQRS, event sourcing, saga pattern, idempotency, gRPC/AIP design) from Saige work — that's ahead of most candidates on system design fundamentals. The gap for "hard interviews" is usually: (1) LeetCode pattern speed/recall under pressure, and (2) structuring a system design answer in 35-45 min instead of just knowing the concepts.

### Timeline (assume ~12-16 weeks before your target interview window)

**Phase 1 (Weeks 1-4): Pattern building — coding**
- Cover one pattern per 2-3 days, not one problem at a time. Patterns: two pointers, sliding window, fast/slow pointers, binary search variants, DFS/BFS (graph + tree), backtracking, heap/priority queue, intervals, top-K, union-find, trie, DP (1D → 2D → on trees/graphs).
- 3-5 problems per pattern: 1 easy (confirm pattern), 2 medium (core reps), 1-2 hard (stress test).
- Track every problem you can't solve in 25 min cold — that's your re-drill list, not a one-and-done.

**Phase 2 (Weeks 5-8): System design foundations**
- Build the standard toolkit: back-of-envelope estimation, load balancing, caching strategies, DB choice (SQL vs NoSQL, sharding, replication), message queues, CDN, consistent hashing.
- Do 1 full design per week end-to-end (URL shortener → rate limiter → chat system → news feed → ride-share → distributed cache). Time-box to 45 min, then review against a reference solution.
- Since you already know CQRS/event sourcing/sagas from Saige, explicitly practice explaining *when* to reach for them vs. simpler approaches — interviewers probe for judgment, not vocabulary.

**Phase 3 (Weeks 9-12): Hard-mode integration**
- Mixed timed sets: 2 mediums back-to-back in 45 min, cold, no pattern hint given.
- Hard-only week: pick 10-15 hard problems across patterns you're weakest on.
- System design: add failure-mode pressure ("now the DB is down," "now traffic 10x's") mid-interview — this is where most candidates crack.
- Add "unfamiliar codebase" reps: grab an open-source repo (or one of your own you haven't touched in months) and (a) implement a small feature that touches 2-3 files, (b) intentionally reintroduce a bug you fix from an old commit and re-find it. This trains the specific skill of *reading and navigating* code fast, which pure LeetCode doesn't touch and which shows up a lot in pair-programming-style onsites.
- Start mock interviews (peer, Pramp-style, or paid) — at least 1/week from here on.

**Phase 4 (Weeks 13-16): Company-specific + polish**
- Target company OA patterns (CodeSignal-style for Capital One TIP, etc.) and behavioral/STAR stories tied to your Saige and side-project work.
- Full mock loops (2-3 coding + 1 system design + 1 behavioral) simulating a real onsite.

### Cadence that actually works
- Daily: 1-2 coding problems, min 5 days/week.
- Weekly: 1 system design deep-dive + 1 mock (once in Phase 3+).
- Always time yourself. Untimed practice doesn't transfer to interview pressure.

---

## Part 2: Practice Tool Spec (no implementation yet)

### Purpose
A lightweight offline-first tracker that closes the loop your memory can't: which patterns you've actually drilled, which problems need re-review, and how your system design reps are progressing — without depending on LeetCode's own (weak) tracking.

### Practice categories (the taxonomy)
Everything below falls into one of these buckets. Concurrency and AI-assisted work are called out separately from plain implementation/debugging because they test distinct skills and are common enough in "hard interview" formats to deserve their own tracking:

1. **Coding patterns** — LeetCode-style, timed, pattern-tagged.
2. **Concurrency / race conditions** — cross-cutting; can show up as either a debug task or an implementation task, but worth tracking on its own since it's consistently where candidates get exposed.
3. **Feature implementation** — multi-file, "add this to an existing codebase."
4. **AI-assisted implementation** — simulating interviews that explicitly allow AI tool use during implementation.
5. **Debugging** — given broken code, find + fix.
6. **System design** — end-to-end design sessions.
7. **Mocks / behavioral** — full-loop simulation.

### Core entities
1. **Problem attempt**
   - `id`, `title`, `platform` (LeetCode/other), `pattern_tags[]`, `difficulty`, `date`, `time_to_solve_min`, `result` (solved-clean / solved-with-hints / failed), `notes`, `review_after` (spaced-repetition date)
2. **Pattern**
   - `name`, `status` (not-started / drilling / solid), `problems_attempted`, `problems_needing_review`
3. **System design session**
   - `id`, `topic`, `date`, `time_spent_min`, `self_rated_score` (1-5), `weak_areas[]` (estimation, data modeling, scaling, tradeoffs, communication), `notes`
4. **Mock interview log**
   - `date`, `type` (coding/system-design/behavioral), `feedback_summary`, `action_items[]`
5. **Implementation task** (multi-file, "build this feature into an existing codebase")
   - `id`, `title`, `language` (Go / C++ / Python), `scope` (single-file/multi-file), `files_touched[]`, `time_spent_min`, `result` (complete / partial / stuck), `what_slowed_you_down` (e.g. navigating unfamiliar code, test setup, API design decisions), `notes`
6. **Debug task** (given broken code, find + fix)
   - `id`, `title`, `language` (Go / C++ / Python), `bug_category`, `time_to_locate_min`, `time_to_fix_min`, `result`, `how_found` (read + reason / added logging / used debugger / test-driven), `notes`
   - Suggested `bug_category` values per language, since the common failure modes differ enough to track separately:
     - **Go**: goroutine leaks, race conditions (unprotected shared state), channel deadlocks, nil pointer dereference, incorrect error handling/swallowed errors, slice aliasing surprises, closure-over-loop-variable bugs.
     - **C++**: dangling pointers/use-after-free, memory leaks, undefined behavior (uninitialized vars, out-of-bounds access), off-by-one, iterator invalidation, object slicing, incorrect const-correctness or reference vs. value semantics.
     - **Python**: mutable default arguments, late-binding closures in loops, type mismatches, off-by-one, incorrect mutability assumptions (shared references), exception swallowing, GIL-related false assumptions about thread safety.
7. **Concurrency task** (race conditions, deadlocks, synchronization — can be framed as either debug or implementation)
   - `id`, `title`, `language` (Go / C++ / Python), `frame` (debug / implement-from-scratch), `primitive_involved` (mutex, channel, atomic, condition variable, goroutine/thread pool, async/await, GIL interaction), `failure_mode` (data race, deadlock, livelock, starvation, missed signal), `detection_method` (manual reasoning / race detector / sanitizer / stress-test loop), `time_to_locate_min`, `time_to_fix_min`, `result`, `notes`
   - This is worth isolating from generic debugging because concurrency bugs are often non-deterministic — "I couldn't repro it" is itself a data point worth tracking, not a failure to hide.
8. **AI-assisted implementation task** (practicing interviews that permit AI tool use)
   - `id`, `title`, `language`, `ai_model_used` (local model name/version), `scope`, `time_spent_min`, `what_ai_helped_with` (boilerplate, API lookup, test generation, debugging suggestions), `what_you_had_to_correct` (where the AI's output was wrong or you had to redirect it), `result`, `notes`
   - The skill being tested here isn't "can you code" — it's "can you supervise and correct an AI's output fast while still explaining your reasoning out loud." Track corrections separately; that's the part interviewers are actually watching for.

### Key features
- **Spaced repetition queue**: surfaces problems you got wrong or slow on, resurfacing at increasing intervals (1 day → 3 days → 1 week → 1 month).
- **Pattern coverage dashboard**: visual gap map — which patterns have zero/few reps, so you're not unconsciously avoiding weak spots (a very common failure mode).
- **Timer built in**: every attempt logs actual time, not estimated — this is the #1 thing self-reporting gets wrong.
- **System design rubric**: a fixed checklist per session (requirements clarification, estimation, high-level design, deep dive, tradeoffs, failure handling) so you self-score consistently instead of vaguely.
- **Weekly review view**: auto-generated summary — problems solved, patterns drilled, mocks completed, streaks — so progress is visible without manual tallying.
- **Implementation/debug tracking**: separate from the LeetCode pattern dashboard, since these test a different muscle (reading unfamiliar code, navigating a codebase, isolating a fault) rather than pattern recall. Track `time_to_locate` vs `time_to_fix` separately for debug tasks — most people are slow at *finding* the bug, not fixing it once found, and that distinction tells you what to actually practice.
- **Offline-first**: local storage/local file (e.g., SQLite or flat JSON) as source of truth, since you don't need this tied to an account or service.

### Suggested build approach (when you're ready to implement)
- Simplest version: a local SQLite DB + a small CLI or terminal TUI (fits your backend-leaning skillset, fast to build, no frontend overhead).
- Upgrade path: add a lightweight local web UI later (e.g., a small Go or Node service + static frontend) if you want dashboards/graphs.
- This is a good fit for Claude Code once you're ready to build — it's a self-contained, scriptable project with a clear data model, similar in shape to your other side-project handoffs.

### What this deliberately does NOT do
- Doesn't replace LeetCode/Pramp/mock platforms — it's a tracking + spaced-repetition layer on top of them.
- Doesn't auto-grade code — you self-assess result and time; keeping that manual keeps you honest and forces reflection.

---

## Part 3: Environment Setup Plan (before building anything)

This is the physical setup you need in place before you can generate exercises against the categories above. Still just planning — no implementation yet.

### 1. Language toolchains
- **Go**: standard toolchain + `delve` (debugger) + built-in race detector (`go test -race` / `go run -race`) — this last one is essential for the concurrency category, since it catches races you'd never spot by reading.
- **C++**: a compiler (clang or gcc) + `gdb` or `lldb` + sanitizers (`-fsanitize=thread` for races/deadlocks, `-fsanitize=address` for memory bugs, `-fsanitize=undefined` for UB). These sanitizers are what make C++ concurrency/memory bugs practice-able instead of guesswork.
- **Python**: standard interpreter + `pdb`/`debugpy` + `faulthandler` for hangs; `threading`/`asyncio` for concurrency exercises (note Python's GIL means "race conditions" here often look different than Go/C++ — more about shared mutable state and ordering than true parallel memory access).

### 2. Repo staging structure
- One working repo per language (`practice-go/`, `practice-cpp/`, `practice-python/`), each with a `main` branch kept clean.
- Feature/bug branches per exercise: e.g. `bug/goroutine-leak-01`, `feature/rate-limiter-multi-file`. Reset from `main` each time so exercises don't contaminate each other.
- A small pool of "host" projects per language sized to require real navigation (not toy single-file demos) — e.g. a basic REST service, a CLI tool, a small concurrent worker pool — into which bugs/features get injected. Reusing 2-3 host projects per language repeatedly (with different injected issues) mirrors "codebase you don't own" better than fresh toy examples every time.

### 3. Local AI model setup (for the AI-assisted category)
- Goal: a **local**, offline model so practice sessions are consistent, private, and not rate-limited or subject to API drift — you want the same "AI collaborator" behavior every rep.
- Your hardware (M4 Max, 36GB unified memory) comfortably fits a strong coding-specialized model:
  - **Primary**: `qwen3-coder:30b` via Ollama — 30B MoE (3.3B active params), ~19GB at Q4_K_M, 256K context. Currently the strongest local coding model for hardware in this range.
  - **Secondary/lighter option**: `devstral:24b` — ~14GB, leaves more headroom, good for comparison or running alongside other apps.
  - **Deliberately weaker model**: keep a smaller model (7-9B range) on hand too. The real value of the AI-assisted category isn't practicing with your best possible tool — it's practicing catching and correcting a *mediocre* AI's mistakes fast while narrating your reasoning, since that's closer to what a sandboxed interview tool is likely to give you. Training only against a great model teaches the wrong instinct.
  - Make sure Ollama is updated (0.19+) — on 32GB+ Macs it now defaults to the MLX backend instead of Metal, roughly doubling generation speed on the same hardware, so an outdated version leaves real performance on the table.
- Keep this local setup separate from your actual paid Claude/Copilot access — the point is to simulate the "AI tool provided by the interviewer" constraint, which is usually sandboxed and weaker than your daily driver.

### 4. Interim tracking (before the real tool exists)
- Until the SQLite tracker is built, a plain spreadsheet or even a dated markdown log covering the same fields (time to locate/fix, category, result) is enough to not lose data — the discipline of logging matters more than the medium right now.

### 5. Timer discipline
- A visible countdown timer (phone, or any timer app) for every single rep, coding or design — this is the cheapest, highest-leverage part of the whole setup and needs zero building.

---

## Part 4: Desktop App + Containerization Architecture (planning only)

Goal: a single app you launch on any machine that gives you a real terminal, a real Go/C++/Python toolchain (debuggers, sanitizers, race detector included), and consistent behavior regardless of host OS.

### Why containerize
- Solves "works on any machine" directly — the practice environment (compilers, debuggers, sanitizers) lives inside container images, not on the host. Windows, Mac, Linux all get identical behavior.
- Lets you reset an exercise instantly (tear down and recreate the container) instead of manually cleaning up a host environment between reps.
- Keeps the toolchain isolated from whatever you already have installed for actual work (Saige, side projects) — no version conflicts.

### Proposed layers
1. **Desktop shell** — a cross-platform app wrapper (Tauri or Electron) that's just a window hosting: an exercise browser/picker, a timer, an embedded code editor pane, and an embedded terminal pane.
   - Tauri: smaller binary, Rust-based backend, lighter resource use.
   - Electron: heavier, but more mature ecosystem for embedding terminals (`xterm.js` + `node-pty`) and editors (Monaco). Given the goal is "works everywhere reliably" over "smallest binary," Electron's maturity is the safer default; Tauri is the leaner alternative if binary size/resource use matters more to you.
2. **Local orchestration layer** — talks to the Docker daemon on the host to spin up/tear down per-exercise containers, mount the relevant repo/branch into them, and pipe a shell session into the desktop terminal pane.
3. **Practice containers** — one image per language (`practice-go`, `practice-cpp`, `practice-python`), each pre-loaded with its toolchain, debugger, and relevant sanitizers/race detector from Part 3. Exercise repos get mounted in as volumes so containers stay generic and reusable.
4. **Local AI model (host-level, not containerized)** — this is a real constraint worth flagging: Docker on macOS doesn't get direct Metal GPU passthrough, so running Ollama *inside* a container on your M4 Max would lose the GPU acceleration that makes it fast. Better plan: run Ollama on the host, expose it on `localhost:11434` as normal, and have the desktop app talk to that local endpoint rather than containerizing the model itself.
5. **Data layer** — the SQLite tracker from Part 2, stored locally on the host (not in a container), since it needs to persist across every exercise regardless of which container ran.

### "Any machine" caveats worth knowing upfront
- Cross-platform means Windows/Mac/Linux, but Docker itself still needs to be installed on the host (Docker Desktop, or Colima/Podman as lighter alternatives) — the app doesn't remove that dependency, it just standardizes what runs inside.
- Multi-architecture matters if you ever run this on both Apple Silicon and an Intel/Linux machine — container images need multi-arch builds (`linux/amd64` + `linux/arm64`) or they won't pull correctly on both.
- The local AI model piece is inherently host-specific (GPU/acceleration differs by machine), so "any machine" for that part really means "any machine, but the AI feels different/slower on weaker hardware" — which arguably isn't a bad thing to practice against anyway, since a real interview's AI tool won't adapt to your hardware either.

---

## Part 5: Container Image Contents + Exercise Lifecycle (planning only)

### What each language image needs to contain

**`practice-go`**
- Go toolchain (matching a recent stable version)
- `delve` for interactive debugging
- Race detector is built into the toolchain (`-race` flag) — no extra install needed
- A small set of common libraries pre-vendored if exercises will use them (e.g. a basic HTTP router, a test framework) so container startup doesn't require a network fetch mid-exercise
- `golangci-lint` or similar, optional — useful if you want static-analysis-catchable bugs as a distinct sub-category later

**`practice-cpp`**
- A compiler (clang preferred — better sanitizer support/diagnostics than gcc for this use case) + `cmake`/`make`
- `gdb` or `lldb`
- Sanitizers enabled: ThreadSanitizer (races/deadlocks), AddressSanitizer (memory), UndefinedBehaviorSanitizer — these need to be compile-time flags baked into the exercise build scripts, not just present on the image
- A minimal test framework (e.g. Catch2 or GoogleTest) if exercises include tests to run

**`practice-python`**
- Recent CPython + `pdb`/`debugpy`
- `faulthandler` enabled by default (catches hangs/deadlocks in async or threaded code)
- `pytest` for test-driven debug exercises
- `threading`/`asyncio` are stdlib, no extra install; if exercises touch multiprocessing, that needs to be explicit since it behaves differently across OSes

**Shared across all three**
- A non-root user inside the container (matches real-world dev environment hygiene, avoids permission headaches with mounted volumes)
- A consistent working directory convention (e.g. `/workspace`) so the desktop app's terminal/editor panes can assume a fixed path regardless of language

### Exercise lifecycle (per rep)

1. **Pick** — you select an exercise from the browser (filtered by category: pattern / concurrency / implementation / AI-assisted / debug / system design / mock).
2. **Stage** — the orchestration layer pulls the relevant host repo + branch (e.g. `bug/goroutine-leak-01`) and mounts it into a fresh container from the matching language image. Container starts clean every time — no state carried over from prior reps.
3. **Timer starts** — the moment the terminal pane becomes interactive, not the moment you clicked "start" (so container boot time doesn't eat into your practice time).
4. **Work** — you use the embedded terminal + editor as normal; for AI-assisted exercises, the app also surfaces a chat pane wired to the local Ollama endpoint.
5. **Stop** — you mark result (solved-clean / solved-with-hints / failed / stuck) and it logs `time_to_locate`/`time_to_fix` or `time_spent` depending on exercise type, plus any notes, into the SQLite tracker.
6. **Teardown** — container is destroyed (not just stopped) so the next rep on the same exercise type starts from a truly clean image. The mounted repo branch on the host resets to its clean state via `git checkout`/`git reset` so bug-injection branches stay reusable.
7. **Review trigger** — if result was "failed" or "solved-with-hints," the exercise gets queued into the spaced-repetition schedule from Part 2 automatically.

### One design question worth deciding early
Whether the editor lives *inside* the container (e.g. a terminal-based editor like `vim`/`nvim` you use over the embedded terminal) or *outside* it (a GUI editor pane in the desktop app that edits files on a mounted volume, while only the run/debug/test commands execute inside the container). The second is more forgiving for fast editing but means the app needs proper file-sync handling; the first is simpler to build but forces you to be comfortable in a terminal editor under time pressure — which, notably, is itself close to how some real onsites are proctored (shared terminal, no fancy IDE).

---

## Part 6: TUI-on-Neovim Architecture + AI Tutor Agent (planning only)

Deciding on a TUI built on Neovim settles Part 5's open question in favor of the terminal-native path — and it's a better fit for your background than a GUI wrapper would've been. This also simplifies Part 4 considerably: no Electron/Tauri layer needed, since a terminal app is cross-platform by nature. The desktop shell essentially disappears — the TUI *is* the app.

### Two ways to build "TUI on top of nvim" — worth deciding between

**Option A: Embed Neovim via its RPC API**
- Neovim can run headless (`nvim --embed`) and exposes a full msgpack-RPC API. A custom TUI app (written in Go, using something like `bubbletea`/`tview`) drives an embedded Neovim instance as its editor pane, alongside separate panes for terminal output and the AI tutor chat.
- Pro: one cohesive binary, full control over layout, panes can talk to each other (e.g. tutor pane can see what file/line the editor pane is on).
- Con: more to build — you're writing a real TUI application, not just composing existing tools.

**Option B: Orchestrate existing terminal tools (tmux + nvim + a chat CLI)**
- A layout is scripted (via `tmux`) with panes: nvim editing files on a mounted volume, a shell pane `docker exec`'d into the practice container for build/run/debug, and a small chat CLI pane talking to the local Ollama endpoint.
- Pro: dramatically less to build — nvim stays exactly as you already use it (your plugins, keybindings, LSP config untouched), tmux handles the pane layout, and the "app" is really just a launch script plus the chat CLI.
- Con: less cohesive — panes don't share context automatically (the tutor doesn't inherently know what file you're on unless you tell it, or unless the script wires that up separately, e.g. by piping `nvim` command history or cursor position into the tutor's context).

Given how much of this project is really about *practice discipline* rather than tool-building for its own sake, Option B is the pragmatic starting point — it gets you practicing sooner, using an editor you already know. Option A becomes worth it later only if the seams in B (context-sharing, exercise lifecycle automation) start actually costing you practice time.

### AI tutor agent — behavior, not just plumbing
The interesting design question isn't "how does the chat pane talk to Ollama" (that's a straightforward local API call) — it's what the tutor is *allowed to do*, which should change by exercise category and phase:

- **Learning phase (Phase 1-2 patterns, first pass on a new topic)**: tutor gives Socratic hints, not answers — nudges toward the right pattern/approach without naming it outright. Closer to how a good mentor behaves than how a code-completion tool behaves.
- **AI-assisted category (Part 2, entity 8)**: tutor behaves like the sandboxed AI an interviewer might actually provide — answers directly when asked, but you're expected to verify/correct it, and your corrections are what get logged.
- **Hard-mode/mock reps (Phase 3-4)**: tutor pane is disabled entirely, simulating a no-AI-allowed round. This needs to be an explicit toggle in the launch script/exercise config, not just "don't open the pane" — otherwise it's too easy to peek.
- **Debug/concurrency categories**: worth considering whether the tutor should be allowed to run tools on your behalf (e.g. suggest `-race` or point you at `tsan` output) versus only reason in text — leaning toward text-only, since interpreting sanitizer/debugger output under pressure is itself part of what you're training.

### Stuck-help calibration — syntax vs. logic are different problems
Worth separating these explicitly, since they're not the same skill and shouldn't be gated the same way:
- **Syntax help is effectively always-on**, independent of the exercise's overall tutor mode (even in hard-mode/mock reps, arguably — real interviews don't dock you for typos, and getting stuck on a missing semicolon instead of the actual bug wastes practice time on the wrong thing). This is a low-stakes capability, not a skill being tested.
- **Logic-level "I'm stuck on the actual bug/approach" help stays gated and hint-first** — the tutor nudges (e.g. "have you checked what happens when two goroutines hit this at the same time?") rather than naming the bug outright, escalating to a direct answer only if you ask again or explicitly request it. This preserves the actual debugging rep instead of routing around it.
- Practically: the debugging environment (Part 3's debuggers/sanitizers) is what should be doing the heavy lifting for *finding* things — a debugger stepping through state, or a race detector pointing at a line — with the tutor's job being to help you *read and reason about* that output, not to replace the tool by just telling you the answer. That keeps "learn to debug" as the actual objective rather than "learn to ask AI what's wrong."

### Where this leaves the data layer
The SQLite tracker from Part 2 stays as-is regardless of which option (A/B) you pick — it's independent of the editor/TUI layer and just needs the exercise lifecycle script (Part 5) to write to it at start/stop of each rep.

---

## Part 7: Automated Verification Layer + Tightened AI Usage Rules (planning only)

This closes the loop that self-assessment alone can't: instead of just marking "solved-clean" by feel, the exercise gets an objective pass/fail from an actual test suite — which also matters because it mirrors how real interviews with AI allowed actually constrain the AI (narrow, directive use — not "build this for me").

### AI usage constraint — redefine the tutor's default mode
Up to now Part 6 treated the tutor as somewhat flexible by phase. Tighten the default (non-learning-phase) behavior to match what real interviews actually permit:
- **Syntax correction only** — fixes typos, wrong API signatures, import errors. Never changes logic.
- **Direct implementation on explicit instruction** — if you say "implement a function that takes X and returns Y using approach Z," it does exactly that, nothing more. No unsolicited refactors, no "I also improved..." additions, no alternate suggestions unless asked.
- **No unsolicited scope creep** — this is the important behavioral constraint to actually enforce (via the tutor's system prompt/instructions), since general-purpose coding assistants default to being helpful in ways that go beyond what was asked. The tutor needs to be told explicitly: do the narrow thing requested, stop, wait.
- **Stuck-help remains available separately** — if you're stuck (not just directing implementation), you can explicitly ask the tutor for a hint or explanation, which is a different mode than the direct-implementation one above. Worth keeping these as distinct, logged separately, so you can see afterward how often you leaned on "explain" vs. "implement this for me."

### Automated verification — new entity: Test suite result
Every exercise category that produces working code (implementation, debug, concurrency, AI-assisted) gets an attached, hidden test suite that runs automatically at "stop," rather than relying only on self-rated result.
- `id`, `exercise_id`, `test_type` (unit / integration / API / SQL-query-diff), `tests_passed`, `tests_total`, `pass_rate`, `failing_test_names[]`, `run_duration_ms`
- This becomes the source of truth for `result` (solved-clean / partial / failed) instead of self-assessment — self-assessment stays as a secondary field for reflection ("did it feel harder than the test results suggest"), not the primary signal.

### What the test suites look like per exercise type
- **API endpoint exercises**: exercise defines expected request/response contracts; test suite fires real requests (e.g. via a lightweight HTTP client) at the running service inside the container and diffs actual vs. expected responses — status codes, payload shape, edge cases (bad input, auth failure, empty results).
- **SQL query exercises**: a fixture database gets loaded into the container; the exercise's query output is diffed against an expected result set (row-for-row or via checksums for larger sets).
- **General implementation/debug exercises**: standard unit tests in the language's native framework (Go's `testing`, C++'s Catch2/GoogleTest, Python's `pytest`), pre-written per exercise, hidden from you until after you submit.
- **Concurrency exercises**: correctness test plus a stress-test loop (run N times, check for flakiness) rather than a single pass/fail — a concurrency fix that only sometimes works needs to be visible as such, not scored as clean.

### Logging AI interactions distinctly
Add to the AI-assisted task entity (Part 2, #8):
- `interaction_mode` per request: `syntax-correction` / `direct-implementation` / `stuck-help`
- `instruction_given` (short text of what you asked for)
- `accepted_as_is` (yes/no/partial) — did you take the AI's output unmodified, or correct it further
This gives you, over time, a real signal on *how* you're using AI — e.g. if "stuck-help" requests trend up in a category, that's a clearer signal of a real knowledge gap than raw pass/fail alone.

### Updated lifecycle (extends Part 5, step 5 "Stop")
5a. **Stop** → automated test suite runs immediately, producing `tests_passed`/`tests_total`.
5b. **Self-assessment** → you still log a subjective result and notes, now explicitly as a secondary/reflective signal.
5c. **Both get written** to SQLite, with test results as the primary driver of whether the exercise gets queued into spaced repetition (Part 2) — a "felt fine but failed 2 hidden tests" outcome should absolutely trigger a review, even if your gut said it went well.

---

## Part 8: Exercise Definition Format (planning only)

Every exercise needs a single definition that the orchestration layer (Part 4/5) can load to know: what category it belongs to, which container image to use, what the tutor is allowed to do, and what hidden tests decide pass/fail. Sketching the fields conceptually (not a real file yet):

### Top-level exercise metadata
- `id`, `title`, `category` (pattern / concurrency / implementation / ai-assisted / debug / system-design)
- `language` (Go / C++ / Python / n/a for system design)
- `difficulty`, `phase` (which Part 1 phase this belongs to — affects tutor defaults)
- `time_limit_min`
- `container_image` (which of the Part 6 images to launch)
- `repo_source` (which host repo + branch to mount, per Part 3's staging structure)

### Tutor permission block (ties Part 6 + Part 7 together)
- `tutor_enabled` (bool) — hard off for hard-mode/mock reps, per Part 6
- `syntax_help_always_on` (bool, default `true`) — kept independent of `tutor_enabled`, since syntax correction isn't the skill being tested (per Part 6's stuck-help calibration); even a "tutor off" hard-mode rep can leave this on unless you specifically want a fully unassisted rep
- `tutor_default_mode` — one of `syntax-correction-only`, `direct-implementation-allowed`, `full-ai-assisted` (matches the AI-assisted category's looser rules), `learning-hints-only` (Socratic mode for early-phase learning exercises)
- `logic_stuck_help` — `hints-first` (default) or `direct-answer` — governs the *separate*, gated behavior for "I'm stuck on the actual bug/approach," as opposed to syntax help above
- `tutor_scope_notes` — free text injected into the tutor's own system prompt for this specific exercise, e.g. "do not suggest the overall algorithm, only correct syntax" — this is what actually enforces the "very direct, no scope creep" behavior from Part 7 at the per-exercise level, rather than relying on a single global rule that might not fit every exercise

### Verification block
- `test_type` (unit / API / SQL-query-diff / stress-test-loop for concurrency)
- `test_suite_path` — where the hidden tests live (mounted into the container but not shown in the editor pane until after submission)
- `pass_threshold` — for partial-credit cases (e.g. "80%+ of tests passing counts as solved-with-hints, not full pass")

### Debug/concurrency-specific fields (when applicable)
- `injected_bug_category` (from Part 2's per-language bug category lists)
- `repro_reliability` — whether the bug reproduces every run or is flaky by design (relevant for concurrency exercises specifically, so a "couldn't repro" result can be checked against whether it's supposed to be flaky)

### Why this matters as its own layer
Without a structured definition per exercise, the tutor's behavior and the test strictness end up being decided ad hoc in the moment, which defeats the purpose of controlled, comparable reps. Having `tutor_scope_notes` per exercise (rather than one global instruction) also means you can deliberately vary how directive the AI is across exercises — useful later if you want to practice a spectrum from "AI barely helps" to "AI does most of the typing," since real interview AI permissions aren't perfectly uniform across companies either.

---

## Part 9: Sandbox / Free-Play Environment (planning only)

Separate from the graded exercise flow entirely — this is for open-ended exploration: trying a new library, poking at a concurrency pattern before attempting it against the clock, seeing how a sanitizer reacts to something, or just tinkering. No timer, no hidden tests, no formal grading.

### Why it needs to be structurally separate, not just "an exercise with grading off"
- Exercise containers (Part 5) are ephemeral by design — clean teardown every rep, so results are comparable. A sandbox needs the opposite property: **state that persists** across sessions, since the point is to keep building on what you tried yesterday.
- Mixing the two modes in one container type risks contaminating graded reps with leftover sandbox state, or making the sandbox annoyingly ephemeral when persistence is the whole point.

### Structure
- **One persistent container per language** (`sandbox-go`, `sandbox-cpp`, `sandbox-python`), same base images as the exercise containers (Part 6) so tooling stays consistent, but backed by a **persistent Docker volume** rather than a git-branch mount that resets.
- **No timer, no hidden test suite** — you're not being scored here.
- **Tutor defaults more open** — `full-ai-assisted` or `learning-hints-only` rather than the tight `syntax-correction-only` default, since exploration benefits from a more collaborative AI rather than a deliberately constrained one. This is the natural place to actually use your local model (or even your regular Claude access) more liberally.
- **Explicit reset, not automatic teardown** — sandboxes accumulate state on purpose, but they'll eventually get messy. A manual `reset sandbox-go` (or similar) command wipes the volume back to a clean base image on demand — you decide when, not the lifecycle script.
- **Not logged into the spaced-repetition system** — since there's no result/test outcome to grade, sandbox sessions don't feed Part 2's tracker the same way. Optional: a lightweight, informal "sandbox notes" log (just freeform text + timestamp) if you want a record of what you explored, without it affecting any metrics.

### How it fits with graded practice
A natural workflow: use the sandbox to first get comfortable with something (e.g. "how does Go's race detector actually flag this pattern") *before* attempting the graded, timed version of a similar exercise. This keeps the graded reps honest as a measure of performance under pressure, while still giving you room to learn without the clock running.

---

## Part 10: Unified App Container with In-Container Pane Splitting (planning only)

This is a real simplification of Parts 4/6, worth adopting: instead of a host-side TUI binary orchestrating separate per-language exercise containers, **the entire environment — editor, terminal, AI chat — runs inside one container**, split into panes/windows via `tmux` from the moment it starts.

### What changes
- **One container, launched with a single command** (`docker run -it <image>`), whose entrypoint immediately starts a `tmux` session pre-arranged with panes: nvim (editor), a shell (terminal/build/debug), and the chat CLI (AI tutor).
- **The "desktop shell" question from Part 4 mostly disappears.** There's no separate host-side app to build — the container *is* the app. This removes Option A vs. B from Part 6 as an open question; it resolves to a version of Option B, just fully inside the container instead of split across host + container.
- **"Works on any machine" reduces to "has Docker installed."** No host TUI binary, no packaging/distribution problem beyond publishing the image. This is the cleanest version of the original goal.

### One image vs. three — a tradeoff worth deciding
Since everything now lives in one container anyway, it's worth reconsidering whether you need three separate language images (Go/C++/Python from Part 6) or one unified image with all three toolchains baked in:
- **One unified image**: simpler startup (one `docker run`, one tmux layout regardless of exercise language), larger image size, less isolation between language environments (not a real concern here since you're not running untrusted code).
- **Three images still**: keeps toolchains cleanly separated, but now needs the entrypoint/tmux logic duplicated (or shared via a common base image) across three images, and exercises need to pick the right one at launch.
- Given the goal is fast iteration and low friction over strict isolation, the unified image is the more practical default — you're the only user, and there's no security boundary being enforced between languages.

### Networking constraint carried over from Part 4
The AI tutor's chat CLI, running inside this container, still needs to reach the local model — which per Part 4 lives on the **host**, not in a container, to keep Metal GPU acceleration on your M4 Max. That means the container needs an explicit route out to the host:
- On Docker Desktop for Mac, `host.docker.internal` resolves to the host automatically — the chat CLI just points there.
- This is a detail worth testing early, since "the tutor pane can't reach the model" is exactly the kind of thing that's silent and confusing to debug later if the networking assumption is wrong.

### How this affects the exercise lifecycle (Part 5) and sandbox (Part 9)
- **Exercise flow**: unchanged in spirit — container still gets destroyed and recreated clean per rep (Part 5, step "Teardown") — it's just one container now instead of potentially separate language-specific ones.
- **Sandbox**: similarly simplifies to one persistent container (or one persistent volume mounted into the unified image) rather than three separate persistent sandbox containers, since the same image handles all three languages.
- **Net effect**: Parts 5 and 9's *logic* (clean teardown vs. persistent volume) stays exactly the same — only the container topology underneath gets simpler.
