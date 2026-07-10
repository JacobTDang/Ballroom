# Implement Trie (Prefix Tree)

A trie (pronounced as "try") or prefix tree is a tree data structure
used to efficiently store and retrieve keys in a dataset of strings.

Implement the `Trie` class:

- `Trie()` initializes the trie object.
- `Insert(word string)` inserts the string `word` into the trie.
- `Search(word string) bool` returns `true` if the string `word` is
  in the trie (i.e., was inserted before), and `false` otherwise.
- `StartsWith(prefix string) bool` returns `true` if there is a
  previously inserted string `word` that has `prefix` as a prefix.

## Example

```
Input:
["Trie", "insert", "search", "search", "startsWith", "insert", "search"]
[[], ["apple"], ["apple"], ["app"], ["app"], ["app"], ["app"]]

Output:
[null, null, true, false, true, null, true]

Explanation:
Trie trie = new Trie();
trie.insert("apple");
trie.search("apple");   // return true
trie.search("app");     // return false
trie.startsWith("app"); // return true
trie.insert("app");
trie.search("app");     // return true
```

## Constraints

- `1 <= word.length, prefix.length <= 2000`
- `word` and `prefix` consist only of lowercase English letters.
- At most `3 * 10^4` calls in total will be made to `Insert`,
  `Search`, and `StartsWith`.
