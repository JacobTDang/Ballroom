# Partition Labels

You are given a string `s`. We want to partition the string into as
many parts as possible so that each letter appears in at most one
part.

Return a list of integers representing the size of these parts, in
the order they appear in `s`.

## Example

```
Input: s = "ababcbacadefegdehijhklij"
Output: [9,7,8]
Explanation: The partition is "ababcbaca", "defegde", "hijhklij".
This is a partition so that each letter appears in at most one part.
A partition like "ababcbacadefegde", "hijhklij" is incorrect, because
it splits s into less parts.
```

## Constraints

- `1 <= s.length <= 500`
- `s` consists of lowercase English letters.
