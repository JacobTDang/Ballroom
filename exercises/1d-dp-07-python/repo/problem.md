# Decode Ways

A message containing letters from `A-Z` can be encoded into numbers
using the mapping:

```
'A' -> "1"
'B' -> "2"
...
'Z' -> "26"
```

Given a string `s` containing only digits, return the number of ways
to decode it. A leading zero in a group (like `"06"`) is never a valid
encoding, and a lone `'0'` cannot be decoded on its own.

## Example

```
Input: s = "226"
Output: 3
Explanation: "226" can be decoded as "BZ" (2 26), "VF" (22 6), or "BBF"
(2 2 6).
```

## Constraints

- `1 <= s.length <= 100`
- `s` consists of digits only
