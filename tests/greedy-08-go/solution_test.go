package main

import "testing"

func TestCheckValidString_Simple(t *testing.T) {
	s := "()"
	if got := CheckValidString(s); got != true {
		t.Errorf("CheckValidString(%q) = %v, want true", s, got)
	}
}

func TestCheckValidString_StarBalances(t *testing.T) {
	s := "(*))"
	if got := CheckValidString(s); got != true {
		t.Errorf("CheckValidString(%q) = %v, want true", s, got)
	}
}

func TestCheckValidString_Unbalanced(t *testing.T) {
	s := "(()"
	if got := CheckValidString(s); got != false {
		t.Errorf("CheckValidString(%q) = %v, want false", s, got)
	}
}

func TestCheckValidString_AllStars(t *testing.T) {
	s := "***"
	if got := CheckValidString(s); got != true {
		t.Errorf("CheckValidString(%q) = %v, want true", s, got)
	}
}

func TestCheckValidString_SingleClose(t *testing.T) {
	s := ")"
	if got := CheckValidString(s); got != false {
		t.Errorf("CheckValidString(%q) = %v, want false", s, got)
	}
}
