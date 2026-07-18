package main

import "testing"

func gridEqual(got, want [][]int) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if len(got[i]) != len(want[i]) {
			return false
		}
		for j := range got[i] {
			if got[i][j] != want[i][j] {
				return false
			}
		}
	}
	return true
}

func TestSnapshotIndependentOfLaterMutation(t *testing.T) {
	g := NewGrid(3, 2)
	g.Set(0, 0, 1)
	g.Set(1, 1, 2)
	g.Set(2, 0, 3)

	snap := g.Snapshot()
	want := [][]int{{1, 0}, {0, 2}, {3, 0}}
	if !gridEqual(snap, want) {
		t.Fatalf("snapshot right after taking it = %v, want %v", snap, want)
	}

	g.Set(0, 0, 111)
	g.Set(2, 1, 999)

	if !gridEqual(snap, want) {
		t.Errorf("snapshot changed after live grid was edited: got %v, want %v", snap, want)
	}
	if g.Get(0, 0) != 111 || g.Get(2, 1) != 999 {
		t.Errorf("live grid did not reflect its own edits: (0,0)=%d (2,1)=%d", g.Get(0, 0), g.Get(2, 1))
	}
}

func TestMultipleSnapshotsAreIndependentOfEachOther(t *testing.T) {
	g := NewGrid(3, 2)
	g.Set(0, 0, 1)
	g.Set(1, 1, 2)
	g.Set(2, 0, 3)

	snap1 := g.Snapshot()

	g.Set(0, 0, 111)
	g.Set(2, 1, 999)

	snap2 := g.Snapshot()

	g.Set(0, 0, 222)

	want1 := [][]int{{1, 0}, {0, 2}, {3, 0}}
	want2 := [][]int{{111, 0}, {0, 2}, {3, 999}}
	if !gridEqual(snap1, want1) {
		t.Errorf("snap1 = %v, want %v", snap1, want1)
	}
	if !gridEqual(snap2, want2) {
		t.Errorf("snap2 = %v, want %v", snap2, want2)
	}
	if g.Get(0, 0) != 222 {
		t.Errorf("live grid (0,0) = %d, want 222", g.Get(0, 0))
	}
}

func TestSnapshotCoversEveryRowNotJustTheFirst(t *testing.T) {
	g := NewGrid(3, 2)
	for r := 0; r < 3; r++ {
		for c := 0; c < 2; c++ {
			g.Set(r, c, r*10+c)
		}
	}

	snap := g.Snapshot()
	g.Set(2, 1, -1) // mutate the LAST row after the snapshot

	if snap[2][1] != 21 {
		t.Errorf("snap[2][1] = %d, want 21 (unaffected by later edit)", snap[2][1])
	}
	if g.Get(2, 1) != -1 {
		t.Errorf("g.Get(2,1) = %d, want -1", g.Get(2, 1))
	}
}
