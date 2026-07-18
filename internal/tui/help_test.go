package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- stageHelp: "?" opens it from every input-taking stage, and Esc/q/?
// returns to wherever it was opened from (issue #242). ---

func TestAppModel_Help_OpensFromEveryInputStageAndReturnsToOrigin(t *testing.T) {
	question := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")}

	stages := []struct {
		name  string
		stage appStage
	}{
		{"stageMain", stageMain},
		{"stageCategories", stageCategories},
		{"stageDSACategories", stageDSACategories},
		{"stageProblems", stageProblems},
		{"stageLanguage", stageLanguage},
		{"stageSearch", stageSearch},
		{"stageStats", stageStats},
		{"stageStatsDetail", stageStatsDetail},
		{"stageSettings", stageSettings},
	}

	for _, c := range stages {
		t.Run(c.name, func(t *testing.T) {
			m := appModel{stage: c.stage}

			opened, cmd := m.Update(question)
			if cmd != nil {
				t.Error("expected no external command — opening help stays inside the same program")
			}
			gotOpened := opened.(appModel)
			if gotOpened.stage != stageHelp {
				t.Fatalf("stage = %v, want stageHelp", gotOpened.stage)
			}
			if gotOpened.helpOrigin != c.stage {
				t.Fatalf("helpOrigin = %v, want %v (the stage help was opened from)", gotOpened.helpOrigin, c.stage)
			}

			// Esc, q, and ? must all return to the origin stage.
			for _, back := range []tea.KeyMsg{
				{Type: tea.KeyEsc},
				{Type: tea.KeyRunes, Runes: []rune("q")},
				question,
			} {
				closed, cmd := gotOpened.Update(back)
				if cmd != nil {
					t.Errorf("%v: expected no external command closing help", back)
				}
				gotClosed := closed.(appModel)
				if gotClosed.stage != c.stage {
					t.Errorf("%v: stage = %v, want %v (the origin stage)", back, gotClosed.stage, c.stage)
				}
			}
		})
	}
}

func TestAppModel_Help_QuestionMarkMidFilterFeedsTheFilterInstead(t *testing.T) {
	// Same carve-out this codebase already gives "q" in stageProblems/
	// stageModelPicker (see updateProblems): once the user has started
	// typing, every rune -- including the ones that double as shortcuts
	// when nothing's typed yet -- feeds the filter instead of firing the
	// shortcut, since it might be part of a real query.
	m := appModel{stage: stageProblems, problemFilter: "arr"}
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	got := newM.(appModel)
	if got.stage != stageProblems {
		t.Fatalf("stage = %v, want stageProblems (still typing, not a help shortcut)", got.stage)
	}
	if got.problemFilter != "arr?" {
		t.Errorf("problemFilter = %q, want %q", got.problemFilter, "arr?")
	}
}

// --- content: the table is the single source of truth, and every key it
// lists must render somewhere in the help screen. ---

func TestHelpSections_NoEmptyKeysOrDescriptions(t *testing.T) {
	if len(helpSections) == 0 {
		t.Fatal("helpSections is empty")
	}
	for _, section := range helpSections {
		if strings.TrimSpace(section.title) == "" {
			t.Error("section with an empty title")
		}
		if len(section.keys) == 0 {
			t.Errorf("section %q has no keys", section.title)
		}
		for _, k := range section.keys {
			if strings.TrimSpace(k.key) == "" || strings.TrimSpace(k.desc) == "" {
				t.Errorf("section %q has an empty key/desc pair: %+v", section.title, k)
			}
		}
	}
}

func TestRenderHelp_ContainsEveryKeyStringFromTheTable(t *testing.T) {
	// Generous size -- this test is about content completeness, not
	// layout at the panel floor (see the wrapping-damage test below for
	// that), so give it enough room that nothing has to wrap.
	m := appModel{stage: stageHelp, helpOrigin: stageMain, width: 220, height: 60}
	view := stripAnsiTUI(m.View())
	for _, section := range helpSections {
		if !strings.Contains(view, section.title) {
			t.Errorf("help view missing section title %q:\n%s", section.title, view)
		}
		for _, k := range section.keys {
			if !strings.Contains(view, k.key) {
				t.Errorf("help view missing key %q from section %q:\n%s", k.key, section.title, view)
			}
		}
	}
}

// TestAppModel_Help_RendersWithoutWrappingDamageAtMinimumPanelSize pins
// the help screen to the same floor panelDimensions clamps every other
// screen to (see dashboard_test.go's TestPanelDimensions_
// ClampsToMinimumOnTinyTerminal) and checks the bordered panel shell
// stays structurally intact: every row the exact same width, and the
// border glyphs still frame it -- proof that however the three-section
// key table has to reflow at the floor, it never blows the box open or
// leaves a ragged edge the way an unbounded word-wrap could (see
// dashboard.go's own comments on this failure mode for the ball+banner
// block).
func TestAppModel_Help_RendersWithoutWrappingDamageAtMinimumPanelSize(t *testing.T) {
	m := appModel{stage: stageHelp, helpOrigin: stageMain, width: 1, height: 1}
	view := m.View()

	wantW, wantH := panelDimensions(1, 1)
	if wantW != minPanelWidth || wantH != minPanelHeight {
		t.Fatalf("test setup assumption broken: panelDimensions(1,1) = (%d,%d), want the (%d,%d) floor",
			wantW, wantH, minPanelWidth, minPanelHeight)
	}

	lines := strings.Split(view, "\n")
	if len(lines) == 0 {
		t.Fatal("empty view")
	}
	for i, line := range lines {
		stripped := stripAnsiTUI(line)
		if w := lipgloss.Width(stripped); w != wantW {
			t.Errorf("line %d width = %d, want %d (panel border broke): %q", i, w, wantW, stripped)
		}
	}

	first := stripAnsiTUI(lines[0])
	if !strings.HasPrefix(first, "╔") || !strings.HasSuffix(first, "╗") {
		t.Errorf("top border row malformed: %q", first)
	}
	last := stripAnsiTUI(lines[len(lines)-1])
	if !strings.HasPrefix(last, "╚") || !strings.HasSuffix(last, "╝") {
		t.Errorf("bottom border row malformed: %q", last)
	}
	for i := 1; i < len(lines)-1; i++ {
		l := stripAnsiTUI(lines[i])
		if !strings.HasPrefix(l, "║") || !strings.HasSuffix(l, "║") {
			t.Errorf("interior row %d missing its side borders: %q", i, l)
		}
	}
}

func TestRenderHelp_ShowsBackHint(t *testing.T) {
	m := appModel{stage: stageHelp, helpOrigin: stageMain, width: 220, height: 60}
	view := stripAnsiTUI(m.View())
	if !strings.Contains(view, "esc") && !strings.Contains(view, "q") {
		t.Errorf("expected a back-hint mentioning esc/q, got:\n%s", view)
	}
}
