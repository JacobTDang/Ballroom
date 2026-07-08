package orchestrator

import (
	"reflect"
	"testing"
)

func TestParseImageIDs_Empty(t *testing.T) {
	got := parseImageIDs("")
	if len(got) != 0 {
		t.Errorf("parseImageIDs(\"\") = %v, want empty", got)
	}
}

func TestParseImageIDs_SingleID(t *testing.T) {
	got := parseImageIDs("sha256:abc123\n")
	want := []string{"sha256:abc123"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseImageIDs = %v, want %v", got, want)
	}
}

func TestParseImageIDs_MultipleIDs(t *testing.T) {
	got := parseImageIDs("abc123\ndef456\nghi789\n")
	want := []string{"abc123", "def456", "ghi789"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseImageIDs = %v, want %v", got, want)
	}
}

func TestParseImageIDs_WhitespaceOnly(t *testing.T) {
	got := parseImageIDs("   \n\n  \n")
	if len(got) != 0 {
		t.Errorf("parseImageIDs(whitespace only) = %v, want empty", got)
	}
}
