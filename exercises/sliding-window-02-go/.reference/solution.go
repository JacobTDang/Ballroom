package main

// LengthOfLongestSubstring returns the length of the longest substring
// of s with no repeating characters.
func LengthOfLongestSubstring(s string) int {
	lastSeen := make(map[byte]int)
	left, best := 0, 0
	for right := 0; right < len(s); right++ {
		c := s[right]
		if idx, ok := lastSeen[c]; ok && idx >= left {
			left = idx + 1
		}
		lastSeen[c] = right
		if window := right - left + 1; window > best {
			best = window
		}
	}
	return best
}
