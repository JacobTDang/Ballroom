# Glob Matcher

`Match(pattern, s)`: `*` matches any run (including empty), `?` matches
exactly one character, `[a-cx]` matches one character from a set (with
ranges). Whole-string matching, built by hand — **no regex library**;
the two-pointer-with-backtracking loop *is* the exercise.

An unclosed `[` makes the pattern invalid: it matches nothing.

## The invariant the tests enforce

Classic star backtracking (`a*b*c` against `aXXbYYc`), empty-string
edges (`*` matches `""`), exact-one semantics for `?`, character-class
ranges and misses, and the unclosed-class rule — a table of exact
cases.

API: `bool Match(const std::string& pattern, const std::string& s)`. `<regex>` is off-limits.

Think: when a `*` match fails later, where do you resume — and what do
you have to remember to resume there?
