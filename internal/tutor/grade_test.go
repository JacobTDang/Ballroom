package tutor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// gradeWorkDir builds a design-session workspace mid-submit: the
// user's solution.md plus the rubric.md the reveal just delivered.
func gradeWorkDir(t *testing.T, withRubric bool) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "solution.md"), []byte("# My design\n\nshard by user id"), 0o644); err != nil {
		t.Fatalf("write solution: %v", err)
	}
	if withRubric {
		if err := os.WriteFile(filepath.Join(dir, "rubric.md"), []byte("- estimates\n- sharding"), 0o644); err != nil {
			t.Fatalf("write rubric: %v", err)
		}
	}
	return dir
}

func TestGradeDesign_ParsesPassVerdict(t *testing.T) {
	mock := newSequencedOllama(t, "VERDICT: pass\nEstimates: strong. Sharding: adequate.")
	cfg := testConfig(mock.URL)
	cfg.WorkDir = gradeWorkDir(t, true)

	verdict, summary, err := GradeDesign(context.Background(), cfg)
	if err != nil {
		t.Fatalf("GradeDesign: %v", err)
	}
	if verdict != "pass" {
		t.Errorf("verdict = %q, want pass", verdict)
	}
	if !strings.Contains(summary, "Estimates: strong") {
		t.Errorf("summary = %q, want the per-dimension text", summary)
	}
}

func TestGradeDesign_ParsesFailVerdict(t *testing.T) {
	mock := newSequencedOllama(t, "VERDICT: fail\nNo estimates at all.")
	cfg := testConfig(mock.URL)
	cfg.WorkDir = gradeWorkDir(t, true)

	verdict, _, err := GradeDesign(context.Background(), cfg)
	if err != nil {
		t.Fatalf("GradeDesign: %v", err)
	}
	if verdict != "fail" {
		t.Errorf("verdict = %q, want fail", verdict)
	}
}

func TestGradeDesign_ToleratesCaseAndLeadingProse(t *testing.T) {
	// Small models rarely obey format instructions perfectly -- accept
	// the verdict line anywhere in the reply, any case.
	mock := newSequencedOllama(t, "Here is my assessment.\nverdict: PASS\nGood work overall.")
	cfg := testConfig(mock.URL)
	cfg.WorkDir = gradeWorkDir(t, true)

	verdict, _, err := GradeDesign(context.Background(), cfg)
	if err != nil {
		t.Fatalf("GradeDesign: %v", err)
	}
	if verdict != "pass" {
		t.Errorf("verdict = %q, want pass", verdict)
	}
}

func TestGradeDesign_UnparseableReplyIsAnError(t *testing.T) {
	// No silent default: a reply without a verdict must error so the
	// submit flow falls back to explicit self-assessment instead of
	// recording a guess.
	mock := newSequencedOllama(t, "The design looks quite good to me overall!")
	cfg := testConfig(mock.URL)
	cfg.WorkDir = gradeWorkDir(t, true)

	if _, _, err := GradeDesign(context.Background(), cfg); err == nil {
		t.Fatal("GradeDesign = nil error for a reply with no verdict, want an error")
	}
}

func TestGradeDesign_MissingRubricIsAnError(t *testing.T) {
	mock := newSequencedOllama(t, "VERDICT: pass")
	cfg := testConfig(mock.URL)
	cfg.WorkDir = gradeWorkDir(t, false)

	if _, _, err := GradeDesign(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "rubric") {
		t.Fatalf("GradeDesign err = %v, want a rubric-missing error", err)
	}
}

func TestGradeDesign_RequestCarriesRubricAndSolution(t *testing.T) {
	mock := newSequencedOllama(t, "VERDICT: pass\nok")
	cfg := testConfig(mock.URL)
	cfg.WorkDir = gradeWorkDir(t, true)

	if _, _, err := GradeDesign(context.Background(), cfg); err != nil {
		t.Fatalf("GradeDesign: %v", err)
	}
	var all strings.Builder
	for _, msg := range mock.allRequests()[0].Messages {
		all.WriteString(msg.Content)
		all.WriteString("\n")
	}
	for _, want := range []string{"shard by user id", "- estimates"} {
		if !strings.Contains(all.String(), want) {
			t.Errorf("grading request missing %q -- the model can't grade what it can't see", want)
		}
	}
}
