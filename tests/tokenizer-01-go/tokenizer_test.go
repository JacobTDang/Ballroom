package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestDenseExpressionExactTokens(t *testing.T) {
	got, err := Tokenize("3+4.5*x")
	if err != nil {
		t.Fatalf("Tokenize: %v", err)
	}
	want := []Token{
		{"number", "3", 0},
		{"op", "+", 1},
		{"number", "4.5", 2},
		{"op", "*", 5},
		{"ident", "x", 6},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tokens = %+v, want %+v", got, want)
	}
}

func TestParensAndIdentifiers(t *testing.T) {
	got, err := Tokenize("price * (1 + tax_rate2)")
	if err != nil {
		t.Fatalf("Tokenize: %v", err)
	}
	want := []Token{
		{"ident", "price", 0},
		{"op", "*", 6},
		{"lparen", "(", 8},
		{"number", "1", 9},
		{"op", "+", 11},
		{"ident", "tax_rate2", 13},
		{"rparen", ")", 22},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("tokens = %+v, want %+v", got, want)
	}
}

func TestSecondDecimalPointIsAnError(t *testing.T) {
	_, err := Tokenize("12..3")
	if err == nil {
		t.Fatal("Tokenize(12..3) succeeded, want an error for the second dot")
	}
	if !strings.Contains(err.Error(), "3") {
		t.Fatalf("error %q should name position 3", err)
	}
}

func TestUnknownCharacterErrorNamesPosition(t *testing.T) {
	_, err := Tokenize("a @ b")
	if err == nil {
		t.Fatal("Tokenize with @ succeeded, want an error")
	}
	if !strings.Contains(err.Error(), "2") {
		t.Fatalf("error %q should name position 2", err)
	}
}

func TestEmptyInputIsEmptyTokenList(t *testing.T) {
	got, err := Tokenize("")
	if err != nil || len(got) != 0 {
		t.Fatalf("Tokenize(\"\") = %v, %v; want empty, nil", got, err)
	}
}
