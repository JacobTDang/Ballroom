#!/usr/bin/env bash
# Tests for the pure/testable parts of chat.sh: build_file_context (#22),
# the code-enforced comprehension check (#23), process_highlights
# (directive parsing + stripping, #24), and apply_highlight's RPC-string
# construction (#24) plus ballroom_highlight.toggle() (#25). Plain bash +
# assertions, no framework — sources chat.sh under its BASH_SOURCE guard
# (so `main`'s interactive loop never runs). Run directly or via CI.
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

run_build_file_context() {
  local workdir max_bytes
  workdir="$1"
  max_bytes="${2:-8000}"
  WORKDIR="$workdir" TUTOR_FILE_CONTEXT_MAX_BYTES="$max_bytes" \
    OLLAMA_HOST="unused" TUTOR_MODEL="unused" PRACTICE_TUTOR_MODE="full-assist" \
    bash -c "source '$SCRIPT_DIR/chat.sh'; build_file_context"
}

run_process_highlights() {
  local text="$1"
  NVIM_SOCKET="" bash -c "source '$SCRIPT_DIR/chat.sh'; process_highlights \"\$1\"" _ "$text"
}

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "--- build_file_context: reading and truncating the active solution file (#22) ---"

# 1: normal file, contents come back exactly
mkdir -p "$tmp/normal"
printf 'package main\n\nfunc uniqueFn() {}\n' >"$tmp/normal/solution.go"
got=$(run_build_file_context "$tmp/normal")
assert_eq "normal file returns exact contents" "$(printf 'package main\n\nfunc uniqueFn() {}')" "$got"

# 2: no solution file at all (sandbox mode) -> empty, no error
mkdir -p "$tmp/empty"
got=$(run_build_file_context "$tmp/empty")
rc=$?
assert_eq "missing file returns empty output" "" "$got"
if [ "$rc" -ne 0 ]; then
  echo "FAIL: missing file should exit 0, got $rc"
  fail_count=$((fail_count + 1))
else
  echo "PASS: missing file exits 0 (degrades gracefully)"
fi

# 3: oversized file gets truncated with a marker, and does NOT exceed the cap
mkdir -p "$tmp/big"
python3 -c "print('x' * 500)" >"$tmp/big/solution.py" 2>/dev/null || \
  perl -e 'print "x" x 500 . "\n"' >"$tmp/big/solution.py"
got=$(run_build_file_context "$tmp/big" 100)
case "$got" in
  *"...[truncated, file is"*"showing first 100"*) echo "PASS: oversized file includes truncation marker" ;;
  *)
    echo "FAIL: oversized file missing truncation marker, got: $got"
    fail_count=$((fail_count + 1))
    ;;
esac
content_before_marker="${got%%$'\n'...\[truncated*}"
if [ "${#content_before_marker}" -le 100 ]; then
  echo "PASS: truncated content does not exceed the byte cap"
else
  echo "FAIL: truncated content is ${#content_before_marker} bytes, want <= 100"
  fail_count=$((fail_count + 1))
fi

# 4: re-reads fresh each call — editing the file between calls changes the result
mkdir -p "$tmp/live"
printf 'first version\n' >"$tmp/live/solution.go"
first=$(run_build_file_context "$tmp/live")
printf 'second version, edited\n' >"$tmp/live/solution.go"
second=$(run_build_file_context "$tmp/live")
if [ "$first" = "$second" ]; then
  echo "FAIL: expected a fresh read to reflect the edit, got the same content both times"
  fail_count=$((fail_count + 1))
else
  echo "PASS: each call re-reads the file (edit between calls is reflected)"
fi

echo
echo "--- wants_comprehension_check: which modes run it (#23) ---"

run_wants_check() {
  bash -c "source '$SCRIPT_DIR/chat.sh'; wants_comprehension_check \"\$1\" && echo true || echo false" _ "$1"
}

for mode in hints-first full-assist unrecognized-mode; do
  got=$(run_wants_check "$mode")
  assert_eq "$mode wants a comprehension check" "true" "$got"
done
got=$(run_wants_check "syntax-only")
assert_eq "syntax-only does not want a comprehension check" "false" "$got"

echo
echo "--- run_comprehension_check: isolated request, not the user's real question (#23) ---"
# The whole point of enforcing this in code (see chat.sh's header comment)
# is that the request asking for the comprehension check must NOT contain
# the user's real, concrete question — that's what a prompt-only approach
# couldn't reliably prevent the model from answering directly instead of
# doing the check. Verify the actual request body chat.sh builds.
fake_curl_dir=$(mktemp -d)
cat >"$fake_curl_dir/curl" <<'FAKECURL'
#!/usr/bin/env bash
# Records the -d request body to $FAKE_CURL_REQ_FILE (overwritten each
# call — tests here only ever make one call at a time) and returns the
# canned response in $FAKE_CURL_RESPONSE_FILE.
payload="" prev=""
for a in "$@"; do
  [ "$prev" = "-d" ] && payload="$a"
  prev="$a"
done
printf '%s' "$payload" >"$FAKE_CURL_REQ_FILE"
cat "$FAKE_CURL_RESPONSE_FILE"
FAKECURL
chmod +x "$fake_curl_dir/curl"

fake_curl_req="$tmp/fake_req.json"
fake_curl_resp="$tmp/fake_resp.json"
printf '%s' '{"message":{"content":"Restated: you want the two-sum indices. What should happen with duplicate values?"}}' >"$fake_curl_resp"

run_comprehension_check_case() {
  PATH="$fake_curl_dir:$PATH" FAKE_CURL_REQ_FILE="$fake_curl_req" FAKE_CURL_RESPONSE_FILE="$fake_curl_resp" \
    PRACTICE_TUTOR_MODE="hints-first" OLLAMA_HOST="http://fake" TUTOR_MODEL="test-model" WORKDIR="$tmp/empty" \
    bash -c "
      source '$SCRIPT_DIR/chat.sh'
      messages=\$(jq -n --arg system \"\$SYSTEM_PROMPT\" '[{role: \"system\", content: \$system}]')
      run_comprehension_check 'this user question should never reach the request'
      echo \"MESSAGES_LEN:\$(jq 'length' <<<\"\$messages\")\"
    "
}

out=$(run_comprehension_check_case)
echo "$out" | grep -v '^MESSAGES_LEN:' # the printed check reply
req_body=$(cat "$fake_curl_req")

case "$req_body" in
  *"this user question should never reach the request"*)
    echo "FAIL: the isolated check request leaked the user's real question — this is exactly what the code-enforced approach must prevent"
    fail_count=$((fail_count + 1))
    ;;
  *) echo "PASS: the isolated check request does not contain the user's real question" ;;
esac

case "$req_body" in
  *"restate the problem"*) echo "PASS: the isolated check request asks for a restatement + clarifying questions" ;;
  *)
    echo "FAIL: the isolated check request is missing the comprehension-check instruction"
    fail_count=$((fail_count + 1))
    ;;
esac

case "$out" in
  *"Restated: you want the two-sum indices"*) echo "PASS: the check's reply is printed for the user to see" ;;
  *)
    echo "FAIL: the check's reply was not printed, got: $out"
    fail_count=$((fail_count + 1))
    ;;
esac

messages_len=$(echo "$out" | grep '^MESSAGES_LEN:' | cut -d: -f2)
assert_eq "the real user message + check reply are both persisted to history" "3" "$messages_len"

rm -rf "$fake_curl_dir"

echo
echo "--- process_highlights: parsing and stripping (#24) ---"

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
    echo "--- ballroom_highlight.toggle(): hide/show without deleting data (#25) ---"
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
