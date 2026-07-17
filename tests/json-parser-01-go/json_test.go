package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestNestedDocumentExactStructure(t *testing.T) {
	input := `{"name": "ada", "age": -3, "tags": ["a", "b"], "meta": {"ok": true, "note": null}, "empty": [], "eo": {}}`
	got, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	want := map[string]any{
		"name": "ada",
		"age":  -3,
		"tags": []any{"a", "b"},
		"meta": map[string]any{"ok": true, "note": nil},
		"empty": []any{},
		"eo":    map[string]any{},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse = %#v, want %#v", got, want)
	}
}

func TestEscapes(t *testing.T) {
	got, err := Parse(`"say \"hi\" and \\"`)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if got != `say "hi" and \` {
		t.Fatalf("Parse = %q, want the unescaped string", got)
	}
}

func TestWhitespaceEverywhere(t *testing.T) {
	got, err := Parse("  { \"a\" :  [ 1 , 2 ]  }  ")
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	want := map[string]any{"a": []any{1, 2}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse = %#v, want %#v", got, want)
	}
}

func TestErrorsNamePositions(t *testing.T) {
	cases := []struct {
		input   string
		wantPos string
	}{
		{`{"a" 1}`, "5"},        // missing colon
		{`"unterminated`, "0"},  // unterminated string
		{`tru`, "0"},            // bad literal
		{`{"a": 1} extra`, "9"}, // trailing garbage
		{`"bad \x escape"`, "5"}, // unsupported escape
	}
	for _, c := range cases {
		_, err := Parse(c.input)
		if err == nil {
			t.Fatalf("Parse(%q) succeeded, want an error", c.input)
		}
		if !strings.Contains(err.Error(), c.wantPos) {
			t.Fatalf("Parse(%q) error %q should name position %s", c.input, err, c.wantPos)
		}
	}
}
