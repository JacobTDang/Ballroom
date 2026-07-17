package main

import "strings"

// Token is one lexeme: its kind (number, ident, op, lparen, rparen),
// its exact text, and the byte position it starts at.
type Token struct {
	Kind string
	Text string
	Pos  int
}

// Tokenize splits input into tokens.
//
// TODO: splitting on spaces calls "3+4" one token and loses every
// position -- and nothing is ever an error.
func Tokenize(input string) ([]Token, error) {
	var tokens []Token
	for _, f := range strings.Fields(input) {
		tokens = append(tokens, Token{Kind: "ident", Text: f})
	}
	return tokens, nil
}
