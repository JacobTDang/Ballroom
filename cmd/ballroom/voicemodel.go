package main

import (
	"fmt"
	"strings"
)

// Dictation accuracy is a model-size choice, not an engine choice.
// whisper-cli already accepts any of these; the app just hardcoded the
// second-smallest one, so the cheapest available accuracy win was a
// bigger file rather than a different transcriber. Swapping engines
// (Parakeet and friends) would mean a new runtime, a new install path,
// and new failure modes for a gain this covers most of.
//
// base stays the default: it is what existing installs already have on
// disk, and silently pulling a larger model because someone upgraded
// would be a rude surprise on a metered connection.

// voiceModel is one selectable whisper model.
type voiceModel struct {
	Name string // as passed to `ballroom config set-voice-model`
	File string // ggml file name
	MB   int    // download size, for the consent prompt
	Note string // what you get for the bytes
}

var voiceModels = []voiceModel{
	{"tiny", "ggml-tiny.en.bin", 75, "fastest, roughest — fine for short commands"},
	{"base", "ggml-base.en.bin", 141, "the default: quick, good enough for dictation"},
	{"small", "ggml-small.en.bin", 465, "noticeably better on technical words and names"},
	{"medium", "ggml-medium.en.bin", 1462, "best available here; slower and a big download"},
}

const defaultVoiceModel = "base"

// lookupVoiceModel resolves a configured name. An unknown name is an
// error rather than a silent fallback: a typo that quietly transcribes
// with the wrong model is worse than being told.
func lookupVoiceModel(name string) (voiceModel, error) {
	if name == "" {
		name = defaultVoiceModel
	}
	for _, m := range voiceModels {
		if strings.EqualFold(m.Name, name) {
			return m, nil
		}
	}
	return voiceModel{}, fmt.Errorf("unknown voice model %q — choose one of %s", name, voiceModelNames())
}

func voiceModelNames() string {
	names := make([]string, 0, len(voiceModels))
	for _, m := range voiceModels {
		names = append(names, m.Name)
	}
	return strings.Join(names, ", ")
}

// modelURL is the download location for a model file.
func (m voiceModel) modelURL() string {
	return "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/" + m.File
}
