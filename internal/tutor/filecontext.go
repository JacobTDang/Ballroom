package tutor

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

// buildFileContext returns the active solution.* file's contents
// (truncated to maxBytes, with a trailing marker if the file is larger),
// or "" if there's no solution file yet (sandbox mode's fresh start) or
// it can't be read — a missing/unreadable file is an expected state, not
// a bug, so this never returns an error. Port of tutor/chat.sh's
// build_file_context, which has the same always-succeeds contract.
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

	info, err := os.Stat(file)
	if err != nil {
		return string(content)
	}
	if info.Size() > int64(maxBytes) {
		return fmt.Sprintf("%s\n...[truncated, file is %d bytes; showing first %d]", content, info.Size(), maxBytes)
	}
	return string(content)
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
