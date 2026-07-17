package main

import (
	"fmt"
	"strings"
)

// Parse: line-oriented with a current-section cursor. Every branch
// either consumes the line's whole meaning or errors with its 1-based
// number -- a config line that parses as nothing is a typo the user
// deserves to hear about, not silence.
func Parse(input string) (map[string]map[string]string, error) {
	result := map[string]map[string]string{}
	section := ""
	result[section] = map[string]string{}

	for n, raw := range strings.Split(input, "\n") {
		line := strings.TrimSpace(raw)
		lineNo := n + 1
		switch {
		case line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";"):
			continue
		case strings.HasPrefix(line, "["):
			if !strings.HasSuffix(line, "]") {
				return nil, fmt.Errorf("ini: line %d: unclosed section header %q", lineNo, line)
			}
			section = strings.TrimSpace(line[1 : len(line)-1])
			if _, ok := result[section]; !ok {
				result[section] = map[string]string{}
			}
		case strings.Contains(line, "="):
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			if key == "" {
				return nil, fmt.Errorf("ini: line %d: empty key", lineNo)
			}
			result[section][key] = strings.TrimSpace(parts[1]) // later key wins
		default:
			return nil, fmt.Errorf("ini: line %d: not a header, comment, or key=value: %q", lineNo, line)
		}
	}
	return result, nil
}
