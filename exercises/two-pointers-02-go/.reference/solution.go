package main

// IsPalindrome reports whether s is a palindrome, considering only
// alphanumeric characters and ignoring case.
func IsPalindrome(s string) bool {
	isAlnum := func(b byte) bool {
		return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
	}
	toLower := func(b byte) byte {
		if b >= 'A' && b <= 'Z' {
			return b - 'A' + 'a'
		}
		return b
	}

	lo, hi := 0, len(s)-1
	for lo < hi {
		for lo < hi && !isAlnum(s[lo]) {
			lo++
		}
		for lo < hi && !isAlnum(s[hi]) {
			hi--
		}
		if toLower(s[lo]) != toLower(s[hi]) {
			return false
		}
		lo++
		hi--
	}
	return true
}
