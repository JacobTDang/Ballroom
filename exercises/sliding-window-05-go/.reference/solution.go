package main

// MinWindow returns the shortest substring of s containing every
// character of t (with duplicates), or "" if no such substring exists.
func MinWindow(s, t string) string {
	if len(t) == 0 || len(s) < len(t) {
		return ""
	}
	need := make(map[byte]int)
	for i := 0; i < len(t); i++ {
		need[t[i]]++
	}
	required := len(need)
	have := 0
	window := make(map[byte]int)

	bestLen := -1
	bestStart := 0
	left := 0
	for right := 0; right < len(s); right++ {
		c := s[right]
		window[c]++
		if cnt, ok := need[c]; ok && window[c] == cnt {
			have++
		}
		for have == required {
			if bestLen == -1 || right-left+1 < bestLen {
				bestLen = right - left + 1
				bestStart = left
			}
			leftChar := s[left]
			window[leftChar]--
			if cnt, ok := need[leftChar]; ok && window[leftChar] < cnt {
				have--
			}
			left++
		}
	}
	if bestLen == -1 {
		return ""
	}
	return s[bestStart : bestStart+bestLen]
}
