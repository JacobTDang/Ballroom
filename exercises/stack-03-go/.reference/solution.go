package main

import "strconv"

// EvalRPN evaluates an arithmetic expression given in Reverse Polish
// Notation and returns the result.
func EvalRPN(tokens []string) int {
	var stack []int
	for _, tok := range tokens {
		switch tok {
		case "+", "-", "*", "/":
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var res int
			switch tok {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				res = a / b
			}
			stack = append(stack, res)
		default:
			n, _ := strconv.Atoi(tok)
			stack = append(stack, n)
		}
	}
	return stack[0]
}
