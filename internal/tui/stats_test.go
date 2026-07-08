package tui

import (
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

var ansiPatternTUI = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripAnsiTUI(s string) string {
	return ansiPatternTUI.ReplaceAllString(s, "")
}

func TestStatsModel_AnyKeyGoesBack(t *testing.T) {
	m := newStatsModel(treeFixture(), nil)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
	if cmd == nil {
		t.Fatal("expected any keypress to return a quit command")
	}
	if !newM.(statsModel).back {
		t.Error("expected back=true")
	}
}

func TestStatsModel_NonKeyMsgIsIgnored(t *testing.T) {
	m := newStatsModel(treeFixture(), nil)
	newM, cmd := m.Update(tickMsg{})
	if cmd != nil {
		t.Error("expected non-key messages to be ignored, not trigger quit")
	}
	if newM.(statsModel).back {
		t.Error("expected back=false for a non-key message")
	}
}

func TestStatsModel_ViewRendersWithoutPanicOnEmptyHistory(t *testing.T) {
	m := newStatsModel(treeFixture(), nil)
	out := m.View()
	if out == "" {
		t.Error("expected non-empty view even with no attempt history")
	}
}

func TestStatsModel_ViewRendersRecentAttempts(t *testing.T) {
	recent := []tracker.Attempt{
		{ID: 2, ExerciseID: "two-pointers-01", Date: "2026-07-08", Result: tracker.ResultPass},
		{ID: 1, ExerciseID: "off-by-one-01-go", Date: "2026-07-07", Result: tracker.ResultFail},
	}
	m := newStatsModel(treeFixture(), recent)
	out := stripAnsiTUI(m.View())
	if !strings.Contains(out, "two-pointers-01") || !strings.Contains(out, "off-by-one-01-go") {
		t.Errorf("expected recent attempts listed in view:\n%s", out)
	}
}
