package tutor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// solutionSnapshot is the "last version the model saw" of the solution
// file, shared by read_solution_file and read_solution_diff: either
// tool advances it, so a diff always answers "what changed since you
// last looked". The baseline is taken when the session's tools are
// built (buildTools), so the first diff shows changes since session
// start. Mutex'd because tool calls run on turn goroutines.
type solutionSnapshot struct {
	mu   sync.Mutex
	last string
}

// swap records current as the newest seen version and returns what it
// replaced.
func (s *solutionSnapshot) swap(current string) (previous string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	previous = s.last
	s.last = current
	return previous
}

// rawSolutionContent is buildFileContext's un-numbered sibling: the
// active solution.* file's raw bytes (capped at maxBytes), or "" when
// missing/unreadable. Raw on purpose -- diffing the numbered render
// would make every line after an insertion look changed.
func rawSolutionContent(workDir string, maxBytes int) string {
	matches, err := filepath.Glob(filepath.Join(workDir, "solution.*"))
	if err != nil || len(matches) == 0 {
		return ""
	}
	f, err := os.Open(matches[0])
	if err != nil {
		return ""
	}
	defer f.Close()
	content, err := io.ReadAll(io.LimitReader(f, int64(maxBytes)))
	if err != nil {
		return ""
	}
	return string(content)
}

// A small line-based unified diff, hand-rolled: the only consumer is
// the read_solution_diff tool, whose audience is a model that just
// needs "what changed" in the familiar -/+ shape -- not worth a
// dependency (and the ask-before-adding-dependencies rule agrees).
// O(n*m) LCS is plenty for solution files.

const diffContextLines = 2

type diffOpKind int

const (
	diffSame diffOpKind = iota
	diffDel
	diffAdd
)

type diffOp struct {
	kind  diffOpKind
	text  string
	aLine int // 1-indexed line in old (valid for same/del)
	bLine int // 1-indexed line in new (valid for same/add)
}

// diffUnified renders the changes from oldText to newText as a unified
// diff with hunk headers and 2 lines of context. Empty string means no
// changes.
func diffUnified(oldText, newText string) string {
	if oldText == newText {
		return ""
	}
	a := splitLines(oldText)
	b := splitLines(newText)

	ops := diffOps(a, b)

	// Group ops into hunks: changes separated by more than 2*context
	// unchanged lines get their own hunk.
	var out strings.Builder
	i := 0
	for i < len(ops) {
		if ops[i].kind == diffSame {
			i++
			continue
		}
		// Found a change: hunk start (back up for leading context).
		start := i - diffContextLines
		if start < 0 {
			start = 0
		}
		// Extend to the last change within merge distance.
		end := i
		lastChange := i
		for end < len(ops) {
			if ops[end].kind != diffSame {
				lastChange = end
				end++
				continue
			}
			// Run of unchanged lines: stop if it exceeds the merge gap.
			gap := 0
			for end+gap < len(ops) && ops[end+gap].kind == diffSame {
				gap++
			}
			if gap > 2*diffContextLines {
				break
			}
			end += gap
		}
		hunkEnd := lastChange + diffContextLines + 1
		if hunkEnd > len(ops) {
			hunkEnd = len(ops)
		}

		out.WriteString(hunkHeader(ops[start:hunkEnd]))
		for _, op := range ops[start:hunkEnd] {
			switch op.kind {
			case diffDel:
				out.WriteString("-" + op.text + "\n")
			case diffAdd:
				out.WriteString("+" + op.text + "\n")
			default:
				out.WriteString(" " + op.text + "\n")
			}
		}
		i = hunkEnd
	}
	return strings.TrimSuffix(out.String(), "\n")
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// diffOps computes the op sequence via a standard LCS backtrack.
func diffOps(a, b []string) []diffOp {
	n, m := len(a), len(b)
	lcs := make([][]int, n+1)
	for i := range lcs {
		lcs[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else if lcs[i+1][j] >= lcs[i][j+1] {
				lcs[i][j] = lcs[i+1][j]
			} else {
				lcs[i][j] = lcs[i][j+1]
			}
		}
	}
	var ops []diffOp
	i, j := 0, 0
	for i < n && j < m {
		switch {
		case a[i] == b[j]:
			ops = append(ops, diffOp{kind: diffSame, text: a[i], aLine: i + 1, bLine: j + 1})
			i++
			j++
		case lcs[i+1][j] >= lcs[i][j+1]:
			ops = append(ops, diffOp{kind: diffDel, text: a[i], aLine: i + 1})
			i++
		default:
			ops = append(ops, diffOp{kind: diffAdd, text: b[j], bLine: j + 1})
			j++
		}
	}
	for ; i < n; i++ {
		ops = append(ops, diffOp{kind: diffDel, text: a[i], aLine: i + 1})
	}
	for ; j < m; j++ {
		ops = append(ops, diffOp{kind: diffAdd, text: b[j], bLine: j + 1})
	}
	return ops
}

// hunkHeader renders the @@ -start,count +start,count @@ line for a
// hunk's op slice.
func hunkHeader(ops []diffOp) string {
	aStart, aCount, bStart, bCount := 0, 0, 0, 0
	for _, op := range ops {
		if op.kind != diffAdd {
			if aStart == 0 {
				aStart = op.aLine
			}
			aCount++
		}
		if op.kind != diffDel {
			if bStart == 0 {
				bStart = op.bLine
			}
			bCount++
		}
	}
	return fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", aStart, aCount, bStart, bCount)
}
