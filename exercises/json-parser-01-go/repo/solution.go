package main

import (
	"fmt"
	"strings"
)

// Parse reads a JSON subset: objects, arrays, strings (\" and \\
// escapes), integers, true/false/null.
//
// TODO: this handles only a flat {"key": "value"} object via string
// splitting -- no nesting, no arrays, no numbers, no errors worth
// the name.
func Parse(input string) (any, error) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "{") || !strings.HasSuffix(input, "}") {
		return nil, fmt.Errorf("only objects supported")
	}
	result := map[string]any{}
	body := strings.Trim(input, "{}")
	for _, pair := range strings.Split(body, ",") {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			k := strings.Trim(strings.TrimSpace(kv[0]), "\"")
			v := strings.Trim(strings.TrimSpace(kv[1]), "\"")
			result[k] = v
		}
	}
	return result, nil
}
