# Dedupe by Identity

`dedupe` is supposed to remove duplicate records from a list -- two
records are duplicates if they have the same key and the same value,
even if they came from different sources and are different objects in
memory. It works when the exact same object appears twice, but records
that are merely *equal* -- same key, same value, different object --
both survive.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
dedupe([Record("a", 1), Record("a", 1), Record("b", 2)])
  => [Record("a", 1), Record("b", 2)]

dedupe([Record("a", 1), Record("a", 2)])
  => [Record("a", 1), Record("a", 2)]   # same key, different value -- not a duplicate
```

## Constraints

- A record's identity for deduping purposes is its (key, value) pair.
- The first occurrence of each distinct (key, value) pair is kept;
  later duplicates are dropped.
- Output order otherwise follows input order.
