package tutor

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// startTestNvim starts a real `nvim --headless --listen <socket>` with
// docker/nvim/'s real config (so the real ballroom_highlight.lua module
// loads, not a stub) and returns the socket path once it's ready to
// accept RPC calls. Self-skips with a clear message if `nvim` isn't on
// PATH — these are integration tests against a real editor, not
// something every environment running `go test` is expected to have
// installed (matches tutor/chat_test.sh's original behavior).
func startTestNvim(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("nvim"); err != nil {
		t.Skip("SKIP: nvim not found on PATH, skipping live-nvim RPC tests")
	}

	configHome := t.TempDir()
	nvimConfigDir := filepath.Join(configHome, "nvim")
	if err := os.MkdirAll(filepath.Join(nvimConfigDir, "lua"), 0o755); err != nil {
		t.Fatalf("mkdir scratch nvim config: %v", err)
	}

	repoNvimDir := repoNvimConfigDir(t)
	copyFile(t, filepath.Join(repoNvimDir, "init.lua"), filepath.Join(nvimConfigDir, "init.lua"))
	copyFile(t, filepath.Join(repoNvimDir, "lua", "ballroom_highlight.lua"), filepath.Join(nvimConfigDir, "lua", "ballroom_highlight.lua"))

	// Unix domain socket paths are capped at ~104 bytes by the OS
	// (sockaddr_un). t.TempDir() nests under the test function's full
	// name (e.g. .../TestApplyHighlight_LegitApostropheNoteSucceeds/001/...),
	// which alone can exceed that limit — bind() then fails silently
	// from this test's point of view (nvim just never creates the
	// socket). Use a short, non-test-name-derived temp dir for the
	// socket specifically to stay well under the limit.
	socketDir, err := os.MkdirTemp("", "ballroom-nvim-test-")
	if err != nil {
		t.Fatalf("create scratch socket dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(socketDir) })
	socket := filepath.Join(socketDir, "s.sock")

	var stderr bytes.Buffer
	cmd := exec.Command("nvim", "--headless", "--listen", socket)
	cmd.Env = append(os.Environ(), "XDG_CONFIG_HOME="+configHome)
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start nvim --headless: %v", err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	})

	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if info, err := os.Stat(socket); err == nil && info.Mode()&os.ModeSocket != 0 {
			return socket
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("nvim --headless never created its RPC socket; stderr: %s", stderr.String())
	return ""
}

func repoNvimConfigDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	abs, err := filepath.Abs(filepath.Join(filepath.Dir(thisFile), "..", "..", "docker", "nvim"))
	if err != nil {
		t.Fatalf("resolve docker/nvim path: %v", err)
	}
	return abs
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read %s: %v", src, err)
	}
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", dst, err)
	}
}

func TestEscapeVimSingleQuoted_DoublesEmbeddedQuotes(t *testing.T) {
	got := escapeVimSingleQuoted("it's a test")
	want := "it''s a test"
	if got != want {
		t.Errorf("escapeVimSingleQuoted(%q) = %q, want %q", "it's a test", got, want)
	}
}

func TestRemoteExpr_EmptySocketReturnsEmptyNoError(t *testing.T) {
	ctx := context.Background()
	out, err := remoteExpr(ctx, "", "1+1")
	if err != nil {
		t.Fatalf("remoteExpr with empty socket: %v", err)
	}
	if out != "" {
		t.Errorf("result = %q, want empty string when socket is unset", out)
	}
}

func TestRemoteExpr_NonexistentSocketReturnsEmptyNoError(t *testing.T) {
	ctx := context.Background()
	out, err := remoteExpr(ctx, filepath.Join(t.TempDir(), "does-not-exist.sock"), "1+1")
	if err != nil {
		t.Fatalf("remoteExpr with nonexistent socket: %v", err)
	}
	if out != "" {
		t.Errorf("result = %q, want empty string when socket doesn't exist", out)
	}
}

func TestApplyHighlight_LegitApostropheNoteSucceeds(t *testing.T) {
	socket := startTestNvim(t)
	ctx := context.Background()

	expr := highlightExpr("test.txt", 1, 1, "it's off by one here")
	out, err := remoteExpr(ctx, socket, expr)
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}
	if out != "ok" {
		t.Errorf("result = %q, want %q — an apostrophe in the note should not break the VimL expression", out, "ok")
	}
}

func TestApplyHighlight_InjectionPayloadDoesNotExecute(t *testing.T) {
	socket := startTestNvim(t)
	ctx := context.Background()

	canary := filepath.Join(t.TempDir(), "canary")
	// Attempt to break out of the single-quoted VimL string and execute
	// arbitrary Lua via a crafted note. If escaping is correct, this is
	// just inert text inside the string literal, not executable.
	payload := "x'); os.execute('touch " + canary + "'); local y = ('"

	expr := highlightExpr("test.txt", 1, 1, payload)
	if _, err := remoteExpr(ctx, socket, expr); err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}

	if _, err := os.Stat(canary); err == nil {
		t.Fatal("injection payload executed — canary file was created")
	}
}

func TestReadCursorPosition_ReturnsPosition(t *testing.T) {
	socket := startTestNvim(t)
	ctx := context.Background()

	out, err := remoteExpr(ctx, socket, cursorPositionExpr())
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}
	if out == "" {
		t.Fatal("expected a JSON result, got empty string")
	}
	// Full shape is verified end-to-end by the tool-level test in
	// tools_test.go; here just confirm the raw RPC call succeeds and
	// returns something JSON-shaped.
	if !strings.HasPrefix(out, "{") {
		t.Errorf("result = %q, want it to look like a JSON object", out)
	}
}
