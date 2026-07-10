package main

var phoneLetters = map[byte]string{
	'2': "abc", '3': "def", '4': "ghi", '5': "jkl",
	'6': "mno", '7': "pqrs", '8': "tuv", '9': "wxyz",
}

// LetterCombinations returns every letter combination that digits
// could represent on a phone keypad.
func LetterCombinations(digits string) []string {
	if len(digits) == 0 {
		return nil
	}
	var res []string
	cur := make([]byte, 0, len(digits))
	var backtrack func(idx int)
	backtrack = func(idx int) {
		if idx == len(digits) {
			res = append(res, string(cur))
			return
		}
		letters := phoneLetters[digits[idx]]
		for i := 0; i < len(letters); i++ {
			cur = append(cur, letters[i])
			backtrack(idx + 1)
			cur = cur[:len(cur)-1]
		}
	}
	backtrack(0)
	return res
}
