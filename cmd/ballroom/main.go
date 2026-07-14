package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/orchestrator"
	"github.com/JacobTDang/Ballroom/internal/session"
	"github.com/JacobTDang/Ballroom/internal/tui"
	"github.com/JacobTDang/Ballroom/internal/tutor"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		exitOnErr(homeCmd())
		return
	}

	switch args[0] {
	case "-h", "--help", "help":
		printUsage(os.Stdout)
	case "home":
		exitOnErr(homeCmd())
	case "practice":
		exitOnErr(practiceCmd(args[1:]))
	case "sandbox":
		exitOnErr(sandboxCmd())
	case "submit":
		exitOnErr(submitCmd())
	case "return":
		exitOnErr(returnCmd())
	case "tutor":
		exitOnErr(tutorCmd())
	case "config":
		exitOnErr(configCmd(args[1:]))
	default:
		fmt.Fprintf(os.Stderr, "ballroom: unknown command %q\n\n", args[0])
		printUsage(os.Stderr)
		os.Exit(1)
	}
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ballroom: %v\n", err)
		os.Exit(1)
	}
}

func printUsage(w *os.File) {
	fmt.Fprint(w, `Ballroom — interview practice CLI

Usage:
  ballroom                    Open the homepage (pick an exercise interactively)
  ballroom home                Same as above
  ballroom practice <id>       Practice a specific exercise by id
  ballroom sandbox             Free practice, no grading, persists across sessions
  ballroom submit              Submit your solution (run this inside an active session)
  ballroom tutor               Start the tutor chat (run this inside an active session)
  ballroom return              Return to the host homepage (run this inside an active session)
  ballroom config set-model <tag>   Set the tutor (worker) model (a local Ollama tag, or
                                     an openrouter:<slug> API model) without opening the TUI
  ballroom config set-orchestrator-model <tag|none>
                                     Set the orchestrator model that routes turns to the
                                     worker model, or "none" to disable routing
  ballroom config set-key <key>     Set the OpenRouter API key used by openrouter: models
  ballroom help | -h | --help  Show this help

Examples:
  ballroom
  ballroom practice two-pointers-01-go
  ballroom sandbox
  ballroom config set-model openrouter:anthropic/claude-3.5-sonnet
  ballroom config set-orchestrator-model openrouter:nvidia/nemotron-3-ultra-550b-a55b:free
  ballroom config set-key sk-...

Reset the sandbox volume:
  docker volume rm ballroom-sandbox
`)
}

// homeCmd shows the full-screen boot check + exercise picker (see
// internal/tui) — the "home base" you return to between sessions rather
// than having to remember exercise ids. The ballroom binary baked into
// the practice image (docker/Dockerfile) means this same code path can
// run either on the host or inside an active session's container; the
// boot screen's preflight checks (CheckDocker, CheckOllama against
// localhost:11434, ...) assume host networking and there's no Docker
// client inside the container, so booting a nested instance there
// doesn't fail cleanly. If we're inside a session, return to the host
// homepage instead of attempting that nested boot.
func homeCmd() error {
	if isSessionContext() {
		return returnToHost()
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return tui.Run(cfg)
}

// isSessionContext reports whether this process is running inside an
// active practice session's container, as opposed to on the host. It
// checks all three session-scoped env vars orchestrator.RunExercise sets
// via `docker run -e` (PRACTICE_WORKSPACE_DIR, PRACTICE_CONTROL_DIR,
// PRACTICE_STARTED_AT) rather than any single one, so a host shell that
// happens to have one of these set for unrelated reasons isn't
// misdetected as a session.
func isSessionContext() bool {
	return os.Getenv("PRACTICE_WORKSPACE_DIR") != "" &&
		os.Getenv("PRACTICE_CONTROL_DIR") != "" &&
		os.Getenv("PRACTICE_STARTED_AT") != ""
}

// returnCmd is `ballroom return`, run from a session's TERMINAL pane to
// get back to the host homepage. Unlike homeCmd, it's only meaningful
// inside a session, so it errors instead of silently falling back to
// tui.Run when there's nothing to return from.
func returnCmd() error {
	if !isSessionContext() {
		return fmt.Errorf("return: not running inside an active practice session (did you mean `ballroom home`?)")
	}
	return returnToHost()
}

// returnToHost ends the practice session so control lands back on the
// host homepage. The container can't reach out and control the host's
// `docker run -it --rm` (orchestrator.RunExercise) directly — no Docker
// client is installed inside it — but that `docker run` is blocking on
// the container's PID 1, which docker/entrypoint.sh sets to `tmux
// attach` after starting the session's tmux server. Killing that server
// tears down every window, which ends the attach client, which exits the
// container, which is what makes `docker run -it --rm` on the host
// return. RunExercise returning is what lets practiceCmd continue on to
// homeCmd and open the homepage picker.
func returnToHost() error {
	if err := exec.Command("tmux", "kill-server").Run(); err != nil {
		return fmt.Errorf("return: tmux kill-server: %w", err)
	}
	return nil
}

func practiceCmd(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: ballroom practice <exercise-id>")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := runExercise(cfg, args[0]); err != nil {
		return err
	}
	// The session container just exited (see returnToHost); land back on
	// the host homepage rather than dropping to a bare shell prompt.
	return homeCmd()
}

func runExercise(cfg config.Config, id string) error {
	ex, err := exercise.Load(cfg.ExercisePath(id))
	if err != nil {
		return fmt.Errorf("unknown exercise %q — run `ballroom help` for usage: %w", id, err)
	}
	return orchestrator.RunExercise(cfg, ex)
}

func sandboxCmd() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	return orchestrator.RunSandbox(cfg)
}

func submitCmd() error {
	startedAtRaw := os.Getenv("PRACTICE_STARTED_AT")
	if startedAtRaw == "" {
		return fmt.Errorf("submit: not running inside a graded exercise session (did you mean to run `ballroom sandbox`?)")
	}
	startedAt, err := time.Parse(time.RFC3339, startedAtRaw)
	if err != nil {
		return fmt.Errorf("submit: parse PRACTICE_STARTED_AT: %w", err)
	}

	cfg := session.Config{
		ControlDir:    os.Getenv("PRACTICE_CONTROL_DIR"),
		WorkspaceDir:  os.Getenv("PRACTICE_WORKSPACE_DIR"),
		TestCommand:   os.Getenv("PRACTICE_TEST_COMMAND"),
		ExerciseID:    os.Getenv("PRACTICE_EXERCISE_ID"),
		Category:      os.Getenv("PRACTICE_CATEGORY"),
		Language:      os.Getenv("PRACTICE_LANGUAGE"),
		StartedAt:     startedAt,
		DBPath:        os.Getenv("PRACTICE_DB_PATH"),
		PollInterval:  200 * time.Millisecond,
		RevealTimeout: 30 * time.Second,
	}

	attempt, err := session.Submit(cfg, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	fmt.Printf("logged attempt #%d\n", attempt.ID)
	return nil
}

// tutorCmd is `ballroom tutor`, launched in the tutor pane by
// docker/entrypoint.sh (env vars below match what it sets — see
// NVIM_SOCKET/OLLAMA_HOST/TUTOR_MODEL/PRACTICE_TUTOR_MODE there, plus
// WORKDIR which every pane shares). Defaults mirror tutor/chat.sh's own
// fallbacks so a standalone run (e.g. local testing outside a real
// session) behaves the same way.
func tutorCmd() error {
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://host.docker.internal:11434"
	}
	model := os.Getenv("TUTOR_MODEL")
	if model == "" {
		model = config.DefaultTutorModel
	}
	mode := os.Getenv("PRACTICE_TUTOR_MODE")
	if mode == "" {
		mode = exercise.TutorModeFullAssist
	}
	workDir := os.Getenv("WORKDIR")
	if workDir == "" {
		workDir = "/workspace"
	}
	maxContextBytes := 8000
	if raw := os.Getenv("TUTOR_FILE_CONTEXT_MAX_BYTES"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil {
			maxContextBytes = n
		}
	}

	cfg := tutor.Config{
		OllamaHost: ollamaHost,
		Model:      model,
		// OrchestratorModel is optional -- an empty value (the default
		// when TUTOR_ORCHESTRATOR_MODEL isn't set) means no routing at
		// all, matching this project's single-model behavior before
		// routing existed (see tutor.Run).
		OrchestratorModel: os.Getenv("TUTOR_ORCHESTRATOR_MODEL"),
		// APIKey is only meaningful when Model or OrchestratorModel is
		// tutor.OpenRouterModelPrefix-prefixed (see agent.go's
		// newChatModel); harmless to always set from the env var
		// regardless, same as OllamaHost being set but unused on that
		// path.
		APIKey:          os.Getenv("OPENROUTER_API_KEY"),
		Mode:            mode,
		WorkDir:         workDir,
		NvimSocket:      os.Getenv("NVIM_SOCKET"),
		MaxContextBytes: maxContextBytes,
	}
	return tutor.Run(context.Background(), cfg, os.Stdin, os.Stdout)
}

// checkToolCallingFn is a var (not a direct call) so tests can
// substitute a fake instead of making a real LLM round-trip — same
// indirection pattern internal/tui/app.go uses for the identical
// reason.
var checkToolCallingFn = tutor.CheckToolCalling

// hostOllamaAddr is where the CLI itself (running on the host, unlike
// tutorCmd's Model/OllamaHost above, which run inside the practice
// container) reaches Ollama — same value as internal/tui/boot.go's own
// unexported ollamaHost const; not worth exporting across packages just
// for this one reuse.
const hostOllamaAddr = "http://localhost:11434"

// configCmd is `ballroom config`, a non-interactive alternative to the
// TUI's Settings tab (internal/tui/app.go) for switching the tutor
// model or OpenRouter API key without opening the picker — useful for
// scripting, or just a faster path when you already know exactly what
// you want to set.
func configCmd(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: ballroom config set-model <tag> | ballroom config set-orchestrator-model <tag|none> | ballroom config set-key <key>")
	}
	switch args[0] {
	case "set-model":
		if len(args) < 2 {
			return fmt.Errorf("usage: ballroom config set-model <tag>")
		}
		return setModelCmd(args[1])
	case "set-orchestrator-model":
		if len(args) < 2 {
			return fmt.Errorf("usage: ballroom config set-orchestrator-model <tag|none>")
		}
		return setOrchestratorModelCmd(args[1])
	case "set-key":
		if len(args) < 2 {
			return fmt.Errorf("usage: ballroom config set-key <key>")
		}
		return setKeyCmd(args[1])
	default:
		return fmt.Errorf("ballroom config: unknown subcommand %q", args[0])
	}
}

// setModelCmd persists tag as the tutor model, preserving the existing
// OpenRouterAPIKey (config.Settings is saved as a whole struct, so
// dropping this would silently wipe a previously-set key — the same
// bug class fixed in the TUI's selectModel when OpenRouter support was
// added). Then, unlike the TUI (which validates tool-calling support
// asynchronously via a background tea.Cmd to stay responsive), checks
// synchronously — a one-shot CLI command has no event loop to keep
// responsive, so there's no reason not to just wait for the real
// answer before printing it.
func setModelCmd(tag string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := config.SaveSettings(cfg.SettingsPath(), config.Settings{
		TutorModel:        tag,
		OpenRouterAPIKey:  cfg.OpenRouterAPIKey,
		OrchestratorModel: cfg.OrchestratorModel,
	}); err != nil {
		return err
	}
	fmt.Printf("tutor model set to %s\n", tag)

	if strings.HasPrefix(tag, tutor.OpenRouterModelPrefix) && cfg.OpenRouterAPIKey == "" {
		fmt.Println("warning: no OpenRouter API key configured yet -- run `ballroom config set-key <key>` or export OPENROUTER_API_KEY")
		return nil
	}

	supported, err := checkToolCallingFn(context.Background(), hostOllamaAddr, tag, cfg.OpenRouterAPIKey)
	switch {
	case err != nil:
		fmt.Printf("warning: checking whether %s supports real tool calling failed: %v\n", tag, err)
	case !supported:
		fmt.Printf("warning: %s may not support real tool calling -- the tutor may not be able to read your code, problem, or tests reliably\n", tag)
	}
	return nil
}

// setKeyCmd persists key as the OpenRouter API key, preserving the
// existing TutorModel (same round-trip concern as setModelCmd, in the
// other direction). Never prompts interactively for anything and never
// echoes the key back -- it's a secret.
func setKeyCmd(key string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := config.SaveSettings(cfg.SettingsPath(), config.Settings{
		TutorModel:        cfg.TutorModel,
		OpenRouterAPIKey:  key,
		OrchestratorModel: cfg.OrchestratorModel,
	}); err != nil {
		return err
	}
	fmt.Println("OpenRouter API key saved")
	return nil
}

// setOrchestratorModelCmd persists tag as the orchestrator model,
// preserving the existing TutorModel and OpenRouterAPIKey (same
// round-trip concern as setModelCmd/setKeyCmd). "none" clears it,
// disabling routing -- internal/tutor.Run treats an empty
// OrchestratorModel as today's single-model behavior.
func setOrchestratorModelCmd(tag string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	orchestratorModel := tag
	if tag == "none" {
		orchestratorModel = ""
	}
	if err := config.SaveSettings(cfg.SettingsPath(), config.Settings{
		TutorModel:        cfg.TutorModel,
		OpenRouterAPIKey:  cfg.OpenRouterAPIKey,
		OrchestratorModel: orchestratorModel,
	}); err != nil {
		return err
	}
	if orchestratorModel == "" {
		fmt.Println("orchestrator model cleared -- routing disabled")
	} else {
		fmt.Printf("orchestrator model set to %s\n", orchestratorModel)
	}
	return nil
}
