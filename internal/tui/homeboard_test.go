package tui

import (
	"strings"
	"testing"
	"time"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

func TestPracticeStreak_Table(t *testing.T) {
	today := time.Date(2026, 7, 16, 15, 0, 0, 0, time.UTC)
	day := func(offset int) string { return today.AddDate(0, 0, offset).Format("2006-01-02") }
	attemptsOn := func(days ...string) []tracker.Attempt {
		var out []tracker.Attempt
		for _, d := range days {
			out = append(out, tracker.Attempt{Date: d})
		}
		return out
	}

	cases := []struct {
		name     string
		attempts []tracker.Attempt
		want     int
	}{
		{"no attempts, no streak", nil, 0},
		{"one attempt today", attemptsOn(day(0)), 1},
		{"three consecutive days ending today", attemptsOn(day(-2), day(-1), day(0)), 3},
		{"streak alive if yesterday practiced but today not yet", attemptsOn(day(-2), day(-1)), 2},
		{"gap breaks the streak", attemptsOn(day(-3), day(-1), day(0)), 2},
		{"last practice two days ago is a dead streak", attemptsOn(day(-4), day(-3), day(-2)), 0},
		{"several attempts one day count once", attemptsOn(day(0), day(0), day(0)), 1},
		{"unparseable dates are ignored", attemptsOn("garbage", day(0)), 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := practiceStreak(c.attempts, today); got != c.want {
				t.Errorf("practiceStreak = %d, want %d", got, c.want)
			}
		})
	}
}

func TestProgressBar_Rendering(t *testing.T) {
	if got := stripAnsiTUI(progressBar(0, 10, 10)); got != "░░░░░░░░░░" {
		t.Errorf("empty bar = %q", got)
	}
	if got := stripAnsiTUI(progressBar(10, 10, 10)); got != "▓▓▓▓▓▓▓▓▓▓" {
		t.Errorf("full bar = %q", got)
	}
	if got := stripAnsiTUI(progressBar(5, 10, 10)); got != "▓▓▓▓▓░░░░░" {
		t.Errorf("half bar = %q", got)
	}
	// One solve out of many must still show one visible cell of
	// progress -- flooring 1/149 to zero would read as "nothing done".
	if got := stripAnsiTUI(progressBar(1, 149, 14)); !strings.HasPrefix(got, "▓") {
		t.Errorf("tiny progress rendered invisible: %q", got)
	}
	if got := stripAnsiTUI(progressBar(3, 0, 10)); got != "░░░░░░░░░░" {
		t.Errorf("zero-total bar must not divide by zero: %q", got)
	}
}

func homeboardFixtureProblems() []catalog.ProblemStatus {
	return []catalog.ProblemStatus{
		{ProblemID: "two-pointers-01", Title: "Two Sum II", Category: "two-pointers", Solved: true,
			Variants: []catalog.ExerciseStatus{{Exercise: exercise.Exercise{Language: "go"}, Attempts: 1, LastResult: tracker.ResultPass, LastAttemptDate: "2026-07-15"}}},
		{ProblemID: "url-shortener-01", Title: "Design Pastebin / Bit.ly", Category: exercise.CategorySystemDesign,
			Variants: []catalog.ExerciseStatus{
				{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}, Attempts: 1, LastResult: tracker.ResultPass, LastAttemptDate: "2026-07-15"},
				{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageInterviewer}},
			}},
		{ProblemID: "disagreement-01", Title: "Disagreement With a Teammate", Category: exercise.CategoryBehavioral,
			Variants: []catalog.ExerciseStatus{{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}}}},
		{ProblemID: "library-api-01", Title: "Design the Library API", Category: exercise.CategoryAPIDesign,
			Variants: []catalog.ExerciseStatus{{Exercise: exercise.Exercise{Kind: exercise.KindDesign, Language: exercise.LanguageCoach}}}},
		{ProblemID: "parking-lot-01", Title: "Design a Parking Lot", Category: exercise.CategoryOODesign,
			Variants: []catalog.ExerciseStatus{{Exercise: exercise.Exercise{Language: "python"}, Attempts: 1, LastResult: tracker.ResultFail, LastAttemptDate: "2020-01-01"}}},
		{ProblemID: "bounded-queue-01", Title: "Bounded Producer-Consumer Queue", Category: exercise.CategoryConcurrency,
			Variants: []catalog.ExerciseStatus{{Exercise: exercise.Exercise{Language: "go"}}}},
		{ProblemID: "bloom-filter-01", Title: "Bloom Filter", Category: exercise.CategoryImplementation,
			Variants: []catalog.ExerciseStatus{{Exercise: exercise.Exercise{Language: "go"}}}},
		{ProblemID: "off-by-one-01", Title: "Off-by-one IndexError in max_of", Category: exercise.CategoryDebug,
			Variants: []catalog.ExerciseStatus{{Exercise: exercise.Exercise{Language: "python"}}}},
	}
}

func TestRenderHomeboard_ShowsTracksDueDailyStreakRecent(t *testing.T) {
	attempts := []tracker.Attempt{
		{ExerciseID: "parking-lot-01-python", Result: tracker.ResultFail, Date: "2020-01-01"},
		{ExerciseID: "two-pointers-01-go", Result: tracker.ResultPass, Date: time.Now().Format("2006-01-02")},
	}

	got := stripAnsiTUI(renderHomeboard(homeboardFixtureProblems(), attempts, time.Now()))

	for _, want := range []string{"DSA", "Debug", "Concurrency", "Implementation", "System Design", "API Design", "OO Design", "Behavioral"} {
		if !strings.Contains(got, want) {
			t.Errorf("homeboard missing the %s track:\n%s", want, got)
		}
	}
	if !strings.Contains(got, "1/1") {
		t.Errorf("homeboard missing solved counts:\n%s", got)
	}
	// url-shortener is mock due; parking-lot's 2020 failure is review due.
	if !strings.Contains(got, "1 mock") || !strings.Contains(got, "1 review") {
		t.Errorf("homeboard missing due counts:\n%s", got)
	}
	if !strings.Contains(got, "streak 1 day") {
		t.Errorf("homeboard missing the streak (one attempt today):\n%s", got)
	}
	if !strings.Contains(got, "today:") {
		t.Errorf("homeboard missing the Daily preview:\n%s", got)
	}
	if !strings.Contains(got, "two-pointers-01-go") {
		t.Errorf("homeboard missing recent attempts:\n%s", got)
	}
}

// TestRenderHomeboard_GaveUpGetsItsOwnThirdGlyph covers issue #238: the
// recent-activity row must distinguish "gave up" from both a genuine
// pass and a genuine fail, not silently collapse into the fail glyph.
func TestRenderHomeboard_GaveUpGetsItsOwnThirdGlyph(t *testing.T) {
	attempts := []tracker.Attempt{
		{ExerciseID: "two-pointers-01-go", Result: tracker.ResultPass, Date: "2026-07-16"},
		{ExerciseID: "off-by-one-01-python", Result: tracker.ResultFail, Date: "2026-07-17"},
		{ExerciseID: "bloom-filter-01-go", Result: tracker.ResultGaveUp, Date: "2026-07-18"},
	}

	got := stripAnsiTUI(renderHomeboard(homeboardFixtureProblems(), attempts, time.Now()))

	if !strings.Contains(got, "~ bloom-filter-01-go") {
		t.Errorf("homeboard missing a distinct gave-up glyph for bloom-filter-01-go:\n%s", got)
	}
	if !strings.Contains(got, "✗ off-by-one-01-python") {
		t.Errorf("homeboard missing the fail glyph for off-by-one-01-python:\n%s", got)
	}
	if !strings.Contains(got, "✓ two-pointers-01-go") {
		t.Errorf("homeboard missing the pass glyph for two-pointers-01-go:\n%s", got)
	}
}

func TestRenderHomeboard_EmptyDataDegradesQuietly(t *testing.T) {
	got := stripAnsiTUI(renderHomeboard(nil, nil, time.Now()))
	if strings.TrimSpace(got) != "" {
		t.Errorf("homeboard with no catalog should render nothing, got:\n%s", got)
	}
}

func TestRenderMain_IncludesHomeboardAndDropsInlineHint(t *testing.T) {
	m := appModel{
		problems:     homeboardFixtureProblems(),
		homeAttempts: []tracker.Attempt{{ExerciseID: "two-pointers-01-go", Result: tracker.ResultPass, Date: time.Now().Format("2006-01-02")}},
	}
	got := stripAnsiTUI(m.renderMain())
	if !strings.Contains(got, "System Design") || !strings.Contains(got, "streak") {
		t.Errorf("renderMain missing the homeboard:\n%s", got)
	}
	if strings.Contains(got, "enter select") {
		t.Errorf("renderMain still carries the inline hint line -- it moved to the panel footer:\n%s", got)
	}
}

func TestView_MainStageHasFooterPinnedAtPanelBottom(t *testing.T) {
	m := appModel{width: 120, height: 40}
	lines := strings.Split(m.View(), "\n")
	// The footer's key hint should sit on the last few rows, above the
	// bottom border.
	tail := stripAnsiTUI(strings.Join(lines[len(lines)-4:], "\n"))
	if !strings.Contains(tail, "enter select") {
		t.Errorf("footer keys not pinned at the panel bottom:\n%s", tail)
	}
}
