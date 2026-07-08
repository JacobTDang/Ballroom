package catalog

import (
	"strings"
	"testing"

	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func TestBanner_ContainsBallroomArtAndTagline(t *testing.T) {
	out := stripAnsi(Banner())
	if !strings.Contains(out, "I N T E R V I E W") || !strings.Contains(out, "P R E P") {
		t.Errorf("banner missing tagline:\n%s", out)
	}
	// Sanity check the ASCII art actually rendered multiple lines, not
	// just the tagline.
	if strings.Count(out, "\n") < 6 {
		t.Errorf("banner looks too short to contain the ASCII art:\n%s", out)
	}
}

func TestStyled_ProducesAnsiCodesByDefault(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	out := styled(colorTeal1, "hello")
	if out == "hello" {
		t.Error("expected styled() to wrap text in ANSI codes when NO_COLOR is unset")
	}
	if !strings.Contains(out, "hello") {
		t.Errorf("styled output lost the original text: %q", out)
	}
}

func TestStyled_RespectsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	out := styled(colorTeal1, "hello")
	if out != "hello" {
		t.Errorf("styled(%q) with NO_COLOR set = %q, want plain %q", "hello", out, "hello")
	}
}

func TestFormatTable_RespectsNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	statuses := []ExerciseStatus{
		{Exercise: fakeExercise("a", "pattern", "go", "A"), LastResult: tracker.ResultPass, Attempts: 1},
	}
	out := FormatTable(statuses)
	if strings.Contains(out, "\x1b[") {
		t.Errorf("expected no ANSI escape codes with NO_COLOR set:\n%q", out)
	}
}
