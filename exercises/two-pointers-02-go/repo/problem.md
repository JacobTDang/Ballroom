# Valid Palindrome

A phrase is a **palindrome** if, after converting all uppercase letters
into lowercase letters and removing all non-alphanumeric characters, it
reads the same forward and backward. Alphanumeric characters include
letters and numbers.

Given a string `s`, return `true` if it is a palindrome, or `false`
otherwise.

## Examples

```
Input: s = "A man, a plan, a canal: Panama"
Output: true
Explanation: "amanaplanacanalpanama" is a palindrome.
```

```
Input: s = "race a car"
Output: false
Explanation: "raceacar" is not a palindrome.
```

```
Input: s = " "
Output: true
Explanation: After removing non-alphanumeric characters, s is an empty
string "". An empty string reads the same forward and backward, so it
is a palindrome.
```

## Constraints

- `1 <= s.length <= 2 * 10^5`
- `s` consists only of printable ASCII characters.
