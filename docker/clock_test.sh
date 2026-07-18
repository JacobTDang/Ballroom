#!/usr/bin/env bash
# Smoke tests for clock.sh's color transitions -- teal (>=10 min left),
# gold (<10 min left), red TIME UP (deadline passed) -- using the
# optional injected-"now" argument so this runs in well under a second
# instead of waiting out a real countdown.
#
# No shell test framework exists in this repo (see README's Development
# section for the Go equivalent); this is a plain pass/fail script run
# directly:
#   docker/clock_test.sh
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CLOCK="$SCRIPT_DIR/clock.sh"

fail=0

assert_contains() {
  local desc="$1" haystack="$2" needle="$3"
  if [[ "$haystack" != *"$needle"* ]]; then
    echo "FAIL: $desc"
    echo "  expected to contain: $needle"
    echo "  got:                 $haystack"
    fail=1
  else
    echo "ok: $desc"
  fi
}

# clock.sh takes deadline/now as container uptime (seconds since this
# container's kernel booted -- /proc/uptime's first field), not wall
# clock epoch seconds (issue #229) -- these are plain small numbers, not
# anything tied to a real clock, which is exactly what makes the
# injected-"now" argument usable for deterministic tests.

# Deadline at uptime 2000s, "now" at 1000s -> 1000s (16:40) remaining, teal.
out="$("$CLOCK" 2000 1000)"
assert_contains "teal when far from deadline" "$out" '#[fg=#2FA6A6]16:40#[default]'

# "now" at 1750s -> 250s (04:10) remaining, under the 600s gold threshold.
out="$("$CLOCK" 2000 1750)"
assert_contains "gold under 10 minutes remaining" "$out" '#[fg=#E8A93C]04:10#[default]'

# The 600s boundary itself: strictly-less-than, so 599 remaining is gold
# and 600 remaining is still the last teal second.
out="$("$CLOCK" 2000 1401)"
assert_contains "599s remaining is gold" "$out" '#[fg=#E8A93C]09:59#[default]'
out="$("$CLOCK" 2000 1400)"
assert_contains "600s remaining is still teal" "$out" '#[fg=#2FA6A6]10:00#[default]'

# "now" at or past the deadline -> TIME UP.
out="$("$CLOCK" 2000 2000)"
assert_contains "deadline exactly reached" "$out" 'TIME UP'
out="$("$CLOCK" 2000 2500)"
assert_contains "deadline passed" "$out" 'TIME UP'

# Fractional /proc/uptime-shaped inputs (the whole point of this fix --
# clock.sh now reads two uptime readings, which carry decimals, not two
# epoch-second integers) must truncate to whole seconds rather than
# erroring out under bash's integer-only arithmetic.
out="$("$CLOCK" 2000.87 1000.12)"
assert_contains "fractional deadline/now truncate to whole seconds" "$out" '#[fg=#2FA6A6]16:40#[default]'

# No injected "now": falls through to reading /proc/uptime for real.
# Deadline of 1 second is guaranteed already passed on any host that's
# been up more than 1 second, so this is deterministic without depending
# on the actual uptime value. Guarded on /proc/uptime existing so this
# script still runs end-to-end on a macOS dev machine (no procfs there);
# the container (this repo's real target) always has it.
if [ -r /proc/uptime ]; then
  out="$("$CLOCK" 1)"
  assert_contains "falls back to real /proc/uptime when now is omitted" "$out" 'TIME UP'
fi

# Usage error when called with no arguments at all.
if out="$("$CLOCK" 2>&1)"; then
  echo "FAIL: clock.sh with no args should exit non-zero"
  fail=1
else
  assert_contains "usage message on missing deadline arg" "$out" 'usage:'
fi

if [ "$fail" -ne 0 ]; then
  echo "clock_test.sh: FAILED"
  exit 1
fi
echo "clock_test.sh: all checks passed"
