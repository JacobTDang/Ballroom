package main

// Partition returns every way to split s into substrings that are
// all palindromes.
func Partition(s string) [][]string {
	isPalindrome := func(str string) bool {
		l, r := 0, len(str)-1
		for l < r {
			if str[l] != str[r] {
				return false
			}
			l++
			r--
		}
		return true
	}

	var res [][]string
	var cur []string
	var backtrack func(start int)
	backtrack = func(start int) {
		if start == len(s) {
			res = append(res, append([]string(nil), cur...))
			return
		}
		for end := start + 1; end <= len(s); end++ {
			sub := s[start:end]
			if isPalindrome(sub) {
				cur = append(cur, sub)
				backtrack(end)
				cur = cur[:len(cur)-1]
			}
		}
	}
	backtrack(0)
	return res
}
