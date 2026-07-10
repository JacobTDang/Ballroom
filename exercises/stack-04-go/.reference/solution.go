package main

// GenerateParenthesis returns every well-formed combination of n pairs
// of parentheses.
func GenerateParenthesis(n int) []string {
	var res []string
	var backtrack func(cur string, open, close int)
	backtrack = func(cur string, open, close int) {
		if len(cur) == 2*n {
			res = append(res, cur)
			return
		}
		if open < n {
			backtrack(cur+"(", open+1, close)
		}
		if close < open {
			backtrack(cur+")", open, close+1)
		}
	}
	backtrack("", 0, 0)
	return res
}
