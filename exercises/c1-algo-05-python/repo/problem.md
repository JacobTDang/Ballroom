# Statement Compression

Archived statement lines often contain long stretches of repeated
characters, so the archiver run-length encodes each line: every maximal
run of a repeated character becomes the character followed by the run
length, except that a run of length 1 keeps just the character — a
count of `1` is never written. If the encoded line is not strictly
shorter than the original, the original line is stored unchanged.

Write `compress(s: str) -> str` returning the stored form of `s`.

## Examples

    compress("aaab")        -> "a3b"
    compress("aabcccccaaa") -> "a2bc5a3"
    compress("abc")         -> "abc"

## Constraints

- 0 <= len(s) <= 100_000
- s contains only lowercase letters
