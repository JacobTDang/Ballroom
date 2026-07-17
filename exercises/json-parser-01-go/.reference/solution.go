package main

import "fmt"

// A recursive-descent parser: one function per grammar rule, each
// consuming exactly its production and leaving p.i just past it.
// Every failure names the byte position -- a parser that guesses is
// worse than no parser.
type parser struct {
	s string
	i int
}

func Parse(input string) (any, error) {
	p := &parser{s: input}
	p.skipSpace()
	v, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	p.skipSpace()
	if p.i != len(p.s) {
		return nil, fmt.Errorf("json: trailing garbage at position %d", p.i)
	}
	return v, nil
}

func (p *parser) skipSpace() {
	for p.i < len(p.s) && (p.s[p.i] == ' ' || p.s[p.i] == '\t' || p.s[p.i] == '\n' || p.s[p.i] == '\r') {
		p.i++
	}
}

func (p *parser) parseValue() (any, error) {
	if p.i >= len(p.s) {
		return nil, fmt.Errorf("json: unexpected end of input at position %d", p.i)
	}
	switch c := p.s[p.i]; {
	case c == '{':
		return p.parseObject()
	case c == '[':
		return p.parseArray()
	case c == '"':
		return p.parseString()
	case c == '-' || (c >= '0' && c <= '9'):
		return p.parseNumber()
	default:
		return p.parseLiteral()
	}
}

func (p *parser) parseObject() (any, error) {
	obj := map[string]any{}
	p.i++ // '{'
	p.skipSpace()
	if p.i < len(p.s) && p.s[p.i] == '}' {
		p.i++
		return obj, nil
	}
	for {
		p.skipSpace()
		if p.i >= len(p.s) || p.s[p.i] != '"' {
			return nil, fmt.Errorf("json: expected object key at position %d", p.i)
		}
		key, err := p.parseString()
		if err != nil {
			return nil, err
		}
		p.skipSpace()
		if p.i >= len(p.s) || p.s[p.i] != ':' {
			return nil, fmt.Errorf("json: expected ':' at position %d", p.i)
		}
		p.i++
		p.skipSpace()
		v, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key.(string)] = v
		p.skipSpace()
		if p.i >= len(p.s) {
			return nil, fmt.Errorf("json: unterminated object at position %d", p.i)
		}
		if p.s[p.i] == ',' {
			p.i++
			continue
		}
		if p.s[p.i] == '}' {
			p.i++
			return obj, nil
		}
		return nil, fmt.Errorf("json: expected ',' or '}' at position %d", p.i)
	}
}

func (p *parser) parseArray() (any, error) {
	var arr []any
	p.i++ // '['
	p.skipSpace()
	if p.i < len(p.s) && p.s[p.i] == ']' {
		p.i++
		return []any{}, nil
	}
	for {
		p.skipSpace()
		v, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, v)
		p.skipSpace()
		if p.i >= len(p.s) {
			return nil, fmt.Errorf("json: unterminated array at position %d", p.i)
		}
		if p.s[p.i] == ',' {
			p.i++
			continue
		}
		if p.s[p.i] == ']' {
			p.i++
			return arr, nil
		}
		return nil, fmt.Errorf("json: expected ',' or ']' at position %d", p.i)
	}
}

func (p *parser) parseString() (any, error) {
	start := p.i
	p.i++ // '"'
	var out []byte
	for p.i < len(p.s) {
		c := p.s[p.i]
		if c == '"' {
			p.i++
			return string(out), nil
		}
		if c == '\\' {
			if p.i+1 >= len(p.s) {
				break
			}
			switch p.s[p.i+1] {
			case '"':
				out = append(out, '"')
			case '\\':
				out = append(out, '\\')
			default:
				return nil, fmt.Errorf("json: unsupported escape at position %d", p.i)
			}
			p.i += 2
			continue
		}
		out = append(out, c)
		p.i++
	}
	return nil, fmt.Errorf("json: unterminated string starting at position %d", start)
}

func (p *parser) parseNumber() (any, error) {
	start := p.i
	if p.s[p.i] == '-' {
		p.i++
	}
	digits := 0
	for p.i < len(p.s) && p.s[p.i] >= '0' && p.s[p.i] <= '9' {
		p.i++
		digits++
	}
	if digits == 0 {
		return nil, fmt.Errorf("json: malformed number at position %d", start)
	}
	n := 0
	neg := p.s[start] == '-'
	for _, c := range p.s[start:p.i] {
		if c == '-' {
			continue
		}
		n = n*10 + int(c-'0')
	}
	if neg {
		n = -n
	}
	return n, nil
}

func (p *parser) parseLiteral() (any, error) {
	rest := p.s[p.i:]
	for lit, v := range map[string]any{"true": true, "false": false, "null": nil} {
		if len(rest) >= len(lit) && rest[:len(lit)] == lit {
			p.i += len(lit)
			return v, nil
		}
	}
	return nil, fmt.Errorf("json: unexpected value at position %d", p.i)
}
