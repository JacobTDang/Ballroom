package tui

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/preflight"
)

// fakeBuildImage substitutes buildImageFn in tests so triggering a live
// build doesn't shell out to real docker. Returns a restore func to defer.
func fakeBuildImage(lineCh <-chan string, errCh <-chan error) func() {
	orig := buildImageFn
	buildImageFn = func(config.Config) (<-chan string, <-chan error) {
		return lineCh, errCh
	}
	return func() { buildImageFn = orig }
}

// noCheckDelay zeroes the pacing delay between checks so tests that drive
// the sequence manually don't actually sleep. Returns a restore func.
func noCheckDelay() func() {
	orig := checkStartDelay
	checkStartDelay = 0
	return func() { checkStartDelay = orig }
}

func TestBootModel_ChecksRunSequentiallyThenReady(t *testing.T) {
	defer noCheckDelay()()

	var order []string
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { order = append(order, "a"); return preflight.Check{Name: "a", OK: true} },
			func() preflight.Check { order = append(order, "b"); return preflight.Check{Name: "b", OK: true} },
		},
	}

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init() should return a command to run the first check")
	}

	// Init batches the first check with the banner's tick command, so
	// unwrap the batch to find the checkDoneMsg among them.
	batch, ok := cmd().(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected Init() to return a batch of commands, got %T", cmd())
	}
	var msg1 tea.Msg
	for _, c := range batch {
		if c == nil {
			continue
		}
		result := c()
		if _, isCheck := result.(checkDoneMsg); isCheck {
			msg1 = result
			break
		}
	}
	if msg1 == nil {
		t.Fatal("expected a checkDoneMsg among Init()'s batched commands")
	}

	newM, cmd2 := m.Update(msg1)
	bm := newM.(bootModel)
	if len(bm.checks) != 1 || bm.checks[0].Name != "a" {
		t.Fatalf("expected 1 check recorded (a), got %+v", bm.checks)
	}
	if bm.ready {
		t.Fatal("should not be ready after only 1 of 2 checks")
	}
	if cmd2 == nil {
		t.Fatal("expected a command to schedule the second check")
	}

	// The next check doesn't run immediately — it's paced behind a
	// startCheckMsg (see checkStartDelay) so checks visibly appear one
	// at a time instead of all resolving within the same frame.
	startMsg, ok := cmd2().(startCheckMsg)
	if !ok {
		t.Fatalf("expected a startCheckMsg to be scheduled, got %T", cmd2())
	}
	newM2a, cmd2b := bm.Update(startMsg)
	bm2a := newM2a.(bootModel)
	if cmd2b == nil {
		t.Fatal("expected startCheckMsg to trigger running the second check")
	}

	msg2 := cmd2b()
	newM2, cmd3 := bm2a.Update(msg2)
	bm2 := newM2.(bootModel)
	if len(bm2.checks) != 2 {
		t.Fatalf("expected 2 checks recorded, got %d", len(bm2.checks))
	}
	if !bm2.ready {
		t.Fatal("expected ready=true once all checks complete")
	}
	if cmd3 != nil {
		t.Error("expected no further command once ready")
	}
	if order[0] != "a" || order[1] != "b" {
		t.Errorf("checks did not run in order: %v", order)
	}
}

func TestBootModel_EnterOnlyQuitsWhenReady(t *testing.T) {
	notReady := bootModel{ready: false}
	_, cmd := notReady.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("enter before ready should be a no-op")
	}

	ready := bootModel{ready: true}
	newM, cmd2 := ready.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd2 == nil {
		t.Fatal("enter when ready should return a quit command")
	}
	if newM.(bootModel).quit {
		t.Error("enter should proceed (quit=false), not request quit")
	}
}

func TestBootModel_QAlwaysRequestsQuit(t *testing.T) {
	m := bootModel{ready: false}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command even before ready")
	}
	if !newM.(bootModel).quit {
		t.Error("expected quit=true after pressing q")
	}
}

func TestBootModel_TickAdvancesPhaseAndReschedules(t *testing.T) {
	m := bootModel{}
	newM, cmd := m.Update(tickMsg{})
	if cmd == nil {
		t.Fatal("expected tick to reschedule another tick command")
	}
	if newM.(bootModel).phase != 1 {
		t.Errorf("phase = %d, want 1", newM.(bootModel).phase)
	}
}

func TestBootModel_ImageNotOkTriggersLiveBuildInsteadOfFailing(t *testing.T) {
	lineCh := make(chan string)
	errCh := make(chan error)
	defer fakeBuildImage(lineCh, errCh)()

	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: "Image", OK: false, Detail: "not built"} },
		},
		checks: []preflight.Check{{Name: "Docker", OK: true}},
	}

	newM, cmd := m.Update(checkDoneMsg(preflight.Check{Name: "Image", OK: false, Detail: "not built"}))
	if cmd == nil {
		t.Fatal("expected a command to wait for build output")
	}
	bm := newM.(bootModel)
	if !bm.building {
		t.Fatal("expected building=true")
	}
	if len(bm.checks) != 1 {
		t.Errorf("image check should not be appended yet while building, got %+v", bm.checks)
	}
	if bm.ready {
		t.Error("should not be ready while the image is still building")
	}
}

func TestBootModel_ImageOkDoesNotTriggerBuild(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: "Image", OK: true} },
		},
		checks: []preflight.Check{{Name: "Docker", OK: true}},
	}
	newM, _ := m.Update(checkDoneMsg(preflight.Check{Name: "Image", OK: true, Detail: "built"}))
	bm := newM.(bootModel)
	if bm.building {
		t.Error("should not start a build when the image check already passed")
	}
	if len(bm.checks) != 2 {
		t.Errorf("expected the image check to be appended normally, got %+v", bm.checks)
	}
}

func TestBootModel_DockerNotReachableSkipsBuildAttempt(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: false} },
			func() preflight.Check { return preflight.Check{Name: "Image", OK: false} },
		},
		checks: []preflight.Check{{Name: "Docker", OK: false}},
	}
	newM, _ := m.Update(checkDoneMsg(preflight.Check{Name: "Image", OK: false, Detail: "not built"}))
	bm := newM.(bootModel)
	if bm.building {
		t.Error("should not attempt a build when Docker itself isn't reachable")
	}
	if len(bm.checks) != 2 {
		t.Errorf("expected the image check to be appended as a plain failure, got %+v", bm.checks)
	}
}

func TestStepID_ExtractsLeadingHashNumberToken(t *testing.T) {
	cases := map[string]string{
		"#5 [2/4] RUN go build ./...": "#5",
		"#12 0.235 exporting layers":  "#12",
		"#1 DONE 0.1s":                "#1",
	}
	for line, want := range cases {
		if got := stepID(line); got != want {
			t.Errorf("stepID(%q) = %q, want %q", line, got, want)
		}
	}
}

func TestStepID_ReturnsEmptyForNonStepLine(t *testing.T) {
	cases := []string{"", "no leading hash", "#", "#abc not digits"}
	for _, line := range cases {
		if got := stepID(line); got != "" {
			t.Errorf("stepID(%q) = %q, want empty", line, got)
		}
	}
}

func TestBootModel_BuildLineGroupsByStepAndCapsEachStepAtThreeLines(t *testing.T) {
	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	m := bootModel{building: true, buildLineCh: lineCh, buildErrCh: errCh}

	for i := 0; i < maxStepLogLines+5; i++ {
		newM, cmd := m.Update(buildLineMsg("#3 some output"))
		if cmd == nil {
			t.Fatal("expected buildLineMsg to keep listening for more output")
		}
		m = newM.(bootModel)
	}
	if len(m.buildSteps) != 1 {
		t.Fatalf("expected all lines from the same step to stay grouped as 1 entry, got %d", len(m.buildSteps))
	}
	if len(m.buildSteps[0].lines) != maxStepLogLines {
		t.Errorf("step lines = %d, want capped at %d", len(m.buildSteps[0].lines), maxStepLogLines)
	}
}

func TestBootModel_BuildLinePersistsPriorStepsWhenNewStepStarts(t *testing.T) {
	m := bootModel{building: true}

	newM, _ := m.Update(buildLineMsg("#1 [1/4] FROM golang"))
	m = newM.(bootModel)
	newM, _ = m.Update(buildLineMsg("#2 [2/4] RUN go build"))
	m = newM.(bootModel)

	if len(m.buildSteps) != 2 {
		t.Fatalf("expected both steps to persist once a new step starts, got %d: %+v", len(m.buildSteps), m.buildSteps)
	}
	if m.buildSteps[0].id != "#1" || m.buildSteps[1].id != "#2" {
		t.Errorf("expected steps in order #1, #2, got %+v", m.buildSteps)
	}
}

func TestBootModel_BuildDoneSuccessAdvancesAsPassingImageCheck(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: "Image", OK: true} },
		},
		checks:     []preflight.Check{{Name: "Docker", OK: true}},
		building:   true,
		buildSteps: []buildStepLog{{id: "#1", lines: []string{"some output"}}},
	}

	newM, cmd := m.Update(buildDoneMsg{err: nil})
	if cmd != nil {
		t.Fatal("expected no further pending checks after the last one resolves")
	}
	bm := newM.(bootModel)
	if bm.building {
		t.Error("expected building=false once the build resolves")
	}
	if len(bm.buildSteps) != 0 {
		t.Error("expected the build log to clear once resolved (panel collapses)")
	}
	if len(bm.checks) != 2 || !bm.checks[1].OK {
		t.Fatalf("expected a passing Image check appended, got %+v", bm.checks)
	}
	if !bm.ready {
		t.Error("expected ready=true once the (now only) remaining check resolves")
	}
}

func TestBootModel_BuildDoneFailureAdvancesAsFailingImageCheck(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: "Image", OK: true} },
		},
		checks:   []preflight.Check{{Name: "Docker", OK: true}},
		building: true,
	}

	newM, _ := m.Update(buildDoneMsg{err: errors.New("boom")})
	bm := newM.(bootModel)
	if bm.building {
		t.Error("expected building=false once the build resolves, even on failure")
	}
	if len(bm.checks) != 2 || bm.checks[1].OK {
		t.Fatalf("expected a failing Image check appended, got %+v", bm.checks)
	}
	if !bm.ready {
		t.Error("expected ready=true so the user isn't stuck — they can see the failure and quit")
	}
}

func TestLastLines_CapsToMostRecentNNonEmptyLines(t *testing.T) {
	out := "one\ntwo\n\nthree\nfour\nfive"
	got := lastLines(out, 3)
	want := []string{"three", "four", "five"}
	if len(got) != len(want) {
		t.Fatalf("lastLines = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("lastLines[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestLastLines_EmptyOutputReturnsNoLines(t *testing.T) {
	if got := lastLines("", 3); len(got) != 0 {
		t.Errorf("lastLines(\"\", 3) = %v, want empty", got)
	}
}

func TestBootModel_RenderRightColumnShowsRealCommandAndOutputForRecentCheck(t *testing.T) {
	m := bootModel{
		checks: []preflight.Check{
			{Name: "Docker daemon", OK: true, Detail: "running", Command: "docker info", Output: "Client:\n Version: 27.0.0\nServer:\n Containers: 3"},
		},
	}
	out := m.renderRightColumn()
	if !strings.Contains(out, "docker info") {
		t.Errorf("expected the real command to be visible, got:\n%s", out)
	}
	if !strings.Contains(out, "Containers: 3") {
		t.Errorf("expected real command output to be visible, got:\n%s", out)
	}
}

func TestBootModel_RenderRightColumnNeverCollapsesEarlierChecks(t *testing.T) {
	m := bootModel{
		checks: []preflight.Check{
			{Name: "Docker daemon", OK: true, Detail: "running", Command: "docker info", Output: "oldest output line"},
			{Name: "Practice image", OK: true, Detail: "built", Command: "docker image inspect x", Output: "sha256:abc"},
			{Name: "Ollama", OK: true, Detail: "reachable", Command: "GET /api/tags", Output: `{"models":[]}`},
			{Name: "Tutor model", OK: true, Detail: "ready", Command: "GET /api/tags", Output: `{"models":[{"name":"x"}]}`},
		},
	}
	out := m.renderRightColumn()
	if !strings.Contains(out, "oldest output line") {
		t.Errorf("expected the oldest check to stay expanded and keep showing its output, got:\n%s", out)
	}
	if !strings.Contains(out, "sha256:abc") {
		t.Errorf("expected a later check to also still show its output, got:\n%s", out)
	}
}

func TestBootModel_RenderRightColumnQueuedCheckShowsNameOnlyNoCommand(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker"} },
			func() preflight.Check { return preflight.Check{Name: "Ollama"} },
		},
		checkNames: []string{"Docker daemon", "Ollama"},
	}
	out := m.renderRightColumn()
	if !strings.Contains(out, "Ollama") {
		t.Errorf("expected the queued check's name to be visible, got:\n%s", out)
	}
	if strings.Contains(out, "$") {
		t.Errorf("expected no command to be shown for a check that hasn't been invoked yet, got:\n%s", out)
	}
}

func TestBootModel_CtrlCAlwaysRequestsQuit(t *testing.T) {
	m := bootModel{ready: false}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected ctrl+c to return a quit command")
	}
	if !newM.(bootModel).quit {
		t.Error("expected quit=true after ctrl+c")
	}
}
