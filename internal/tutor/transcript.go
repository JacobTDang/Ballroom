package tutor

import (
	"fmt"
	"os"
	"path/filepath"
)

// Transcript export: every completed turn is appended to the session's
// transcript file(s) as it happens -- no exit hook, so a crashed pane
// or killed container loses nothing already written. Two destinations
// in a real session (see cmd/ballroom's tutorCmd): the workspace copy
// the user can open mid-session, and a mirror under the persistent
// /data mount, because the workspace temp dir is deleted the moment
// the session ends. Raw markdown, matching history's clean-pairs-only
// semantics -- failed turns and tool scaffolding are not conversation.

// appendTranscriptTurn appends one (user, reply) exchange to every
// path, creating parent directories as needed. Returns the first
// error; the caller decides how loudly to surface it (the pane warns
// once and keeps going -- see turnCompleteMsg).
func appendTranscriptTurn(paths []string, userMessage, reply string) error {
	entry := fmt.Sprintf("## you\n\n%s\n\n## tutor\n\n%s\n\n", userMessage, reply)
	for _, path := range paths {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("tutor: transcript: %w", err)
		}
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("tutor: transcript: %w", err)
		}
		_, werr := f.WriteString(entry)
		cerr := f.Close()
		if werr != nil {
			return fmt.Errorf("tutor: transcript: %w", werr)
		}
		if cerr != nil {
			return fmt.Errorf("tutor: transcript: %w", cerr)
		}
	}
	return nil
}
