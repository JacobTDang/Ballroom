# Mask Account Numbers

Support agents see customer account numbers on screen, but policy says
only the final four characters may ever be displayed. Anything too
short to mask meaningfully is shown as-is.

Write `mask(accounts: list[str]) -> list[str]` returning a new list in
which every string of 5 or more characters has all but its last 4
characters replaced with `*`. Strings shorter than 5 characters are
kept unchanged.

## Examples

    mask(["123456789"])         -> ["*****6789"]
    mask(["9604", "10011001"])  -> ["9604", "****1001"]
    mask([])                    -> []

## Constraints

- 0 <= len(accounts) <= 1_000
- 1 <= len(account) <= 32
- every character of every account is a digit
