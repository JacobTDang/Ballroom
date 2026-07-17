package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestFullDocument(t *testing.T) {
	input := `top = level
# a comment
; another comment

[server]
host = localhost
port = 8080
host = example.com

[client]
retries=3
`
	got, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	want := map[string]map[string]string{
		"": {"top": "level"},
		// later host wins
		"server": {"host": "example.com", "port": "8080"},
		"client": {"retries": "3"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse = %v, want %v", got, want)
	}
}

func TestWhitespaceTrimmedEverywhere(t *testing.T) {
	got, err := Parse("  spaced key   =   spaced value  ")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got[""]["spaced key"] != "spaced value" {
		t.Fatalf("got %v, want trimmed key and value", got)
	}
}

func TestMalformedLineErrorsWithLineNumber(t *testing.T) {
	_, err := Parse("ok = 1\nnot a valid line\nok2 = 2")
	if err == nil {
		t.Fatal("Parse accepted a line that is not a header, comment, or key=value")
	}
	if !strings.Contains(err.Error(), "2") {
		t.Fatalf("error %q should name line 2", err)
	}
}

func TestUnclosedSectionErrorsWithLineNumber(t *testing.T) {
	_, err := Parse("[server]\nkey = v\n[broken")
	if err == nil {
		t.Fatal("Parse accepted an unclosed section header")
	}
	if !strings.Contains(err.Error(), "3") {
		t.Fatalf("error %q should name line 3", err)
	}
}
