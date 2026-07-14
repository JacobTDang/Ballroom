package tutor

import (
	"encoding/json"
	"testing"
)

func TestParseFallbackToolCall_WellFormedCall(t *testing.T) {
	call, matched, err := parseFallbackToolCall(`{"name": "read_solution_file", "arguments": {"a": 1}}`)
	if err != nil {
		t.Fatalf("parseFallbackToolCall: %v", err)
	}
	if !matched {
		t.Fatal("matched = false, want true for a well-formed call")
	}
	if call.Name != "read_solution_file" {
		t.Errorf("Name = %q, want %q", call.Name, "read_solution_file")
	}
	if string(call.Arguments) != `{"a": 1}` {
		t.Errorf("Arguments = %s, want %s", call.Arguments, `{"a": 1}`)
	}
}

func TestParseFallbackToolCall_FencedWithJSONTag(t *testing.T) {
	content := "```json\n{\"name\": \"highlight_lines\", \"arguments\": {\"file\": \"solution.go\"}}\n```"
	call, matched, err := parseFallbackToolCall(content)
	if err != nil {
		t.Fatalf("parseFallbackToolCall: %v", err)
	}
	if !matched {
		t.Fatal("matched = false, want true")
	}
	if call.Name != "highlight_lines" {
		t.Errorf("Name = %q, want %q", call.Name, "highlight_lines")
	}
}

func TestParseFallbackToolCall_FencedWithoutLanguageTag(t *testing.T) {
	content := "```\n{\"name\": \"read_test_output\", \"arguments\": {}}\n```"
	call, matched, err := parseFallbackToolCall(content)
	if err != nil {
		t.Fatalf("parseFallbackToolCall: %v", err)
	}
	if !matched {
		t.Fatal("matched = false, want true")
	}
	if call.Name != "read_test_output" {
		t.Errorf("Name = %q, want %q", call.Name, "read_test_output")
	}
}

// TestParseFallbackToolCall_ParametersAlias covers the real observed
// drift pattern for this exact class of model behavior in this codebase
// -- leakedToolCallPattern's own doc comment and agent_test.go's fixture
// both show real leaked native tool calls using "parameters", not
// "arguments". The fallback parser accepts both defensively.
func TestParseFallbackToolCall_ParametersAlias(t *testing.T) {
	call, matched, err := parseFallbackToolCall(`{"name": "read_cursor_position", "parameters": {"x": 1}}`)
	if err != nil {
		t.Fatalf("parseFallbackToolCall: %v", err)
	}
	if !matched {
		t.Fatal("matched = false, want true")
	}
	if string(call.Arguments) != `{"x": 1}` {
		t.Errorf("Arguments = %s, want the \"parameters\" value", call.Arguments)
	}
}

func TestParseFallbackToolCall_NoArgumentsKeyDefaultsToEmptyObject(t *testing.T) {
	call, matched, err := parseFallbackToolCall(`{"name": "read_problem_statement"}`)
	if err != nil {
		t.Fatalf("parseFallbackToolCall: %v", err)
	}
	if !matched {
		t.Fatal("matched = false, want true")
	}
	if string(call.Arguments) != "{}" {
		t.Errorf("Arguments = %s, want {} for a tool with no explicit arguments", call.Arguments)
	}
}

func TestParseFallbackToolCall_PrefacedWithProseStillParses(t *testing.T) {
	content := "Let me check that.\n\n{\"name\": \"read_solution_file\", \"arguments\": {}}"
	call, matched, err := parseFallbackToolCall(content)
	if err != nil {
		t.Fatalf("parseFallbackToolCall: %v", err)
	}
	if !matched {
		t.Fatal("matched = false, want true -- a preface shouldn't prevent finding the call")
	}
	if call.Name != "read_solution_file" {
		t.Errorf("Name = %q, want %q", call.Name, "read_solution_file")
	}
}

func TestParseFallbackToolCall_NestedBracesInArguments(t *testing.T) {
	content := `{"name": "x", "arguments": {"nested": {"a": 1}, "note": "has } a brace"}}`
	call, matched, err := parseFallbackToolCall(content)
	if err != nil {
		t.Fatalf("parseFallbackToolCall: %v", err)
	}
	if !matched {
		t.Fatal("matched = false, want true")
	}
	var args map[string]any
	if err := json.Unmarshal(call.Arguments, &args); err != nil {
		t.Fatalf("Arguments didn't round-trip as valid JSON: %v (got %s)", err, call.Arguments)
	}
}

func TestParseFallbackToolCall_MalformedJSONThatLooksLikeAnAttempt(t *testing.T) {
	// Missing closing brace -- the hint still matches (name+arguments
	// shape is present), so this must be reported as a failed attempt
	// (matched=true, err!=nil), not silently treated as a final answer.
	content := `{"name": "read_solution_file", "arguments": {}`
	_, matched, err := parseFallbackToolCall(content)
	if !matched {
		t.Fatal("matched = false, want true -- this looks like an attempted call")
	}
	if err == nil {
		t.Fatal("expected an error for malformed JSON that still matches the call shape")
	}
}

func TestParseFallbackToolCall_PlainFinalAnswerIsNotACall(t *testing.T) {
	content := "Your solution looks correct. The time complexity is O(n)."
	_, matched, err := parseFallbackToolCall(content)
	if err != nil {
		t.Fatalf("parseFallbackToolCall: %v", err)
	}
	if matched {
		t.Error("matched = true, want false for plain prose with no tool-call shape")
	}
}

func TestParseFallbackToolCall_JSONWithoutNameFieldIsNotACall(t *testing.T) {
	// Some other JSON shape entirely -- not a tool-call attempt, since it
	// never matches the hint pattern at all.
	content := `{"foo": "bar"}`
	_, matched, err := parseFallbackToolCall(content)
	if err != nil {
		t.Fatalf("parseFallbackToolCall: %v", err)
	}
	if matched {
		t.Error("matched = true, want false -- no \"name\" field means this was never a call attempt")
	}
}

func TestExtractFirstJSONObject_FindsBalancedObject(t *testing.T) {
	obj, ok := extractFirstJSONObject(`prefix {"a": {"b": 1}} suffix`)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	if string(obj) != `{"a": {"b": 1}}` {
		t.Errorf("obj = %s, want %s", obj, `{"a": {"b": 1}}`)
	}
}

func TestExtractFirstJSONObject_IgnoresBracesInsideStrings(t *testing.T) {
	obj, ok := extractFirstJSONObject(`{"note": "a } brace and a { brace"}`)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	want := `{"note": "a } brace and a { brace"}`
	if string(obj) != want {
		t.Errorf("obj = %s, want %s", obj, want)
	}
}

func TestExtractFirstJSONObject_HandlesEscapedQuotes(t *testing.T) {
	obj, ok := extractFirstJSONObject(`{"note": "she said \"hi }\" to me"}`)
	if !ok {
		t.Fatal("ok = false, want true")
	}
	var decoded map[string]string
	if err := json.Unmarshal(obj, &decoded); err != nil {
		t.Fatalf("extracted object didn't parse: %v (got %s)", err, obj)
	}
}

func TestExtractFirstJSONObject_NoObjectReturnsFalse(t *testing.T) {
	_, ok := extractFirstJSONObject("no braces here at all")
	if ok {
		t.Error("ok = true, want false")
	}
}

func TestExtractFirstJSONObject_UnbalancedReturnsFalse(t *testing.T) {
	_, ok := extractFirstJSONObject(`{"name": "x", "arguments": {}`)
	if ok {
		t.Error("ok = true, want false for an object that never closes")
	}
}
