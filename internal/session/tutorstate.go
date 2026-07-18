package session

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// tutorStateFile is the well-known dotfile internal/tutor's tutor pane
// writes into the workspace after every turn (see model.go's
// writeTutorState) -- same file-based handoff lastTestResultFile uses
// in the other direction, and same local-duplication convention this
// package already follows for that constant (see internal/tutor/
// filecontext.go's own doc comment).
const tutorStateFile = ".ballroom-tutor-state.json"

// tutorState mirrors the JSON shape internal/tutor writes.
type tutorState struct {
	HintsUsed int    `json:"hints_used"`
	TutorMode string `json:"tutor_mode"`
	Model     string `json:"model"`
}

// readTutorState returns the tutor pane's current assistance counters
// from workDir, degrading to the zero value whenever the file is
// missing (sandbox mode, or a submit/reference that happens before the
// tutor pane has written anything yet) or malformed. Unlike
// readLastTestResult, a bad tutor-state file is never a real error: it
// describes an attempt's tutor-assistance metadata, not the attempt's
// own graded outcome, so it must never fail (or even warn on) the
// submit/reference it's attached to.
func readTutorState(workDir string) tutorState {
	data, err := os.ReadFile(filepath.Join(workDir, tutorStateFile))
	if err != nil {
		return tutorState{}
	}
	var s tutorState
	if err := json.Unmarshal(data, &s); err != nil {
		return tutorState{}
	}
	return s
}
