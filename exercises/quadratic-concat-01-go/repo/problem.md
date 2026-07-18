# Slow log rebuild: `BuildLog`

A log shipper reads a page of log chunks from storage newest-chunk-first
(a common "tail the log backward" read pattern). `BuildLog` is supposed
to return the full log text in the opposite order -- oldest chunk first,
so it reads top-to-bottom in the order things actually happened -- by
concatenating the chunks directly (each chunk already carries its own
formatting, so there's no separator to insert).

It returns the right text on a small page. On a full page from a busy
service it takes far longer than it has any right to.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
BuildLog(["c", "b", "a"])  => "abc"
BuildLog(["single"])       => "single"
BuildLog([])                => ""
```

## Constraints

- Chunks may repeat or be empty strings.
- Must stay fast from an empty page up to a page of hundreds of
  thousands of chunks.
