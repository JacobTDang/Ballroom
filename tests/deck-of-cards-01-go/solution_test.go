package main

import "testing"

func TestNewDeckHas52Cards(t *testing.T) {
	if got := NewDeck().Size(); got != 52 {
		t.Errorf("Size() = %d, want 52", got)
	}
}

func TestFirstCardIsTheAceOfClubs(t *testing.T) {
	got := NewDeck().Deal(1)
	if len(got) != 1 || got[0] != (Card{"clubs", 1}) {
		t.Errorf("Deal(1) = %v, want [{clubs 1}]", got)
	}
}

func TestDealFollowsCanonicalOrder(t *testing.T) {
	d := NewDeck()
	d.Deal(1)
	got := d.Deal(3)
	want := []Card{{"clubs", 2}, {"clubs", 3}, {"clubs", 4}}
	if len(got) != 3 || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Errorf("Deal(3) = %v, want %v", got, want)
	}
}

func TestSuitsChangeAtTheRightBoundary(t *testing.T) {
	d := NewDeck()
	d.Deal(13)
	got := d.Deal(1)
	if len(got) != 1 || got[0] != (Card{"diamonds", 1}) {
		t.Errorf("14th card = %v, want {diamonds 1}", got)
	}
}

func TestDealReducesSize(t *testing.T) {
	d := NewDeck()
	d.Deal(5)
	if got := d.Size(); got != 47 {
		t.Errorf("Size() = %d after dealing 5, want 47", got)
	}
}

func TestAll52CardsAreUnique(t *testing.T) {
	seen := map[Card]bool{}
	for _, c := range NewDeck().Deal(52) {
		if seen[c] {
			t.Fatalf("card %v dealt twice", c)
		}
		seen[c] = true
	}
	if len(seen) != 52 {
		t.Errorf("dealt %d unique cards, want 52", len(seen))
	}
}

func TestDealPastEmptyReturnsRemainder(t *testing.T) {
	d := NewDeck()
	d.Deal(50)
	if got := d.Deal(5); len(got) != 2 {
		t.Errorf("Deal(5) with 2 left returned %d cards, want 2", len(got))
	}
	if got := d.Size(); got != 0 {
		t.Errorf("Size() = %d, want 0", got)
	}
	if got := d.Deal(1); len(got) != 0 {
		t.Errorf("Deal(1) on empty deck returned %d cards, want 0", len(got))
	}
}

func TestResetRestoresCanonicalDeck(t *testing.T) {
	d := NewDeck()
	d.Deal(30)
	d.Reset()
	if got := d.Size(); got != 52 {
		t.Errorf("Size() = %d after Reset, want 52", got)
	}
	if got := d.Deal(1); len(got) != 1 || got[0] != (Card{"clubs", 1}) {
		t.Errorf("first card after Reset = %v, want {clubs 1}", got)
	}
}

func TestShuffleKeepsExactlyTheRemainingCards(t *testing.T) {
	d := NewDeck()
	dealt := d.Deal(10)
	d.Shuffle()
	rest := d.Deal(52)
	if len(rest) != 42 {
		t.Fatalf("dealt %d cards after shuffle, want 42", len(rest))
	}
	seen := map[Card]bool{}
	for _, c := range dealt {
		seen[c] = true
	}
	for _, c := range rest {
		if seen[c] {
			t.Fatalf("card %v appeared twice across deal+shuffle", c)
		}
		seen[c] = true
	}
	if len(seen) != 52 {
		t.Errorf("deal+shuffled rest covered %d unique cards, want 52", len(seen))
	}
}

func TestShuffleDoesNotChangeSize(t *testing.T) {
	d := NewDeck()
	d.Deal(7)
	d.Shuffle()
	if got := d.Size(); got != 45 {
		t.Errorf("Size() = %d after shuffle, want 45", got)
	}
}
