package catalog

import (
	"strings"
	"testing"
)

func TestStyled_ProducesAnsiCodesByDefault(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	out := styled(colorTeal, "hello")
	if out == "hello" {
		t.Error("expected styled() to wrap text in ANSI codes when NO_COLOR is unset")
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("styled output lost the original text: %q", out)
	}
}

func TestStyled_RespectsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	out := styled(colorTeal, "hello")
	if out != "hello" {
		t.Errorf("styled(%q) with NO_COLOR set = %q, want plain %q", "hello", out, "hello")
	}
}

func TestFormatSummary_RespectsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	statuses := []ExerciseStatus{
		{Exercise: fakeExercise("a", "pattern", "go", "A"), LastResult: "pass", Attempts: 1},
	}
	out := FormatSummary(statuses)
	if strings.Contains(out, "\x1b[") {
		t.Errorf("expected no ANSI escape codes with NO_COLOR set:\n%q", out)
	}
}
