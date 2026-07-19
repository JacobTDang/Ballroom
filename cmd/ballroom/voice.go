package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
)

// ballroom voice: host-side speech-to-text into the running session's
// tutor pane. Records the mic until Enter (ffmpeg/avfoundation),
// transcribes locally (whisper-cpp), and types -- never sends -- the
// text into the tutor pane via docker exec + tmux send-keys, so the
// user reviews it before pressing Enter themselves. Everything runs on
// the host: the container never needs microphone access.
//
// Dependencies (user-approved): ffmpeg and whisper-cpp, both via
// Homebrew; the whisper model file is downloaded on first use with
// explicit consent. macOS built-in dictation (press fn twice) remains
// the zero-setup alternative and needs none of this.

// tutorPaneTarget is the tmux pane the transcript is typed into --
// session:window.pane per docker/entrypoint.sh's layout (pane 1 is the
// tutor chat).
const tutorPaneTarget = "practice:MAIN.1"

func voiceCmd(args []string) error {
	fromWav := ""
	if len(args) == 2 && args[0] == "--from-wav" {
		// Transcribe an existing recording instead of the mic --
		// debugging aid, and how this command is verified end to end
		// without microphone hardware.
		fromWav = args[1]
	} else if len(args) != 0 {
		return fmt.Errorf("usage: ballroom voice [--from-wav <file.wav>]")
	}

	whisper, err := exec.LookPath("whisper-cli")
	if err != nil {
		return fmt.Errorf("voice: whisper-cli not found -- install it with `brew install whisper-cpp`")
	}
	if fromWav == "" {
		if _, err := exec.LookPath("ffmpeg"); err != nil {
			return fmt.Errorf("voice: ffmpeg not found -- install it with `brew install ffmpeg`")
		}
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	model, err := lookupVoiceModel(cfg.VoiceModel)
	if err != nil {
		return err
	}
	modelPath := filepath.Join(cfg.DataDir, "whisper", model.File)
	if err := ensureWhisperModel(modelPath, model); err != nil {
		return err
	}

	container, err := findSessionContainer()
	if err != nil {
		return err
	}

	wav := fromWav
	if wav == "" {
		wav, err = recordUntilEnter()
		if err != nil {
			return err
		}
		defer os.Remove(wav)
	}

	fmt.Println("transcribing...")
	text, err := transcribe(whisper, modelPath, wav)
	if err != nil {
		return err
	}
	if text == "" {
		return fmt.Errorf("voice: transcription came back empty -- was anything said?")
	}

	fmt.Printf("\n> %s\n\n", text)
	if err := typeIntoTutorPane(container, text); err != nil {
		return err
	}
	fmt.Println("typed into the tutor pane -- review it there and press Enter to send.")
	return nil
}

// ensureWhisperModel downloads the model on first use, with explicit
// consent -- a 142 MB download should never happen silently.
func ensureWhisperModel(path string, model voiceModel) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	fmt.Printf("whisper %s model not found\ndownload %s (~%d MB, %s) from Hugging Face now? [y/N]: ",
		model.Name, model.File, model.MB, model.Note)
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() || !strings.EqualFold(strings.TrimSpace(scanner.Text()), "y") {
		return fmt.Errorf("voice: model download declined -- nothing was downloaded")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("voice: create model dir: %w", err)
	}
	curl := exec.Command("curl", "-L", "--fail", "--progress-bar", "-o", path, model.modelURL())
	curl.Stdout = os.Stdout
	curl.Stderr = os.Stderr
	if err := curl.Run(); err != nil {
		os.Remove(path) // never leave a truncated model behind
		return fmt.Errorf("voice: model download failed: %w", err)
	}
	return nil
}

// findSessionContainer locates the one running practice container via
// the orchestrator's session label.
func findSessionContainer() (string, error) {
	out, err := exec.Command("docker", "ps", "--latest", "-q", "--filter", orchestrator.SessionContainerFilter).Output()
	if err != nil {
		return "", fmt.Errorf("voice: docker ps: %w", err)
	}
	id := strings.TrimSpace(string(out))
	if id == "" {
		return "", fmt.Errorf("voice: no running practice session found -- start one first (sessions launched by an older ballroom binary lack the container label; relaunch after updating)")
	}
	return id, nil
}

// recordUntilEnter captures the default microphone to a temp wav
// (16 kHz mono, what whisper wants) until the user presses Enter.
func recordUntilEnter() (string, error) {
	f, err := os.CreateTemp("", "ballroom-voice-*.wav")
	if err != nil {
		return "", fmt.Errorf("voice: temp wav: %w", err)
	}
	f.Close()

	// avfoundation ":0" is macOS's default audio input; -y overwrites
	// the just-created temp file.
	rec := exec.Command("ffmpeg", "-hide_banner", "-loglevel", "error",
		"-f", "avfoundation", "-i", ":0", "-ar", "16000", "-ac", "1", "-y", f.Name())
	rec.Stderr = os.Stderr
	if err := rec.Start(); err != nil {
		os.Remove(f.Name())
		return "", fmt.Errorf("voice: start recording: %w", err)
	}

	fmt.Println("recording -- press Enter to stop...")
	bufio.NewScanner(os.Stdin).Scan()

	// SIGINT lets ffmpeg finalize the wav header; Kill would corrupt it.
	if err := rec.Process.Signal(syscall.SIGINT); err != nil {
		rec.Process.Kill()
	}
	if err := rec.Wait(); err != nil {
		// ffmpeg exits non-zero on SIGINT (255) even after writing a
		// valid file -- only a missing/empty output is a real failure.
		if info, statErr := os.Stat(f.Name()); statErr != nil || info.Size() == 0 {
			os.Remove(f.Name())
			return "", fmt.Errorf("voice: recording produced no audio (microphone permission for your terminal app?): %w", err)
		}
	}
	return f.Name(), nil
}

// transcribe runs whisper-cli and returns the cleaned transcript.
func transcribe(whisper, modelPath, wav string) (string, error) {
	out, err := exec.Command(whisper, "-m", modelPath, "-f", wav, "-nt", "-np").Output()
	if err != nil {
		return "", fmt.Errorf("voice: whisper-cli: %w", err)
	}
	return strings.Join(strings.Fields(string(out)), " "), nil
}

// typeIntoTutorPane literally types text into the tutor pane (-l =
// literal, no key-name interpretation) without pressing Enter -- the
// user confirms in the pane, so a mis-transcription never fires a turn.
func typeIntoTutorPane(container, text string) error {
	cmd := exec.Command("docker", "exec", container, "tmux", "send-keys", "-t", tutorPaneTarget, "-l", text)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("voice: typing into the tutor pane: %v: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}
