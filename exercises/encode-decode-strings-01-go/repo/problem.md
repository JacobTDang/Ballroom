# Encode and Decode Strings

Design an algorithm to encode a list of strings to a single string. The
encoded string is then decoded back to the original list of strings.

Your encode/decode pair must round-trip **exactly** — every character in
every string must survive, including digits, spaces, commas, and
whatever delimiter you choose to use internally. A string in the input
list might itself contain your delimiter character; your encoding must
not be confused by that.

## Example

```
Input: ["neet","code","love","you"]
Output after encode -> decode: ["neet","code","love","you"]
```

There's no single "correct" wire format for what encode produces — only
that `decode(encode(strs)) == strs` for any list of strings, including
ones containing whatever delimiter character your own encoding uses.

## Constraints

- `1 <= strs.length <= 200`
- `0 <= strs[i].length <= 200`
- `strs[i]` may contain any character, including digits and whatever
  delimiter character your own encoding uses.
