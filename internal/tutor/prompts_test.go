package tutor

import (
	"testing"

	"github.com/JacobTDang/Ballroom/internal/exercise"
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

func TestWantsComprehensionCheck(t *testing.T) {
	cases := []struct {
		mode string
		want bool
	}{
		{exercise.TutorModeSyntaxOnly, false},
		{exercise.TutorModeHintsFirst, true},
		{exercise.TutorModeFullAssist, true},
		{"unrecognized-mode", true},
	}
	for _, c := range cases {
		if got := wantsComprehensionCheck(c.mode); got != c.want {
			t.Errorf("wantsComprehensionCheck(%q) = %v, want %v", c.mode, got, c.want)
		}
	}
}
