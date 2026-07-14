package main

import "testing"

func TestMultiplyStrings_SingleDigits(t *testing.T) {
	if got := MultiplyStrings("2", "3"); got != "6" {
		t.Errorf(`MultiplyStrings("2", "3") = %q, want "6"`, got)
	}
}

func TestMultiplyStrings_MultiDigits(t *testing.T) {
	if got := MultiplyStrings("123", "456"); got != "56088" {
		t.Errorf(`MultiplyStrings("123", "456") = %q, want "56088"`, got)
	}
}

func TestMultiplyStrings_ZeroFirst(t *testing.T) {
	if got := MultiplyStrings("0", "12345"); got != "0" {
		t.Errorf(`MultiplyStrings("0", "12345") = %q, want "0"`, got)
	}
}

func TestMultiplyStrings_AllNines(t *testing.T) {
	if got := MultiplyStrings("999", "999"); got != "998001" {
		t.Errorf(`MultiplyStrings("999", "999") = %q, want "998001"`, got)
	}
}

func TestMultiplyStrings_OneByOne(t *testing.T) {
	if got := MultiplyStrings("1", "1"); got != "1" {
		t.Errorf(`MultiplyStrings("1", "1") = %q, want "1"`, got)
	}
}

func TestMultiplyStrings_BothZero(t *testing.T) {
	if got := MultiplyStrings("0", "0"); got != "0" {
		t.Errorf(`MultiplyStrings("0", "0") = %q, want "0"`, got)
	}
}

func TestMultiplyStrings_IdentityByLarger(t *testing.T) {
	if got := MultiplyStrings("1", "999"); got != "999" {
		t.Errorf(`MultiplyStrings("1", "999") = %q, want "999"`, got)
	}
}

func TestMultiplyStrings_LargerMultiDigit(t *testing.T) {
	if got := MultiplyStrings("12345", "67890"); got != "838102050" {
		t.Errorf(`MultiplyStrings("12345", "67890") = %q, want "838102050"`, got)
	}
}

func TestMultiplyStrings_SingleDigitByMultiDigit(t *testing.T) {
	if got := MultiplyStrings("9", "123"); got != "1107" {
		t.Errorf(`MultiplyStrings("9", "123") = %q, want "1107"`, got)
	}
}
