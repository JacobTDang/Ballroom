package main

import "fmt"

// Token is one lexeme: its kind (number, ident, op, lparen, rparen),
// its exact text, and the byte position it starts at.
type Token struct {
	Kind string
	Text string
	Pos  int
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }
func isAlpha(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// Tokenize: one pass, one switch on the current byte -- each branch
// consumes a maximal token and records where it started. The error
// cases (a second dot, an unknown byte) name the exact position,
// because "invalid input" without a position is a debugging session.
func Tokenize(input string) ([]Token, error) {
	var tokens []Token
	i := 0
	for i < len(input) {
		c := input[i]
		switch {
		case c == ' ' || c == '\t' || c == '\n':
			i++
		case isDigit(c):
			start := i
			sawDot := false
			for i < len(input) && (isDigit(input[i]) || input[i] == '.') {
				if input[i] == '.' {
					if sawDot {
						return nil, fmt.Errorf("tokenize: second decimal point at position %d", i)
					}
					sawDot = true
				}
				i++
			}
			tokens = append(tokens, Token{Kind: "number", Text: input[start:i], Pos: start})
		case isAlpha(c):
			start := i
			for i < len(input) && (isAlpha(input[i]) || isDigit(input[i])) {
				i++
			}
			tokens = append(tokens, Token{Kind: "ident", Text: input[start:i], Pos: start})
		case c == '+' || c == '-' || c == '*' || c == '/':
			tokens = append(tokens, Token{Kind: "op", Text: string(c), Pos: i})
			i++
		case c == '(':
			tokens = append(tokens, Token{Kind: "lparen", Text: "(", Pos: i})
			i++
		case c == ')':
			tokens = append(tokens, Token{Kind: "rparen", Text: ")", Pos: i})
			i++
		default:
			return nil, fmt.Errorf("tokenize: unexpected character %q at position %d", c, i)
		}
	}
	return tokens, nil
}
