package main

import "testing"

func TestClimbStairs_One(t *testing.T) {
	if got := ClimbStairs(1); got != 1 {
		t.Errorf("ClimbStairs(1) = %d, want 1", got)
	}
}

func TestClimbStairs_Two(t *testing.T) {
	if got := ClimbStairs(2); got != 2 {
		t.Errorf("ClimbStairs(2) = %d, want 2", got)
	}
}

func TestClimbStairs_Three(t *testing.T) {
	if got := ClimbStairs(3); got != 3 {
		t.Errorf("ClimbStairs(3) = %d, want 3", got)
	}
}

func TestClimbStairs_Five(t *testing.T) {
	if got := ClimbStairs(5); got != 8 {
		t.Errorf("ClimbStairs(5) = %d, want 8", got)
	}
}

func TestClimbStairs_Four(t *testing.T) {
	if got := ClimbStairs(4); got != 5 {
		t.Errorf("ClimbStairs(4) = %d, want 5", got)
	}
}

func TestClimbStairs_Ten(t *testing.T) {
	if got := ClimbStairs(10); got != 89 {
		t.Errorf("ClimbStairs(10) = %d, want 89", got)
	}
}

func TestClimbStairs_BoundaryMax(t *testing.T) {
	if got := ClimbStairs(45); got != 1836311903 {
		t.Errorf("ClimbStairs(45) = %d, want 1836311903", got)
	}
}
