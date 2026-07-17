# JSON Subset Parser

A real recursive-descent parser for a JSON subset: objects, arrays,
strings (with `\"` and `\\` escapes only), integers (with negatives),
`true`, `false`, `null`, and whitespace between tokens.

Every error names its **byte position**: unterminated strings, missing
colons, bad literals (`tru`), unsupported escapes, and trailing garbage
after the value are all errors — a parser that guesses is worse than no
parser.

## The invariant the tests enforce

Exact parsed structure for nested documents, each escape rule, each
error case with its position, and the trailing-garbage rule.

API: `Parse(input string) (any, error)` — objects are `map[string]any`, arrays `[]any`, strings `string`, integers `int`, booleans `bool`, null `nil`. The encoding/json package is off-limits.

Think: one function per grammar rule (`parseValue`, `parseObject`,
`parseArray`, `parseString`, `parseNumber`), each leaving the cursor
just past what it consumed.
