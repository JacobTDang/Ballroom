# Word Ladder

You are given two words, `beginWord` and `endWord`, and a dictionary
`wordList`. Return the number of words in the shortest transformation
sequence from `beginWord` to `endWord`, such that:

- Only one letter can be changed at a time.
- Each transformed word must exist in `wordList`, including
  `endWord`.

Return `0` if there is no such transformation sequence. Note that
`beginWord` does not need to be in `wordList`.

## Example

```
Input: beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log","cog"]
Output: 5
Explanation: hit -> hot -> dot -> dog -> cog, which is 5 words.

Input: beginWord = "hit", endWord = "cog", wordList = ["hot","dot","dog","lot","log"]
Output: 0
Explanation: endWord "cog" is not in wordList, so no valid transformation exists.
```

## Constraints

- `1 <= beginWord.length <= 10`
- `endWord.length == beginWord.length`
- `1 <= wordList.length <= 5000`
- `wordList[i].length == beginWord.length`
- `beginWord`, `endWord`, and `wordList[i]` consist of lowercase
  English letters.
- `beginWord != endWord`
- All words in `wordList` are unique.
