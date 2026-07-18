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

func TestImageNeedsRebuild(t *testing.T) {
	cases := []struct {
		name  string
		label string
		want  string
		out   bool
	}{
		{"label matches current content hash", "abc123", "abc123", false},
		{"label is a stale content hash", "abc123", "def456", true},
		{"label empty (image predates this label)", "", "abc123", true},
		{"both empty", "", "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := imageNeedsRebuild(tc.label, tc.want)
			if got != tc.out {
				t.Errorf("imageNeedsRebuild(%q, %q) = %v, want %v", tc.label, tc.want, got, tc.out)
			}
		})
	}
}
