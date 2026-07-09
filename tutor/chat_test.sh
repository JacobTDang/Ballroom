#!/usr/bin/env bash
# Tests for the pure/testable parts of chat.sh (build_file_context).
# Plain bash + assertions, no framework — sources chat.sh under its
# BASH_SOURCE guard (so `main`'s interactive loop never runs) and drives
# the function directly against real temp files. Run directly or via CI.
set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
fail_count=0

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

run_case() {
  local desc="$1"
  shift
  local workdir max_bytes
  workdir="$1"
  max_bytes="${2:-8000}"
  WORKDIR="$workdir" TUTOR_FILE_CONTEXT_MAX_BYTES="$max_bytes" \
    OLLAMA_HOST="unused" TUTOR_MODEL="unused" PRACTICE_TUTOR_MODE="full-assist" \
    bash -c "source '$SCRIPT_DIR/chat.sh'; build_file_context"
}

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

# 1: normal file, contents come back exactly
mkdir -p "$tmp/normal"
printf 'package main\n\nfunc uniqueFn() {}\n' >"$tmp/normal/solution.go"
got=$(run_case "normal file" "$tmp/normal")
assert_eq "normal file returns exact contents" "$(printf 'package main\n\nfunc uniqueFn() {}')" "$got"

# 2: no solution file at all (sandbox mode) -> empty, no error
mkdir -p "$tmp/empty"
got=$(run_case "no file" "$tmp/empty")
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
got=$(run_case "oversized file" "$tmp/big" 100)
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
first=$(run_case "before edit" "$tmp/live")
printf 'second version, edited\n' >"$tmp/live/solution.go"
second=$(run_case "after edit" "$tmp/live")
if [ "$first" = "$second" ]; then
  echo "FAIL: expected a fresh read to reflect the edit, got the same content both times"
  fail_count=$((fail_count + 1))
else
  echo "PASS: each call re-reads the file (edit between calls is reflected)"
fi

echo
if [ "$fail_count" -eq 0 ]; then
  echo "All tests passed."
  exit 0
else
  echo "$fail_count test(s) failed."
  exit 1
fi
