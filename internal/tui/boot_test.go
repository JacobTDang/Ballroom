package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/preflight"
)

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
