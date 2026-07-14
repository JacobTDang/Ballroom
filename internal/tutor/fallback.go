package tutor

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// fallbackToolCall is one parsed tool-call attempt from a
// jsonFallbackToolCalling model's reply.
type fallbackToolCall struct {
	Name      string
	Arguments json.RawMessage
}

// fallbackToolCallHint pre-filters before a full parse: cheap way to
// decide "this reply is at least attempting a tool call" before paying
// for a brace-balanced scan. Accepts both "arguments" (the shape
// jsonFallbackInstruction asks for) and "parameters" (the shape real
// leaked native tool calls have actually used in this codebase -- see
// leakedToolCallPattern's doc comment and agent_test.go's own fixture)
// as aliases.
var fallbackToolCallHint = regexp.MustCompile(`\{"name"\s*:\s*"[a-zA-Z_][a-zA-Z0-9_]*"\s*,\s*"(arguments|parameters)"\s*:`)

// rawFallbackCall is the wire shape parseFallbackToolCall decodes into
// before folding the arguments/parameters alias down to a single field.
type rawFallbackCall struct {
	Name       string          `json:"name"`
	Arguments  json.RawMessage `json:"arguments"`
	Parameters json.RawMessage `json:"parameters"`
}

// parseFallbackToolCall reports whether content is a tool-call attempt
// and, if so, parses it. Tries the whole trimmed (and code-fence-
// stripped, if the entire content is one) reply as a call first -- the
// well-behaved case jsonFallbackInstruction asks for. Only if that fails
// AND fallbackToolCallHint matches somewhere does it fall back to a
// brace-balanced scan for the first {...} object, tolerating a preface
// the model wasn't supposed to add.
//
// matched=false (err always nil) means content is a final answer, not a
// call attempt at all -- the common case once the model is done calling
// tools. matched=true, err!=nil means it looks like an attempted call
// but doesn't parse -- the caller should treat this as a failed attempt
// (corrective retry), not silently fall through to treating content as
// the final answer.
//
// Known limitation, accepted rather than chased: a final answer that
// happens to describe or quote a tool-call-shaped JSON example (e.g.
// explaining what a call "would look like") can be misread as a real
// attempt. jsonFallbackInstruction tells the model its entire reply must
// be ONLY the JSON when calling a tool specifically to keep this
// vanishingly rare in practice, matching this codebase's own established
// lesson (see toolsInstruction's doc comment) that chasing every edge
// case with more parsing complexity tends to regress reliability more
// than it fixes.
func parseFallbackToolCall(content string) (fallbackToolCall, bool, error) {
	unfenced := stripJSONFence(strings.TrimSpace(content))
	if call, err := decodeFallbackCall(unfenced); err == nil {
		return call, true, nil
	}

	if !fallbackToolCallHint.MatchString(content) {
		return fallbackToolCall{}, false, nil
	}

	obj, ok := extractFirstJSONObject(content)
	if !ok {
		return fallbackToolCall{}, true, fmt.Errorf("tutor: reply looks like a tool call attempt but contains no valid JSON object")
	}
	call, err := decodeFallbackCall(string(obj))
	if err != nil {
		return fallbackToolCall{}, true, err
	}
	return call, true, nil
}

// stripJSONFence removes a ``` or ```json fence wrapping the entire
// (already-trimmed) string s, if present -- a model asked to reply with
// ONLY JSON sometimes wraps it in a code fence anyway. Returns s
// unchanged if it isn't fenced.
func stripJSONFence(s string) string {
	if !strings.HasPrefix(s, "```") {
		return s
	}
	body := strings.TrimPrefix(s, "```")
	body = strings.TrimPrefix(body, "json")
	body = strings.TrimPrefix(body, "\n")
	body = strings.TrimSuffix(strings.TrimSpace(body), "```")
	return strings.TrimSpace(body)
}

// decodeFallbackCall parses s as a rawFallbackCall and folds its
// arguments/parameters alias down to a single field, defaulting to an
// empty object for a tool that takes no arguments. Returns an error if s
// isn't valid JSON or has no "name".
func decodeFallbackCall(s string) (fallbackToolCall, error) {
	var raw rawFallbackCall
	if err := json.Unmarshal([]byte(s), &raw); err != nil {
		return fallbackToolCall{}, fmt.Errorf("tutor: parse fallback tool call: %w", err)
	}
	if raw.Name == "" {
		return fallbackToolCall{}, fmt.Errorf("tutor: fallback tool call has no \"name\"")
	}
	args := raw.Arguments
	if len(args) == 0 {
		args = raw.Parameters
	}
	if len(args) == 0 {
		args = json.RawMessage("{}")
	}
	return fallbackToolCall{Name: raw.Name, Arguments: args}, nil
}

// extractFirstJSONObject brace-scans s for the first top-level {...}
// span, honoring string literals and escaped quotes so a "}" or "{"
// inside a quoted value (e.g. a tool argument's own text) doesn't
// desync the scan. Returns ok=false if s has no balanced object.
func extractFirstJSONObject(s string) (json.RawMessage, bool) {
	start := strings.IndexByte(s, '{')
	if start == -1 {
		return nil, false
	}

	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if inString {
			switch {
			case escaped:
				escaped = false
			case c == '\\':
				escaped = true
			case c == '"':
				inString = false
			}
			continue
		}
		switch c {
		case '"':
			inString = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return json.RawMessage(s[start : i+1]), true
			}
		}
	}
	return nil, false
}
