# Backwards Ties

`sort_leaderboard` orders leaderboard entries by score, highest first;
when two entries tie on score, it's supposed to break the tie
alphabetically by name so the board is deterministic. Score ordering
works fine. Name ordering doesn't: whenever two or more entries land
on the same score, they come out in the wrong order relative to each
other.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
sort_leaderboard([("bob", 60), ("dan", 100), ("cara", 75), ("amy", 90)])
  => [("dan", 100), ("amy", 90), ("cara", 75), ("bob", 60)]

sort_leaderboard([("erin", 90), ("cara", 75), ("amy", 90), ("bob", 90)])
  => [("amy", 90), ("bob", 90), ("erin", 90), ("cara", 75)]
```

## Constraints

- Scores are integers and may be negative.
- The result is a full ordering: highest score first; among entries
  with equal scores, names in ascending alphabetical order.
