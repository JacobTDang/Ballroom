package main

// Match: the classic two-pointer loop. On '*' remember both positions
// (starIdx, starMatch); on any later mismatch, back up to just after
// the star and let it swallow one more character. That pair IS the
// backtracking state -- no recursion needed, O(len(p)*len(s)) worst
// case.
func Match(pattern, s string) bool {
	p, i := 0, 0
	starIdx, starMatch := -1, 0

	for i < len(s) {
		if p < len(pattern) {
			switch pattern[p] {
			case '*':
				starIdx, starMatch = p, i
				p++
				continue
			case '?':
				p++
				i++
				continue
			case '[':
				ok, next, valid := matchClass(pattern, p, s[i])
				if !valid {
					return false // unclosed class: invalid pattern
				}
				if ok {
					p = next
					i++
					continue
				}
			default:
				if pattern[p] == s[i] {
					p++
					i++
					continue
				}
			}
		}
		// Mismatch: if a star is behind us, let it swallow one more.
		if starIdx == -1 {
			return false
		}
		starMatch++
		p = starIdx + 1
		i = starMatch
	}

	for p < len(pattern) && pattern[p] == '*' {
		p++
	}
	return p == len(pattern)
}

// matchClass matches s[i] against the class starting at pattern[p]
// (which is '['). Returns whether it matched, the index just past the
// closing ']', and whether the class was well-formed.
func matchClass(pattern string, p int, c byte) (matched bool, next int, valid bool) {
	q := p + 1
	for q < len(pattern) && pattern[q] != ']' {
		if q+2 < len(pattern) && pattern[q+1] == '-' && pattern[q+2] != ']' {
			if pattern[q] <= c && c <= pattern[q+2] {
				matched = true
			}
			q += 3
		} else {
			if pattern[q] == c {
				matched = true
			}
			q++
		}
	}
	if q >= len(pattern) {
		return false, 0, false // unclosed class
	}
	return matched, q + 1, true
}
