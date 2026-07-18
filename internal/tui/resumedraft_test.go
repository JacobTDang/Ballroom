package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/draft"
	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// seedDraft writes a draft for exerciseID under a temp data dir and
// returns the dir, so the resume-prompt tests exercise the real
// internal/draft round-trip rather than a hand-built fixture.
func seedDraft(t *testing.T, exerciseID, content string) string {
	t.Helper()
	dataDir := t.TempDir()
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("write workspace solution: %v", err)
	}
	if _, err := draft.Snapshot(dataDir, exerciseID, workspace); err != nil {
		t.Fatalf("seed draft: %v", err)
	}
	return dataDir
}

func draftModel(t *testing.T, dataDir string, problem catalog.ProblemStatus) appModel {
	t.Helper()
	m := appModel{
		cfg:             config.Config{DataDir: dataDir},
		stage:           stageLanguage,
		selectedProblem: problem,
	}
	return m
}

// TestLaunchExercise_WithDraftEntersResumeStage: the whole point of the
// prompt — a saved draft must never be silently overwritten or silently
// resumed. Both the manual language pick and the default-language fast
// path funnel through launchExercise, so this covers both.
func TestLaunchExercise_WithDraftEntersResumeStage(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	ex := problem.Variants[0].Exercise
	dataDir := seedDraft(t, ex.ID, "package main\n\nfunc twoSum() {}\n")

	m := draftModel(t, dataDir, problem)
	got, cmd := m.launchExercise(ex)
	next := got.(appModel)

	if next.stage != stageResumeDraft {
		t.Fatalf("stage = %v, want stageResumeDraft", next.stage)
	}
	if cmd != nil {
		t.Error("launchExercise with a draft should not quit yet")
	}
	if next.pendingExercise.ID != ex.ID {
		t.Errorf("pendingExercise = %q, want %q", next.pendingExercise.ID, ex.ID)
	}
	if len(next.pendingDraft.Preview) == 0 {
		t.Error("pendingDraft has no preview lines to show the user")
	}
}

func TestLaunchExercise_WithoutDraftLaunchesImmediately(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	ex := problem.Variants[0].Exercise

	m := draftModel(t, t.TempDir(), problem)
	got, cmd := m.launchExercise(ex)
	next := got.(appModel)

	if next.stage == stageResumeDraft {
		t.Error("no draft exists, want no resume prompt")
	}
	if next.outcome != outcomeRunExercise || next.exerciseToRun.ID != ex.ID {
		t.Errorf("outcome/exercise = %v/%q, want run %q", next.outcome, next.exerciseToRun.ID, ex.ID)
	}
	if cmd == nil {
		t.Error("want tea.Quit to hand off to the session")
	}
}

// TestResolveLanguageStage_DefaultLanguageStillPrompts pins the trap
// called out in planning: the default-language fast path bypassed the
// language screen entirely, so a user with a default set would never
// have seen the prompt.
func TestResolveLanguageStage_DefaultLanguageStillPrompts(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	ex := problem.Variants[0].Exercise
	dataDir := seedDraft(t, ex.ID, "package main\n")

	m := draftModel(t, dataDir, problem)
	m.cfg.DefaultLanguage = exercise.LanguageGo

	got, _ := m.resolveLanguageStage()
	if next := got.(appModel); next.stage != stageResumeDraft {
		t.Fatalf("stage = %v, want stageResumeDraft even on the default-language fast path", next.stage)
	}
}

func TestUpdateResumeDraft_ResumeCarriesDraftDir(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	ex := problem.Variants[0].Exercise
	dataDir := seedDraft(t, ex.ID, "package main\n")

	m := draftModel(t, dataDir, problem)
	started, _ := m.launchExercise(ex)
	m = started.(appModel)

	got, cmd := m.updateResumeDraft(tea.KeyMsg{Type: tea.KeyEnter})
	next := got.(appModel)

	if next.outcome != outcomeRunExercise {
		t.Fatalf("outcome = %v, want outcomeRunExercise", next.outcome)
	}
	if next.draftDirToUse != draft.Dir(dataDir, ex.ID) {
		t.Errorf("draftDirToUse = %q, want the draft dir", next.draftDirToUse)
	}
	if cmd == nil {
		t.Error("want tea.Quit after choosing resume")
	}
}

// TestUpdateResumeDraft_StartFreshArchivesAndLaunchesClean: "start
// fresh" must archive rather than delete — the user gets the starter,
// but yesterday's work is still recoverable on disk.
func TestUpdateResumeDraft_StartFreshArchivesAndLaunchesClean(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	ex := problem.Variants[0].Exercise
	dataDir := seedDraft(t, ex.ID, "package main\n\n// my work\n")

	m := draftModel(t, dataDir, problem)
	started, _ := m.launchExercise(ex)
	m = started.(appModel)

	m.resumeCursor = 1 // "start fresh"
	got, cmd := m.updateResumeDraft(tea.KeyMsg{Type: tea.KeyEnter})
	next := got.(appModel)

	if next.draftDirToUse != "" {
		t.Errorf("draftDirToUse = %q, want empty so the starter is used", next.draftDirToUse)
	}
	if next.outcome != outcomeRunExercise || cmd == nil {
		t.Error("start fresh should still launch the session")
	}
	if _, ok, _ := draft.Load(dataDir, ex.ID); ok {
		t.Error("draft still loads after start-fresh, want it archived")
	}
	archived, err := filepath.Glob(filepath.Join(draft.Dir(dataDir, ex.ID), "previous.*"))
	if err != nil || len(archived) == 0 {
		t.Error("start fresh destroyed the draft instead of archiving it")
	}
}

func TestUpdateResumeDraft_EscBacksOutWithoutTouchingTheDraft(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	ex := problem.Variants[0].Exercise
	dataDir := seedDraft(t, ex.ID, "package main\n")

	m := draftModel(t, dataDir, problem)
	started, _ := m.launchExercise(ex)
	m = started.(appModel)

	got, cmd := m.updateResumeDraft(tea.KeyMsg{Type: tea.KeyEsc})
	next := got.(appModel)

	if next.stage != stageLanguage {
		t.Errorf("stage = %v, want back to stageLanguage", next.stage)
	}
	if next.outcome == outcomeRunExercise || cmd != nil {
		t.Error("Esc should not launch anything")
	}
	if _, ok, _ := draft.Load(dataDir, ex.ID); !ok {
		t.Error("Esc destroyed the draft")
	}
}

// TestProblemHasDraft_* and TestRenderProblems_DraftMarker_* cover
// issue #255: the picker should mark rows that have a saved draft, so
// the resume prompt (above) is never a surprise. problemHasDraft must
// check cheaply (draft.Exists -- directory/glob presence only, no file
// content reads) since renderProblems calls it once per visible row on
// every render.

func TestProblemHasDraft_TrueWhenAnyVariantHasADraft(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	ex := problem.Variants[0].Exercise
	dataDir := seedDraft(t, ex.ID, "package main\n")

	if !problemHasDraft(dataDir, problem) {
		t.Error("expected problemHasDraft to be true for a problem whose variant has a saved draft")
	}
}

func TestProblemHasDraft_FalseWhenNoVariantHasADraft(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	if problemHasDraft(t.TempDir(), problem) {
		t.Error("expected problemHasDraft to be false when nothing was ever snapshotted")
	}
}

func TestRenderProblems_DraftMarker_ShownWhenDraftExists(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	ex := problem.Variants[0].Exercise
	dataDir := seedDraft(t, ex.ID, "package main\n")

	m := appModel{stage: stageProblems, category: "two-pointers", cfg: config.Config{DataDir: dataDir}, categoryProblems: []catalog.ProblemStatus{problem}}
	out := stripAnsiTUI(m.renderProblems())
	if !strings.Contains(out, "· draft") {
		t.Errorf("renderProblems missing the draft marker for a problem with a saved draft, got:\n%s", out)
	}
}

func TestRenderProblems_DraftMarker_AbsentWhenNoDraft(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")

	m := appModel{stage: stageProblems, category: "two-pointers", cfg: config.Config{DataDir: t.TempDir()}, categoryProblems: []catalog.ProblemStatus{problem}}
	out := stripAnsiTUI(m.renderProblems())
	if strings.Contains(out, "· draft") {
		t.Errorf("renderProblems should not show a draft marker with no saved draft, got:\n%s", out)
	}
}

func TestRenderResumeDraft_ShowsAgeAndPreview(t *testing.T) {
	problem := codingProblem("two-sum-01", "Two Sum", "easy")
	ex := problem.Variants[0].Exercise
	dataDir := seedDraft(t, ex.ID, "package main\n\nfunc twoSum(nums []int) []int {\n\treturn nil\n}\n")

	m := draftModel(t, dataDir, problem)
	started, _ := m.launchExercise(ex)
	m = started.(appModel)

	out := stripAnsiTUI(m.renderResumeDraft())
	for _, want := range []string{"Two Sum", "func twoSum", "Resume", "Start fresh"} {
		if !strings.Contains(out, want) {
			t.Errorf("renderResumeDraft missing %q:\n%s", want, out)
		}
	}
}
