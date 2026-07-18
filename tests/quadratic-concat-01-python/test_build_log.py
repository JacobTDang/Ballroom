import time

from solution import build_log


# Tier 1: exact small-n correctness.
def test_reverses_three_chunks():
    assert build_log(["c", "b", "a"]) == "abc"


def test_single_chunk():
    assert build_log(["single"]) == "single"


def test_empty_input():
    assert build_log([]) == ""


def test_four_chunks():
    assert build_log(["d", "c", "b", "a"]) == "abcd"


def test_multi_char_chunks():
    assert build_log(["world", "hello"]) == "helloworld"


def test_repeated_chunks():
    assert build_log(["x", "x", "x"]) == "xxx"


# Tier 2: an in-test stopwatch, not a harness timeout or a pytest
# plugin. n=200,000 chunks of 15 bytes each (~3MB total copying for a
# correct, linear-time build):
#   fixed:     ~3,000,000 bytes copied total -- well under 100ms even
#              at CPython's interpreted speed -- >=100x headroom under
#              the 10s bound.
#   quadratic: ~15 * 200,000^2 / 2 =~ 3*10^11 bytes (~300GB) of
#              copying -- tens of seconds (20-40s observed) -- >=2x
#              OVER the 10s bound.
# The 10s cutoff sits comfortably between those two, so this cannot
# flake based on machine speed; it only distinguishes O(n) from
# O(n^2). time.monotonic() is immune to wall-clock adjustments.
def test_large_input_finishes_within_ten_seconds():
    n = 200_000
    chunks = [f"{i:014d};" for i in range(n)]  # 15 bytes each

    start = time.monotonic()
    result = build_log(chunks)
    elapsed = time.monotonic() - start

    assert len(result) == n * 15
    assert result[:15] == chunks[-1]
    assert result[-15:] == chunks[0]
    assert elapsed < 10, f"build_log took {elapsed:.1f}s on {n} chunks -- looks quadratic"
