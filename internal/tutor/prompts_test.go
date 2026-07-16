package tutor

import (
	"strings"
	"testing"

	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/cloudwego/eino/schema"
)

func TestSystemPromptForMode_DiffersByMode(t *testing.T) {
	syntaxOnly := systemPromptForMode(exercise.TutorModeSyntaxOnly)
	hintsFirst := systemPromptForMode(exercise.TutorModeHintsFirst)
	fullAssist := systemPromptForMode(exercise.TutorModeFullAssist)

	for name, p := range map[string]string{"syntax-only": syntaxOnly, "hints-first": hintsFirst, "full-assist": fullAssist} {
		if p == "" {
			t.Errorf("mode %s: system prompt is empty", name)
		}
	}
	if syntaxOnly == hintsFirst {
		t.Error("syntax-only and hints-first got the same system prompt — modes must differ")
	}
	if syntaxOnly == fullAssist {
		t.Error("syntax-only and full-assist got the same system prompt — modes must differ")
	}
	if hintsFirst == fullAssist {
		t.Error("hints-first and full-assist got the same system prompt — modes must differ")
	}
}

func TestSystemPromptForMode_UnknownModeDefaultsToFullAssist(t *testing.T) {
	got := systemPromptForMode("some-unrecognized-mode")
	want := systemPromptForMode(exercise.TutorModeFullAssist)
	if got != want {
		t.Errorf("unknown mode prompt = %q, want it to match full-assist's prompt %q", got, want)
	}
}

func TestModePrompts_HasEntryForEveryKnownMode(t *testing.T) {
	for _, mode := range []string{exercise.TutorModeSyntaxOnly, exercise.TutorModeHintsFirst, exercise.TutorModeFullAssist} {
		if modePrompts[mode] == "" {
			t.Errorf("modePrompts[%q] is empty, want a real prompt body", mode)
		}
	}
}

func TestToolsInstructions_HasNativeEntry(t *testing.T) {
	got := toolsInstructions[nativeToolCalling]
	if got == "" {
		t.Fatal("toolsInstructions[nativeToolCalling] is empty")
	}
	if got != toolsInstruction {
		t.Errorf("toolsInstructions[nativeToolCalling] = %q, want it to match the toolsInstruction constant", got)
	}
}

func TestToolsInstructions_HasJSONFallbackEntry(t *testing.T) {
	got := toolsInstructions[jsonFallbackToolCalling]
	if got == "" {
		t.Fatal("toolsInstructions[jsonFallbackToolCalling] is empty")
	}
	if got == toolsInstructions[nativeToolCalling] {
		t.Error("jsonFallbackToolCalling's instruction is identical to native's -- it needs fundamentally different protocol text, not a copy")
	}
}

// TestPersonaPromptForMode_ExcludesToolsText locks in why
// personaPromptForMode exists at all: newTutorModel seeds history with
// this (not systemPromptForMode), because the tools instruction is
// strategy-dependent and can change per role per turn once routing is
// enabled, unlike the persona text which is fixed for the session. If
// this ever accidentally included tools text again, a fallback-strategy
// role would see native-style instructions baked permanently into
// history with no way to override them per turn.
func TestPersonaPromptForMode_ExcludesToolsText(t *testing.T) {
	for _, mode := range []string{exercise.TutorModeSyntaxOnly, exercise.TutorModeHintsFirst, exercise.TutorModeFullAssist} {
		got := personaPromptForMode(mode)
		if got == "" {
			t.Errorf("personaPromptForMode(%q) is empty", mode)
		}
		if strings.Contains(got, toolsInstruction) {
			t.Errorf("personaPromptForMode(%q) contains the native toolsInstruction text, want persona-only", mode)
		}
		if strings.Contains(got, jsonFallbackInstruction) {
			t.Errorf("personaPromptForMode(%q) contains jsonFallbackInstruction text, want persona-only", mode)
		}
	}
}

func TestPersonaPromptForMode_UnknownModeDefaultsToFullAssist(t *testing.T) {
	got := personaPromptForMode("some-unrecognized-mode")
	want := personaPromptForMode(exercise.TutorModeFullAssist)
	if got != want {
		t.Errorf("unknown mode persona = %q, want it to match full-assist's %q", got, want)
	}
}

func TestSystemPromptForMode_EqualsPersonaPlusNativeToolsInstruction(t *testing.T) {
	for _, mode := range []string{exercise.TutorModeSyntaxOnly, exercise.TutorModeHintsFirst, exercise.TutorModeFullAssist} {
		got := systemPromptForMode(mode)
		want := toolsInstructions[nativeToolCalling] + personaPromptForMode(mode)
		if got != want {
			t.Errorf("systemPromptForMode(%q) = %q, want %q", mode, got, want)
		}
	}
}

func TestPrependToolsPrompt_NativeStrategyOmitsToolCatalog(t *testing.T) {
	msgs := []*schema.Message{schema.UserMessage("hello")}
	got := prependToolsPrompt(nativeToolCalling, "CATALOG_MARKER", msgs)

	if len(got) != 2 {
		t.Fatalf("got %d messages, want 2 (prepended system message + original user message)", len(got))
	}
	if got[0].Role != schema.System {
		t.Errorf("got[0].Role = %v, want schema.System", got[0].Role)
	}
	if got[0].Content != toolsInstructions[nativeToolCalling] {
		t.Errorf("got[0].Content = %q, want exactly the native toolsInstructions entry", got[0].Content)
	}
	if strings.Contains(got[0].Content, "CATALOG_MARKER") {
		t.Error("native strategy must not include the tool catalog text -- tool schemas are already bound via the real API, not described in prose")
	}
	if got[1] != msgs[0] {
		t.Error("original message was not preserved as the second element")
	}
}

func TestPrependToolsPrompt_JSONFallbackStrategyIncludesToolCatalog(t *testing.T) {
	msgs := []*schema.Message{schema.UserMessage("hello")}
	got := prependToolsPrompt(jsonFallbackToolCalling, "CATALOG_MARKER", msgs)

	if len(got) != 2 {
		t.Fatalf("got %d messages, want 2", len(got))
	}
	if !strings.Contains(got[0].Content, jsonFallbackInstruction) {
		t.Error("expected the fallback protocol instruction text")
	}
	if !strings.Contains(got[0].Content, "CATALOG_MARKER") {
		t.Error("expected the tool catalog text to be appended for the fallback strategy, since the model has no other way to learn tool schemas")
	}
}

func TestPrependToolsPrompt_DoesNotMutateInputSlice(t *testing.T) {
	original := []*schema.Message{schema.UserMessage("hello")}
	_ = prependToolsPrompt(nativeToolCalling, "", original)

	if len(original) != 1 {
		t.Fatalf("caller's slice was mutated: len = %d, want 1", len(original))
	}
	if original[0].Role != schema.User {
		t.Error("caller's slice was mutated")
	}
}

func TestWantsComprehensionCheck(t *testing.T) {
	cases := []struct {
		mode string
		want bool
	}{
		{exercise.TutorModeSyntaxOnly, false},
		{exercise.TutorModeHintsFirst, true},
		{exercise.TutorModeFullAssist, true},
		// The interviewer must NOT restate the problem or ask the
		// clarifying questions -- that's the candidate's own step 1 of
		// the design method, and doing it for them defeats the drill.
		{exercise.TutorModeInterviewer, false},
		{exercise.TutorModeDesignCoach, true},
		// The behavioral interviewer opens by asking the question --
		// a restate-and-clarify check would talk over that opener.
		{exercise.TutorModeBehavioralInterviewer, false},
		{exercise.TutorModeStoryCoach, true},
		{"unrecognized-mode", true},
	}
	for _, c := range cases {
		if got := wantsComprehensionCheck(c.mode); got != c.want {
			t.Errorf("wantsComprehensionCheck(%q) = %v, want %v", c.mode, got, c.want)
		}
	}
}

func TestModePrompts_DesignModesHaveDistinctPersonas(t *testing.T) {
	interviewer, ok := modePrompts[exercise.TutorModeInterviewer]
	if !ok {
		t.Fatal("modePrompts has no interviewer entry")
	}
	coach, ok := modePrompts[exercise.TutorModeDesignCoach]
	if !ok {
		t.Fatal("modePrompts has no design-coach entry")
	}
	if interviewer == coach {
		t.Error("interviewer and design-coach personas are identical")
	}
	// The interviewer persona's defining constraints: candidate drives
	// step 1, and no solutions are volunteered.
	for _, want := range []string{"restate", "rubric"} {
		if !strings.Contains(strings.ToLower(interviewer), want) {
			t.Errorf("interviewer persona doesn't mention %q:\n%s", want, interviewer)
		}
	}
	if !strings.Contains(strings.ToLower(coach), "solution.md") {
		t.Errorf("design-coach persona should direct writing into solution.md:\n%s", coach)
	}
}
