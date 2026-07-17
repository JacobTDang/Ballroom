package tutor

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiffUnified_Table(t *testing.T) {
	cases := []struct {
		name       string
		old, new   string
		wantParts  []string
		wantAbsent []string
	}{
		{
			"insertion",
			"a\nb\nc",
			"a\nb\nx\nc",
			[]string{"+x", " b", " c"},
			[]string{"+a", "-a"},
		},
		{
			"deletion",
			"a\nb\nc",
			"a\nc",
			[]string{"-b"},
			[]string{"+b"},
		},
		{
			"modification is a remove plus an add",
			"count := 0\nreturn count",
			"count := 1\nreturn count",
			[]string{"-count := 0", "+count := 1", " return count"},
			nil,
		},
		{
			"from empty shows all adds",
			"",
			"a\nb",
			[]string{"+a", "+b"},
			[]string{"-"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := diffUnified(c.old, c.new)
			for _, want := range c.wantParts {
				if !strings.Contains(got, want+"\n") && !strings.HasSuffix(got, want) {
					t.Errorf("diff missing %q:\n%s", want, got)
				}
			}
			for _, absent := range c.wantAbsent {
				for _, line := range strings.Split(got, "\n") {
					if line == absent {
						t.Errorf("diff has unexpected line %q:\n%s", absent, got)
					}
				}
			}
			if !strings.Contains(got, "@@") {
				t.Errorf("diff has no hunk header:\n%s", got)
			}
		})
	}
}

func TestDiffUnified_IdenticalInputsAreEmpty(t *testing.T) {
	if got := diffUnified("a\nb", "a\nb"); got != "" {
		t.Errorf("diff of identical content = %q, want empty", got)
	}
}

// TestReadSolutionDiffTool_Lifecycle drives the real tool: the session
// baseline is taken when the tools are built, each diff (or full read)
// advances the snapshot, and no-change reads say so.
func TestReadSolutionDiffTool_Lifecycle(t *testing.T) {
	dir := t.TempDir()
	write := func(content string) {
		if err := os.WriteFile(filepath.Join(dir, "solution.go"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("package main\n\nfunc solve() int {\n\treturn 0\n}\n")

	cfg := Config{WorkDir: dir, MaxContextBytes: 8000}
	tools, err := BuildTools(cfg)
	if err != nil {
		t.Fatalf("BuildTools: %v", err)
	}
	ctx := context.Background()
	diffTool := findTool(ctx, tools, "read_solution_diff")
	if diffTool == nil {
		t.Fatal("read_solution_diff not in the tool set")
	}
	readTool := findTool(ctx, tools, "read_solution_file")

	// InvokableRun returns the JSON-marshaled readFileOutput -- decode
	// to assert on the real content, not its escaped form.
	runDiff := func() string {
		t.Helper()
		out, err := diffTool.InvokableRun(ctx, "{}")
		if err != nil {
			t.Fatalf("read_solution_diff: %v", err)
		}
		var parsed readFileOutput
		if err := json.Unmarshal([]byte(out), &parsed); err != nil {
			t.Fatalf("unmarshal tool output %q: %v", out, err)
		}
		return parsed.Content
	}

	// No edits yet: nothing has changed since the session-start baseline.
	if got := runDiff(); !strings.Contains(got, "no changes") {
		t.Errorf("pre-edit diff = %q, want a no-changes note", got)
	}

	// Edit, then diff: only the change shows.
	write("package main\n\nfunc solve() int {\n\treturn 42\n}\n")
	if got := runDiff(); !strings.Contains(got, "-\treturn 0") || !strings.Contains(got, "+\treturn 42") {
		t.Errorf("diff after edit missing the change:\n%s", got)
	}

	// Consecutive diff with no further edits: back to no-changes (the
	// previous diff advanced the snapshot).
	if got := runDiff(); !strings.Contains(got, "no changes") {
		t.Errorf("repeat diff = %q, want no-changes (snapshot must advance)", got)
	}

	// A full read also advances the snapshot: edit, read, then diff
	// shows nothing new.
	write("package main\n\nfunc solve() int {\n\treturn 7\n}\n")
	if _, err := readTool.InvokableRun(ctx, "{}"); err != nil {
		t.Fatalf("read: %v", err)
	}
	if got := runDiff(); !strings.Contains(got, "no changes") {
		t.Errorf("diff after full read = %q, want no-changes (read advances the snapshot too)", got)
	}
}
