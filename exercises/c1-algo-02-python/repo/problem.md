# Balanced Statement Brackets

A statement renderer wraps adjustments in parentheses, but an upstream
bug sometimes drops or duplicates characters, leaving strings of `(`
and `)` that no longer balance. Before rendering, the system deletes as
few characters as possible so that the remaining string is balanced:
every `(` is closed by a later `)`, and every `)` is opened by an
earlier `(`.

Write `min_removals(s: str) -> int` returning the minimum number of
characters that must be removed to leave a balanced string.

## Examples

    min_removals("(()")  -> 1
    min_removals("())(") -> 2
    min_removals("()()") -> 0

## Constraints

- 0 <= len(s) <= 100_000
- s contains only the characters '(' and ')'
