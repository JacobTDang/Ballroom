package main

// CountSubstrings returns the number of palindromic substrings in s,
// counting substrings at different positions separately even if they
// contain the same characters.
func CountSubstrings(s string) int {
	count := 0
	for i := 0; i < len(s); i++ {
		count += countExpansions(s, i, i)
		count += countExpansions(s, i, i+1)
	}
	return count
}

// countExpansions grows outward from the center between indices l and
// r, counting one palindrome for every successful expansion.
func countExpansions(s string, l, r int) int {
	count := 0
	for l >= 0 && r < len(s) && s[l] == s[r] {
		count++
		l--
		r++
	}
	return count
}
