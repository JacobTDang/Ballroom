// Command verify-exercises walks exercises/* and, for each one, checks
// that its hidden tests actually discriminate correct from incorrect code
// (see internal/verify). Run from the repo root, or pass it as the sole
// argument.
package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/JacobTDang/Ballroom/internal/verify"
)

func main() {
	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	ok, err := run(root, os.Stdout)
	if err != nil {
		fmt.Fprintln(os.Stderr, "verify-exercises:", err)
		os.Exit(1)
	}
	if !ok {
		os.Exit(1)
	}
}

// listExerciseDirs returns the ids (subdirectory names) under
// exercisesRoot that actually contain an exercise.json, skipping
// _template and anything else that isn't a real exercise.
func listExerciseDirs(exercisesRoot string) ([]string, error) {
	entries, err := os.ReadDir(exercisesRoot)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", exercisesRoot, err)
	}

	var ids []string
	for _, e := range entries {
		if !e.IsDir() || e.Name() == "_template" {
			continue
		}
		if _, err := os.Stat(filepath.Join(exercisesRoot, e.Name(), "exercise.json")); err != nil {
			continue
		}
		ids = append(ids, e.Name())
	}
	return ids, nil
}

// run verifies every exercise under root/exercises against
// root/tests/<id>, writing a PASS/FAIL line (plus failure detail) per
// exercise to out. ok is false if any exercise failed either check.
func run(root string, out io.Writer) (ok bool, err error) {
	exercisesRoot := filepath.Join(root, "exercises")
	testsRoot := filepath.Join(root, "tests")

	ids, err := listExerciseDirs(exercisesRoot)
	if err != nil {
		return false, err
	}

	ok = true
	for _, id := range ids {
		exDir := filepath.Join(exercisesRoot, id)
		testsDir := filepath.Join(testsRoot, id)

		result, verr := verify.Exercise(exDir, testsDir)
		if verr != nil {
			ok = false
			fmt.Fprintf(out, "FAIL %s: %v\n", id, verr)
			continue
		}
		if !result.OK() {
			ok = false
			fmt.Fprintf(out, "FAIL %s\n", id)
			if !result.StarterFailed {
				fmt.Fprintf(out, "  starter unexpectedly PASSED (hidden tests don't discriminate):\n%s\n", result.StarterOutput)
			}
			if !result.ReferencePassed {
				fmt.Fprintf(out, "  reference solution FAILED:\n%s\n", result.ReferenceOutput)
			}
			continue
		}
		fmt.Fprintf(out, "PASS %s\n", id)
	}
	return ok, nil
}
