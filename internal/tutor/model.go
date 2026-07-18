package tutor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	agentopt "github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	template "github.com/cloudwego/eino/utils/callbacks"
)

// tutorStateFile is the well-known dotfile the tutor pane writes into
// the workspace so internal/session's submit/reference commands can
// attach this session's tutor-assistance footprint to a tracker
// attempt without talking to this package directly -- same file-based
// handoff lastTestResultFile (filecontext.go) uses in the other
// direction (session writing, tutor reading). Duplicated as a literal
// in internal/session rather than imported, matching this codebase's
// established local-duplication convention (see filecontext.go's own
// doc comment on lastTestResultFile).
const tutorStateFile = ".ballroom-tutor-state.json"

// tutorState is the JSON shape written to tutorStateFile.
type tutorState struct {
	HintsUsed int    `json:"hints_used"`
	TutorMode string `json:"tutor_mode"`
	Model     string `json:"model"`
}

// writeTutorState persists the session's current tutor-assistance
// counters to workDir -- best-effort and silent: a write failure is
// not something a mid-conversation turn should ever fail or even
// visibly warn about (the user is mid-flow), matching this file's
// other silent-degradation paths (e.g. remoteExpr's own no-editor-
// attached case). A missing/malformed read on the session side already
// degrades to zero/empty, so a failed write here just means the
// eventual attempt records nothing extra, never a broken submit.
//
// An empty workDir (most of this package's own tests, which build a
// Config via testConfig without setting WorkDir) is a deliberate no-op
// rather than writing into the process's current working directory --
// see TestWriteTutorState_EmptyWorkDirIsANoOp.
func writeTutorState(workDir string, s tutorState) {
	if workDir == "" {
		return
	}
	data, err := json.Marshal(s)
	if err != nil {
		return
	}
	_ = os.WriteFile(filepath.Join(workDir, tutorStateFile), data, 0o644)
}

// minTextareaRows/minViewportRows are floors, not targets — the textarea
// always shows at least one row even when empty, and the viewport always
// keeps at least a few rows of conversation visible even when the
// textarea has grown to its cap (see recomputeLayout).
const (
	minTextareaRows = 1
	minViewportRows = 3
)

// tutorModel is the tutor pane's bubbletea Model — replaces
// internal/tutor/scrollbox.go's hand-rolled ANSI positioning entirely.
// Mirrors internal/tui/app.go's own Model/Update/View architecture,
// which has needed zero manual cursor/escape-sequence code all session,
// unlike the hand-rolled approach this replaces.
type tutorModel struct {
	viewport viewport.Model
	textarea textarea.Model
	width    int
	height   int

	// displayBlocks is what the viewport shows: each submitted message
	// echoed immediately (before its reply arrives), each reply, and
	// pre-styled notes (tool-usage summaries, error notes) -- including
	// a failed turn's honest fallback message, which is shown but
	// deliberately never added to history (see submit/Update's
	// turnCompleteMsg case). Decoupled from history on purpose: history
	// is model context only, this is display only. Blocks hold raw
	// content and are styled per frame at the viewport's current width
	// (renderBlock, refreshViewport) -- styling at append time baked a
	// width into the transcript, which is exactly what fixed-width
	// constructs like the editor cards can't survive a resize with.
	displayBlocks []displayBlock

	// ctx is stored on the model (unusual for Go, normally a function
	// parameter) because bubbletea's Update(msg) signature has nowhere
	// else to thread it through for the lifetime of the program -- same
	// ctx Run(ctx, cfg, ...) already received from its own caller.
	ctx context.Context
	cfg Config

	workerAgent          *react.Agent
	orchestratorAgent    *react.Agent
	workerCM             model.ToolCallingChatModel
	orchestratorCM       model.ToolCallingChatModel
	workerEndpoint       string
	orchestratorEndpoint string
	routingEnabled       bool

	// workerStrategy/orchestratorStrategy are each role's detected
	// toolCallingStrategy -- orchestratorStrategy is only meaningful
	// when routingEnabled. startTurn's callRole dispatch reads these to
	// decide whether a role's Generate call goes through the real
	// react.Agent or runFallbackToolLoop (fallback.go).
	//
	// Detection is asynchronous: Init()'s detectStrategies command makes
	// the (up to two) live probe calls off the UI thread and delivers a
	// strategiesDetectedMsg. It used to run blocking inside newTutorModel
	// -- a real bug found live: with slow free-tier models the whole
	// pane, chat box included, stayed blank until both probes returned,
	// because Run()'s tea.Program hadn't even started. strategiesDetected
	// flips true when the message lands; a submit arriving before then is
	// held in queuedSubmits and started by the strategiesDetectedMsg case.
	workerStrategy       toolCallingStrategy
	orchestratorStrategy toolCallingStrategy
	strategiesDetected   bool
	queuedSubmits        []queuedSubmit

	// tools/toolCatalogText are built once per session (buildTools +
	// renderToolCatalog) and reused by both roles -- the tool set is
	// identical regardless of which role answers a given turn.
	// toolCatalogText only matters for a jsonFallbackToolCalling role
	// (see prependToolsPrompt, prompts.go); a native role never sees it.
	tools           []tool.BaseTool
	toolCatalogText string

	// history holds only the system prompt plus clean (user, assistant)
	// pairs -- never a failed turn's fallback message, never
	// tool-call scaffolding -- exactly matching the pre-rewrite Run()
	// loop's own history semantics.
	history                   []*schema.Message
	comprehensionCheckPending bool
	helpRequestCount          int
	turnInFlight              bool

	// turnCancel cancels the in-flight turn's derived context (see
	// startTurn's ctx, cancel := context.WithTimeout(m.ctx, turnTimeout))
	// -- nil whenever no turn is in flight, and also nil while a turn is
	// still queued behind strategy detection (queuedSubmit) and hasn't
	// actually started a goroutine yet. Set by submit/the
	// strategiesDetectedMsg case; cleared once the turn is no longer
	// this model's concern, either by cancelInFlightTurn or by
	// turnCompleteMsg's own handling of a natural completion.
	turnCancel context.CancelFunc

	// pendingUserText is the in-flight turn's submitted line, held here
	// purely so a cancel or failure can hand it back to the user instead
	// of discarding it (see recoverPendingText) -- set the moment a turn
	// starts (submit), cleared the moment it's recovered one way or the
	// other. A successful completion also clears it without recovering
	// anything: success has nothing left to hand back.
	pendingUserText string

	// turnSeq tags every message a turn produces (activityEventMsg,
	// streamTextMsg, turnCompleteMsg all carry the seq they were started
	// under -- see startTurn/waitForActivityEvent) with a monotonic
	// generation number. Bumped both when a new turn starts AND when the
	// current one is cancelled (cancelInFlightTurn) -- cancelling must
	// retire the generation immediately even though the cancelled
	// goroutine is still unwinding in the background, since ctx
	// cancellation aborts the in-flight request but not instantaneously,
	// and by the time that happens bubbletea's runtime may already have
	// moved on to a different turn, or none at all.
	//
	// Update's three turn-message cases drop anything whose seq doesn't
	// match m.turnSeq before touching any state. That's what makes a
	// cancelled turn's eventual straggling result harmless once it does
	// land: bubbletea's runtime started that wait the moment submit
	// returned it, so a cancel genuinely cannot prevent it from
	// eventually being delivered to Update -- only rendering it inert
	// once it gets there.
	turnSeq int

	// activeCalls/pulsePhase drive the live tool-call activity region
	// (see activityView, buildActivityChannelOption) -- activeCalls is
	// only ever non-nil while turnInFlight; pulsePhase free-runs for the
	// program's whole lifetime (see pulseTickMsg) rather than being
	// started/stopped per turn, which is only visually relevant while
	// turnInFlight (plus the aurora fade window just after -- see
	// auroraFade) anyway and avoids any start/stop bookkeeping.
	activeCalls []activityCall
	pulsePhase  int

	// turnStartedAt/turnSettledAt anchor the thinking aurora's two
	// ramps (aurora.go): the bloom-in measures from turnStartedAt, the
	// fade-out from turnSettledAt. Both may be backdated at the
	// transition (see submit and the turnCompleteMsg case) so the glow
	// always continues from its current level rather than jumping.
	// Zero turnSettledAt means no turn has ever completed, so the
	// aurora has never had a reason to exist.
	turnStartedAt time.Time
	turnSettledAt time.Time

	// streamingText is the in-flight turn's partial reply so far (full
	// accumulated text, not a delta -- see pushLatestStreamText), shown
	// as a provisional styled block below displayLines while the turn
	// runs (refreshViewport) and cleared when it completes; the final
	// reply is appended through the exact same turnCompleteMsg path as a
	// non-streamed turn, so the settled display is byte-identical.
	// streamPainting flips true at the first painted chunk and makes the
	// aurora begin its fade-out mid-turn (auroraFade) -- the reply is
	// visibly arriving, so the "thinking" glow should already be dying.
	streamingText  string
	streamPainting bool

	// transcriptWarned makes transcript-export failures warn exactly
	// once (see the turnCompleteMsg case) instead of once per turn.
	transcriptWarned bool
}

// displayBlock is one transcript entry, held raw and styled per frame
// (renderBlock) at the viewport's current width -- see the
// displayBlocks field's doc comment for why styling can't happen at
// append time anymore.
type displayBlock struct {
	kind displayBlockKind
	raw  string
}

type displayBlockKind int

const (
	// blockNote is pre-styled content rendered exactly as appended --
	// tool-usage summaries, error notes, fallback messages. These were
	// styled for a specific width at append time before this refactor
	// too; nothing regressed, they just don't re-flow.
	blockNote displayBlockKind = iota
	// blockUser is one line the user submitted, echoed as an
	// accent-bar block (renderUserBlock).
	blockUser
	// blockTutor is one tutor reply's raw markdown, styled by
	// styleMarkdown at render time.
	blockTutor
)

// renderBlock styles one transcript block for the given content width.
func renderBlock(b displayBlock, width int) string {
	switch b.kind {
	case blockUser:
		return renderUserBlock(b.raw, width)
	case blockTutor:
		return styleMarkdown(b.raw, width)
	default:
		return b.raw
	}
}

// renderUserBlock renders one submitted user message as a block with a
// pink accent bar down its left edge — the bar alone is the speaker
// signal (no "you" label), in the user's own accent color, distinct
// from the tutor's teal. The text is wrapped here at width−2 (bar +
// space) so every emitted line is already ≤ width — the invariant that
// keeps refreshViewport's outer word-wrap from ever re-breaking a line
// and orphaning a continuation row without its bar. Built from raw
// escapes rather than a lipgloss border style for the same reason
// statusbar.go uses them: lipgloss colors can vanish under go test's
// TTY-less termenv profile, and this block is pinned by
// string-assertion tests.
func renderUserBlock(raw string, width int) string {
	bar := ansiFg(panePink) + "│" + mdColorReset + " "
	text := raw
	if w := width - 2; w > 0 {
		text = lipgloss.NewStyle().Width(w).Render(raw)
	}
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = bar + line
	}
	return strings.Join(lines, "\n")
}

// newTutorLayoutOnly builds a model with just the textarea/viewport
// wiring — no agents, no config, no network. Used by this file's own
// pure-layout tests (resize, dynamic growth, Enter-submits-not-newline)
// that have no need to exercise real turn logic. newTutorModel (below)
// is what Run() and every turn-logic test actually uses.
func newTutorLayoutOnly() tutorModel {
	ta := textarea.New()
	ta.Placeholder = "Ask a question..."
	ta.Prompt = "› "
	ta.FocusedStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(paneTeal))
	ta.BlurredStyle.Prompt = lipgloss.NewStyle().Foreground(lipgloss.Color(paneInputRule))
	ta.ShowLineNumbers = false
	ta.KeyMap.InsertNewline.SetEnabled(false)
	ta.Focus()

	vp := viewport.New(0, minViewportRows)
	vp.Style = viewportContentStyle

	return tutorModel{viewport: vp, textarea: ta}
}

// newTutorModel builds the real tutor model: the layout from
// newTutorLayoutOnly, plus every piece of setup Run()'s old for-loop
// used to do once at the top — building tools, the worker chat
// model/agent, and (when cfg.OrchestratorModel is set) the orchestrator
// chat model/agent, seeding history with the mode's system prompt, and
// deciding whether the first message wants a comprehension check. Same
// construction logic as before, just returned as a Model instead of
// consumed inline by a for-loop.
func newTutorModel(ctx context.Context, cfg Config) (tutorModel, error) {
	m := newTutorLayoutOnly()
	m.ctx = ctx
	m.cfg = cfg

	tools, err := buildTools(cfg)
	if err != nil {
		return tutorModel{}, err
	}
	m.tools = tools
	m.toolCatalogText, err = renderToolCatalog(ctx, tools)
	if err != nil {
		return tutorModel{}, err
	}

	m.workerCM, err = newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
	if err != nil {
		return tutorModel{}, err
	}
	m.workerAgent, err = newAgent(ctx, m.workerCM, tools)
	if err != nil {
		return tutorModel{}, err
	}
	m.workerEndpoint = providerEndpoint(cfg.Model, cfg.OllamaHost)

	m.routingEnabled = cfg.OrchestratorModel != ""
	if m.routingEnabled {
		m.orchestratorEndpoint = providerEndpoint(cfg.OrchestratorModel, cfg.OllamaHost)
		m.orchestratorCM, err = newChatModel(ctx, cfg.OrchestratorModel, cfg.OllamaHost, cfg.APIKey)
		if err != nil {
			return tutorModel{}, err
		}
		m.orchestratorAgent, err = newAgent(ctx, m.orchestratorCM, tools)
		if err != nil {
			return tutorModel{}, err
		}
	}
	// No banner block in the transcript -- the session's identity
	// (model, mode, endpoint, exit hint) lives in the fixed header line
	// above the viewport now (headerView), where it stays visible
	// instead of scrolling away with the conversation.

	// Strategy detection deliberately does NOT happen here -- see the
	// strategiesDetected field's doc comment. Construction must make
	// zero probe network calls so Run()'s tea.Program starts (and the
	// chat box paints) immediately; Init()'s detectStrategies command
	// delivers the strategies asynchronously.

	// Seeded with persona text only, not systemPromptForMode's combined
	// persona+native-tools string: the tools instruction is strategy-
	// dependent and (once routing is enabled) can differ turn to turn
	// depending on which role answers, so it's prepended fresh per call
	// by startTurn (via prependToolsPrompt) instead of being baked once
	// into the session's fixed history.
	m.history = []*schema.Message{schema.SystemMessage(personaPromptForMode(cfg.Mode))}
	m.comprehensionCheckPending = wantsComprehensionCheck(cfg.Mode)
	m.refreshViewport()

	// Written once here (hints_used at its true starting value, zero)
	// so a submit or `ballroom reference` that happens before the user
	// ever asks the tutor anything still finds real tutor_mode/model
	// values instead of nothing at all -- see writeTutorState's doc
	// comment.
	writeTutorState(cfg.WorkDir, tutorState{TutorMode: cfg.Mode, Model: cfg.Model})

	return m, nil
}

// RunOneTurn builds a tutor session and submits exactly one message,
// synchronously draining the turn to completion -- for headless callers
// (cmd/tutor-eval's grounding checks) that need one turn's real result
// without driving a full interactive tea.Program.
//
// This deliberately does NOT go through Run()/tea.Program: a real
// interactive terminal never types its next input (including Ctrl-D)
// before seeing the current turn's reply, but a scripted byte stream fed
// through tea.WithInput has no such guarantee -- bubbletea delivers
// queued input immediately, so a trailing Ctrl-D appended right after a
// message would very likely reach Update() and quit the program before
// the turn's async goroutine (a real network round trip) ever finishes,
// racing the very thing the caller wants to observe. Calling
// newTutorModel + Update(KeyEnter) directly and draining the returned
// tea.Cmd chain synchronously (same pattern as this package's own
// submitAndRun test helper) sidesteps that race entirely: there is no
// second input source that can outrun the turn.
func RunOneTurn(ctx context.Context, cfg Config, message string) (reply string, err error) {
	m, err := newTutorModel(ctx, cfg)
	if err != nil {
		return "", err
	}
	// Detection is async in a real session (Init's command); headless
	// there is no tea runtime to deliver it, and blocking here is fine --
	// this caller wants one turn's result, not a responsive UI. Without
	// this the submit below would queue forever.
	newModel, _ := m.Update(m.detectStrategies()())
	m = newModel.(tutorModel)
	m.textarea.SetValue(message)
	newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newModel.(tutorModel)
	if cmd == nil {
		return "", fmt.Errorf("tutor: empty message produced no turn")
	}
	for {
		msg := cmd()
		newModel, cmd = m.Update(msg)
		m = newModel.(tutorModel)
		if tc, ok := msg.(turnCompleteMsg); ok {
			if tc.err != nil {
				return "", tc.err
			}
			return tc.reply.Content, nil
		}
		if cmd == nil {
			return "", fmt.Errorf("tutor: turn ended without a result")
		}
	}
}

func (m tutorModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, pulseTickCmd(), m.detectStrategies())
}

// strategiesDetectedMsg delivers both roles' async strategy-detection
// results (see the strategiesDetected field's doc comment).
// orchestrator is only meaningful for a routing-enabled session.
type strategiesDetectedMsg struct {
	worker       toolCallingStrategy
	orchestrator toolCallingStrategy
}

// detectStrategies probes each configured role's tool-calling strategy
// concurrently -- worker always, orchestrator only when routing is
// enabled -- off the UI thread, as a tea.Cmd from Init(). detectStrategy
// never fails (any probe error resolves to jsonFallbackToolCalling, the
// strategy that can actually work when the check itself can't complete
// -- see its doc comment), so this always delivers a usable result.
func (m tutorModel) detectStrategies() tea.Cmd {
	ctx, cfg, routingEnabled := m.ctx, m.cfg, m.routingEnabled
	return func() tea.Msg {
		var msg strategiesDetectedMsg
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			msg.worker = detectStrategy(ctx, cfg.OllamaHost, cfg.Model, cfg.APIKey)
		}()
		if routingEnabled {
			wg.Add(1)
			go func() {
				defer wg.Done()
				msg.orchestrator = detectStrategy(ctx, cfg.OllamaHost, cfg.OrchestratorModel, cfg.APIKey)
			}()
		}
		wg.Wait()
		return msg
	}
}

func (m tutorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.recomputeLayout()
		// Re-wrap the already-displayed conversation to the new width --
		// refreshViewport bakes wrapping into the content at SetContent
		// time (see its doc comment), so a resize after a long message
		// is already showing needs this to re-flow it, not just resize
		// the viewport's own frame.
		m.refreshViewport()
		return m, tea.ClearScreen

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlD {
			// bubbletea does not quit on its own when the underlying
			// io.Reader hits EOF (confirmed directly from tty.go's
			// readLoop: io.EOF is explicitly excluded from the errors it
			// forwards as a shutdown) -- unlike the old bufio.Scanner
			// loop, which exited naturally when Scan() returned false.
			// This is the explicit replacement, matching the banner's
			// own "Ctrl-D to exit." text.
			return m, tea.Quit
		}
		if msg.Type == tea.KeyCtrlC {
			// Ctrl-C is otherwise unhandled by bubbletea entirely: the
			// program runs the terminal in raw mode, so there is no
			// SIGINT to catch -- this KeyMsg is the only signal a Ctrl-C
			// press ever produces. While a turn is running it cancels
			// that turn (same as Esc, just below); idle, it quits --
			// matching every other program in this codebase (see
			// internal/tui/app.go's own "ctrl+c" -> tea.Quit, and this
			// file's own Ctrl-D case above).
			if m.turnInFlight {
				return m.cancelInFlightTurn()
			}
			return m, tea.Quit
		}
		if msg.Type == tea.KeyEsc && m.turnInFlight {
			// Esc is only claimed here, while a turn is in flight --
			// idle, it falls through to the textarea below unchanged,
			// same as before this existed.
			return m.cancelInFlightTurn()
		}
		if msg.Type == tea.KeyEnter {
			return m.submit()
		}
		if msg.Type == tea.KeyPgUp || msg.Type == tea.KeyPgDown {
			// Dedicated to scrolling the conversation history -- a real
			// bug found live: the textarea swallowed every key
			// unconditionally, so there was no way to scroll up and see
			// earlier messages at all once they'd scrolled past the top
			// of the viewport (refreshViewport's GotoBottom always pins
			// to the latest content). PgUp/PgDn specifically (not arrow
			// keys, and not bubbles/viewport's own default single-letter
			// vim bindings like "j"/"k") because they're never used for
			// normal text editing, so routing them to the viewport
			// instead of the textarea can never eat real typing.
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		m.recomputeLayout()
		return m, cmd

	case tea.MouseMsg:
		// Mouse wheel scrolling -- never conflicts with typing, so every
		// mouse event goes straight to the viewport. Requires
		// tea.WithMouseCellMotion() in Run() for the terminal to
		// actually report these.
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd

	case activityEventMsg:
		if msg.seq != m.turnSeq {
			// Stale: this turn was cancelled (or superseded) before it
			// finished -- see turnSeq's doc comment for why bubbletea's
			// runtime can still deliver this regardless, and why dropping
			// it here (rather than re-arming a wait nobody asked for
			// anymore) is the correct response.
			return m, nil
		}
		m.activeCalls = msg.calls
		m.recomputeLayout()
		return m, waitForActivityEvent(msg.seq, msg.activityCh, msg.streamCh, msg.doneCh)

	case streamTextMsg:
		if msg.seq != m.turnSeq {
			return m, nil
		}
		if !m.streamPainting && msg.text != "" {
			// First painted chunk: the reply is visibly arriving, so the
			// thinking glow starts dying now rather than at turn
			// completion -- backdated the same way turnCompleteMsg does,
			// so the fade continues from the glow's current level.
			glowLevel := m.auroraFade()
			m.streamPainting = true
			m.turnSettledAt = time.Now().Add(-auroraFadeOutLead(glowLevel))
		}
		m.streamingText = msg.text
		m.refreshViewport()
		return m, waitForActivityEvent(msg.seq, msg.activityCh, msg.streamCh, msg.doneCh)

	case pulseTickMsg:
		// Free-running for the program's whole lifetime (see the doc
		// comment on tutorModel.pulsePhase) -- always re-arms, never
		// stops, cheap when idle since it's only rendered while
		// turnInFlight.
		m.pulsePhase++
		return m, pulseTickCmd()

	case strategiesDetectedMsg:
		m.workerStrategy = msg.worker
		m.orchestratorStrategy = msg.orchestrator
		m.strategiesDetected = true
		// Start every submit that arrived while detection was still in
		// flight -- each queued line already echoed and engaged the
		// thinking state in submit(); only the actual turn was held back,
		// because startTurn snapshots m and needs the strategies resolved.
		var cmds []tea.Cmd
		for _, q := range m.queuedSubmits {
			activityCh, streamCh, doneCh, cancel := startTurn(m, q.line, q.checkComprehension, q.seq)
			m.turnCancel = cancel
			cmds = append(cmds, waitForActivityEvent(q.seq, activityCh, streamCh, doneCh))
		}
		m.queuedSubmits = nil
		return m, tea.Batch(cmds...)

	case turnCompleteMsg:
		if msg.seq != m.turnSeq {
			// Stale: see turnSeq's doc comment -- either this turn was
			// cancelled and (maybe) replaced by a newer one already, or
			// it was cancelled and nothing has started since. Either way
			// its result is no longer this model's concern; the cancel
			// already showed its own note and recovered the user's text
			// at cancel time (cancelInFlightTurn).
			return m, nil
		}
		// Read the glow's level before flipping turnInFlight, then
		// backdate the settle time so the fade-out resumes from that
		// level -- a reply landing mid-bloom must not flash the glow
		// up to full brightness before it dies.
		glowLevel := m.auroraFade()
		m.turnInFlight = false
		m.turnCancel = nil
		m.turnSettledAt = time.Now().Add(-auroraFadeOutLead(glowLevel))
		// The provisional streamed block (if any) is done: the final
		// reply below replaces it on success, and a failed turn must not
		// leave a partial reply looking like a real one.
		m.streamingText = ""
		m.streamPainting = false
		calls := m.activeCalls
		m.activeCalls = nil
		m.helpRequestCount = msg.helpRequestCount
		// Kept in sync with helpRequestCount every turn (success or
		// failure -- msg.helpRequestCount already reflects either
		// outcome, see startTurn) so session.Submit/session.Reference
		// always see this session's latest count, not just its value
		// at startup.
		writeTutorState(m.cfg.WorkDir, tutorState{HintsUsed: m.helpRequestCount, TutorMode: m.cfg.Mode, Model: m.cfg.Model})
		m.recomputeLayout()
		// toolUsageSummary leaves a permanent record of which tools this
		// turn used -- the live activity region above is about to
		// disappear entirely now that turnInFlight is false, so without
		// this the conversation history would show no trace a tool was
		// ever called, only the final reply. Applies on both success and
		// failure: a turn can call tools and still fail on the final
		// reply, and that's still worth showing. Empty (a no-op append)
		// for a turn that made no tool calls at all.
		if summary := toolUsageSummary(calls, m.activityContentWidth()); summary != "" {
			m.displayBlocks = append(m.displayBlocks, displayBlock{kind: blockNote, raw: summary})
		}
		// routingWarning and a turn failure's real error detail both used
		// to go to a raw fmt.Fprintf(m.stderr, ...) call -- a real bug
		// found live: a real interactive session has stderr and stdout on
		// the very same tty, so that write bypassed bubbletea's renderer
		// entirely and visibly corrupted the alt-screen frame (stray text,
		// e.g. an eino error's own "node path: [chat]" detail, landing
		// wherever the cursor happened to be, never cleared by the next
		// redraw). Rendered into displayLines instead, both now go through
		// the same safe pipeline as everything else on screen.
		if msg.routingWarning != "" {
			m.displayBlocks = append(m.displayBlocks, displayBlock{kind: blockNote, raw: activityErrorNote(msg.routingWarning)})
		}
		if msg.err != nil {
			// classifyTurnError (history.go) gives a context-window
			// overflow, a per-turn timeout, or a cancellation their own
			// wording instead of the generic connectivity phrasing below
			// -- none of those three actually failed to REACH the model
			// the way "could not reach" implies. Empty means none of
			// those matched, so this keeps today's original wording for
			// a genuine connectivity/provider-rejection failure.
			detail := classifyTurnError(msg.err)
			if detail == "" {
				detail = fmt.Sprintf("could not reach %s: %v", msg.endpoint, msg.err)
			}
			m.displayBlocks = append(m.displayBlocks, displayBlock{kind: blockNote, raw: turnFailedFallbackReply})
			m.displayBlocks = append(m.displayBlocks, displayBlock{kind: blockNote, raw: activityErrorNote(m.recoverPendingText(detail))})
			m.refreshViewport()
			return m, nil
		}
		m.pendingUserText = ""
		// trimHistory (history.go) keeps history from growing without
		// bound across a long session -- resent in full every turn, it
		// would both slow every request down and, eventually, exceed the
		// provider's own hard context-length limit outright. Applied
		// here, right as a new pair is added, rather than only at
		// request-build time, so m.history itself never carries more
		// than the budget either -- there's no other reason to keep the
		// untrimmed tail around once a pair has aged out of it.
		m.history = trimHistory(append(m.history, schema.UserMessage(msg.userMessage), schema.AssistantMessage(msg.reply.Content, nil)), historyBudgetChars)
		// The block keeps the raw reply -- styleMarkdown runs per frame
		// in renderBlock, and history above keeps the raw text too,
		// because that goes back to the model as context and escape
		// codes there would pollute every later turn's prompt.
		m.displayBlocks = append(m.displayBlocks, displayBlock{kind: blockTutor, raw: msg.reply.Content})
		if len(m.cfg.TranscriptPaths) > 0 {
			if err := appendTranscriptTurn(m.cfg.TranscriptPaths, msg.userMessage, msg.reply.Content); err != nil && !m.transcriptWarned {
				// Warn once, then stay quiet -- a broken transcript disk
				// must not turn every turn into an error banner, and the
				// session itself is unaffected.
				m.transcriptWarned = true
				m.displayBlocks = append(m.displayBlocks, displayBlock{kind: blockNote, raw: activityErrorNote(fmt.Sprintf("transcript export failed (turns continue unrecorded): %v", err))})
			}
		}
		m.refreshViewport()
		return m, nil
	}
	return m, nil
}

// recomputeLayout resizes the textarea and viewport to fit the current
// terminal dimensions and the textarea's current content — called on
// every resize AND every keystroke (per bubbles/textarea and
// bubbles/viewport's own design: neither auto-grows or reacts to resize
// itself, the composing model always owns this). This per-keystroke
// recompute is the actual mechanism that makes input growth dynamic,
// not something built from scratch.
//
// The textarea's height is capped at half the terminal height (a floor
// of minTextareaRows, a ceiling that scales with the real terminal size
// rather than a fixed row count) so a pathological single huge message
// can't starve the viewport entirely — beyond the cap, bubbles/textarea
// scrolls *within its own bounded region*, which is the structural fix
// over the old hand-rolled box: View() always renders exactly as many
// rows as it's told to, so there is no way for overflow to land outside
// the box and corrupt anything else, unlike raw cooked-mode wrapping.
func (m *tutorModel) recomputeLayout() {
	m.textarea.SetWidth(max(m.width-textareaBoxStyle.GetHorizontalFrameSize(), 1))

	maxTaRows := max(m.height/2, minTextareaRows)
	taRows := estimatedTextareaRows(m.textarea.Value(), m.textarea.Width())
	taRows = min(max(taRows, minTextareaRows), maxTaRows)
	m.textarea.SetHeight(taRows)

	// activityRows makes room for activityView's output -- zero whenever
	// no turn is in flight (see activityView), so the activity region
	// costs nothing when idle, unlike the old design's permanently
	// reserved 5 rows. Recomputed here (not just on resize) because this
	// also runs from the activityEventMsg/turnCompleteMsg cases, where
	// len(m.activeCalls) or turnInFlight itself just changed. Each call
	// costs its own header row plus however many indented output rows
	// activityOutputLines actually produces for it (0 while still
	// running, up to activityOutputPreviewLines once it has a result) --
	// not a flat 1-per-call, now that a completed call's output renders
	// on its own indented lines instead of squeezed onto the header.
	activityRows := 0
	if m.turnInFlight {
		activityRows = 1 // status line
		w := m.activityContentWidth()
		for _, c := range m.activeCalls {
			activityRows += 1 + len(activityOutputLines(c, w))
		}
	}

	m.viewport.Width = m.width
	m.viewport.Height = max(m.height-statusBarHeight-taRows-textareaBoxStyle.GetVerticalFrameSize()-activityRows, minViewportRows)
}

// activityContentWidth is the actual text width available inside the
// activity region once viewportContentStyle's left padding is
// subtracted — shared by recomputeLayout's row-count accounting and
// activityView's rendering so they can never disagree about how much
// room a call's header/output preview has, which used to let the
// activity region's content overflow the pane by viewportContentStyle's
// own padding width (the same class of bug refreshViewport's word-wrap
// fixes for the main conversation).
func (m tutorModel) activityContentWidth() int {
	return max(m.width-viewportContentStyle.GetHorizontalFrameSize(), 0)
}

// estimatedTextareaRows estimates how many visual rows value will wrap
// to at the given width — deliberately from the raw value text, not
// textarea.LineCount() (which only counts explicit newlines). A real bug
// this is fixing: a single long line that wraps purely from exceeding
// the box's width never grew LineCount() at all, so a height computed
// from it never grew either — exactly the scenario that used to overflow
// past the old hand-rolled box's last row and corrupt the terminal. Pure
// and testable without a real terminal. width <= 0 is treated as 1 to
// keep the division defined during startup/resize edge cases before a
// real size is known.
func estimatedTextareaRows(value string, width int) int {
	if width <= 0 {
		width = 1
	}
	lines := strings.Split(value, "\n")
	rows := 0
	for _, line := range lines {
		w := utf8.RuneCountInString(line)
		rows += max(1, (w+width-1)/width) // ceil division
	}
	return max(rows, 1)
}

// refreshViewport re-renders displayLines into the viewport and scrolls
// to the bottom — called any time displayLines changes (a submit-echo,
// a reply, a failure fallback) so the newest content is always what's
// visible, and also on every resize (see Update's WindowSizeMsg case) to
// re-flow already-displayed content to the new width.
//
// The content is word-wrapped here, before SetContent, rather than left
// to the viewport itself — a real bug found live: bubbles/viewport does
// NOT wrap long lines. It only ever splits on explicit "\n", and any
// line wider than the viewport gets hard-cut at the frame edge
// (visibleLines' ansi.Cut) with the rest silently discarded, not shown
// on a wrapped row below — exactly what a long assistant reply (one
// long unbroken line, no embedded newlines) hit. lipgloss.Style.Render
// wraps to its Width via cellbuf.Wrap, which is reused here for the
// same reason recomputeLayout uses estimatedTextareaRows: real word
// wrapping, not truncation.
func (m *tutorModel) refreshViewport() {
	w := m.viewport.Width - m.viewport.Style.GetHorizontalFrameSize()
	rendered := make([]string, 0, len(m.displayBlocks)+1)
	for _, b := range m.displayBlocks {
		rendered = append(rendered, renderBlock(b, w))
	}
	if m.turnInFlight && m.streamingText != "" {
		// The in-flight turn's provisional partial reply, styled the
		// same way its final version will be (styleMarkdown renders an
		// unterminated code fence fine, so a reply cut mid-block still
		// displays) -- never appended to displayBlocks itself: the
		// turnCompleteMsg case owns the permanent append.
		rendered = append(rendered, styleMarkdown(m.streamingText, w))
	}
	content := strings.Join(rendered, "\n\n")
	if w > 0 {
		content = lipgloss.NewStyle().Width(w).Render(content)
	}
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
}

// turnCancelledNote is what a user-initiated cancel (Esc or Ctrl-C while
// a turn is in flight, see cancelInFlightTurn) shows -- distinct
// wording, AND distinct color (dimSpan's neutral gray, not
// activityErrorNote's red), from turnFailedFallbackReply/the failure
// detail below it, so a deliberate cancel never reads like the model or
// network failed on its own.
const turnCancelledNote = "Turn cancelled."

// cancelInFlightTurn handles Esc or Ctrl-C while m.turnInFlight (see
// Update's tea.KeyMsg case): cancels the turn and restores the pane to
// idle SYNCHRONOUSLY, without waiting for the turn's own goroutine to
// actually unwind. That's a deliberate choice, not an oversight -- the
// issue asks for control back "immediately", and ctx cancellation
// propagates to an in-flight HTTP request promptly but not
// instantaneously (a request already past its network call, deep in
// react.Agent's own graph bookkeeping between steps, might not observe
// ctx.Done() for another moment). Bumping turnSeq here (not just
// clearing turnCancel) is what makes that eventual straggling result
// harmless once it does land -- see the field's doc comment.
//
// turnCancel may be nil here (a turn queued behind strategy detection,
// see queuedSubmit, never got as far as starting a goroutine at all) --
// guarded rather than assumed, since either way the rest of this
// (clearing turnInFlight, dropping the queue, recovering the text) still
// applies.
func (m tutorModel) cancelInFlightTurn() (tea.Model, tea.Cmd) {
	if m.turnCancel != nil {
		m.turnCancel()
	}
	m.turnCancel = nil
	m.turnInFlight = false
	m.turnSeq++
	m.queuedSubmits = nil

	glowLevel := m.auroraFade()
	m.turnSettledAt = time.Now().Add(-auroraFadeOutLead(glowLevel))
	m.streamingText = ""
	m.streamPainting = false
	m.activeCalls = nil
	m.recomputeLayout()

	m.displayBlocks = append(m.displayBlocks, displayBlock{kind: blockNote, raw: dimSpan(m.recoverPendingText(turnCancelledNote))})
	m.refreshViewport()
	return m, nil
}

// recoverPendingText hands m.pendingUserText back to the user rather
// than ever silently discarding it (see the field's doc comment): back
// into the textarea when it's empty -- the common case, nothing lost --
// or, if the user has since typed ahead (a fresh, unsent draft already
// sits there), appended onto note instead so it stays recoverable by
// copy rather than clobbering what they're mid-typing. Always clears
// pendingUserText: once shown here, either way, there's nothing left to
// recover a second time.
func (m *tutorModel) recoverPendingText(note string) string {
	text := m.pendingUserText
	m.pendingUserText = ""
	if text == "" {
		return note
	}
	if strings.TrimSpace(m.textarea.Value()) == "" {
		m.textarea.SetValue(text)
		return note
	}
	return fmt.Sprintf("%s (you've typed ahead, so your message wasn't put back -- copy it from here: %q)", note, text)
}

// submit handles Enter: empty input is a no-op (nothing to send). A real
// message resets the textarea (growth immediately collapses back down —
// recomputeLayout runs again right after), echoes into the viewport
// immediately (the reply can take many seconds), and starts the turn on
// its own goroutine — mirroring internal/tui/boot.go's own
// buildImageFn/pullModelFn pattern: the goroutine-starting call happens
// directly here (not itself wrapped in a tea.Cmd), returning channels
// that waitForActivityEvent's tea.Cmd drains.
//
// checkComprehension is snapshotted here, before comprehensionCheckPending
// is cleared on m below — matching the old Run() loop's own behavior
// exactly: the flag clears unconditionally on the very first message,
// whether the check succeeds or fails (see startTurn's comment on what
// happens on failure).
//
// helpRequestCount is deliberately NOT incremented here -- a successful
// comprehension check consumes this line without it ever counting as a
// hints-first "help request" (matching the old Run() loop's own
// placement: its helpRequestCount++ sat after the check's early
// `continue`, so it only ever ran for a genuine normal turn). Since
// whether this line ends up on the check path or the normal path is
// only known once startTurn's goroutine runs, the actual increment
// happens there, and the resulting count comes back on turnCompleteMsg
// for Update to apply.
func (m tutorModel) submit() (tea.Model, tea.Cmd) {
	line := strings.TrimSpace(m.textarea.Value())
	if line == "" {
		return m, nil
	}
	checkComprehension := m.comprehensionCheckPending
	m.comprehensionCheckPending = false

	m.textarea.Reset()
	// Read any still-fading glow before flipping turnInFlight, then
	// backdate the start time so the bloom resumes from that level --
	// resubmitting during a fade-out must not blink the glow down to
	// zero before it gathers again.
	glowLevel := m.auroraFade()
	m.turnInFlight = true
	m.turnStartedAt = time.Now().Add(-auroraFadeInLead(glowLevel))
	m.streamingText = ""
	m.streamPainting = false
	// A new generation starts here, whether the turn can actually start
	// its goroutine yet or is about to be queued below -- either way
	// this line's text is now "in flight" from the user's perspective,
	// and turnSeq/pendingUserText need to track it from this point so a
	// cancel arriving before detection resolves still has something
	// concrete to cancel and recover (see cancelInFlightTurn).
	m.turnSeq++
	m.pendingUserText = line
	m.recomputeLayout()
	m.displayBlocks = append(m.displayBlocks, displayBlock{kind: blockUser, raw: line})
	m.refreshViewport()

	// Strategy detection still in flight (a submit typed faster than the
	// probes resolved): hold the turn -- everything visible (echo,
	// thinking state) is already engaged above, and the
	// strategiesDetectedMsg case starts the queued turn the moment the
	// strategies land. startTurn can't run yet because it snapshots m,
	// strategies included.
	if !m.strategiesDetected {
		m.queuedSubmits = append(m.queuedSubmits, queuedSubmit{line: line, checkComprehension: checkComprehension, seq: m.turnSeq})
		return m, nil
	}

	activityCh, streamCh, doneCh, cancel := startTurn(m, line, checkComprehension, m.turnSeq)
	m.turnCancel = cancel
	return m, waitForActivityEvent(m.turnSeq, activityCh, streamCh, doneCh)
}

// queuedSubmit is one message submitted before strategy detection
// resolved -- see the strategiesDetected field's doc comment.
// checkComprehension is captured at submit time because submit consumes
// comprehensionCheckPending immediately, whether or not the turn can
// start yet. seq is submit's own turnSeq at the time this was queued,
// carried through to startTurn once the strategiesDetectedMsg case
// actually starts it (see turnSeq's doc comment).
type queuedSubmit struct {
	line               string
	checkComprehension bool
	seq                int
}

// turnCompleteMsg carries one turn's final outcome — whether it went
// through the comprehension-check path or a normal turn (see startTurn),
// the result-handling shape is identical either way: on success, persist
// (userMessage, reply) to history and show the reply; on failure, show
// turnFailedFallbackReply and persist nothing (see Update's case for
// this message). helpRequestCount is always the count Update should
// adopt regardless of outcome -- unchanged from the pre-turn snapshot
// for a successful comprehension check, incremented for any real
// (non-check) turn attempt, success or failure (see startTurn).
// routingWarning is non-empty only when routing was enabled and the
// handoff decision itself failed (defaulting to the specialist, per
// decideHandoff) -- the turn still completes normally either way, this
// is just visibility into why the decision defaulted. seq is the
// generation this turn was started under (see tutorModel.turnSeq's doc
// comment) -- Update drops this message outright if it no longer
// matches the model's current turnSeq.
type turnCompleteMsg struct {
	seq              int
	reply            *schema.Message
	err              error
	endpoint         string
	userMessage      string
	helpRequestCount int
	routingWarning   string
}

// activityEventMsg carries one live snapshot of the turn's tool calls so
// far (see activityFeed.currentCalls) — pushed by buildActivityChannelOption's
// eino callbacks as they fire, delivered here by waitForActivityEvent.
// Carries its own source channels so Update can re-arm the wait without
// needing to store them as model fields (they're per-turn, ephemeral).
// seq is turnCompleteMsg's own field of the same name, threaded through
// unchanged by waitForActivityEvent on every re-arm.
type activityEventMsg struct {
	seq        int
	calls      []activityCall
	activityCh <-chan []activityCall
	streamCh   <-chan string
	doneCh     <-chan turnCompleteMsg
}

// streamTextMsg carries the in-flight turn's accumulated partial reply
// (see stream.go's pushLatestStreamText) — the streaming counterpart of
// activityEventMsg, delivered by the same waitForActivityEvent wait and
// carrying the same channels (and seq, see activityEventMsg's doc
// comment) so Update re-arms identically.
type streamTextMsg struct {
	seq        int
	text       string
	activityCh <-chan []activityCall
	streamCh   <-chan string
	doneCh     <-chan turnCompleteMsg
}

// pulseTickMsg drives the fading-dot animation (see activityDotColor) —
// free-running for the program's whole lifetime rather than
// started/stopped per turn (see tutorModel.pulsePhase's doc comment).
type pulseTickMsg struct{}

func pulseTickCmd() tea.Cmd {
	return tea.Tick(activityPulseInterval, func(time.Time) tea.Msg { return pulseTickMsg{} })
}

// waitForActivityEvent mirrors internal/tui/boot.go's waitForBuildLine:
// blocks on the turn's live channels — activity snapshots and streamed
// partial text — forwarding each event as it arrives; once activityCh
// closes (the turn's goroutine has finished — see startTurn's deferred
// closes, which close streamCh first) it reads the buffered final
// result from doneCh instead. One single re-armed command rather than
// one per channel, so every caller that drains a turn synchronously
// (RunOneTurn, the test helpers) keeps its one-cmd-chain invariant.
//
// seq is stamped onto every activityEventMsg/streamTextMsg this produces
// (turnCompleteMsg already carries its own, stamped by startTurn) so
// Update can recognize a stale turn's messages -- see tutorModel.turnSeq's
// doc comment.
func waitForActivityEvent(seq int, activityCh <-chan []activityCall, streamCh <-chan string, doneCh <-chan turnCompleteMsg) tea.Cmd {
	return func() tea.Msg {
		for {
			select {
			case calls, ok := <-activityCh:
				if !ok {
					// streamCh is already closed (startTurn's defer order),
					// so nothing more can paint: the turn is done.
					return <-doneCh
				}
				return activityEventMsg{seq: seq, calls: calls, activityCh: activityCh, streamCh: streamCh, doneCh: doneCh}
			case text, ok := <-streamCh:
				if !ok {
					// Stream finished but the turn hasn't: disable this arm
					// (a nil channel never fires) and keep waiting.
					streamCh = nil
					continue
				}
				return streamTextMsg{seq: seq, text: text, activityCh: activityCh, streamCh: streamCh, doneCh: doneCh}
			}
		}
	}
}

// turnTimeout bounds one whole turn -- comprehension check, routing
// decision, and the main model call together (see startTurn) -- unlike
// ollamaRequestTimeout (agent.go), which only bounds a single HTTP
// request and can't stop a react.Agent turn that legitimately makes
// several of those in sequence (each individually within budget) from
// still running for tens of minutes; reactMaxStep=30 steps at up to
// ollamaRequestTimeout=120s each is a long time with no escape but
// Ctrl-D. A var, not a const, so tests can shrink it rather than
// waiting out the real duration -- same pattern as ollamaRequestTimeout.
var turnTimeout = 5 * time.Minute

// startTurn runs one submitted line's whole turn — comprehension check
// (if checkComprehension), routing decision (if m.routingEnabled), and
// the actual model call — on its own goroutine, exactly mirroring the
// old Run() loop's own sequencing. Returns the three channels
// waitForActivityEvent drains, plus the derived context's cancel func so
// the caller can store it for a later Esc/Ctrl-C (see
// tutorModel.turnCancel); m is a snapshot (Go closures capture by value
// here since m is passed as a parameter, not a live reference), which is
// why helpRequestCount/comprehensionCheckPending are already resolved by
// submit before this is called.
//
// seq is stamped onto every turnCompleteMsg this produces so Update can
// recognize a stale result once cancelling has moved the model on to a
// different (or no) generation -- see tutorModel.turnSeq's doc comment.
//
// Every network call in the goroutine below uses ctx (derived from
// m.ctx here), never m.ctx directly -- that's the actual mechanism that
// makes both the timeout and a manual cancel (via the returned
// CancelFunc) reach the request in flight, comprehension check, routing
// decision, or main call alike.
func startTurn(m tutorModel, line string, checkComprehension bool, seq int) (<-chan []activityCall, <-chan string, <-chan turnCompleteMsg, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(m.ctx, turnTimeout)
	activityCh := make(chan []activityCall, 32)
	streamCh := make(chan string, 1)
	doneCh := make(chan turnCompleteMsg, 1)

	go func() {
		// Always released once the turn ends, one way or another --
		// natural completion, timeout, or an external cancel via the
		// returned CancelFunc racing to get here first (cancel is
		// idempotent, so whichever fires first wins harmlessly).
		defer cancel()
		// LIFO: streamCh closes first, so by the time waitForActivityEvent
		// sees activityCh closed (its read-doneCh signal), no more partial
		// text can arrive.
		defer close(activityCh)
		defer close(streamCh)

		feed := &activityFeed{}
		activityOpt := buildActivityChannelOption(feed, activityCh)

		// onTextFor gates streaming per role by that role's own model --
		// worker and orchestrator can differ (e.g. an OpenRouter worker
		// with a local Ollama orchestrator must stream only worker turns).
		onTextFor := func(roleModel string) func(string) {
			if !streamingEnabled(roleModel) {
				return nil
			}
			return func(text string) { pushLatestStreamText(streamCh, text) }
		}

		if checkComprehension {
			checkAgent, checkCM, checkStrategy, checkModel := m.workerAgent, m.workerCM, m.workerStrategy, m.cfg.Model
			if m.routingEnabled {
				checkAgent, checkCM, checkStrategy, checkModel = m.orchestratorAgent, m.orchestratorCM, m.orchestratorStrategy, m.cfg.OrchestratorModel
			}
			checkMessages := prependToolsPrompt(checkStrategy, m.toolCatalogText, comprehensionCheckMessages(m.history, m.cfg.WorkDir, line))
			reply, err := callRole(ctx, checkStrategy, checkAgent, checkCM, m.tools, checkMessages, feed, activityCh, activityOpt, onTextFor(checkModel))
			if err == nil {
				// helpRequestCount stays at its pre-turn snapshot -- a
				// successful comprehension check never counts as a
				// hints-first "help request" (see submit's doc comment).
				doneCh <- turnCompleteMsg{seq: seq, reply: reply, userMessage: line, helpRequestCount: m.helpRequestCount}
				return
			}
			// Couldn't reach the provider for the check -- fall through
			// and handle this same message as a normal turn instead of
			// silently dropping it, exactly like the old Run() loop.
		}

		turnAgent, turnCM, turnStrategy, turnEndpoint, turnModel := m.workerAgent, m.workerCM, m.workerStrategy, m.workerEndpoint, m.cfg.Model
		routingWarning := ""
		if m.routingEnabled {
			handoff, err := decideHandoff(ctx, m.orchestratorCM, line)
			if err != nil {
				// Doesn't abort the turn -- decideHandoff already
				// defaulted to handoff (true) on this same error, so the
				// turn still gets answered by the specialist; this is
				// just visibility into why (rendered by Update, not
				// written directly here -- see turnCompleteMsg.routingWarning's
				// doc comment for why a direct write is unsafe).
				routingWarning = fmt.Sprintf("routing decision failed, defaulting to handoff: %v", err)
			}
			if !handoff {
				turnAgent, turnCM, turnStrategy, turnEndpoint, turnModel = m.orchestratorAgent, m.orchestratorCM, m.orchestratorStrategy, m.orchestratorEndpoint, m.cfg.OrchestratorModel
			}
		}

		// A real (non-check) turn attempt always counts as a help
		// request, success or failure -- matching the old Run() loop's
		// own placement of helpRequestCount++ (unconditional, right
		// before this same call).
		newHelpRequestCount := m.helpRequestCount + 1
		turn := turnMessages(m.cfg.Mode, newHelpRequestCount, line)
		if note := interviewClockNote(m.cfg.Mode, m.cfg.StartedAt, m.cfg.TimeLimitMin, time.Now()); note != nil {
			turn = append([]*schema.Message{note}, turn...)
		}
		requestMessages := prependToolsPrompt(turnStrategy, m.toolCatalogText, append(append([]*schema.Message{}, m.history...), turn...))
		reply, err := callRole(ctx, turnStrategy, turnAgent, turnCM, m.tools, requestMessages, feed, activityCh, activityOpt, onTextFor(turnModel))
		if err != nil {
			doneCh <- turnCompleteMsg{seq: seq, err: err, endpoint: turnEndpoint, userMessage: line, helpRequestCount: newHelpRequestCount, routingWarning: routingWarning}
			return
		}
		doneCh <- turnCompleteMsg{seq: seq, reply: reply, userMessage: line, helpRequestCount: newHelpRequestCount, routingWarning: routingWarning}
	}()

	return activityCh, streamCh, doneCh, cancel
}

// buildActivityChannelOption wires the eino callback machinery
// react.BuildAgentCallback/utils/callbacks.ToolCallbackHandler give real
// OnStart/OnEnd/OnError events for, pushing activityFeed snapshots onto
// a channel for the bubbletea Update loop to pick up (see startTurn,
// Update's activityEventMsg case) instead of writing directly to a
// terminal, which is how this package's now-deleted hand-rolled ANSI box
// (internal/tutor/scrollbox.go) rendered them.
func buildActivityChannelOption(feed *activityFeed, activityCh chan<- []activityCall) agentopt.AgentOption {
	toolHandler := &template.ToolCallbackHandler{
		OnStart: func(ctx context.Context, info *callbacks.RunInfo, input *tool.CallbackInput) context.Context {
			argsPreview := ""
			if input != nil {
				argsPreview = truncateLine(input.ArgumentsInJSON, activityArgsPreviewMax)
			}
			feed.started(compose.GetToolCallID(ctx), info.Name, argsPreview)
			pushActivity(feed, activityCh)
			return ctx
		},
		OnEnd: func(ctx context.Context, info *callbacks.RunInfo, output *tool.CallbackOutput) context.Context {
			resultPreview := ""
			if output != nil {
				resultPreview = truncateLine(output.Response, activityResultPreviewMax)
			}
			feed.finished(compose.GetToolCallID(ctx), resultPreview)
			pushActivity(feed, activityCh)
			return ctx
		},
		OnError: func(ctx context.Context, info *callbacks.RunInfo, err error) context.Context {
			feed.failed(compose.GetToolCallID(ctx), truncateLine(err.Error(), activityResultPreviewMax))
			pushActivity(feed, activityCh)
			return ctx
		},
	}
	handler := react.BuildAgentCallback(nil, toolHandler)
	return agentopt.WithComposeOptions(compose.WithCallbacks(handler))
}

// comprehensionCheckMessages builds one comprehension check's request —
// extracted as a pure function from the old runComprehensionCheck (now
// dead code, removed alongside the rest of the pre-rewrite Run() loop)
// so it stays testable without needing a real agent. The problem
// statement is injected directly as ephemeral context rather than
// calling read_problem_statement (see comprehensionCheckInstruction's
// own doc comment in prompts.go for why).
func comprehensionCheckMessages(history []*schema.Message, workDir, userFirstMessage string) []*schema.Message {
	checkMessages := append([]*schema.Message{}, history...)
	if problem := readProblemStatement(workDir); problem != "" {
		checkMessages = append(checkMessages, schema.SystemMessage("The exercise's problem statement:\n\n"+problem))
	}
	checkMessages = append(checkMessages, schema.SystemMessage(comprehensionCheckInstruction), schema.UserMessage(userFirstMessage))
	return checkMessages
}

// activityView renders the live tool-call activity region: a pulsing
// status dot, then one header line per active call plus (once it has a
// result) its own indented output preview beneath it, via activity.go's
// pulsedStatusLine/pulsedCallLines. Empty whenever no turn is in flight,
// so the region costs zero rows when idle — an improvement over the old
// design's permanently reserved 5 rows (see recomputeLayout's
// activityRows, which must stay in sync with how many lines this
// actually produces).
func (m tutorModel) activityView() string {
	if !m.turnInFlight {
		return ""
	}
	w := m.activityContentWidth()
	lines := make([]string, 0, activityToolLines*(activityOutputPreviewLines+1)+1)
	lines = append(lines, pulsedStatusLine(m.pulsePhase, w))
	for _, c := range m.activeCalls {
		lines = append(lines, pulsedCallLines(c, m.pulsePhase, w)...)
	}
	return viewportContentStyle.Render(strings.Join(lines, "\n"))
}

func (m tutorModel) View() string {
	parts := []string{m.viewport.View()}
	if av := m.activityView(); av != "" {
		parts = append(parts, av)
	}
	parts = append(parts, textareaBoxStyle.Render(m.textarea.View()), m.statusBarView())
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	fade := m.auroraFade()
	if fade <= 0 {
		return content
	}
	t := float64(m.pulsePhase) * activityPulseInterval.Seconds()
	return overlayAurora(content, m.width, m.height, t, auroraBrightness*fade)
}
