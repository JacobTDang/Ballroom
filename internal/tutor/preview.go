package tutor

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// RunPreview runs the pane on canned fixture content with no network,
// agents, or session — the visual harness for restyling work
// (cmd/tutor-preview). Every visual state is reachable from the
// keyboard so a tmux capture-pane can eyeball them:
//
//	ctrl+t  toggle a fake in-flight turn (activity region + aurora)
//	ctrl+s  toggle fake streaming text during the fake turn
//	enter   echo the typed line as a user block (no model call)
//	ctrl+d  quit
func RunPreview() error {
	m := newTutorLayoutOnly()
	m.cfg = Config{Model: "preview-model", Mode: exercise.TutorModeHintsFirst}
	m.workerEndpoint = "preview (no network)"
	m.helpRequestCount = 2
	// Detection is skipped entirely (previewModel.Init never calls
	// tutorModel.Init, whose strategy probe would hit the nil agents),
	// so mark it resolved for anything that consults the flag.
	m.strategiesDetected = true
	m.displayBlocks = []displayBlock{
		{kind: blockUser, raw: "how should I approach two sum? my brute force keeps timing out on the big cases and I am not sure what the next step is"},
		{kind: blockTutor, raw: previewReply},
		{kind: blockNote, raw: toolUsageSummary(previewCalls(time.Now().Add(-3*time.Second)), 76)},
	}
	m.refreshViewport()

	p := tea.NewProgram(previewModel{inner: m}, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tutor preview: %w", err)
	}
	return nil
}

// previewReply exercises every markdown construct the pane styles:
// bold, inline code, a header, a list (raw until the renderer learns
// lists), and a fenced block that renders as an editor card.
const previewReply = "### narrowing it down\n\n" +
	"Before reaching for code, think about what the **inner loop** is\n" +
	"actually looking for on each pass — a specific value, right?\n\n" +
	"- what would let you ask \"have I seen this value\" in O(1)?\n" +
	"- what do you need to remember alongside each value?\n\n" +
	"> hint: `target - nums[i]` is a single number, not a search.\n\n" +
	"Here is the shape of the scan (not the whole answer):\n\n" +
	"```python\nseen = {}\nfor i, x in enumerate(nums):\n    want = target - x\n    # ... your move\n```\n"

// previewCalls is the settled tool-usage fixture; completedAt is
// backdated so the summary's dots render at rest color, same as the
// real turnCompleteMsg append.
func previewCalls(done time.Time) []activityCall {
	return []activityCall{
		{callID: "1", name: "read_problem_statement", status: "done", detail: "Two Sum — given an array of integers...", completedAt: done},
		{callID: "2", name: "read_solution_file", status: "done", detail: "def two_sum(nums, target):\\n    for i in range(len(nums)):", completedAt: done},
	}
}

// previewStreamText is the fake in-flight partial reply (ctrl+s): an
// unterminated fence, deliberately, so the bottomless-card streaming
// state is reachable.
const previewStreamText = "Let me look at what changed in your loop.\n\n```python\nfor i, x in enumerate(nums):\n    want = target - x"

// previewModel wraps tutorModel to intercept the harness keys and, on
// enter, echo the draft as a user block instead of starting a real
// turn (newTutorLayoutOnly has no agents; a real submit would panic).
type previewModel struct {
	inner tutorModel
}

// Init deliberately bypasses tutorModel.Init: its strategy-detection
// probe calls the real chat models, which the layout-only preview
// doesn't have. The cursor blink is the only piece worth keeping.
func (p previewModel) Init() tea.Cmd { return textarea.Blink }

func (p previewModel) View() string { return p.inner.View() }

func (p previewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.Type {
		case tea.KeyCtrlT:
			p.inner.turnInFlight = !p.inner.turnInFlight
			if p.inner.turnInFlight {
				p.inner.activeCalls = previewCalls(time.Time{})
				p.inner.activeCalls[1].status = "running"
				p.inner.recomputeLayout()
				p.inner.refreshViewport()
				return p, pulseTickCmd()
			}
			p.inner.turnSettledAt = time.Now()
			p.inner.activeCalls = nil
			p.inner.streamingText = ""
			p.inner.recomputeLayout()
			p.inner.refreshViewport()
			return p, nil
		case tea.KeyCtrlS:
			if p.inner.streamingText == "" {
				p.inner.streamingText = previewStreamText
			} else {
				p.inner.streamingText = ""
			}
			p.inner.refreshViewport()
			return p, nil
		case tea.KeyEnter:
			if line := strings.TrimSpace(p.inner.textarea.Value()); line != "" {
				p.inner.displayBlocks = append(p.inner.displayBlocks, displayBlock{kind: blockUser, raw: line})
				p.inner.textarea.Reset()
				p.inner.recomputeLayout()
				p.inner.refreshViewport()
			}
			return p, nil
		}
	}
	inner, cmd := p.inner.Update(msg)
	p.inner = inner.(tutorModel)
	return p, cmd
}
