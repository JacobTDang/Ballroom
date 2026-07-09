#!/usr/bin/env bash
# Tests for the pure/testable parts of chat.sh: build_file_context (#22)
# and the code-enforced comprehension check (#23). Plain bash + assertions,
# no framework — sources chat.sh under its BASH_SOURCE guard (so `main`'s
# interactive loop never runs). Run directly or via CI.
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

run_build_file_context() {
  local workdir max_bytes
  workdir="$1"
  max_bytes="${2:-8000}"
  WORKDIR="$workdir" TUTOR_FILE_CONTEXT_MAX_BYTES="$max_bytes" \
    OLLAMA_HOST="unused" TUTOR_MODEL="unused" PRACTICE_TUTOR_MODE="full-assist" \
    bash -c "source '$SCRIPT_DIR/chat.sh'; build_file_context"
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
if [ "$fail_count" -eq 0 ]; then
  echo "All tests passed."
  exit 0
else
  echo "$fail_count test(s) failed."
  exit 1
fi
