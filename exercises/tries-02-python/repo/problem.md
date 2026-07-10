# Design Add and Search Words Data Structure

Design a data structure that supports adding new words and finding
if a string matches any previously added string.

Implement the `WordDictionary` class:

- `WordDictionary()` initializes the object.
- `AddWord(word string)` adds `word` to the data structure, it can
  be matched later.
- `Search(word string) bool` returns `true` if there is any string in
  the data structure that matches `word` or `false` otherwise. `word`
  may contain dots `'.'` where dots can be matched with any letter.

## Example

```
Input:
["WordDictionary","addWord","addWord","addWord","search","search","search","search"]
[[],["bad"],["dad"],["mad"],["pad"],["bad"],[".ad"],["b.."]]

Output:
[null,null,null,null,false,true,true,true]

Explanation:
WordDictionary wordDictionary = new WordDictionary();
wordDictionary.addWord("bad");
wordDictionary.addWord("dad");
wordDictionary.addWord("mad");
wordDictionary.search("pad"); // return false
wordDictionary.search("bad"); // return true
wordDictionary.search(".ad"); // return true
wordDictionary.search("b.."); // return true
```

## Constraints

- `1 <= word.length <= 25`
- `word` in `AddWord` consists of lowercase English letters.
- `word` in `Search` consists of `'.'` or lowercase English letters.
- There will be at most `2` dots in `word` for `Search` queries.
- At most `10^4` calls will be made to `AddWord` and `Search`.
