package tutor

import (
	"context"
	"fmt"
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

// minTextareaRows/minViewportRows are floors, not targets — the textarea
// always shows at least one row even when empty, and the viewport always
// keeps at least a few rows of conversation visible even when the
// textarea has grown to its cap (see recomputeLayout).
const (
	minTextareaRows = 1
	minViewportRows = 3
)

// textareaBoxStyle separates the input from the scrolling conversation
// above it with a single top rule, styled closer to Claude Code's own
// CLI input — replacing an earlier full rounded box on all four sides.
// The full box's colored left edge read visually as a persistent
// vertical "sidebar" running down the pane, a real complaint from live
// use ("remove the side bar ... it will look a lot nicer"); a lone top
// border keeps the same visual separation from the conversation above
// without that vertical bar. PaddingLeft(1) keeps the "> " prompt
// roughly aligned with the conversation's own left inset
// (viewportContentStyle's PaddingLeft(2)) instead of sitting flush
// against the pane's bare left edge.
var textareaBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true, false, false, false).
	BorderForeground(lipgloss.Color("#2FA6A6")).
	PaddingLeft(1)

// viewportContentStyle left-pads the scrolling conversation area — a
// real readability complaint found live: text printed flush against the
// pane's own left edge read as cramped. A lipgloss primitive instead of
// prefixing every print call site by hand (the old architecture's only
// option).
var viewportContentStyle = lipgloss.NewStyle().PaddingLeft(2)

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

	// displayLines is what the viewport actually shows: the banner, each
	// submitted message echoed immediately (before its reply arrives),
	// and each reply -- including a failed turn's honest fallback
	// message, which is shown but deliberately never added to history
	// (see submit/Update's turnCompleteMsg case). Decoupled from history
	// on purpose: history is model context only, this is display only.
	displayLines []string
	banner       string

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
}

// newTutorLayoutOnly builds a model with just the textarea/viewport
// wiring — no agents, no config, no network. Used by this file's own
// pure-layout tests (resize, dynamic growth, Enter-submits-not-newline)
// that have no need to exercise real turn logic. newTutorModel (below)
// is what Run() and every turn-logic test actually uses.
func newTutorLayoutOnly() tutorModel {
	ta := textarea.New()
	ta.Placeholder = "Ask a question..."
	ta.Prompt = "> "
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
		m.banner = fmt.Sprintf("tutor (worker=%s, orchestrator=%s, mode=%s) — connected to %s / %s. Ctrl-D to exit.", cfg.Model, cfg.OrchestratorModel, cfg.Mode, m.workerEndpoint, m.orchestratorEndpoint)
	} else {
		m.banner = fmt.Sprintf("tutor (%s, mode=%s) — connected to %s. Ctrl-D to exit.", cfg.Model, cfg.Mode, m.workerEndpoint)
	}

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
	m.displayLines = []string{m.banner}
	m.refreshViewport()

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
		m.activeCalls = msg.calls
		m.recomputeLayout()
		return m, waitForActivityEvent(msg.activityCh, msg.doneCh)

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
			activityCh, doneCh := startTurn(m, q.line, q.checkComprehension)
			cmds = append(cmds, waitForActivityEvent(activityCh, doneCh))
		}
		m.queuedSubmits = nil
		return m, tea.Batch(cmds...)

	case turnCompleteMsg:
		// Read the glow's level before flipping turnInFlight, then
		// backdate the settle time so the fade-out resumes from that
		// level -- a reply landing mid-bloom must not flash the glow
		// up to full brightness before it dies.
		glowLevel := m.auroraFade()
		m.turnInFlight = false
		m.turnSettledAt = time.Now().Add(-auroraFadeOutLead(glowLevel))
		calls := m.activeCalls
		m.activeCalls = nil
		m.helpRequestCount = msg.helpRequestCount
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
			m.displayLines = append(m.displayLines, summary)
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
			m.displayLines = append(m.displayLines, activityErrorNote(msg.routingWarning))
		}
		if msg.err != nil {
			m.displayLines = append(m.displayLines, turnFailedFallbackReply)
			m.displayLines = append(m.displayLines, activityErrorNote(fmt.Sprintf("could not reach %s: %v", msg.endpoint, msg.err)))
			m.refreshViewport()
			return m, nil
		}
		m.history = append(m.history, schema.UserMessage(msg.userMessage), schema.AssistantMessage(msg.reply.Content, nil))
		// styleMarkdown is display-only: history above keeps the raw
		// reply, because that text goes back to the model as context and
		// escape codes there would pollute every later turn's prompt.
		m.displayLines = append(m.displayLines, styleMarkdown(msg.reply.Content))
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
	m.viewport.Height = max(m.height-taRows-textareaBoxStyle.GetVerticalFrameSize()-activityRows, minViewportRows)
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
	content := strings.Join(m.displayLines, "\n\n")
	if w := m.viewport.Width - m.viewport.Style.GetHorizontalFrameSize(); w > 0 {
		content = lipgloss.NewStyle().Width(w).Render(content)
	}
	m.viewport.SetContent(content)
	m.viewport.GotoBottom()
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
	m.recomputeLayout()
	m.displayLines = append(m.displayLines, userEchoPrefix+line)
	m.refreshViewport()

	// Strategy detection still in flight (a submit typed faster than the
	// probes resolved): hold the turn -- everything visible (echo,
	// thinking state) is already engaged above, and the
	// strategiesDetectedMsg case starts the queued turn the moment the
	// strategies land. startTurn can't run yet because it snapshots m,
	// strategies included.
	if !m.strategiesDetected {
		m.queuedSubmits = append(m.queuedSubmits, queuedSubmit{line: line, checkComprehension: checkComprehension})
		return m, nil
	}

	activityCh, doneCh := startTurn(m, line, checkComprehension)
	return m, waitForActivityEvent(activityCh, doneCh)
}

// queuedSubmit is one message submitted before strategy detection
// resolved -- see the strategiesDetected field's doc comment.
// checkComprehension is captured at submit time because submit consumes
// comprehensionCheckPending immediately, whether or not the turn can
// start yet.
type queuedSubmit struct {
	line               string
	checkComprehension bool
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
// is just visibility into why the decision defaulted.
type turnCompleteMsg struct {
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
type activityEventMsg struct {
	calls      []activityCall
	activityCh <-chan []activityCall
	doneCh     <-chan turnCompleteMsg
}

// pulseTickMsg drives the fading-dot animation (see activityDotColor) —
// free-running for the program's whole lifetime rather than
// started/stopped per turn (see tutorModel.pulsePhase's doc comment).
type pulseTickMsg struct{}

func pulseTickCmd() tea.Cmd {
	return tea.Tick(activityPulseInterval, func(time.Time) tea.Msg { return pulseTickMsg{} })
}

// waitForActivityEvent mirrors internal/tui/boot.go's waitForBuildLine
// exactly: blocks on activityCh, forwarding each snapshot as it arrives;
// once that channel closes (the turn's goroutine has finished — see
// startTurn's deferred close) it reads the buffered final result from
// doneCh instead.
func waitForActivityEvent(activityCh <-chan []activityCall, doneCh <-chan turnCompleteMsg) tea.Cmd {
	return func() tea.Msg {
		calls, ok := <-activityCh
		if ok {
			return activityEventMsg{calls: calls, activityCh: activityCh, doneCh: doneCh}
		}
		return <-doneCh
	}
}

// startTurn runs one submitted line's whole turn — comprehension check
// (if checkComprehension), routing decision (if m.routingEnabled), and
// the actual model call — on its own goroutine, exactly mirroring the
// old Run() loop's own sequencing. Returns the two channels
// waitForActivityEvent drains; m is a snapshot (Go closures capture by
// value here since m is passed as a parameter, not a live reference),
// which is why helpRequestCount/comprehensionCheckPending are already
// resolved by submit before this is called.
func startTurn(m tutorModel, line string, checkComprehension bool) (<-chan []activityCall, <-chan turnCompleteMsg) {
	activityCh := make(chan []activityCall, 32)
	doneCh := make(chan turnCompleteMsg, 1)

	go func() {
		defer close(activityCh)

		feed := &activityFeed{}
		activityOpt := buildActivityChannelOption(feed, activityCh)

		if checkComprehension {
			checkAgent, checkCM, checkStrategy := m.workerAgent, m.workerCM, m.workerStrategy
			if m.routingEnabled {
				checkAgent, checkCM, checkStrategy = m.orchestratorAgent, m.orchestratorCM, m.orchestratorStrategy
			}
			checkMessages := prependToolsPrompt(checkStrategy, m.toolCatalogText, comprehensionCheckMessages(m.history, m.cfg.WorkDir, line))
			reply, err := callRole(m.ctx, checkStrategy, checkAgent, checkCM, m.tools, checkMessages, feed, activityCh, activityOpt)
			if err == nil {
				// helpRequestCount stays at its pre-turn snapshot -- a
				// successful comprehension check never counts as a
				// hints-first "help request" (see submit's doc comment).
				doneCh <- turnCompleteMsg{reply: reply, userMessage: line, helpRequestCount: m.helpRequestCount}
				return
			}
			// Couldn't reach the provider for the check -- fall through
			// and handle this same message as a normal turn instead of
			// silently dropping it, exactly like the old Run() loop.
		}

		turnAgent, turnCM, turnStrategy, turnEndpoint := m.workerAgent, m.workerCM, m.workerStrategy, m.workerEndpoint
		routingWarning := ""
		if m.routingEnabled {
			handoff, err := decideHandoff(m.ctx, m.orchestratorCM, line)
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
				turnAgent, turnCM, turnStrategy, turnEndpoint = m.orchestratorAgent, m.orchestratorCM, m.orchestratorStrategy, m.orchestratorEndpoint
			}
		}

		// A real (non-check) turn attempt always counts as a help
		// request, success or failure -- matching the old Run() loop's
		// own placement of helpRequestCount++ (unconditional, right
		// before this same call).
		newHelpRequestCount := m.helpRequestCount + 1
		requestMessages := prependToolsPrompt(turnStrategy, m.toolCatalogText, append(append([]*schema.Message{}, m.history...), turnMessages(m.cfg.Mode, newHelpRequestCount, line)...))
		reply, err := callRole(m.ctx, turnStrategy, turnAgent, turnCM, m.tools, requestMessages, feed, activityCh, activityOpt)
		if err != nil {
			doneCh <- turnCompleteMsg{err: err, endpoint: turnEndpoint, userMessage: line, helpRequestCount: newHelpRequestCount, routingWarning: routingWarning}
			return
		}
		doneCh <- turnCompleteMsg{reply: reply, userMessage: line, helpRequestCount: newHelpRequestCount, routingWarning: routingWarning}
	}()

	return activityCh, doneCh
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
	parts = append(parts, textareaBoxStyle.Render(m.textarea.View()))
	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	fade := m.auroraFade()
	if fade <= 0 {
		return content
	}
	t := float64(m.pulsePhase) * activityPulseInterval.Seconds()
	return overlayAurora(content, m.width, m.height, t, auroraBrightness*fade)
}
