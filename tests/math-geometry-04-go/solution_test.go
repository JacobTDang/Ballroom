package main

import "testing"

func TestIsHappy_19(t *testing.T) {
	if got := IsHappy(19); got != true {
		t.Errorf("IsHappy(19) = %v, want true", got)
	}
}

func TestIsHappy_2(t *testing.T) {
	if got := IsHappy(2); got != false {
		t.Errorf("IsHappy(2) = %v, want false", got)
	}
}

func TestIsHappy_1(t *testing.T) {
	if got := IsHappy(1); got != true {
		t.Errorf("IsHappy(1) = %v, want true", got)
	}
}

func TestIsHappy_7(t *testing.T) {
	if got := IsHappy(7); got != true {
		t.Errorf("IsHappy(7) = %v, want true", got)
	}
}

func TestIsHappy_4(t *testing.T) {
	if got := IsHappy(4); got != false {
		t.Errorf("IsHappy(4) = %v, want false", got)
	}
}

func TestIsHappy_100(t *testing.T) {
	if got := IsHappy(100); got != true {
		t.Errorf("IsHappy(100) = %v, want true", got)
	}
}

func TestIsHappy_3(t *testing.T) {
	if got := IsHappy(3); got != false {
		t.Errorf("IsHappy(3) = %v, want false", got)
	}
}

func TestIsHappy_986(t *testing.T) {
	if got := IsHappy(986); got != false {
		t.Errorf("IsHappy(986) = %v, want false", got)
	}
}

func TestIsHappy_BoundaryLarge(t *testing.T) {
	if got := IsHappy(2147483647); got != false {
		t.Errorf("IsHappy(2147483647) = %v, want false", got)
	}
}
