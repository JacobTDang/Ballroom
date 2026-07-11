package tutor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// escapeVimSingleQuoted doubles every embedded single quote so s is safe
// to interpolate into a VimL single-quoted string literal — the only
// escape VimL single-quoted strings need (nothing else is special inside
// them). Port of tutor/chat.sh's apply_highlight escaping.
func escapeVimSingleQuoted(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// highlightExpr builds the --remote-expr call into
// docker/nvim/lua/ballroom_highlight.lua's add_highlight, with file/note
// safely escaped — both are model-controlled, so injection-safety here
// is load-bearing (see nvimrpc_test.go's live-nvim injection tests).
func highlightExpr(file string, start, end int, note string) string {
	return fmt.Sprintf(
		"v:lua.require('ballroom_highlight').add_highlight('%s', %d, %d, '%s')",
		escapeVimSingleQuoted(file), start, end, escapeVimSingleQuoted(note),
	)
}

// cursorPositionExpr builds the --remote-expr call returning the
// currently-focused window's cursor position as JSON. Entirely static —
// no model/user-controlled value is ever interpolated into it — so
// unlike highlightExpr there is no injection surface and no escaping is
// needed.
func cursorPositionExpr() string {
	return `json_encode({'file': expand('%:t'), 'line': line('.'), 'col': col('.'), 'total_lines': line('$')})`
}

// remoteExpr evaluates expr in the nvim instance listening on socket via
// `nvim --server socket --remote-expr expr` and returns its raw string
// result. Returns ("", nil) — not an error — when socket is empty or
// isn't a real Unix socket: a missing/unreachable editor pane is an
// expected state in some runs (e.g. sandbox mode without an editor
// attached yet), not a bug. A genuine RPC failure (nvim reachable but the
// call itself errored) IS returned as an error, so it flows through
// utils.WrapToolWithErrorHandler back to the model as a visible failure
// rather than silently no-op'ing.
func remoteExpr(ctx context.Context, socket, expr string) (string, error) {
	if socket == "" {
		return "", nil
	}
	info, err := os.Stat(socket)
	if err != nil || info.Mode()&os.ModeSocket == 0 {
		return "", nil
	}

	cmd := exec.CommandContext(ctx, "nvim", "--server", socket, "--remote-expr", expr)
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("tutor: nvim remote-expr failed: %v: %s", err, errOut.String())
	}
	return strings.TrimRight(out.String(), "\n"), nil
}
