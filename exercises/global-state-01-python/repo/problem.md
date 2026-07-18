# A report that remembers too much: `GenerateReport`

`GenerateReport` formats a list of low-stock item names into report
lines, one per item ("LOW STOCK: <item>"), in the order given.

Call it once and the report looks right. Call it again right after
with different items, and the new report still contains lines from the
first call. Call it a third time with an item from the very first
call, and that line shows up more than once.

## Task

Find and fix the bug. Don't change the function's signature.

## Examples

```
GenerateReport(["apples"])   => ["LOW STOCK: apples"]
GenerateReport(["bananas"])  => ["LOW STOCK: bananas"]
```

Each call's report should only ever contain lines for that call's
items -- regardless of what was reported before it.

## Constraints

- Items may repeat within a single call.
- The function may be called any number of times, in any sequence.
