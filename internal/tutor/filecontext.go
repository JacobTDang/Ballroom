package tutor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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
