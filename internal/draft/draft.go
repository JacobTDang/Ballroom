// Package draft persists in-progress solution files somewhere that
// survives session exit: data/.drafts/<exercise-id>/, written host-side
// while a session runs (see internal/orchestrator.RunExercise and its
// SnapshotLoop), independent of the disposable per-session workspace
// temp dir that gets deleted the moment the session ends (issue #221 --
// previously nothing persisted a session's actual code anywhere).
//
// A leaf package: every function takes plain dataDir/exerciseID/
// workspaceDir strings, the same paths internal/config and
// internal/exercise already resolve, rather than importing either --
// matching this codebase's convention of small local duplication over
// a cross-package dependency (see internal/tutor/transcript.go).
package draft

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// metaFileName is the sidecar written alongside a draft's solution
// files in its draft directory.
const metaFileName = "meta.json"

// previewLineCount bounds how many non-blank lines Load collects for
// Draft.Preview -- enough to recognize a draft at a glance without
// dumping the whole file.
const previewLineCount = 5

// Meta is the sidecar JSON written next to a snapshot's solution
// files, recording when it was taken and a cheap size summary -- e.g.
// for a future "resume draft? saved 3m ago, 42 lines" prompt, without
// having to read the solution files themselves.
type Meta struct {
	SavedAt    string `json:"saved_at"` // RFC3339, UTC
	ExerciseID string `json:"exercise_id"`
	Bytes      int    `json:"bytes"`
	Lines      int    `json:"lines"`
}

// Draft is a loaded solution draft: where its files live on disk, the
// snapshot's metadata, and a short preview for showing the user what
// they'd be resuming.
type Draft struct {
	ExerciseID string
	Dir        string
	Files      []string
	Meta       Meta
	Preview    []string
}

// draftsDirName is dot-prefixed deliberately: drafts of Go exercises are
// real .go files, and a visible data/drafts/ would put a package main
// with no func main() into the module tree -- enough to break a plain
// `go build ./...` from the repo root as soon as any Go exercise has
// been practiced. Go's tooling skips dot-directories, so the drafts stay
// invisible to it while remaining ordinary files for everything else.
const draftsDirName = ".drafts"

// Dir returns the directory under dataDir where exerciseID's solution
// draft is stored. The directory itself is created lazily by the first
// successful Snapshot; callers must not assume it exists.
func Dir(dataDir, exerciseID string) string {
	return filepath.Join(dataDir, draftsDirName, exerciseID)
}

// Snapshot copies workspaceDir's current solution.* files (the same
// top-level glob docker/entrypoint.sh uses to find the file to open in
// the editor) into exerciseID's draft directory under dataDir, plus a
// meta.json sidecar. Designed for cheap periodic polling (every few
// seconds -- see internal/orchestrator.SnapshotLoop): when the
// workspace's solution files are already byte-identical to what's
// saved, it writes nothing and returns (false, nil) rather than
// rewriting unchanged bytes on every tick.
//
// Also returns (false, nil), not an error, when the workspace has no
// solution.* files at all -- an expected state (sandbox mode, or a
// session that hasn't reached the editor yet), not a bug.
func Snapshot(dataDir, exerciseID, workspaceDir string) (bool, error) {
	workspaceFiles, err := solutionFiles(workspaceDir)
	if err != nil {
		return false, fmt.Errorf("draft: glob workspace solution files: %w", err)
	}
	if len(workspaceFiles) == 0 {
		return false, nil
	}

	contents, err := readFiles(workspaceFiles)
	if err != nil {
		return false, fmt.Errorf("draft: read workspace solution files: %w", err)
	}
	newHash := hashContents(contents)

	dir := Dir(dataDir, exerciseID)
	existingFiles, err := solutionFiles(dir)
	if err != nil {
		return false, fmt.Errorf("draft: glob existing draft files: %w", err)
	}
	existing, err := readFiles(existingFiles)
	if err != nil {
		return false, fmt.Errorf("draft: read existing draft files: %w", err)
	}
	if hashContents(existing) == newHash {
		return false, nil
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, fmt.Errorf("draft: create draft dir: %w", err)
	}

	// Drop any solution.* left over from a previous snapshot that the
	// current workspace no longer has, so the draft dir never mixes
	// files from two different snapshots.
	for _, f := range existingFiles {
		if _, keep := contents[filepath.Base(f)]; !keep {
			if err := os.Remove(f); err != nil {
				return false, fmt.Errorf("draft: remove stale %s: %w", f, err)
			}
		}
	}

	var totalBytes, totalLines int
	for _, data := range contents {
		totalBytes += len(data)
		totalLines += countLines(data)
	}
	for name, data := range contents {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o644); err != nil {
			return false, fmt.Errorf("draft: write %s: %w", name, err)
		}
	}

	meta := Meta{
		SavedAt:    time.Now().UTC().Format(time.RFC3339),
		ExerciseID: exerciseID,
		Bytes:      totalBytes,
		Lines:      totalLines,
	}
	metaData, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return false, fmt.Errorf("draft: marshal meta: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, metaFileName), metaData, 0o644); err != nil {
		return false, fmt.Errorf("draft: write meta: %w", err)
	}
	return true, nil
}

// Load reads exerciseID's saved draft from dataDir. A missing or empty
// draft directory -- never snapshotted, or archived with nothing saved
// since -- is reported as (_, false, nil), not an error: that's the
// ordinary state for any exercise not currently in progress.
func Load(dataDir, exerciseID string) (Draft, bool, error) {
	dir := Dir(dataDir, exerciseID)

	files, err := solutionFiles(dir)
	if err != nil {
		return Draft{}, false, fmt.Errorf("draft: glob %s: %w", dir, err)
	}
	if len(files) == 0 {
		return Draft{}, false, nil
	}

	metaData, err := os.ReadFile(filepath.Join(dir, metaFileName))
	if err != nil {
		return Draft{}, false, fmt.Errorf("draft: read meta: %w", err)
	}
	var meta Meta
	if err := json.Unmarshal(metaData, &meta); err != nil {
		return Draft{}, false, fmt.Errorf("draft: parse meta %s: %w", filepath.Join(dir, metaFileName), err)
	}

	preview, err := previewLines(files, previewLineCount)
	if err != nil {
		return Draft{}, false, fmt.Errorf("draft: build preview: %w", err)
	}

	return Draft{
		ExerciseID: exerciseID,
		Dir:        dir,
		Files:      files,
		Meta:       meta,
		Preview:    preview,
	}, true, nil
}

// Archive rotates exerciseID's current draft files to previous.<ext>
// (overwriting any earlier previous.<ext> -- one generation of
// history, not a stack) and removes the now-stale meta.json, so a
// caller can start a session fresh from the exercise's pristine
// starter without ever destroying the last draft outright: it's still
// on disk as the previous generation. A no-op, not an error, when
// there's no current draft to archive.
func Archive(dataDir, exerciseID string) error {
	dir := Dir(dataDir, exerciseID)

	files, err := solutionFiles(dir)
	if err != nil {
		return fmt.Errorf("draft: glob %s: %w", dir, err)
	}
	if len(files) == 0 {
		return nil
	}

	for _, f := range files {
		ext := strings.TrimPrefix(filepath.Base(f), "solution")
		dest := filepath.Join(dir, "previous"+ext)
		if err := os.Rename(f, dest); err != nil {
			return fmt.Errorf("draft: archive %s: %w", f, err)
		}
	}

	metaPath := filepath.Join(dir, metaFileName)
	if err := os.Remove(metaPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("draft: remove stale meta: %w", err)
	}
	return nil
}

// solutionFiles returns the sorted, absolute paths of every solution.*
// file directly inside dir -- the same top-level glob
// docker/entrypoint.sh uses to find the file to open in the editor
// (and internal/tutor's buildFileContext/diff.go use to find the
// active file to read). A dir that doesn't exist yet yields an empty,
// nil-error slice -- Glob simply matches nothing.
func solutionFiles(dir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "solution.*"))
	if err != nil {
		return nil, err
	}
	files := make([]string, 0, len(matches))
	for _, m := range matches {
		info, err := os.Stat(m)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			continue
		}
		files = append(files, m)
	}
	sort.Strings(files)
	return files, nil
}

// readFiles reads every path into a basename -> bytes map.
func readFiles(paths []string) (map[string][]byte, error) {
	out := make(map[string][]byte, len(paths))
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		out[filepath.Base(p)] = data
	}
	return out, nil
}

// hashContents produces a deterministic content hash over a
// basename->bytes set, independent of Go's randomized map iteration
// order, so two snapshots with the same file set and bytes always hash
// equal -- the cheap short-circuit Snapshot polls on every tick.
func hashContents(files map[string][]byte) string {
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	h := sha256.New()
	for _, name := range names {
		h.Write([]byte(name))
		h.Write([]byte{0})
		h.Write(files[name])
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

// countLines counts newline-delimited lines the same way
// tutor.numberLines does (bufio.Scanner, so a trailing newline doesn't
// produce a spurious extra empty line).
func countLines(data []byte) int {
	n := 0
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		n++
	}
	return n
}

// previewLines reads up to n non-blank lines (whitespace-trimmed,
// blank lines skipped) from files in order, for a human-scannable
// "here's what your draft looks like" preview.
func previewLines(files []string, n int) ([]string, error) {
	lines := make([]string, 0, n)
	for _, f := range files {
		if len(lines) >= n {
			break
		}
		data, err := os.ReadFile(f)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", f, err)
		}
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			lines = append(lines, line)
			if len(lines) >= n {
				break
			}
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("scan %s: %w", f, err)
		}
	}
	return lines, nil
}
