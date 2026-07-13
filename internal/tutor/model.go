package tutor

import (
	"context"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

// minTextareaRows/minViewportRows are floors, not targets — the textarea
// always shows at least one row even when empty, and the viewport always
// keeps at least a few rows of conversation visible even when the
// textarea has grown to its cap (see recomputeLayout).
const (
	minTextareaRows    = 1
	minViewportRows    = 3
	textareaBorderRows = 2 // top + bottom border, added by textareaBoxStyle
	textareaBorderCols = 2 // left + right border, same style
)

// textareaBoxStyle replaces the old hand-rolled box borders
// (internal/tutor/scrollbox.go, deleted alongside this rewrite) with a
// lipgloss border — same teal accent used elsewhere in this project's
// palette (docker/tmux.conf, internal/tui/styles.go).
var textareaBoxStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#2FA6A6"))

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

	// stderr is the same real process stream Run's caller passed in --
	// routing-decision-failed logging (see startTurnCmd) still goes here
	// directly, on the turn's own goroutine, entirely independent of
	// bubbletea's own stdout-only rendering.
	stderr io.Writer

	workerAgent          *react.Agent
	orchestratorAgent    *react.Agent
	orchestratorCM       model.ToolCallingChatModel
	workerEndpoint       string
	orchestratorEndpoint string
	routingEnabled       bool

	// history holds only the system prompt plus clean (user, assistant)
	// pairs -- never a failed turn's fallback message, never
	// tool-call scaffolding -- exactly matching the pre-rewrite Run()
	// loop's own history semantics.
	history                   []*schema.Message
	comprehensionCheckPending bool
	helpRequestCount          int
	turnInFlight              bool
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
func newTutorModel(ctx context.Context, cfg Config, stderr io.Writer) (tutorModel, error) {
	m := newTutorLayoutOnly()
	m.ctx = ctx
	m.cfg = cfg
	m.stderr = stderr

	tools, err := buildTools(cfg)
	if err != nil {
		return tutorModel{}, err
	}

	workerCM, err := newChatModel(ctx, cfg.Model, cfg.OllamaHost, cfg.APIKey)
	if err != nil {
		return tutorModel{}, err
	}
	m.workerAgent, err = newAgent(ctx, workerCM, tools)
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

	m.history = []*schema.Message{schema.SystemMessage(systemPromptForMode(cfg.Mode))}
	m.comprehensionCheckPending = wantsComprehensionCheck(cfg.Mode)
	m.displayLines = []string{m.banner}
	m.refreshViewport()

	return m, nil
}

func (m tutorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m tutorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.recomputeLayout()
		return m, tea.ClearScreen

	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			return m.submit()
		}
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		m.recomputeLayout()
		return m, cmd

	case turnCompleteMsg:
		m.turnInFlight = false
		if msg.err != nil {
			fmt.Fprintf(m.stderr, "tutor: could not reach %s: %v\n", msg.endpoint, msg.err)
			m.displayLines = append(m.displayLines, turnFailedFallbackReply)
			m.refreshViewport()
			return m, nil
		}
		m.history = append(m.history, schema.UserMessage(msg.userMessage), schema.AssistantMessage(msg.reply.Content, nil))
		m.displayLines = append(m.displayLines, msg.reply.Content)
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
	m.textarea.SetWidth(max(m.width-textareaBorderCols, 1))

	maxTaRows := max(m.height/2, minTextareaRows)
	taRows := estimatedTextareaRows(m.textarea.Value(), m.textarea.Width())
	taRows = min(max(taRows, minTextareaRows), maxTaRows)
	m.textarea.SetHeight(taRows)

	m.viewport.Width = m.width
	m.viewport.Height = max(m.height-taRows-textareaBorderRows, minViewportRows)
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
// visible.
func (m *tutorModel) refreshViewport() {
	m.viewport.SetContent(strings.Join(m.displayLines, "\n\n"))
	m.viewport.GotoBottom()
}

// submit handles Enter: empty input is a no-op (nothing to send). A real
// message resets the textarea (growth immediately collapses back down —
// recomputeLayout runs again right after), echoes into the viewport
// immediately (the reply can take many seconds), and starts the turn as
// a tea.Cmd — mirroring internal/tui/boot.go's own "kick off background
// work, delivered back via a tea.Msg" pattern, not a blocking call
// inside Update itself.
//
// checkComprehension is snapshotted here, before comprehensionCheckPending
// is cleared on m below — matching the old Run() loop's own behavior
// exactly: the flag clears unconditionally on the very first message,
// whether the check succeeds or fails (see startTurnCmd's comment on
// what happens on failure).
func (m tutorModel) submit() (tea.Model, tea.Cmd) {
	line := strings.TrimSpace(m.textarea.Value())
	if line == "" {
		return m, nil
	}
	checkComprehension := m.comprehensionCheckPending
	m.comprehensionCheckPending = false

	m.textarea.Reset()
	m.recomputeLayout()
	m.turnInFlight = true
	m.helpRequestCount++
	m.displayLines = append(m.displayLines, "> "+line)
	m.refreshViewport()

	return m, startTurnCmd(m, line, checkComprehension)
}

// turnCompleteMsg carries one turn's final outcome — whether it went
// through the comprehension-check path or a normal turn (see
// startTurnCmd), the result-handling shape is identical either way: on
// success, persist (userMessage, reply) to history and show the reply;
// on failure, show turnFailedFallbackReply and persist nothing (see
// Update's case for this message).
type turnCompleteMsg struct {
	reply       *schema.Message
	err         error
	endpoint    string
	userMessage string
}

// startTurnCmd runs one submitted line's whole turn — comprehension
// check (if checkComprehension), routing decision (if
// m.routingEnabled), and the actual model call — on its own goroutine,
// exactly mirroring the old Run() loop's own sequencing, just wrapped as
// a tea.Cmd instead of inline for-loop code. m is a snapshot (bubbletea
// Cmds capture their closure's values at creation, not a live
// reference), which is why helpRequestCount/comprehensionCheckPending
// are already resolved by submit before this is built.
func startTurnCmd(m tutorModel, line string, checkComprehension bool) tea.Cmd {
	return func() tea.Msg {
		if checkComprehension {
			checkAgent := m.workerAgent
			if m.routingEnabled {
				checkAgent = m.orchestratorAgent
			}
			checkMessages := comprehensionCheckMessages(m.history, m.cfg.WorkDir, line)
			reply, err := generateWithLeakRetry(m.ctx, checkAgent, checkMessages)
			if err == nil {
				return turnCompleteMsg{reply: reply, userMessage: line}
			}
			// Couldn't reach the provider for the check -- fall through
			// and handle this same message as a normal turn instead of
			// silently dropping it, exactly like the old Run() loop.
		}

		turnAgent, turnEndpoint := m.workerAgent, m.workerEndpoint
		if m.routingEnabled {
			handoff, err := decideHandoff(m.ctx, m.orchestratorCM, line)
			if err != nil {
				// Doesn't abort the turn -- decideHandoff already
				// defaulted to handoff (true) on this same error, so the
				// turn still gets answered by the specialist; this is
				// just visibility into why.
				fmt.Fprintf(m.stderr, "tutor: routing decision failed, defaulting to handoff: %v\n", err)
			}
			if !handoff {
				turnAgent, turnEndpoint = m.orchestratorAgent, m.orchestratorEndpoint
			}
		}

		requestMessages := append(append([]*schema.Message{}, m.history...), turnMessages(m.cfg.Mode, m.helpRequestCount, line)...)
		reply, err := generateWithLeakRetry(m.ctx, turnAgent, requestMessages)
		if err != nil {
			return turnCompleteMsg{err: err, endpoint: turnEndpoint, userMessage: line}
		}
		return turnCompleteMsg{reply: reply, userMessage: line}
	}
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

func (m tutorModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), textareaBoxStyle.Render(m.textarea.View()))
}
