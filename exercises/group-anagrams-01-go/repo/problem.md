# Group Anagrams

Given an array of strings `strs`, group the anagrams together. You can
return the answer in **any order** — both which group comes first, and
the order of strings within each group.

An anagram is a string formed by rearranging the letters of another
string, using all the original letters exactly once.

## Examples

```
Input: strs = ["eat","tea","tan","ate","nat","bat"]
Output: [["bat"],["nat","tan"],["ate","eat","tea"]]
```

```
Input: strs = [""]
Output: [[""]]
```

```
Input: strs = ["a"]
Output: [["a"]]
```

## Constraints

- `1 <= strs.length <= 10^4`
- `0 <= strs[i].length <= 100`
- `strs[i]` consists of lowercase English letters.
