package tutor

import (
	"context"
	"os"
	"testing"
)

// TestMain defaults checkToolCallingForSession to "always native, no
// real network round-trip" for this package's entire test run. Without
// this, newTutorModel's strategy detection (added for JSON-fallback
// tool calling) would make an extra, unscripted probe request against
// the same mock Ollama server nearly every existing test already uses
// (newSequencedOllama), consuming one of that mock's scripted replies
// and desyncing the rest -- request-count assertions like "expected
// exactly 2 requests" (model_test.go's routing tests) would break for
// reasons unrelated to what they're actually testing. Tests that
// specifically exercise strategy detection override this var again
// locally, the same save/restore pattern internal/tui/app_test.go's
// fakeCheckToolCalling already uses.
func TestMain(m *testing.M) {
	checkToolCallingForSession = func(context.Context, string, string, string) (bool, error) {
		return true, nil
	}
	os.Exit(m.Run())
}
