package main

import "testing"

func TestIsNStraightHand_Classic(t *testing.T) {
	hand := []int{1, 2, 3, 6, 2, 3, 4, 7, 8}
	if got := IsNStraightHand(hand, 3); got != true {
		t.Errorf("IsNStraightHand(%v, 3) = %v, want true", hand, got)
	}
}

func TestIsNStraightHand_NotDivisible(t *testing.T) {
	hand := []int{1, 2, 3, 4, 5}
	if got := IsNStraightHand(hand, 4); got != false {
		t.Errorf("IsNStraightHand(%v, 4) = %v, want false", hand, got)
	}
}

func TestIsNStraightHand_MissingCard(t *testing.T) {
	hand := []int{1, 2, 3, 4, 5, 7}
	if got := IsNStraightHand(hand, 3); got != false {
		t.Errorf("IsNStraightHand(%v, 3) = %v, want false", hand, got)
	}
}

func TestIsNStraightHand_GroupSizeOne(t *testing.T) {
	hand := []int{5, 5, 5}
	if got := IsNStraightHand(hand, 1); got != true {
		t.Errorf("IsNStraightHand(%v, 1) = %v, want true", hand, got)
	}
}

func TestIsNStraightHand_ExactSingleGroup(t *testing.T) {
	hand := []int{1, 2, 3}
	if got := IsNStraightHand(hand, 3); got != true {
		t.Errorf("IsNStraightHand(%v, 3) = %v, want true", hand, got)
	}
}
