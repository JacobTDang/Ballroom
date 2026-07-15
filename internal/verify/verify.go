// Package verify checks that an exercise's hidden tests actually
// discriminate correct code from incorrect code: the shipped starter must
// FAIL test_command against the hidden tests, and a known-correct
// reference solution (exercises/<id>/.reference/, never mounted into a
// real session — see internal/orchestrator.PrepareWorkspace, which only
// ever copies RepoPath) must PASS the same tests. Without this, a
// vacuous or too-weak test could ship silently, at any of 150+ problems.
package verify

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/JacobTDang/Ballroom/internal/exercise"
)

// referenceDirName is the sibling directory of repo/, checked into git but
// never mounted into a session's workspace, holding a known-correct
// solution for this verification only.
const referenceDirName = ".reference"

// Result is one exercise's verification outcome.
type Result struct {
	ExerciseID string

	StarterFailed bool // true if the shipped starter fails the hidden tests, as expected
	StarterOutput string

	ReferencePassed bool // true if the .reference/ solution passes the hidden tests
	ReferenceOutput string
}

// OK reports whether both checks came back as expected: the starter
// failed and the reference passed.
func (r Result) OK() bool {
	return r.StarterFailed && r.ReferencePassed
}

// Exercise verifies one exercise directory (containing exercise.json).
// testsDir is the sibling hidden-tests directory (tests/<id>/), copied in
// exactly as orchestrator.WaitAndReveal does at real submit time.
//
// A design-kind exercise has nothing to run, so its bar is structural:
// repo/problem.md (the design prompt), repo/solution.md (the template
// the user fills in), and tests/<id>/rubric.md (the hidden rubric the
// reveal delivers) must all exist. Each missing piece is an error, not
// a Result flag -- catalog.List silently skips broken exercises and a
// missing rubric would otherwise only surface live as a 30-second
// submit timeout, so this is the loud backstop.
func Exercise(exerciseDir, testsDir string) (Result, error) {
	ex, err := exercise.Load(filepath.Join(exerciseDir, "exercise.json"))
	if err != nil {
		return Result{}, err
	}
	result := Result{ExerciseID: ex.ID}

	if ex.Kind == exercise.KindDesign {
		for _, required := range []string{
			filepath.Join(ex.RepoPath, "problem.md"),
			filepath.Join(ex.RepoPath, "solution.md"),
			filepath.Join(testsDir, "rubric.md"),
		} {
			if info, err := os.Stat(required); err != nil || info.IsDir() {
				return Result{}, fmt.Errorf("verify: %s: design exercise missing %s", ex.ID, required)
			}
		}
		// Both flags true so OK() reports success through the same
		// Result shape coding exercises use -- the outputs say what was
		// actually checked.
		result.StarterFailed = true
		result.ReferencePassed = true
		result.StarterOutput = "(design exercise: structural check)"
		result.ReferenceOutput = "(design exercise: structural check)"
		return result, nil
	}

	refDir := filepath.Join(exerciseDir, referenceDirName)
	if info, err := os.Stat(refDir); err != nil || !info.IsDir() {
		return Result{}, fmt.Errorf("verify: %s: missing %s (no reference solution to verify against)", ex.ID, refDir)
	}

	starterDir, cleanup, err := buildWorkspace(ex.RepoPath, testsDir)
	if err != nil {
		return Result{}, err
	}
	defer cleanup()
	starterOut, starterErr := runTestCommand(ex.TestCommand, starterDir)
	result.StarterFailed = starterErr != nil
	result.StarterOutput = starterOut

	refWorkDir, cleanup2, err := buildWorkspace(ex.RepoPath, testsDir)
	if err != nil {
		return Result{}, err
	}
	defer cleanup2()
	if err := copyTree(refDir, refWorkDir); err != nil {
		return Result{}, fmt.Errorf("verify: %s: overlay reference solution: %w", ex.ID, err)
	}
	refOut, refErr := runTestCommand(ex.TestCommand, refWorkDir)
	result.ReferencePassed = refErr == nil
	result.ReferenceOutput = refOut

	return result, nil
}

// buildWorkspace assembles a temp dir the same way a real session's
// workspace ends up looking once hidden tests are revealed: repoPath's
// contents, then testsDir's contents copied on top.
func buildWorkspace(repoPath, testsDir string) (dir string, cleanup func(), err error) {
	dir, err = os.MkdirTemp("", "verify-workspace-")
	if err != nil {
		return "", nil, fmt.Errorf("verify: create workspace: %w", err)
	}
	cleanup = func() { os.RemoveAll(dir) }

	if err := copyTree(repoPath, dir); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("verify: copy repo into workspace: %w", err)
	}
	if err := copyTree(testsDir, dir); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("verify: copy tests into workspace: %w", err)
	}
	return dir, cleanup, nil
}

// runTestCommand mirrors internal/session's runTestCommand exactly — same
// shell invocation, same working-dir semantics — so a verify pass here
// means the same thing a real submission's result would.
func runTestCommand(testCommand, dir string) (output string, err error) {
	cmd := exec.Command("sh", "-c", testCommand)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// copyTree recursively copies the contents of src into dst, creating dst
// (and any subdirectories) as needed, overwriting existing files —
// mirrors internal/orchestrator's unexported copyTree/copyFile.
func copyTree(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
