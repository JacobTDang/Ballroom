# Alien Dictionary

There is a new alien language that uses the English alphabet. However, the
order among the letters is unknown to you.

You are given a list of strings `words` from the alien language's
dictionary, where the strings in `words` are **sorted lexicographically**
by the rules of this new language.

Derive the order of letters in this language, and return it as a string
of unique letters, with any valid ordering acceptable. If there is no
solution (the input is invalid, or the constraints don't define a valid
ordering), return `""`.

## Example 1

```
Input: words = ["wrt","wrf","er","ett","rftt"]
Output: "wertf"
```

## Example 2

```
Input: words = ["z","x"]
Output: "zx"
```

## Example 3

```
Input: words = ["z","x","z"]
Output: ""
Explanation: The order is invalid, so return "".
```

## Constraints

- `1 <= words.length <= 100`
- `1 <= words[i].length <= 100`
- `words[i]` consists of only lowercase English letters.
