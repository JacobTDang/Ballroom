package main

// IsAnagram returns true if t is an anagram of s.
func IsAnagram(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	counts := make(map[rune]int)
	for _, c := range s {
		counts[c]++
	}
	for _, c := range t {
		counts[c]--
		if counts[c] < 0 {
			return false
		}
	}
	return true
}
