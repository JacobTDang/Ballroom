package main

// CheckInclusion reports whether s2 contains a permutation of s1 as a
// contiguous substring.
func CheckInclusion(s1, s2 string) bool {
	if len(s1) > len(s2) {
		return false
	}
	var need, window [26]int
	for i := 0; i < len(s1); i++ {
		need[s1[i]-'a']++
		window[s2[i]-'a']++
	}
	if need == window {
		return true
	}
	for i := len(s1); i < len(s2); i++ {
		window[s2[i]-'a']++
		window[s2[i-len(s1)]-'a']--
		if need == window {
			return true
		}
	}
	return false
}
