package main

import (
	"strings"
	"testing"
)

func TestLookupVoiceModel_DefaultsToBase(t *testing.T) {
	m, err := lookupVoiceModel("")
	if err != nil {
		t.Fatalf("empty name: %v", err)
	}
	if m.Name != defaultVoiceModel {
		t.Errorf("default = %q, want %q", m.Name, defaultVoiceModel)
	}
	// Existing installs already have this file; changing the default
	// would pull a fresh multi-hundred-megabyte download on upgrade.
	if m.File != "ggml-base.en.bin" {
		t.Errorf("default file = %q, want the model existing installs already have", m.File)
	}
}

func TestLookupVoiceModel_CaseInsensitive(t *testing.T) {
	m, err := lookupVoiceModel("SMALL")
	if err != nil {
		t.Fatalf("SMALL: %v", err)
	}
	if m.Name != "small" {
		t.Errorf("got %q, want small", m.Name)
	}
}

// TestLookupVoiceModel_UnknownIsAnError: a typo that quietly falls back
// would transcribe with a model the user didn't choose and never say so.
func TestLookupVoiceModel_UnknownIsAnError(t *testing.T) {
	_, err := lookupVoiceModel("enormous")
	if err == nil {
		t.Fatal("unknown model was accepted")
	}
	if !strings.Contains(err.Error(), "enormous") {
		t.Errorf("error %q should name the bad value", err)
	}
	for _, want := range []string{"tiny", "base", "small", "medium"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q should list the valid choice %q", err, want)
		}
	}
}

func TestVoiceModels_URLsAndSizes(t *testing.T) {
	for _, m := range voiceModels {
		if !strings.HasPrefix(m.modelURL(), "https://") || !strings.HasSuffix(m.modelURL(), m.File) {
			t.Errorf("%s: bad URL %q", m.Name, m.modelURL())
		}
		if m.MB <= 0 {
			t.Errorf("%s: size %d MB, want a real number for the consent prompt", m.Name, m.MB)
		}
		if m.Note == "" {
			t.Errorf("%s: no note explaining what the extra bytes buy", m.Name)
		}
	}
}

func TestVoiceModels_OrderedSmallestFirst(t *testing.T) {
	for i := 1; i < len(voiceModels); i++ {
		if voiceModels[i].MB <= voiceModels[i-1].MB {
			t.Errorf("%s (%d MB) is not larger than %s (%d MB) — the list should read smallest to largest",
				voiceModels[i].Name, voiceModels[i].MB, voiceModels[i-1].Name, voiceModels[i-1].MB)
		}
	}
}
