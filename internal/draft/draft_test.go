package draft

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestDir(t *testing.T) {
	got := Dir("/data", "two-pointers-01-go")
	want := filepath.Join("/data", ".drafts", "two-pointers-01-go")
	if got != want {
		t.Errorf("Dir(...) = %q, want %q", got, want)
	}
}

func TestSnapshot_WritesFilesAndMeta(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()

	content := "package main\n\nfunc TwoSum() {}\n"
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte(content), 0o644); err != nil {
		t.Fatalf("seed workspace: %v", err)
	}

	wrote, err := Snapshot(dataDir, "two-pointers-01-go", workspace)
	if err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	if !wrote {
		t.Fatal("expected Snapshot to report a write on the first snapshot")
	}

	draftDir := Dir(dataDir, "two-pointers-01-go")
	got, err := os.ReadFile(filepath.Join(draftDir, "solution.go"))
	if err != nil {
		t.Fatalf("expected solution.go in draft dir: %v", err)
	}
	if string(got) != content {
		t.Errorf("draft solution.go = %q, want %q", got, content)
	}

	metaRaw, err := os.ReadFile(filepath.Join(draftDir, "meta.json"))
	if err != nil {
		t.Fatalf("expected meta.json in draft dir: %v", err)
	}
	var meta Meta
	if err := json.Unmarshal(metaRaw, &meta); err != nil {
		t.Fatalf("parse meta.json: %v", err)
	}
	if meta.ExerciseID != "two-pointers-01-go" {
		t.Errorf("meta.ExerciseID = %q, want %q", meta.ExerciseID, "two-pointers-01-go")
	}
	if meta.Bytes != len(content) {
		t.Errorf("meta.Bytes = %d, want %d", meta.Bytes, len(content))
	}
	if meta.Lines != 3 {
		t.Errorf("meta.Lines = %d, want 3", meta.Lines)
	}
	if _, err := time.Parse(time.RFC3339, meta.SavedAt); err != nil {
		t.Errorf("meta.SavedAt = %q not RFC3339: %v", meta.SavedAt, err)
	}
}

func TestSnapshot_UnchangedContentDoesNotRewrite(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()

	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("seed workspace: %v", err)
	}

	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatalf("first Snapshot: %v", err)
	}

	draftFile := filepath.Join(Dir(dataDir, "ex-1"), "solution.go")
	before, err := os.Stat(draftFile)
	if err != nil {
		t.Fatalf("stat draft file: %v", err)
	}

	// Sleep so a spurious rewrite would actually move the mtime enough
	// for the assertion below to catch it.
	time.Sleep(10 * time.Millisecond)

	wrote, err := Snapshot(dataDir, "ex-1", workspace)
	if err != nil {
		t.Fatalf("second Snapshot: %v", err)
	}
	if wrote {
		t.Error("expected Snapshot to report no write when content is unchanged")
	}

	after, err := os.Stat(draftFile)
	if err != nil {
		t.Fatalf("stat draft file after second snapshot: %v", err)
	}
	if !before.ModTime().Equal(after.ModTime()) {
		t.Errorf("draft file was rewritten despite unchanged content: mtime %v -> %v", before.ModTime(), after.ModTime())
	}
}

func TestSnapshot_ChangedContentDoesRewrite(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()
	solutionPath := filepath.Join(workspace, "solution.go")

	if err := os.WriteFile(solutionPath, []byte("v1"), 0o644); err != nil {
		t.Fatalf("seed workspace: %v", err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatalf("first Snapshot: %v", err)
	}

	if err := os.WriteFile(solutionPath, []byte("v2"), 0o644); err != nil {
		t.Fatalf("edit workspace: %v", err)
	}
	wrote, err := Snapshot(dataDir, "ex-1", workspace)
	if err != nil {
		t.Fatalf("second Snapshot: %v", err)
	}
	if !wrote {
		t.Fatal("expected Snapshot to report a write when content changed")
	}

	got, err := os.ReadFile(filepath.Join(Dir(dataDir, "ex-1"), "solution.go"))
	if err != nil {
		t.Fatalf("read draft file: %v", err)
	}
	if string(got) != "v2" {
		t.Errorf("draft content = %q, want %q", got, "v2")
	}
}

func TestSnapshot_NoSolutionFilesIsNotAnError(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir() // no solution.* seeded

	wrote, err := Snapshot(dataDir, "ex-1", workspace)
	if err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	if wrote {
		t.Error("expected no write when the workspace has no solution.* files")
	}
	if _, err := os.Stat(Dir(dataDir, "ex-1")); !os.IsNotExist(err) {
		t.Errorf("expected no draft dir created, stat err = %v", err)
	}
}

func TestSnapshot_MultipleSolutionFiles(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()

	if err := os.WriteFile(filepath.Join(workspace, "solution.hpp"), []byte("class X {};\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "solution.cpp"), []byte("#include \"solution.hpp\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Snapshot(dataDir, "ex-cpp", workspace); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}

	draftDir := Dir(dataDir, "ex-cpp")
	for _, name := range []string{"solution.hpp", "solution.cpp"} {
		if _, err := os.Stat(filepath.Join(draftDir, name)); err != nil {
			t.Errorf("expected %s in draft dir: %v", name, err)
		}
	}
}

func TestSnapshot_DropsStaleFileNoLongerInWorkspace(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()

	if err := os.WriteFile(filepath.Join(workspace, "solution.hpp"), []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "solution.cpp"), []byte("v1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-cpp", workspace); err != nil {
		t.Fatalf("first Snapshot: %v", err)
	}

	if err := os.Remove(filepath.Join(workspace, "solution.cpp")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "solution.hpp"), []byte("v2"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-cpp", workspace); err != nil {
		t.Fatalf("second Snapshot: %v", err)
	}

	draftDir := Dir(dataDir, "ex-cpp")
	if _, err := os.Stat(filepath.Join(draftDir, "solution.cpp")); !os.IsNotExist(err) {
		t.Errorf("expected stale solution.cpp removed from draft dir, stat err = %v", err)
	}
	got, err := os.ReadFile(filepath.Join(draftDir, "solution.hpp"))
	if err != nil {
		t.Fatalf("expected solution.hpp still present: %v", err)
	}
	if string(got) != "v2" {
		t.Errorf("solution.hpp = %q, want %q", got, "v2")
	}
}

func TestLoad_MissingDirReturnsNotFoundWithoutError(t *testing.T) {
	dataDir := t.TempDir()

	d, found, err := Load(dataDir, "never-practiced")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if found {
		t.Errorf("expected found=false for a never-snapshotted exercise, got draft %+v", d)
	}
}

func TestLoad_ReturnsFilesAndMeta(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}

	d, found, err := Load(dataDir, "ex-1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !found {
		t.Fatal("expected found=true")
	}
	if d.ExerciseID != "ex-1" {
		t.Errorf("ExerciseID = %q, want ex-1", d.ExerciseID)
	}
	if len(d.Files) != 1 || filepath.Base(d.Files[0]) != "solution.go" {
		t.Errorf("Files = %v, want exactly one .../solution.go", d.Files)
	}
	if d.Meta.ExerciseID != "ex-1" {
		t.Errorf("Meta.ExerciseID = %q, want ex-1", d.Meta.ExerciseID)
	}
}

func TestLoad_PreviewSkipsBlankLines(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()
	content := "package main\n\n\nfunc TwoSum(nums []int, target int) []int {\n\n\treturn nil\n}\n"
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}

	d, found, err := Load(dataDir, "ex-1")
	if err != nil || !found {
		t.Fatalf("Load: found=%v err=%v", found, err)
	}

	for _, line := range d.Preview {
		if strings.TrimSpace(line) == "" {
			t.Errorf("Preview contains a blank line: %q in %v", line, d.Preview)
		}
	}
	want := []string{"package main", "func TwoSum(nums []int, target int) []int {", "return nil", "}"}
	if len(d.Preview) != len(want) {
		t.Fatalf("Preview = %v, want %v", d.Preview, want)
	}
	for i, w := range want {
		if d.Preview[i] != w {
			t.Errorf("Preview[%d] = %q, want %q", i, d.Preview[i], w)
		}
	}
}

func TestLoad_PreviewCapsAtLimit(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()
	var b strings.Builder
	for i := 0; i < 20; i++ {
		b.WriteString("line\n")
	}
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}

	d, found, err := Load(dataDir, "ex-1")
	if err != nil || !found {
		t.Fatalf("Load: found=%v err=%v", found, err)
	}
	if len(d.Preview) == 0 || len(d.Preview) > 10 {
		t.Errorf("Preview has %d lines, want a small bounded preview (>0, <=10)", len(d.Preview))
	}
}

func TestArchive_PreservesOldFileAndLeavesDraftDirLoadable(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()
	original := "package main\n\nfunc Old() {}\n"
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}

	if err := Archive(dataDir, "ex-1"); err != nil {
		t.Fatalf("Archive: %v", err)
	}

	draftDir := Dir(dataDir, "ex-1")
	got, err := os.ReadFile(filepath.Join(draftDir, "previous.go"))
	if err != nil {
		t.Fatalf("expected previous.go to preserve the archived draft: %v", err)
	}
	if string(got) != original {
		t.Errorf("previous.go = %q, want %q", got, original)
	}
	if _, err := os.Stat(filepath.Join(draftDir, "solution.go")); !os.IsNotExist(err) {
		t.Errorf("expected solution.go rotated away after Archive, stat err = %v", err)
	}

	// The draft dir must still be safely Load-able -- no current draft
	// (nothing new saved since the archive), not an error.
	_, found, err := Load(dataDir, "ex-1")
	if err != nil {
		t.Fatalf("Load after Archive: %v", err)
	}
	if found {
		t.Error("expected no current draft immediately after Archive")
	}

	// A fresh snapshot afterward must still work normally, proving
	// Archive didn't leave the directory in a broken state, and must
	// never disturb the archived previous.go -- that's the whole point
	// of "start fresh never destroys anything".
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte("fresh start"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatalf("Snapshot after Archive: %v", err)
	}
	d, found, err := Load(dataDir, "ex-1")
	if err != nil || !found {
		t.Fatalf("Load after fresh snapshot: found=%v err=%v", found, err)
	}
	if d.Meta.Bytes != len("fresh start") {
		t.Errorf("fresh draft Bytes = %d, want %d", d.Meta.Bytes, len("fresh start"))
	}
	stillThere, err := os.ReadFile(filepath.Join(draftDir, "previous.go"))
	if err != nil || string(stillThere) != original {
		t.Errorf("previous.go was disturbed by the later snapshot: %q (err %v), want %q", stillThere, err, original)
	}
}

func TestArchive_NoOpWhenNoCurrentDraft(t *testing.T) {
	dataDir := t.TempDir()
	if err := Archive(dataDir, "never-practiced"); err != nil {
		t.Errorf("Archive on a never-practiced exercise should be a no-op, got %v", err)
	}
}

func TestArchive_RotationOverwritesEarlierPrevious(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()
	solutionPath := filepath.Join(workspace, "solution.go")

	if err := os.WriteFile(solutionPath, []byte("gen1"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatal(err)
	}
	if err := Archive(dataDir, "ex-1"); err != nil {
		t.Fatalf("first Archive: %v", err)
	}

	if err := os.WriteFile(solutionPath, []byte("gen2"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatal(err)
	}
	if err := Archive(dataDir, "ex-1"); err != nil {
		t.Fatalf("second Archive: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(Dir(dataDir, "ex-1"), "previous.go"))
	if err != nil {
		t.Fatalf("read previous.go: %v", err)
	}
	if string(got) != "gen2" {
		t.Errorf("previous.go = %q, want the most recently archived generation %q", got, "gen2")
	}
}

// TestExists_* cover the picker's draft marker (issue #255): Exists must
// be a cheap existence check (no file content reads -- see the doc
// comment) that agrees exactly with Load's ok result, so the marker
// never promises a resume prompt that doesn't actually appear.

func TestExists_FalseWhenNeverSnapshotted(t *testing.T) {
	dataDir := t.TempDir()
	if Exists(dataDir, "never-practiced") {
		t.Error("expected Exists to be false for an exercise with no draft dir at all")
	}
}

func TestExists_TrueAfterSnapshot(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}

	if !Exists(dataDir, "ex-1") {
		t.Error("expected Exists to be true right after a snapshot")
	}
}

// TestExists_FalseAfterArchive: Archive rotates solution.* to
// previous.* and leaves the draft directory itself in place -- Exists
// must track "is there a current, resumable draft", not mere directory
// presence, or the picker marker would promise a resume prompt for a
// problem that was deliberately started fresh.
func TestExists_FalseAfterArchive(t *testing.T) {
	dataDir := t.TempDir()
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, "solution.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(dataDir, "ex-1", workspace); err != nil {
		t.Fatalf("Snapshot: %v", err)
	}
	if err := Archive(dataDir, "ex-1"); err != nil {
		t.Fatalf("Archive: %v", err)
	}

	if Exists(dataDir, "ex-1") {
		t.Error("expected Exists to be false after Archive left only previous.* behind")
	}
}

func TestExists_DoesNotReadFileContents(t *testing.T) {
	// A directory containing a solution.* whose bytes would fail to
	// parse as anything meaningful must still report true -- Exists
	// only checks presence, never opens/parses the file for content.
	dataDir := t.TempDir()
	dir := Dir(dataDir, "ex-1")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "solution.go"), []byte{0xff, 0x00, 0xfe}, 0o644); err != nil {
		t.Fatal(err)
	}

	if !Exists(dataDir, "ex-1") {
		t.Error("expected Exists to report true regardless of file content")
	}
}
