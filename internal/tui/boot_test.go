package tui

import (
	"errors"
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

func TestBootModel_ChecksRunSequentiallyThenReady(t *testing.T) {
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
		t.Fatal("expected a command to run the second check")
	}

	msg2 := cmd2()
	newM2, cmd3 := bm.Update(msg2)
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

func TestBootModel_BuildLineAppendsAndCapsLogAndKeepsWaiting(t *testing.T) {
	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	m := bootModel{building: true, buildLineCh: lineCh, buildErrCh: errCh}

	for i := 0; i < maxBuildLogLines+3; i++ {
		newM, cmd := m.Update(buildLineMsg("line"))
		if cmd == nil {
			t.Fatal("expected buildLineMsg to keep listening for more output")
		}
		m = newM.(bootModel)
	}
	if len(m.buildLines) != maxBuildLogLines {
		t.Errorf("buildLines len = %d, want capped at %d", len(m.buildLines), maxBuildLogLines)
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
		buildLines: []string{"some output"},
	}

	newM, cmd := m.Update(buildDoneMsg{err: nil})
	if cmd != nil {
		t.Fatal("expected no further pending checks after the last one resolves")
	}
	bm := newM.(bootModel)
	if bm.building {
		t.Error("expected building=false once the build resolves")
	}
	if len(bm.buildLines) != 0 {
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
