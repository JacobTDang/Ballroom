package tutor

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
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

// TestRemoteExpr_TimesOutIfNvimHangs simulates a blocking nvim prompt
// (e.g. a swap-file recovery dialog) — without remoteExprTimeout, this
// would hang indefinitely and, by extension, freeze the whole
// synchronous tutor turn loop. Uses a stand-in shell script (via the
// nvimCommand package var) instead of a real nvim instance, which can't
// easily be made to hang on demand.
func TestRemoteExpr_TimesOutIfNvimHangs(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("SKIP: needs a POSIX shell script standing in for nvim")
	}

	scriptDir := t.TempDir()
	script := filepath.Join(scriptDir, "fake-nvim")
	// `exec sleep 30`, not a bare `sleep 30` -- exec replaces the shell
	// process image instead of forking a child. A bare `sleep 30` would
	// fork, and killing the (now-dead) parent shell on timeout would
	// leave the orphaned sleep still holding the stdout/stderr pipes
	// open, so cmd.Wait() would block on it regardless -- not
	// representative of the real nvim binary, which doesn't fork like
	// this for a --remote-expr client call.
	if err := os.WriteFile(script, []byte("#!/bin/sh\nexec sleep 30\n"), 0o755); err != nil {
		t.Fatalf("write fake nvim script: %v", err)
	}

	socketDir, err := os.MkdirTemp("", "ballroom-nvim-timeout-test-")
	if err != nil {
		t.Fatalf("create scratch socket dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(socketDir) })
	socket := filepath.Join(socketDir, "s.sock")
	l, err := net.Listen("unix", socket)
	if err != nil {
		t.Fatalf("create fake socket: %v", err)
	}
	defer l.Close()

	origCommand, origTimeout := nvimCommand, remoteExprTimeout
	nvimCommand = script
	remoteExprTimeout = 200 * time.Millisecond
	t.Cleanup(func() { nvimCommand, remoteExprTimeout = origCommand, origTimeout })

	start := time.Now()
	_, err = remoteExpr(context.Background(), socket, "1+1")
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected an error when nvim hangs past remoteExprTimeout")
	}
	if elapsed > 5*time.Second {
		t.Errorf("remoteExpr took %v to return, want it to time out quickly (~%v)", elapsed, remoteExprTimeout)
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

// noteCountExpr calls ballroom_highlight's own note_count() rather than
// inspecting extmarks/namespaces directly from the test -- same reason
// add_highlight/toggle/clear_all are all reached via v:lua.require(...)
// instead of poking nvim's buffer state by hand: it exercises the
// module's real public surface, not a reimplementation of it.
const noteCountExpr = `v:lua.require('ballroom_highlight').note_count()`

func TestApplyHighlight_DoesNotPaintBackgroundHighlightBlock(t *testing.T) {
	socket := startTestNvim(t)
	ctx := context.Background()

	expr := highlightExpr("test.txt", 1, 1, "a note")
	if _, err := remoteExpr(ctx, socket, expr); err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}

	extmarkCountExpr := `len(nvim_buf_get_extmarks(bufnr('%'), nvim_create_namespace('ballroom_tutor'), 0, -1, {}))`
	out, err := remoteExpr(ctx, socket, extmarkCountExpr)
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}
	if out != "1" {
		t.Errorf("extmark count = %s, want 1 (just the note's virtual text, no background highlight block)", out)
	}
}

func TestFocusGained_ClearsNotesLeftWhileAway(t *testing.T) {
	socket := startTestNvim(t)
	ctx := context.Background()

	expr := highlightExpr("test.txt", 1, 1, "a note")
	if _, err := remoteExpr(ctx, socket, expr); err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}

	before, err := remoteExpr(ctx, socket, noteCountExpr)
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}
	if before != "1" {
		t.Fatalf("note_count() = %s before FocusGained, want 1", before)
	}

	if _, err := remoteExpr(ctx, socket, `execute('doautocmd FocusGained')`); err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}

	after, err := remoteExpr(ctx, socket, noteCountExpr)
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}
	if after != "0" {
		t.Errorf("note_count() = %s after FocusGained, want 0 -- notes should disappear once the user is back", after)
	}
}

func TestFocusLost_LeavesNotesInPlace(t *testing.T) {
	socket := startTestNvim(t)
	ctx := context.Background()

	expr := highlightExpr("test.txt", 1, 1, "a note")
	if _, err := remoteExpr(ctx, socket, expr); err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}

	if _, err := remoteExpr(ctx, socket, `execute('doautocmd FocusLost')`); err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}

	out, err := remoteExpr(ctx, socket, noteCountExpr)
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}
	if out != "1" {
		t.Errorf("note_count() = %s after FocusLost, want 1 -- a note left while away must survive until the user actually returns", out)
	}
}

// noteVirtLines fetches the single note extmark's virt_lines -- each
// entry is one on-screen row, itself a list of [text, hl_group] chunks
// -- via a raw nvim_buf_get_extmarks call, mirroring how noteCountExpr
// reaches into real nvim state that the ballroom_highlight.lua module
// doesn't expose a dedicated getter for (unlike note_count()): wrapping
// is a rendering detail worth locking in directly against nvim's own
// extmark storage, not worth adding module surface just to test.
func noteVirtLines(t *testing.T, ctx context.Context, socket string) [][]string {
	t.Helper()
	extmarksExpr := `json_encode(nvim_buf_get_extmarks(bufnr('%'), nvim_create_namespace('ballroom_tutor'), 0, -1, {'details': v:true}))`
	out, err := remoteExpr(ctx, socket, extmarksExpr)
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}

	var marks []json.RawMessage
	if err := json.Unmarshal([]byte(out), &marks); err != nil {
		t.Fatalf("unmarshal extmarks %q: %v", out, err)
	}
	if len(marks) != 1 {
		t.Fatalf("got %d extmarks, want exactly 1", len(marks))
	}

	var tuple []json.RawMessage
	if err := json.Unmarshal(marks[0], &tuple); err != nil {
		t.Fatalf("unmarshal extmark tuple: %v", err)
	}
	if len(tuple) != 4 {
		t.Fatalf("extmark tuple has %d elements, want 4 (id, row, col, details)", len(tuple))
	}

	var details struct {
		VirtLines [][][2]string `json:"virt_lines"`
	}
	if err := json.Unmarshal(tuple[3], &details); err != nil {
		t.Fatalf("unmarshal extmark details %s: %v", tuple[3], err)
	}

	rows := make([][]string, len(details.VirtLines))
	for i, row := range details.VirtLines {
		for _, chunk := range row {
			rows[i] = append(rows[i], chunk[0])
		}
	}
	return rows
}

// TestApplyHighlight_WrapsLongNoteAcrossMultipleLinesInsteadOfCuttingItOff
// covers the "it cuts off and it's hard to see" report: a single-line
// eol virtual text extmark never wraps in nvim -- it just gets clipped
// at the window's right edge -- so any note longer than the editor
// pane's width was silently unreadable. Notes now render as their own
// virtual lines below the highlighted line, word-wrapped to the window
// width, so a long note is fully visible instead of cut off mid-sentence.
func TestApplyHighlight_WrapsLongNoteAcrossMultipleLinesInsteadOfCuttingItOff(t *testing.T) {
	socket := startTestNvim(t)
	ctx := context.Background()

	longNote := "This line has a typo: 'Falsee' should be 'True', and the function logic is inverted: when you find a duplicate you should return True, not False"
	expr := highlightExpr("test.txt", 1, 1, longNote)
	if _, err := remoteExpr(ctx, socket, expr); err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}

	rows := noteVirtLines(t, ctx, socket)
	if len(rows) < 2 {
		t.Fatalf("got %d virtual-text rows, want >= 2 -- a long note should wrap across multiple lines", len(rows))
	}

	winWidth, err := remoteExpr(ctx, socket, "winwidth(0)")
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}
	maxWidth, err := strconv.Atoi(winWidth)
	if err != nil {
		t.Fatalf("parse winwidth %q: %v", winWidth, err)
	}

	var rebuilt strings.Builder
	for i, row := range rows {
		text := strings.Join(row, "")
		if len(text) > maxWidth {
			t.Errorf("row %d is %d chars, want <= window width %d: %q", i, len(text), maxWidth, text)
		}
		if i > 0 {
			rebuilt.WriteByte(' ')
		}
		rebuilt.WriteString(strings.TrimSpace(text))
	}
	if !strings.Contains(rebuilt.String(), "inverted") {
		t.Errorf("rebuilt wrapped text = %q, want it to still contain the full note (nothing dropped, just wrapped)", rebuilt.String())
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

func TestOpenInSplitExpr_EscapesEmbeddedQuotes(t *testing.T) {
	got := openInSplitExpr("/tmp/it's/solution.py")
	want := "execute('vsplit ' . fnameescape('/tmp/it''s/solution.py'))"
	if got != want {
		t.Errorf("openInSplitExpr = %q, want %q", got, want)
	}
}

func TestOpenInSplit_EmptySocketReturnsNilNoError(t *testing.T) {
	// Same graceful-degradation contract as remoteExpr itself: no editor
	// pane attached is an expected state (e.g. driven headless), not a
	// command failure -- the revealed file is still on disk regardless.
	if err := OpenInSplit(context.Background(), "", "/workspace/reference/solution.py"); err != nil {
		t.Fatalf("OpenInSplit with empty socket: %v", err)
	}
}

func TestOpenInSplit_OpensFileAsNewSplitAndFocusesIt(t *testing.T) {
	socket := startTestNvim(t)
	ctx := context.Background()

	target := filepath.Join(t.TempDir(), "solution.py")
	if err := os.WriteFile(target, []byte("print('reference')\n"), 0o644); err != nil {
		t.Fatalf("write target file: %v", err)
	}

	if err := OpenInSplit(ctx, socket, target); err != nil {
		t.Fatalf("OpenInSplit: %v", err)
	}

	winCount, err := remoteExpr(ctx, socket, "winnr('$')")
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}
	if winCount != "2" {
		t.Errorf("winnr('$') = %s, want 2 (the original window plus the new split)", winCount)
	}

	// The split takes focus and shows the target file -- not left behind
	// in the original (now background) window.
	focused, err := remoteExpr(ctx, socket, "fnamemodify(bufname('%'), ':t')")
	if err != nil {
		t.Fatalf("remoteExpr: %v", err)
	}
	if focused != "solution.py" {
		t.Errorf("focused buffer = %q, want %q", focused, "solution.py")
	}
}
