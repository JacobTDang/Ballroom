package mock

import (
	"os"
	"path/filepath"
	"testing"
)

func sitting(started string, solved int, completed bool) Sitting {
	s := Sitting{
		StartedAt:  started,
		ProblemIDs: [4]string{"c1-warmup-01", "c1-warmup-02", "c1-grid-01", "c1-algo-01"},
		Completed:  completed,
	}
	for i := range s.Outcomes {
		if i < solved {
			s.Outcomes[i] = OutcomeSolved
		} else {
			s.Outcomes[i] = OutcomeSkipped
		}
	}
	return s
}

func TestAppendAndListSittings(t *testing.T) {
	dir := t.TempDir()
	if err := AppendSitting(dir, sitting("2026-08-16T10:00:00Z", 3, true)); err != nil {
		t.Fatalf("append: %v", err)
	}
	if err := AppendSitting(dir, sitting("2026-08-17T10:00:00Z", 2, true)); err != nil {
		t.Fatalf("append: %v", err)
	}
	got, err := ListSittings(dir)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 || got[0].Solved() != 3 || got[1].Solved() != 2 {
		t.Fatalf("got %d sittings, solved %v", len(got), got)
	}
}

func TestListSittings_MissingFileIsEmpty(t *testing.T) {
	got, err := ListSittings(t.TempDir())
	if err != nil || got != nil {
		t.Fatalf("want (nil, nil), got (%v, %v)", got, err)
	}
}

func TestListSittings_CorruptLineErrors(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "mocks.jsonl"), []byte("{not json\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ListSittings(dir); err == nil {
		t.Fatal("corrupt line must error, not be skipped silently")
	}
}

func TestHistoryLine(t *testing.T) {
	if got := HistoryLine(nil); got != "" {
		t.Errorf("empty history = %q, want empty", got)
	}
	line := HistoryLine([]Sitting{
		sitting("2026-08-15T10:00:00Z", 3, true),
		sitting("2026-08-16T10:00:00Z", 2, true),
	})
	want := "2 sittings · best 3/4 · last 2/4"
	if line != want {
		t.Errorf("HistoryLine = %q, want %q", line, want)
	}
}
