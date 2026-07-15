package tutor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// lastTestResultFile is the well-known dotfile internal/session's
// writeLastTestResult writes into the workspace after every submit.
// Duplicated as a literal here (rather than importing internal/session
// for one constant) matching this codebase's convention of local
// duplication for small pieces over a cross-package dependency — see
// internal/verify's copyTree comment for the same rationale.
const lastTestResultFile = ".ballroom-last-test-result.json"

// lastTestResult mirrors the JSON shape internal/session writes.
type lastTestResult struct {
	Result      string    `json:"result"`
	Output      string    `json:"output"`
	TestCommand string    `json:"test_command"`
	RecordedAt  time.Time `json:"recorded_at"`
}

// readLastTestResult returns the most recent submit's test result from
// workDir, or (zero, false, nil) if none exists yet — haven't submitted
// this session, or sandbox mode (which never runs test_command) — an
// expected state, not an error. A malformed file (exists but won't
// parse) IS a real error, unlike a missing file.
func readLastTestResult(workDir string) (lastTestResult, bool, error) {
	data, err := os.ReadFile(filepath.Join(workDir, lastTestResultFile))
	if os.IsNotExist(err) {
		return lastTestResult{}, false, nil
	}
	if err != nil {
		return lastTestResult{}, false, fmt.Errorf("tutor: read last test result: %w", err)
	}

	var result lastTestResult
	if err := json.Unmarshal(data, &result); err != nil {
		return lastTestResult{}, false, fmt.Errorf("tutor: parse last test result: %w", err)
	}
	return result, true, nil
}

// numberLines prefixes each line of content with its 1-indexed line
// number, tab-separated — the same numbering highlight_lines' start/end
// arguments expect. Without this, the only way the model can figure out
// which line is which is by counting through raw, unnumbered prose-
// formatted code — a well-known LLM weak spot, and the direct cause of
// highlight_lines calls landing on the wrong line. Uses bufio.Scanner
// (not strings.Split on "\n") specifically so a trailing newline doesn't
// produce a spurious extra numbered empty line at the end, matching how
// `cat -n` counts lines.
func numberLines(content string) string {
	if content == "" {
		return ""
	}
	var b strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	n := 0
	for scanner.Scan() {
		n++
		if n > 1 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "%d\t%s", n, scanner.Text())
	}
	return b.String()
}

// buildFileContext returns the active solution.* file's contents, each
// line numbered (see numberLines), truncated to maxBytes with a trailing
// marker if the file is larger, or "" if there's no solution file yet
// (sandbox mode's fresh start) or it can't be read — a missing/unreadable
// file is an expected state, not a bug, so this never returns an error.
func buildFileContext(workDir string, maxBytes int) string {
	matches, err := filepath.Glob(filepath.Join(workDir, "solution.*"))
	if err != nil || len(matches) == 0 {
		return ""
	}
	file := matches[0]

	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	content, err := io.ReadAll(io.LimitReader(f, int64(maxBytes)))
	if err != nil {
		return ""
	}
	numbered := numberLines(string(content))

	info, err := os.Stat(file)
	if err != nil {
		return numbered
	}
	if info.Size() > int64(maxBytes) {
		return fmt.Sprintf("%s\n...[truncated, file is %d bytes; showing first %d]", numbered, info.Size(), maxBytes)
	}
	return numbered
}

// readProblemStatement returns problem.md's contents from workDir, or ""
// if it doesn't exist (e.g. sandbox mode) or can't be read — same
// graceful-degradation contract as buildFileContext.
func readProblemStatement(workDir string) string {
	data, err := os.ReadFile(filepath.Join(workDir, "problem.md"))
	if err != nil {
		return ""
	}
	return string(data)
}

// readRubric returns rubric.md's contents from workDir, or "" if it
// doesn't exist -- which is the normal state for a design session's
// whole working phase: the rubric is hidden content (tests/<id>/) that
// only lands in the workspace after M-q submit reveals it. Same
// graceful-degradation contract as readProblemStatement.
func readRubric(workDir string) string {
	data, err := os.ReadFile(filepath.Join(workDir, "rubric.md"))
	if err != nil {
		return ""
	}
	return string(data)
}
