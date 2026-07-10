package main

// CharacterReplacement returns the length of the longest substring of
// s that can be made to contain only one repeating letter after at
// most k character replacements.
func CharacterReplacement(s string, k int) int {
	var count [26]int
	left, maxFreq, best := 0, 0, 0
	for right := 0; right < len(s); right++ {
		count[s[right]-'A']++
		if count[s[right]-'A'] > maxFreq {
			maxFreq = count[s[right]-'A']
		}
		for right-left+1-maxFreq > k {
			count[s[left]-'A']--
			left++
		}
		if window := right - left + 1; window > best {
			best = window
		}
	}
	return best
}
