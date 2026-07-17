package main

import "strings"

// Parse reads an INI document into section -> key -> value.
//
// TODO: no sections, no comments, no errors -- every line is
// split on '=' into the "" section, and malformed lines are silently
// skipped.
func Parse(input string) (map[string]map[string]string, error) {
	result := map[string]map[string]string{"": {}}
	for _, line := range strings.Split(input, "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[""][parts[0]] = parts[1]
		}
	}
	return result, nil
}
