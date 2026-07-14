package main

// LongestPalindrome returns the longest palindromic substring of s. If
// several substrings share the maximum length, the first one found
// scanning left to right is returned.
func LongestPalindrome(s string) string {
	if len(s) == 0 {
		return ""
	}
	start, end := 0, 0
	for i := 0; i < len(s); i++ {
		len1 := expandAroundCenter(s, i, i)
		len2 := expandAroundCenter(s, i, i+1)
		maxLen := len1
		if len2 > maxLen {
			maxLen = len2
		}
		if maxLen > end-start+1 {
			start = i - (maxLen-1)/2
			end = i + maxLen/2
		}
	}
	return s[start : end+1]
}

// expandAroundCenter grows outward from the center between indices l
// and r (l == r for an odd-length center, r == l+1 for an even-length
// center) and returns the length of the palindrome found.
func expandAroundCenter(s string, l, r int) int {
	for l >= 0 && r < len(s) && s[l] == s[r] {
		l--
		r++
	}
	return r - l - 1
}
