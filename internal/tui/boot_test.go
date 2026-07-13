package tui

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/preflight"
	"github.com/JacobTDang/Ballroom/internal/tutor"
)

// fakeBuildImage substitutes buildImageFn in tests so triggering a live
// build doesn't shell out to real docker. Returns a restore func to defer.
func fakeBuildImage(lineCh <-chan string, errCh <-chan error) func() {
	orig := buildImageFn
	buildImageFn = func(config.Config) (<-chan string, <-chan error) {
		return lineCh, errCh
	}
	return func() { buildImageFn = orig }
}

// fakePullModel substitutes pullModelFn in tests so triggering a live
// fallback model pull doesn't make a real HTTP call to Ollama. Returns a
// restore func to defer.
func fakePullModel(lineCh <-chan string, errCh <-chan error) func() {
	orig := pullModelFn
	pullModelFn = func(string, string) (<-chan string, <-chan error) {
		return lineCh, errCh
	}
	return func() { pullModelFn = orig }
}

// noCheckDelay zeroes the pacing delay between checks so tests that drive
// the sequence manually don't actually sleep. Returns a restore func.
func noCheckDelay() func() {
	orig := checkStartDelay
	checkStartDelay = 0
	return func() { checkStartDelay = orig }
}

func TestBootModel_ChecksRunSequentiallyThenReady(t *testing.T) {
	defer noCheckDelay()()

	var order []string
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { order = append(order, "a"); return preflight.Check{Name: "a", OK: true} },
			func() preflight.Check { order = append(order, "b"); return preflight.Check{Name: "b", OK: true} },
		},
	}

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init() should return a command to run the first check")
	}

	// Init batches the first check with the banner's tick command, so
	// unwrap the batch to find the checkDoneMsg among them.
	batch, ok := cmd().(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected Init() to return a batch of commands, got %T", cmd())
	}
	var msg1 tea.Msg
	for _, c := range batch {
		if c == nil {
			continue
		}
		result := c()
		if _, isCheck := result.(checkDoneMsg); isCheck {
			msg1 = result
			break
		}
	}
	if msg1 == nil {
		t.Fatal("expected a checkDoneMsg among Init()'s batched commands")
	}

	newM, cmd2 := m.Update(msg1)
	bm := newM.(bootModel)
	if len(bm.checks) != 1 || bm.checks[0].Name != "a" {
		t.Fatalf("expected 1 check recorded (a), got %+v", bm.checks)
	}
	if bm.ready {
		t.Fatal("should not be ready after only 1 of 2 checks")
	}
	if cmd2 == nil {
		t.Fatal("expected a command to schedule the second check")
	}

	// The next check doesn't run immediately — it's paced behind a
	// startCheckMsg (see checkStartDelay) so checks visibly appear one
	// at a time instead of all resolving within the same frame.
	startMsg, ok := cmd2().(startCheckMsg)
	if !ok {
		t.Fatalf("expected a startCheckMsg to be scheduled, got %T", cmd2())
	}
	newM2a, cmd2b := bm.Update(startMsg)
	bm2a := newM2a.(bootModel)
	if cmd2b == nil {
		t.Fatal("expected startCheckMsg to trigger running the second check")
	}

	msg2 := cmd2b()
	newM2, cmd3 := bm2a.Update(msg2)
	bm2 := newM2.(bootModel)
	if len(bm2.checks) != 2 {
		t.Fatalf("expected 2 checks recorded, got %d", len(bm2.checks))
	}
	if !bm2.ready {
		t.Fatal("expected ready=true once all checks complete")
	}
	if cmd3 != nil {
		t.Error("expected no further command once ready")
	}
	if order[0] != "a" || order[1] != "b" {
		t.Errorf("checks did not run in order: %v", order)
	}
}

func TestBootModel_EnterOnlyQuitsWhenReady(t *testing.T) {
	notReady := bootModel{ready: false}
	_, cmd := notReady.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("enter before ready should be a no-op")
	}

	ready := bootModel{ready: true}
	newM, cmd2 := ready.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd2 == nil {
		t.Fatal("enter when ready should return a quit command")
	}
	if newM.(bootModel).quit {
		t.Error("enter should proceed (quit=false), not request quit")
	}
}

func TestBootModel_QAlwaysRequestsQuit(t *testing.T) {
	m := bootModel{ready: false}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	if cmd == nil {
		t.Fatal("expected q to return a quit command even before ready")
	}
	if !newM.(bootModel).quit {
		t.Error("expected quit=true after pressing q")
	}
}

func TestBootModel_TickAdvancesPhaseAndReschedules(t *testing.T) {
	m := bootModel{}
	newM, cmd := m.Update(tickMsg{})
	if cmd == nil {
		t.Fatal("expected tick to reschedule another tick command")
	}
	if newM.(bootModel).phase != 1 {
		t.Errorf("phase = %d, want 1", newM.(bootModel).phase)
	}
}

func TestBootModel_DockerOkAlwaysTriggersLiveBuildForImageCheck(t *testing.T) {
	lineCh := make(chan string)
	errCh := make(chan error)
	defer fakeBuildImage(lineCh, errCh)()

	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameImage} },
		},
		checks: []preflight.Check{{Name: "Docker", OK: true}},
	}

	newM, cmd := m.Update(checkDoneMsg(preflight.Check{Name: preflight.CheckNameImage}))
	if cmd == nil {
		t.Fatal("expected a command to wait for build output")
	}
	bm := newM.(bootModel)
	if !bm.building {
		t.Fatal("expected building=true")
	}
	if len(bm.checks) != 1 {
		t.Errorf("image check should not be appended yet while building, got %+v", bm.checks)
	}
	if bm.ready {
		t.Error("should not be ready while the image is still building")
	}
}

func TestBootModel_ImageAlreadyBuiltStillTriggersLiveBuild(t *testing.T) {
	// Docker's own layer cache makes a rebuild of an unchanged image
	// fast (every step shows CACHED) — this is the only way to show
	// real build output either way, so it always runs, even when the
	// image already exists.
	lineCh := make(chan string)
	errCh := make(chan error)
	defer fakeBuildImage(lineCh, errCh)()

	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check {
				return preflight.Check{Name: preflight.CheckNameImage, OK: true, Detail: "built"}
			},
		},
		checks: []preflight.Check{{Name: "Docker", OK: true}},
	}
	newM, cmd := m.Update(checkDoneMsg(preflight.Check{Name: preflight.CheckNameImage, OK: true, Detail: "built"}))
	if cmd == nil {
		t.Fatal("expected a command to wait for build output even when the image check already passed")
	}
	bm := newM.(bootModel)
	if !bm.building {
		t.Error("expected building=true so the (cached) build output is still shown")
	}
}

func TestBootModel_DockerNotReachableSkipsBuildAttempt(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: false} },
			func() preflight.Check { return preflight.Check{Name: "Image", OK: false} },
		},
		checks: []preflight.Check{{Name: "Docker", OK: false}},
	}
	newM, _ := m.Update(checkDoneMsg(preflight.Check{Name: "Image", OK: false, Detail: "not built"}))
	bm := newM.(bootModel)
	if bm.building {
		t.Error("should not attempt a build when Docker itself isn't reachable")
	}
	if len(bm.checks) != 2 {
		t.Errorf("expected the image check to be appended as a plain failure, got %+v", bm.checks)
	}
}

func TestBootModel_ModelCheckFailureWithOllamaOKTriggersLivePullFallback(t *testing.T) {
	lineCh := make(chan string)
	errCh := make(chan error)
	defer fakePullModel(lineCh, errCh)()

	m := bootModel{
		cfg: config.Config{TutorModel: "deepseek-coder-v2:16b-lite-instruct-q4_K_M"},
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameImage, OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameOllama, OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameModel, OK: false} },
		},
		checks: []preflight.Check{
			{Name: "Docker", OK: true},
			{Name: preflight.CheckNameImage, OK: true},
			{Name: preflight.CheckNameOllama, OK: true},
		},
	}

	newM, cmd := m.Update(checkDoneMsg(preflight.Check{Name: preflight.CheckNameModel, OK: false, Detail: "deepseek-coder-v2:... not pulled"}))
	if cmd == nil {
		t.Fatal("expected a command to wait for pull output")
	}
	bm := newM.(bootModel)
	if !bm.pullingModel {
		t.Fatal("expected pullingModel=true")
	}
	if len(bm.checks) != 3 {
		t.Errorf("model check should not be appended yet while pulling, got %+v", bm.checks)
	}
	if bm.ready {
		t.Error("should not be ready while the fallback model is still pulling")
	}
}

func TestBootModel_ModelCheckFailureWithOllamaNotOKSkipsPullFallback(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameImage, OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameOllama, OK: false} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameModel, OK: false} },
		},
		checks: []preflight.Check{
			{Name: "Docker", OK: true},
			{Name: preflight.CheckNameImage, OK: true},
			{Name: preflight.CheckNameOllama, OK: false},
		},
	}
	newM, _ := m.Update(checkDoneMsg(preflight.Check{Name: preflight.CheckNameModel, OK: false, Detail: "can't reach Ollama to check"}))
	bm := newM.(bootModel)
	if bm.pullingModel {
		t.Error("should not attempt a fallback pull when Ollama itself isn't reachable")
	}
	if len(bm.checks) != 4 {
		t.Errorf("expected the model check to be appended as a plain failure, got %+v", bm.checks)
	}
}

func TestBootModel_ModelCheckOKSkipsPullFallback(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameImage, OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameOllama, OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameModel, OK: true} },
		},
		checks: []preflight.Check{
			{Name: "Docker", OK: true},
			{Name: preflight.CheckNameImage, OK: true},
			{Name: preflight.CheckNameOllama, OK: true},
		},
	}
	newM, _ := m.Update(checkDoneMsg(preflight.Check{Name: preflight.CheckNameModel, OK: true, Detail: "ready"}))
	bm := newM.(bootModel)
	if bm.pullingModel {
		t.Error("should not attempt a fallback pull when the model check already passed")
	}
	if len(bm.checks) != 4 || !bm.checks[3].OK {
		t.Errorf("expected a passing model check appended, got %+v", bm.checks)
	}
}

// TestNewBootModel_OpenRouterTutorModelSkipsLocalOllamaModelCheck is a
// regression test for a real bug found live: preflight.CheckModel only
// ever checks Ollama's own local /api/tags -- it has no OpenRouter
// awareness at all, so it always reported an OpenRouter-prefixed worker
// model as "not pulled". That always triggered the pull-fallback path
// below (Ollama itself was reachable), which pulled config.DefaultTutorModel
// and overwrote cfg.TutorModel for the rest of the session -- and that
// corrupted value then got persisted back to settings.json the next
// time *any* setting was saved, permanently reverting the user's real
// OpenRouter pick back to the local default. An OpenRouter model was
// never a candidate for a local Ollama pull in the first place, so the
// check itself must never run for one.
func TestNewBootModel_OpenRouterTutorModelSkipsLocalOllamaModelCheck(t *testing.T) {
	openRouterModel := tutor.OpenRouterModelPrefix + "nvidia/nemotron-3-ultra-550b-a55b:free"
	m := newBootModel(config.Config{TutorModel: openRouterModel})

	if len(m.pending) != 4 {
		t.Fatalf("pending has %d checks, want 4", len(m.pending))
	}
	check := m.pending[3]()
	if check.Name != preflight.CheckNameModel {
		t.Fatalf("check.Name = %q, want %q", check.Name, preflight.CheckNameModel)
	}
	if !check.OK {
		t.Errorf("check.OK = false for an OpenRouter model, want true -- it was never a candidate for a local Ollama pull check")
	}
}

// TestBootModel_OpenRouterTutorModelNeverTriggersPullFallbackOrCorruptsCfg
// exercises the same bug end to end through Update(), confirming the
// downstream pull-fallback path never fires and cfg.TutorModel survives
// the boot sequence unchanged.
func TestBootModel_OpenRouterTutorModelNeverTriggersPullFallbackOrCorruptsCfg(t *testing.T) {
	openRouterModel := tutor.OpenRouterModelPrefix + "nvidia/nemotron-3-ultra-550b-a55b:free"
	m := newBootModel(config.Config{TutorModel: openRouterModel})
	m.checks = []preflight.Check{
		{Name: "Docker", OK: true},
		{Name: preflight.CheckNameImage, OK: true},
		{Name: preflight.CheckNameOllama, OK: true},
	}

	check := m.pending[3]() // the real check function this model uses, not a fake
	newM, cmd := m.Update(checkDoneMsg(check))
	bm := newM.(bootModel)

	if bm.pullingModel {
		t.Error("an OpenRouter tutor model must never trigger the local-Ollama pull fallback")
	}
	if bm.cfg.TutorModel != openRouterModel {
		t.Errorf("cfg.TutorModel = %q, want it left untouched at %q", bm.cfg.TutorModel, openRouterModel)
	}
	// This was the last of the 4 checks, so advance() sets ready=true and
	// returns a nil cmd (nothing left to schedule) -- not evidence of a
	// stuck/failed advance.
	if !bm.ready {
		t.Error("expected the check to advance normally and reach ready")
	}
	if cmd != nil {
		t.Errorf("expected a nil cmd (no more checks queued after the last one), got %v", cmd)
	}
	if len(bm.checks) != 4 || !bm.checks[3].OK {
		t.Errorf("expected a passing model check appended, got %+v", bm.checks)
	}
}

func TestBootModel_PullLineCapsAtThreeLinesTotal(t *testing.T) {
	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	m := bootModel{pullingModel: true, pullLineCh: lineCh, pullErrCh: errCh}

	for i := 0; i < maxOutputLines+5; i++ {
		newM, cmd := m.Update(pullLineMsg(fmt.Sprintf("line %d", i)))
		if cmd == nil {
			t.Fatal("expected pullLineMsg to keep listening for more output")
		}
		m = newM.(bootModel)
	}
	if len(m.pullLines) != maxOutputLines {
		t.Errorf("pullLines = %d, want capped at %d", len(m.pullLines), maxOutputLines)
	}
	if m.pullLines[len(m.pullLines)-1] != fmt.Sprintf("line %d", maxOutputLines+4) {
		t.Errorf("expected the most recent line to be kept, got %+v", m.pullLines)
	}
}

func TestBootModel_PullDoneSuccessAdvancesAsPassingModelCheckAndSwitchesCfg(t *testing.T) {
	m := bootModel{
		cfg: config.Config{TutorModel: "deepseek-coder-v2:16b-lite-instruct-q4_K_M"},
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameModel, OK: true} },
		},
		checks:           []preflight.Check{{Name: "Docker", OK: true}},
		pullingModel:     true,
		pullFallbackFrom: "deepseek-coder-v2:16b-lite-instruct-q4_K_M",
		pullLines:        []string{"pulling manifest", "success"},
	}

	newM, cmd := m.Update(pullDoneMsg{err: nil})
	if cmd != nil {
		t.Fatal("expected no further pending checks after the last one resolves")
	}
	bm := newM.(bootModel)
	if bm.pullingModel {
		t.Error("expected pullingModel=false once the pull resolves")
	}
	if len(bm.pullLines) != 0 {
		t.Error("expected the pull log to clear once resolved (panel collapses)")
	}
	if len(bm.checks) != 2 || !bm.checks[1].OK {
		t.Fatalf("expected a passing model check appended, got %+v", bm.checks)
	}
	if bm.cfg.TutorModel != config.DefaultTutorModel {
		t.Errorf("cfg.TutorModel = %q, want it switched to the default %q for this session", bm.cfg.TutorModel, config.DefaultTutorModel)
	}
	if !bm.ready {
		t.Error("expected ready=true once the (now only) remaining check resolves")
	}
}

func TestBootModel_PullDoneFailureAdvancesAsFailingModelCheckAndLeavesCfgUnchanged(t *testing.T) {
	m := bootModel{
		cfg: config.Config{TutorModel: "deepseek-coder-v2:16b-lite-instruct-q4_K_M"},
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: preflight.CheckNameModel, OK: true} },
		},
		checks:           []preflight.Check{{Name: "Docker", OK: true}},
		pullingModel:     true,
		pullFallbackFrom: "deepseek-coder-v2:16b-lite-instruct-q4_K_M",
	}

	newM, _ := m.Update(pullDoneMsg{err: errors.New("boom")})
	bm := newM.(bootModel)
	if bm.pullingModel {
		t.Error("expected pullingModel=false once the pull resolves, even on failure")
	}
	if len(bm.checks) != 2 || bm.checks[1].OK {
		t.Fatalf("expected a failing model check appended, got %+v", bm.checks)
	}
	if bm.cfg.TutorModel != "deepseek-coder-v2:16b-lite-instruct-q4_K_M" {
		t.Errorf("expected cfg.TutorModel to stay unchanged on pull failure, got %q", bm.cfg.TutorModel)
	}
	if !bm.ready {
		t.Error("expected ready=true so the user isn't stuck")
	}
}

func TestBootModel_RenderRightColumnShowsPullingModelPanel(t *testing.T) {
	m := bootModel{
		checks:       []preflight.Check{{Name: "Docker", OK: true}},
		pullingModel: true,
		pullLines:    []string{"pulling manifest", "downloading (42%)"},
	}
	out := m.renderRightColumn()
	if !strings.Contains(out, "pulling manifest") {
		t.Errorf("expected live pull output visible, got:\n%s", out)
	}
	if !strings.Contains(out, "downloading (42%)") {
		t.Errorf("expected live pull output visible, got:\n%s", out)
	}
}

func TestStepID_ExtractsLeadingHashNumberToken(t *testing.T) {
	cases := map[string]string{
		"#5 [2/4] RUN go build ./...": "#5",
		"#12 0.235 exporting layers":  "#12",
		"#1 DONE 0.1s":                "#1",
	}
	for line, want := range cases {
		if got := stepID(line); got != want {
			t.Errorf("stepID(%q) = %q, want %q", line, got, want)
		}
	}
}

func TestStepID_ReturnsEmptyForNonStepLine(t *testing.T) {
	cases := []string{"", "no leading hash", "#", "#abc not digits"}
	for _, line := range cases {
		if got := stepID(line); got != "" {
			t.Errorf("stepID(%q) = %q, want empty", line, got)
		}
	}
}

func TestBootModel_BuildLineGroupsByStepAndCapsEachStepAtThreeLines(t *testing.T) {
	lineCh := make(chan string, 1)
	errCh := make(chan error, 1)
	m := bootModel{building: true, buildLineCh: lineCh, buildErrCh: errCh}

	for i := 0; i < maxStepLogLines+5; i++ {
		newM, cmd := m.Update(buildLineMsg("#3 some output"))
		if cmd == nil {
			t.Fatal("expected buildLineMsg to keep listening for more output")
		}
		m = newM.(bootModel)
	}
	if len(m.buildSteps) != 1 {
		t.Fatalf("expected all lines from the same step to stay grouped as 1 entry, got %d", len(m.buildSteps))
	}
	if len(m.buildSteps[0].lines) != maxStepLogLines {
		t.Errorf("step lines = %d, want capped at %d", len(m.buildSteps[0].lines), maxStepLogLines)
	}
}

func TestBootModel_BuildLinePersistsPriorStepsWhenNewStepStarts(t *testing.T) {
	m := bootModel{building: true}

	newM, _ := m.Update(buildLineMsg("#1 [1/4] FROM golang"))
	m = newM.(bootModel)
	newM, _ = m.Update(buildLineMsg("#2 [2/4] RUN go build"))
	m = newM.(bootModel)

	if len(m.buildSteps) != 2 {
		t.Fatalf("expected both steps to persist once a new step starts, got %d: %+v", len(m.buildSteps), m.buildSteps)
	}
	if m.buildSteps[0].id != "#1" || m.buildSteps[1].id != "#2" {
		t.Errorf("expected steps in order #1, #2, got %+v", m.buildSteps)
	}
}

func TestBootModel_BuildDoneSuccessAdvancesAsPassingImageCheck(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: "Image", OK: true} },
		},
		checks:     []preflight.Check{{Name: "Docker", OK: true}},
		building:   true,
		buildSteps: []buildStepLog{{id: "#1", lines: []string{"some output"}}},
	}

	newM, cmd := m.Update(buildDoneMsg{err: nil})
	if cmd != nil {
		t.Fatal("expected no further pending checks after the last one resolves")
	}
	bm := newM.(bootModel)
	if bm.building {
		t.Error("expected building=false once the build resolves")
	}
	if len(bm.buildSteps) != 0 {
		t.Error("expected the build log to clear once resolved (panel collapses)")
	}
	if len(bm.checks) != 2 || !bm.checks[1].OK {
		t.Fatalf("expected a passing Image check appended, got %+v", bm.checks)
	}
	if !bm.ready {
		t.Error("expected ready=true once the (now only) remaining check resolves")
	}
}

func TestBootModel_BuildDoneFailureAdvancesAsFailingImageCheck(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker", OK: true} },
			func() preflight.Check { return preflight.Check{Name: "Image", OK: true} },
		},
		checks:   []preflight.Check{{Name: "Docker", OK: true}},
		building: true,
	}

	newM, _ := m.Update(buildDoneMsg{err: errors.New("boom")})
	bm := newM.(bootModel)
	if bm.building {
		t.Error("expected building=false once the build resolves, even on failure")
	}
	if len(bm.checks) != 2 || bm.checks[1].OK {
		t.Fatalf("expected a failing Image check appended, got %+v", bm.checks)
	}
	if !bm.ready {
		t.Error("expected ready=true so the user isn't stuck — they can see the failure and quit")
	}
}

func TestLastLines_CapsToMostRecentNNonEmptyLines(t *testing.T) {
	out := "one\ntwo\n\nthree\nfour\nfive"
	got := lastLines(out, 3)
	want := []string{"three", "four", "five"}
	if len(got) != len(want) {
		t.Fatalf("lastLines = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("lastLines[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestLastLines_EmptyOutputReturnsNoLines(t *testing.T) {
	if got := lastLines("", 3); len(got) != 0 {
		t.Errorf("lastLines(\"\", 3) = %v, want empty", got)
	}
}

func TestLastBuildLines_CapsAcrossAllStepsNotPerStep(t *testing.T) {
	steps := []buildStepLog{
		{id: "#1", lines: []string{"a", "b", "c"}},
		{id: "#2", lines: []string{"d", "e", "f"}},
		{id: "#3", lines: []string{"g", "h", "i"}},
	}
	got := lastBuildLines(steps, 3)
	want := []string{"g", "h", "i"}
	if len(got) != len(want) {
		t.Fatalf("lastBuildLines = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("lastBuildLines[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestBootModel_RenderRightColumnCapsBuildLogToThreeLinesTotal(t *testing.T) {
	m := bootModel{
		checks:   []preflight.Check{{Name: "Docker", OK: true}},
		building: true,
		buildSteps: []buildStepLog{
			{id: "#1", lines: []string{"step one line"}},
			{id: "#2", lines: []string{"step two line"}},
			{id: "#3", lines: []string{"step three line"}},
			{id: "#4", lines: []string{"step four line"}},
		},
	}
	out := m.renderRightColumn()
	if strings.Contains(out, "step one line") {
		t.Errorf("expected the build log to cap at the last 3 lines total across all steps, got:\n%s", out)
	}
	if !strings.Contains(out, "step four line") {
		t.Errorf("expected the most recent build line to be visible, got:\n%s", out)
	}
}

func TestBootModel_RenderRightColumnShowsRealCommandAndOutputForRecentCheck(t *testing.T) {
	m := bootModel{
		checks: []preflight.Check{
			{Name: "Docker daemon", OK: true, Detail: "running", Command: "docker info", Output: "Client:\n Version: 27.0.0\nServer:\n Containers: 3"},
		},
	}
	out := m.renderRightColumn()
	if !strings.Contains(out, "docker info") {
		t.Errorf("expected the real command to be visible, got:\n%s", out)
	}
	if !strings.Contains(out, "Containers: 3") {
		t.Errorf("expected real command output to be visible, got:\n%s", out)
	}
}

func TestBootModel_RenderRightColumnNeverCollapsesEarlierChecks(t *testing.T) {
	m := bootModel{
		checks: []preflight.Check{
			{Name: "Docker daemon", OK: true, Detail: "running", Command: "docker info", Output: "oldest output line"},
			{Name: "Practice image", OK: true, Detail: "built", Command: "docker image inspect x", Output: "sha256:abc"},
			{Name: "Ollama", OK: true, Detail: "reachable", Command: "GET /api/tags", Output: `{"models":[]}`},
			{Name: "Tutor model", OK: true, Detail: "ready", Command: "GET /api/tags", Output: `{"models":[{"name":"x"}]}`},
		},
	}
	out := m.renderRightColumn()
	if !strings.Contains(out, "oldest output line") {
		t.Errorf("expected the oldest check to stay expanded and keep showing its output, got:\n%s", out)
	}
	if !strings.Contains(out, "sha256:abc") {
		t.Errorf("expected a later check to also still show its output, got:\n%s", out)
	}
}

func TestBootModel_RenderRightColumnQueuedCheckShowsNameOnlyNoCommand(t *testing.T) {
	m := bootModel{
		pending: []func() preflight.Check{
			func() preflight.Check { return preflight.Check{Name: "Docker"} },
			func() preflight.Check { return preflight.Check{Name: "Ollama"} },
		},
		checkNames: []string{"Docker daemon", "Ollama"},
	}
	out := m.renderRightColumn()
	if !strings.Contains(out, "Ollama") {
		t.Errorf("expected the queued check's name to be visible, got:\n%s", out)
	}
	if strings.Contains(out, "$") {
		t.Errorf("expected no command to be shown for a check that hasn't been invoked yet, got:\n%s", out)
	}
}

func TestBootModel_CtrlCAlwaysRequestsQuit(t *testing.T) {
	m := bootModel{ready: false}
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected ctrl+c to return a quit command")
	}
	if !newM.(bootModel).quit {
		t.Error("expected quit=true after ctrl+c")
	}
}
