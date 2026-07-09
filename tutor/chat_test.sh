#!/usr/bin/env bash
# Tests for the pure/testable parts of chat.sh: process_highlights (directive
# parsing + stripping) and apply_highlight's RPC-string construction. Plain
# bash + assertions, no framework — sources chat.sh under its BASH_SOURCE
# guard (so `main`'s interactive loop never runs). Run directly or via CI.
#
# apply_highlight's own RPC call is exercised against a REAL nvim --listen
# instance when `nvim` is available, specifically to keep the single-quote
# escaping regression-tested against real Vimscript parsing rather than
# just asserting on the constructed string — this escaping is what stands
# between untrusted (LLM-controlled) directive content and nvim's
# `--remote-expr` evaluator, so a change here is security-relevant, not
# just cosmetic. When `nvim` isn't available, those cases are skipped
# with a clear note rather than silently passing.
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fail_count=0
skip_count=0

assert_eq() {
  local desc="$1" want="$2" got="$3"
  if [ "$want" != "$got" ]; then
    echo "FAIL: $desc"
    echo "  want: $(printf '%q' "$want")"
    echo "  got:  $(printf '%q' "$got")"
    fail_count=$((fail_count + 1))
  else
    echo "PASS: $desc"
  fi
}

run_process_highlights() {
  local text="$1"
  NVIM_SOCKET="" bash -c "source '$SCRIPT_DIR/chat.sh'; process_highlights \"\$1\"" _ "$text"
}

echo "--- process_highlights: parsing and stripping ---"

# 1: a well-formed directive is stripped from the displayed reply
reply='Take a look at this. <<<highlight file=solution.go line=10-14 note="off-by-one here">>> Does that help?'
got=$(run_process_highlights "$reply")
assert_eq "well-formed directive is stripped from the reply" \
  "Take a look at this.  Does that help?" "$got"

# 2: multiple directives across a reply all get stripped
reply='First point <<<highlight file=a.go line=1 note="one">>> and second <<<highlight file=a.go line=2-3 note="two">>> done.'
got=$(run_process_highlights "$reply")
assert_eq "multiple directives are all stripped" \
  "First point  and second  done." "$got"

# 3: a malformed directive (missing note) is left in place, not silently eaten,
# and does not crash the function (set -e safe)
reply='See <<<highlight file=solution.go line=10-14>>> for details.'
got=$(run_process_highlights "$reply")
assert_eq "malformed directive (no note) is not matched/stripped" "$reply" "$got"

# 4: plain text with no directive at all passes through unchanged
reply='Nothing special about this reply.'
got=$(run_process_highlights "$reply")
assert_eq "reply with no directive passes through unchanged" "$reply" "$got"

# 5: empty reply doesn't crash
got=$(run_process_highlights "")
assert_eq "empty reply returns empty" "" "$got"

# 6: a non-numeric line value never matches the directive regex at all —
# this is what guarantees apply_highlight's start/end (interpolated
# unquoted into the RPC expr) can only ever be pure digits in the normal
# flow, not attacker-controlled expression text.
reply='See <<<highlight file=solution.go line=ten note="bad">>> for details.'
got=$(run_process_highlights "$reply")
assert_eq "non-numeric line value is never matched/stripped" "$reply" "$got"

echo
echo "--- apply_highlight: escaping/RPC safety (requires a real nvim, see below) ---"
# Deliberately not testing this by re-deriving the constructed expr string
# in the test harness itself — a second implementation of the same
# escaping logic can drift from (or subtly differ from, e.g. via an extra
# quoting layer) the real one and either mask a real bug or manufacture a
# fake one. Since this escaping is what stands between untrusted
# (LLM-controlled) directive content and nvim's --remote-expr evaluator,
# what actually matters is real behavior against a real nvim instance —
# see the live round-trip checks below.

if command -v nvim >/dev/null 2>&1; then
  echo
  echo "--- apply_highlight: live nvim RPC round-trip (nvim found, running real checks) ---"
  sock="$(mktemp -u)"
  workdir="$(mktemp -d)"
  # A bare `nvim --headless` won't have docker/nvim/lua/ballroom_highlight.lua
  # loaded — it only gets required via docker/nvim/init.lua (see that file
  # and the Dockerfile's COPY of docker/nvim/ into the real container's
  # config dir). Point a scratch XDG_CONFIG_HOME at the real config here so
  # this test exercises the actual module, not a bare nvim where every RPC
  # call would fail on "module not found" — a failure apply_highlight's own
  # error handling swallows, which would silently make these checks pass
  # for the wrong reason.
  nvim_config_home="$(mktemp -d)"
  mkdir -p "$nvim_config_home/nvim"
  cp -r "$SCRIPT_DIR/../docker/nvim/." "$nvim_config_home/nvim/"
  printf 'package main\n\nfunc main() {}\n' >"$workdir/solution.go"
  rm -f "$sock"
  XDG_CONFIG_HOME="$nvim_config_home" nvim --headless --listen "$sock" "$workdir/solution.go" >/tmp/chat_test_nvim.log 2>&1 &
  nvim_pid=$!
  # shellcheck disable=SC2064 # intentional: expand these now, not at trap time
  trap "kill '$nvim_pid' 2>/dev/null; rm -rf '$workdir' '$sock' '$nvim_config_home'" EXIT

  for _ in 1 2 3 4 5; do
    [ -S "$sock" ] && break
    sleep 0.2
  done

  if [ -S "$sock" ]; then
    # Legit call must succeed and actually paint the extmark — proves the
    # escaping doesn't just neutralize attacks, it also doesn't break
    # normal usage.
    NVIM_SOCKET="$sock" bash -c "source '$SCRIPT_DIR/chat.sh'; apply_highlight 'solution.go' 1 1 \"it's a legit note\"" >/dev/null 2>&1
    extmark_count=$(nvim --server "$sock" --remote-expr \
      "luaeval('#vim.api.nvim_buf_get_extmarks(0, vim.api.nvim_create_namespace(\"ballroom_tutor\"), 0, -1, {})')" 2>/dev/null)
    if [ "${extmark_count:-0}" -ge 1 ]; then
      echo "PASS: legit apostrophe-containing call actually renders a highlight/note"
    else
      echo "FAIL: legit call did not render any extmark (got count=${extmark_count:-<none>})"
      fail_count=$((fail_count + 1))
    fi

    # Malicious payloads must NOT execute — verified by absence of a
    # canary file that would only exist if injected Lua/Vimscript ran.
    canary="$workdir/PWNED"
    NVIM_SOCKET="$sock" bash -c "source '$SCRIPT_DIR/chat.sh'; apply_highlight \"x') os.execute('touch $canary') --\" 1 1 note" >/dev/null 2>&1 || true
    if [ -e "$canary" ]; then
      echo "FAIL: injection payload executed — found $canary"
      fail_count=$((fail_count + 1))
    else
      echo "PASS: injection payload via file field did not execute"
    fi

    echo
    echo "--- ballroom_highlight.toggle(): hide/show without deleting data (issue #25) ---"
    extmark_count() {
      nvim --server "$sock" --remote-expr \
        "luaeval('#vim.api.nvim_buf_get_extmarks(0, vim.api.nvim_create_namespace(\"ballroom_tutor\"), 0, -1, {})')" 2>/dev/null
    }
    toggle() {
      nvim --server "$sock" --remote-expr "v:lua.require('ballroom_highlight').toggle()" 2>&1
    }

    before=$(extmark_count)
    state=$(toggle)
    after_hide=$(extmark_count)
    if [ "$state" = "hidden" ] && [ "${after_hide:-1}" -eq 0 ]; then
      echo "PASS: toggle() hides rendering (extmarks cleared, state=hidden)"
    else
      echo "FAIL: toggle() did not hide correctly (state=$state, extmark_count=$after_hide)"
      fail_count=$((fail_count + 1))
    fi

    state=$(toggle)
    after_show=$(extmark_count)
    if [ "$state" = "shown" ] && [ "$after_show" -eq "$before" ]; then
      echo "PASS: toggling back on replays the exact same stored note(s) (extmark count restored, state=shown)"
    else
      echo "FAIL: toggle() back on did not restore correctly (state=$state, extmark_count=$after_show, want=$before)"
      fail_count=$((fail_count + 1))
    fi

    # Add a highlight while hidden: must be stored, not rendered, until
    # toggled back on — this is the "does not delete the underlying
    # notes/data" acceptance criterion, verified against a second note
    # rather than just the first.
    toggle >/dev/null # -> hidden
    NVIM_SOCKET="$sock" bash -c "source '$SCRIPT_DIR/chat.sh'; apply_highlight 'solution.go' 2 2 'added while hidden'" >/dev/null 2>&1
    count_while_hidden=$(extmark_count)
    toggle >/dev/null # -> shown
    count_after_reshow=$(extmark_count)
    if [ "${count_while_hidden:-1}" -eq 0 ] && [ "$count_after_reshow" -gt "$before" ]; then
      echo "PASS: a note added while hidden stays hidden, then appears once toggled back on"
    else
      echo "FAIL: note-added-while-hidden case wrong (hidden count=$count_while_hidden, reshow count=$count_after_reshow, want > $before)"
      fail_count=$((fail_count + 1))
    fi
  else
    echo "SKIP: nvim RPC socket never came up, skipping live RPC checks"
    skip_count=$((skip_count + 1))
  fi
else
  echo
  echo "SKIP: nvim not found in PATH — skipping live RPC round-trip checks"
  skip_count=$((skip_count + 1))
fi

echo
if [ "$skip_count" -gt 0 ]; then
  echo "$skip_count check group(s) skipped (nvim unavailable)."
fi
if [ "$fail_count" -eq 0 ]; then
  echo "All tests passed."
  exit 0
else
  echo "$fail_count test(s) failed."
  exit 1
fi
