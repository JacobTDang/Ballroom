package tui

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/draft"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/preflight"
	"github.com/JacobTDang/Ballroom/internal/tracker"
	"github.com/JacobTDang/Ballroom/internal/tutor"
)

// listModelsFn and checkModelFn are vars (not direct calls) so tests can
// substitute fakes instead of making real HTTP calls to Ollama — same
// indirection pattern as buildImageFn in boot.go.
var listModelsFn = preflight.ListModels
var checkModelFn = preflight.CheckModel

// checkToolCallingFn is checkModelFn's counterpart for the deeper
// "does this model actually make real tool calls" question — a real
// Ollama round-trip (see tutor.CheckToolCalling), so it always runs as
// a background tea.Cmd (checkToolCallingCmd) after a pick completes
// rather than blocking selection on it.
var checkToolCallingFn = tutor.CheckToolCalling

// claimsToolSupportFn is the cheap metadata counterpart to
// checkToolCallingFn -- a var for the same test-substitution reason.
var claimsToolSupportFn = tutor.ClaimsToolSupport

// suggestedModels are known-good tutor model tags surfaced in the picker
// even before they're pulled locally, so they're discoverable without
// already knowing their exact Ollama tag. Selecting one that isn't
// pulled yet goes through the same not-pulled warn-then-confirm flow as
// typing an arbitrary tag (see handleModelEnter).
var suggestedModels = []string{
	config.DeepSeekCoderV2LiteModel,
	config.Qwen25Coder14BModel,
}

// suggestedOpenRouterModels are free-tier OpenRouter model slugs (no
// OpenRouterModelPrefix — stageAPIModelEntry adds it on select) verified
// live, on 2026-07-12 (poolside/laguna-xs-2.1:free added and verified
// 2026-07-13), via tutor.CheckToolCalling to actually make real
// tool calls, not just declared "tools"-capable in OpenRouter's
// /models metadata (that distinction mattered in this codebase before —
// see config.Qwen25Coder14BModel's doc comment for a model that claims
// tool support but doesn't deliver it end-to-end). Kept short and
// re-verified rather than exhaustive: OpenRouter's free-tier catalog
// changes over time, and an unverified entry here would be actively
// misleading in a list whose whole point is "known to work".
var suggestedOpenRouterModels = []string{
	"openai/gpt-oss-120b:free",
	"openai/gpt-oss-20b:free",
	"nvidia/nemotron-3-ultra-550b-a55b:free",
	"nvidia/nemotron-3-super-120b-a12b:free",
	"nvidia/nemotron-3-nano-30b-a3b:free",
	"nvidia/nemotron-nano-9b-v2:free",
	"poolside/laguna-xs-2.1:free",
}

// browsableModels merges local (locally pulled) with suggestedModels,
// deduplicated — a suggested tag already pulled just appears once, as a
// normal local entry.
func browsableModels(local []string) []string {
	seen := make(map[string]bool, len(local))
	out := make([]string, 0, len(local)+len(suggestedModels))
	for _, name := range local {
		seen[name] = true
		out = append(out, name)
	}
	for _, name := range suggestedModels {
		if !seen[name] {
			out = append(out, name)
		}
	}
	return out
}

// modelsLoadedMsg carries the result of querying Ollama's /api/tags for
// the locally pulled model list.
type modelsLoadedMsg struct {
	models []string
	err    error
}

// toolCallingCheckMsg carries the result of a background
// tutor.CheckToolCalling run kicked off by selectModel. model is the
// tag it was checked against — the Update handler drops the result if
// it no longer matches cfg.TutorModel, so a check for a pick the user
// has since replaced with another one can't clobber the current
// warning state.
type toolCallingCheckMsg struct {
	model     string
	supported bool
	err       error
}

// checkToolCallingCmd runs tutor.CheckToolCalling in the background —
// a real LLM round-trip, so this must never block Update the way
// checkModelFn's cheap /api/tags lookup can. ollamaHost is unused (but
// still passed) when model is OpenRouterModelPrefix-prefixed — see
// newChatModel's provider branch, which apiKey and ollamaHost are both
// threaded through to.
func checkToolCallingCmd(model, apiKey string) tea.Cmd {
	return func() tea.Msg {
		supported, err := checkToolCallingFn(context.Background(), ollamaHost, model, apiKey)
		return toolCallingCheckMsg{model: model, supported: supported, err: err}
	}
}

// toolSupportMsg carries the metadata pre-check for the picker's rows:
// the set of models whose provider metadata explicitly DENIES tool
// support. Models the metadata couldn't be fetched for (or doesn't
// list) are simply absent -- no signal is never rendered as a denial.
// Advisory only: metadata can rule a model out but not in (a model has
// claimed tools in metadata while failing 6/6 live probes), so
// checkToolCallingFn's live probe on selection stays the ground truth.
type toolSupportMsg struct {
	noTools map[string]bool
}

// toolSupportCmd checks each listed model's metadata in the background.
// All local Ollama calls (the picker lists /api/tags results), so a
// short overall deadline keeps a wedged daemon from stalling the badge
// forever -- on timeout the remaining models just carry no marker.
func toolSupportCmd(models []string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		noTools := map[string]bool{}
		for _, name := range models {
			if claims, known := claimsToolSupportFn(ctx, ollamaHost, name, ""); known && !claims {
				noTools[name] = true
			}
		}
		return toolSupportMsg{noTools: noTools}
	}
}

// filterModels returns the models whose name contains filter
// (case-insensitive). An empty filter matches everything.
func filterModels(models []string, filter string) []string {
	if filter == "" {
		return models
	}
	lower := strings.ToLower(filter)
	var out []string
	for _, name := range models {
		if strings.Contains(strings.ToLower(name), lower) {
			out = append(out, name)
		}
	}
	return out
}

// appStage is which screen the merged program is currently showing —
// every stage renders inside the same renderDashboardPanel shell (disco
// ball + animated banner + bordered panel) so navigating between them
// never tears down and relaunches a differently-styled program.
type appStage int

const (
	stageMain appStage = iota
	stageCategories
	stageDSACategories
	stageProblems
	stageLanguage
	// stageResumeDraft asks whether to resume a saved solution draft or
	// start from the exercise's starter -- see resumedraft.go. Entered
	// from launchExercise, which every launch path funnels through.
	stageResumeDraft
	stageStats
	// stageSettings is the Settings menu entry's landing screen: a
	// Worker model / Orchestrator model role choice (see settingsTarget)
	// — which one is being edited then decides where stageProviderChoice
	// and everything past it write their result (selectModel).
	stageSettings
	// stageProviderChoice is the Local (Ollama) / API (OpenRouter)
	// choice for whichever role stageSettings just picked. Routes to
	// stageModelPicker (Local) or stageAPIModelEntry (API) -- plus a
	// third "None (disable routing)" option, only when the role being
	// edited is the orchestrator (there must always be a worker model,
	// so that option never appears for it).
	stageProviderChoice
	stageModelPicker
	// stageAPIModelEntry asks for a bare OpenRouter model slug (no
	// OpenRouterModelPrefix needed — provider is already established by
	// the point this shows) when stageProviderChoice -> API is chosen.
	// Enter delegates to the same selectModelOrPromptForKey
	// handleModelEnter already uses for a directly-typed openrouter: tag
	// in the local picker.
	stageAPIModelEntry
	// stageOpenRouterKeyEntry asks for an OpenRouter API key the first
	// time an OpenRouterModelPrefix-prefixed model is picked with none
	// available yet (settings.json nor OPENROUTER_API_KEY) — see
	// handleModelEnter/selectModelOrPromptForKey.
	stageOpenRouterKeyEntry
)

// settingsTarget tracks which Config field stageProviderChoice and
// everything past it (stageModelPicker, stageAPIModelEntry,
// stageOpenRouterKeyEntry) are editing — set by updateSettings when the
// user picks Worker or Orchestrator from stageSettings' top-level list,
// read by selectModel to decide whether to write cfg.TutorModel or
// cfg.OrchestratorModel.
type settingsTarget int

const (
	settingsTargetWorker settingsTarget = iota
	settingsTargetOrchestrator
)

// appOutcome is what Run's caller should do once the program exits.
// Sandbox and a chosen language variant both hand the terminal to
// `docker run -it` (orchestrator.RunSandbox/RunExercise) — bubbletea can't
// render inside that external interactive process, so those two are the
// only cases that actually tear the program down.
type appOutcome int

const (
	outcomeNone appOutcome = iota
	outcomeRunExercise
	outcomeRunSandbox
)

// catalogListFn and recentAttemptsFn are vars (not direct calls) so tests
// can substitute fakes instead of touching a real exercises dir / sqlite
// db — same indirection pattern as listModelsFn/checkModelFn/buildImageFn.
var catalogListFn = catalog.List
var recentAttemptsFn = recentAttempts
var allAttemptsFn = allAttempts

// recentAttemptsLimit caps how many recent attempts Stats shows.
const recentAttemptsLimit = 10

// menuChoice is one of the main menu options.
type menuChoice int

const (
	menuPractice menuChoice = iota
	menuDaily
	menuSandbox
	menuStats
	menuSettings
)

var menuLabels = []string{"Practice", "Daily", "Sandbox", "Stats", "Settings"}

var menuDescriptions = []string{
	"Pick a category and work through exercises",
	"Today's pick — due or unsolved, same one all day",
	"Free practice, no grading, persists across sessions",
	"See your progress across categories",
	"Choose your worker and orchestrator models — local (Ollama) or API (OpenRouter)",
}

// menuRightColWidth is the fixed content width of the right column —
// wide enough for the longest line (the keybinding hint) — so the
// selected row's highlight reads as a full-width bar rather than a
// tight box around just the label text.
const menuRightColWidth = 54

// menuRowHighlight matches cursorRowStyle exactly, so the main menu's
// selected row and every other stage's selected row read as the same
// highlight — reused directly instead of redefining an identical style.
var (
	menuSubtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8B8680"))
	menuRowHighlight  = cursorRowStyle
)

// appResume tells newAppModel where to pick back up after a docker
// handoff (Sandbox or an exercise session) returns — stageMain has no
// resume data, stageProblems resumes straight into the category just
// practiced instead of making the user re-pick it from the main menu.
// launchErr carries a session that failed to even start (e.g. Docker
// went down between the boot check and picking a problem) back into
// the resumed screen instead of it being printed to a stderr line the
// very next alt-screen program immediately wipes (issue #230) -- RunApp
// seeds appModel.err from it so renderProblems shows it the same way
// every other screen-level failure renders (see renderFriendlyError).
type appResume struct {
	stage     appStage
	category  string
	launchErr error
}

// appModel is the single bubbletea program behind the whole menu tree:
// main menu, Practice's category/problem/language chain, Stats, and the
// model picker. Every stage keeps its own cursor/data fields (rather than
// reusing one generic set) so switching stages can never leave a stale
// cursor pointing at the wrong list.
type appModel struct {
	cfg   config.Config
	stage appStage

	phase         int
	width, height int
	quit          bool
	err           error

	// stageMain
	cursor int
	// toolCallingWarning is set once a background tutor.CheckToolCalling
	// (kicked off by selectModel) reports the currently selected model
	// doesn't make real tool calls — shown on the main menu since the
	// check resolves after the picker has already returned there.
	toolCallingWarning string

	// stageCategories / source data for stageProblems
	problems       []catalog.ProblemStatus
	categories     []string
	categoryCursor int

	// The home dashboard's data (renderHomeboard) -- problems above is
	// shared with the pickers; homeAttempts is the full attempt history,
	// loaded once at construction. homeLoadErr keeps a load failure
	// visible as a dim notice instead of silently hiding the dashboard.
	homeAttempts []tracker.Attempt
	homeLoadErr  error

	// stageDSACategories — second-level picker shown when the top-level
	// selection is the grouped DSA entry, listing the NeetCode roadmap
	// subcategories (Arrays & Hashing, Two Pointers, ...) it collapses.
	dsaCategories     []string
	dsaCategoryCursor int

	// stageProblems
	category         string
	categoryProblems []catalog.ProblemStatus
	problemCursor    int
	// problemFilter is the picker's live type-to-search text — same UX
	// as the model picker's filter (see updateModelPicker's KeyRunes
	// case): letters feed it, backspace trims it, and it resets on
	// every entry into stageProblems. The cursor indexes
	// visibleProblems(), not categoryProblems, while a filter is live.
	problemFilter string

	// stageLanguage
	selectedProblem catalog.ProblemStatus
	langCursor      int

	// stageResumeDraft — the exercise waiting on the resume/fresh
	// answer, the draft it would resume, and which row is selected.
	// draftDirToUse is the answer, read by the caller after Quit: the
	// draft dir to overlay, or empty for the pristine starter.
	pendingExercise exercise.Exercise
	pendingDraft    draft.Draft
	resumeCursor    int
	draftDirToUse   string

	// stageStats
	statsStatuses   []catalog.ExerciseStatus
	statsRecent     []tracker.Attempt
	statsWeakDims   []catalog.DimensionWeakness
	statsCodingWeak []catalog.CategoryWeakness

	// stageSettings: settingsCursor is 0 = Worker model, 1 = Orchestrator
	// model there; reused by stageProviderChoice as 0 = Local (Ollama),
	// 1 = API (OpenRouter), 2 = None (disable routing, orchestrator only)
	// — the two stages are never on screen at once, so one cursor field
	// covers both without ambiguity.
	settingsCursor  int
	settingsEditing settingsTarget // which Config field stageProviderChoice onward writes to (see selectModel)

	// stageAPIModelEntry: apiModelInput both filters apiModelFiltered
	// (derived from suggestedOpenRouterModels, same shape as
	// modelFilter/modelFiltered below) and, when nothing matches, is the
	// custom slug typed directly.
	apiModelInput    string
	apiModelFiltered []string
	apiModelCursor   int

	// stageModelPicker
	modelLoading  bool
	modelLoadErr  error
	localModels   []string        // exactly what Ollama reports pulled, for isLocalModel
	models        []string        // localModels + suggestedModels, deduplicated
	modelNoTools  map[string]bool // metadata explicitly denies tools (toolSupportMsg)
	modelFilter   string
	modelFiltered []string
	modelCursor   int
	modelWarning  string

	// modelPendingDownloadTag etc. mirror bootModel's pullingModel/
	// pullLines fields in boot.go — same pullLineMsg/pullDoneMsg/
	// waitForPullLine/pullModelFn/maxOutputLines machinery, reused as-is
	// (same package) for a live download triggered from the picker
	// itself instead of boot's own fallback.
	modelPendingDownloadTag string
	modelDownloading        bool
	modelDownloadTarget     string
	modelDownloadLines      []string
	modelDownloadLineCh     <-chan string
	modelDownloadErrCh      <-chan error

	// stageOpenRouterKeyEntry
	openRouterPendingModel string // the openrouter: tag waiting on a key before it can be selected
	openRouterKeyInput     string

	// outcome is read by Run() once the program exits.
	outcome       appOutcome
	exerciseToRun exercise.Exercise
}

// newAppModel builds the starting model for one Program lifetime. resume
// lets a fresh program (after a docker handoff) pick back up at
// stageProblems for the category just practiced instead of dropping back
// to the main menu.
func newAppModel(cfg config.Config, resume appResume) appModel {
	m := appModel{cfg: cfg}
	// The home dashboard's data: local disk + sqlite reads, same work
	// loadPractice does on Enter -- done once up front so the menu is a
	// status page from the first frame (and fresh again after every
	// session, since the launcher rebuilds the model on return). Errors
	// don't block the menu; renderMain shows them as a dim notice.
	if statuses, err := catalogListFn(m.cfg); err != nil {
		m.homeLoadErr = err
	} else {
		m.problems = catalog.GroupByProblem(statuses)
	}
	if attempts, err := allAttemptsFn(m.cfg); err != nil {
		m.homeLoadErr = err
	} else {
		m.homeAttempts = attempts
	}
	if resume.stage == stageProblems {
		m = m.loadPractice()
		if m.err == nil {
			m.category = resume.category
			m.categoryProblems = filterByCategory(m.problems, resume.category)
			m.problemCursor = 0
			m.problemFilter = ""
			m.stage = stageProblems
		}
	}
	return m
}

func (m appModel) Init() tea.Cmd {
	return tickCmd()
}

func (m appModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, tea.ClearScreen

	case tickMsg:
		m.phase++
		return m, tickCmd()

	case modelsLoadedMsg:
		m.modelLoading = false
		m.modelLoadErr = msg.err
		m.localModels = msg.models
		m.models = browsableModels(msg.models)
		m.modelFiltered = filterModels(m.models, m.modelFilter)
		if msg.err != nil || len(msg.models) == 0 {
			return m, nil
		}
		// Only the locally pulled models are worth a metadata call --
		// suggested-but-not-pulled entries have no /api/show record.
		return m, toolSupportCmd(msg.models)

	case toolSupportMsg:
		m.modelNoTools = msg.noTools
		return m, nil

	case pullLineMsg:
		line := string(msg)
		m.modelDownloadLines = append(m.modelDownloadLines, line)
		if len(m.modelDownloadLines) > maxOutputLines {
			m.modelDownloadLines = m.modelDownloadLines[len(m.modelDownloadLines)-maxOutputLines:]
		}
		return m, waitForPullLine(m.modelDownloadLineCh, m.modelDownloadErrCh)

	case pullDoneMsg:
		m.modelDownloading = false
		if msg.err != nil {
			m.modelWarning = fmt.Sprintf("couldn't download that model — %v", msg.err)
			m.modelDownloadLines = nil
			return m, nil
		}
		return m.selectModel(m.modelDownloadTarget)

	case toolCallingCheckMsg:
		if msg.model != m.cfg.TutorModel {
			// The user picked something else again before this
			// resolved — stale, ignore rather than clobbering the
			// warning state for whatever is actually selected now.
			return m, nil
		}
		switch {
		case msg.err != nil:
			// A real bug found live: msg.err was captured on the struct
			// but never actually looked at here, so a hard rejection
			// (e.g. Ollama returning 400 "does not support tools" for a
			// model picked without real tool-calling support) produced
			// no warning at all — the picker went silent, and the
			// problem only surfaced once inside a live tutor session as
			// a confusing "could not reach" error that read like a
			// network/Docker problem. An error here is at least as
			// informative as the "ran cleanly but didn't call it" case
			// below, often more so (a hard rejection is more certain
			// than an inferred non-support), so it must warn too.
			m.toolCallingWarning = fmt.Sprintf("checking whether %s supports real tool calling failed: %v — the tutor may not work correctly with this model. Pick a different model from the Model menu if this causes problems.", msg.model, msg.err)
		case !msg.supported:
			m.toolCallingWarning = fmt.Sprintf("%s may not support real tool calling — the tutor may not be able to read your code, problem, or tests reliably. Pick a different model from the Model menu if this causes problems.", msg.model)
		}
		return m, nil

	case tea.KeyMsg:
		switch m.stage {
		case stageMain:
			return m.updateMain(msg)
		case stageCategories:
			return m.updateCategories(msg)
		case stageDSACategories:
			return m.updateDSACategories(msg)
		case stageProblems:
			return m.updateProblems(msg)
		case stageLanguage:
			return m.updateLanguage(msg)
		case stageResumeDraft:
			return m.updateResumeDraft(msg)
		case stageStats:
			return m.updateStats(msg)
		case stageSettings:
			return m.updateSettings(msg)
		case stageProviderChoice:
			return m.updateProviderChoice(msg)
		case stageModelPicker:
			return m.updateModelPicker(msg)
		case stageAPIModelEntry:
			return m.updateAPIModelEntry(msg)
		case stageOpenRouterKeyEntry:
			return m.updateOpenRouterKeyEntry(msg)
		}
	}
	return m, nil
}

// --- stageMain ---

func (m appModel) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(menuLabels)-1 {
			m.cursor++
		}
	case "1", "2", "3", "4", "5":
		n, _ := strconv.Atoi(msg.String())
		m.cursor = n - 1
	case "enter":
		return m.chooseMain()
	case "q", "ctrl+c":
		m.quit = true
		return m, tea.Quit
	}
	return m, nil
}

func (m appModel) chooseMain() (tea.Model, tea.Cmd) {
	switch menuChoice(m.cursor) {
	case menuPractice:
		return m.loadPractice(), nil
	case menuDaily:
		return m.loadDaily().resolveLanguageStage()
	case menuSandbox:
		m.outcome = outcomeRunSandbox
		return m, tea.Quit
	case menuStats:
		return m.loadStats(), nil
	case menuSettings:
		return m.loadSettings()
	}
	return m, nil
}

// loadPractice loads every problem (synchronous local disk + sqlite
// reads, same as the pre-merge runPracticeLoop did before every category
// picker launch) and derives the category list from it, in
// first-encountered order — catalog.List already sorts by categoryOrder,
// so that ordering carries through here too.
func (m appModel) loadPractice() appModel {
	statuses, err := catalogListFn(m.cfg)
	if err != nil {
		m.err = err
		return m
	}
	m.err = nil
	m.problems = catalog.GroupByProblem(statuses)
	m.categories = distinctCategories(m.problems)
	m.categoryCursor = 0
	m.stage = stageCategories
	return m
}

// loadDaily jumps straight to the language/style picker for today's
// DailyPick -- the whole point of the entry is skipping the
// category-and-problem navigation. The whole back-navigation context
// is still set up (category problems with the cursor on the pick, the
// category list, and the DSA topic list when the pick is in a grouped
// category) so backing out with q walks the normal picker screens
// instead of empty ones -- a real gap found in review: the first
// version populated only categoryProblems, and Daily -> q -> q
// rendered an empty topic list and an empty category list.
func (m appModel) loadDaily() appModel {
	statuses, err := catalogListFn(m.cfg)
	if err != nil {
		m.err = err
		return m
	}
	m.err = nil
	m.problems = catalog.GroupByProblem(statuses)
	pick, ok := catalog.DailyPick(m.problems, time.Now())
	if !ok {
		m.err = fmt.Errorf("tui: no problems available for a daily pick")
		return m
	}
	m.selectedProblem = pick
	m.category = pick.Category
	m.categories = distinctCategories(m.problems)
	m.categoryCursor = 0
	if catalog.IsGroupedCategory(pick.Category) {
		m.dsaCategories = distinctDSASubcategories(m.problems)
		m.dsaCategoryCursor = 0
	}
	m.categoryProblems = filterByCategory(m.problems, pick.Category)
	m.problemCursor = 0
	m.problemFilter = ""
	for i, p := range m.categoryProblems {
		if p.ProblemID == pick.ProblemID {
			m.problemCursor = i
			break
		}
	}
	m.langCursor = 0
	m.stage = stageLanguage
	return m
}

func (m appModel) loadStats() appModel {
	statuses, err := catalogListFn(m.cfg)
	if err != nil {
		m.err = err
		return m
	}
	recent, err := recentAttemptsFn(m.cfg, recentAttemptsLimit)
	if err != nil {
		m.err = err
		return m
	}
	all, err := allAttemptsFn(m.cfg)
	if err != nil {
		m.err = err
		return m
	}
	m.statsWeakDims = catalog.WeakDimensions(all)
	m.statsCodingWeak = catalog.CodingWeakSpots(all, codingWeakSpotMinAttempts)
	m.err = nil
	m.statsStatuses = statuses
	m.statsRecent = recent
	m.stage = stageStats
	return m
}

// loadSettings enters stageSettings, the Worker/Orchestrator role
// choice — unlike loadModelPicker/loadStats, this needs no data fetch,
// so it's synchronous (no tea.Cmd).
func (m appModel) loadSettings() (tea.Model, tea.Cmd) {
	m.stage = stageSettings
	m.settingsCursor = 0
	return m, nil
}

func (m appModel) loadModelPicker() (tea.Model, tea.Cmd) {
	m.stage = stageModelPicker
	m.modelLoading = true
	m.modelLoadErr = nil
	m.modelFilter = ""
	m.modelFiltered = nil
	m.modelCursor = 0
	m.modelWarning = ""
	m.modelPendingDownloadTag = ""
	m.modelDownloading = false
	m.modelDownloadLines = nil
	return m, func() tea.Msg {
		models, err := listModelsFn(ollamaHost)
		return modelsLoadedMsg{models: models, err: err}
	}
}

// loadAPIModelEntry enters stageAPIModelEntry seeded with
// suggestedOpenRouterModels — unlike loadModelPicker's Ollama tags,
// this list is static, so no async fetch (and no loading state) is
// needed.
func (m appModel) loadAPIModelEntry() appModel {
	m.stage = stageAPIModelEntry
	m.apiModelInput = ""
	m.apiModelFiltered = suggestedOpenRouterModels
	m.apiModelCursor = 0
	return m
}

// distinctCategories collects the top-level practice-picker entries —
// every NeetCode roadmap subcategory collapses into a single "dsa" entry
// via catalog.TopLevelGroup, so DSA shows once no matter how many
// subcategories have problems in them. Sorted by catalog.CategoryRank
// (rather than left at first-encountered order) so DSA always sorts to
// its own taxonomy position even though no exercise carries the literal
// "dsa" category anymore — every DSA problem lives under a specific
// subcategory like "two-pointers", which alone would otherwise anchor
// the group whenever its first subcategory in the list happens to be.
func distinctCategories(problems []catalog.ProblemStatus) []string {
	var categories []string
	seen := make(map[string]bool)
	for _, p := range problems {
		group := catalog.TopLevelGroup(p.Category)
		if !seen[group] {
			seen[group] = true
			categories = append(categories, group)
		}
	}
	sort.Slice(categories, func(i, j int) bool {
		return catalog.CategoryRank(categories[i]) < catalog.CategoryRank(categories[j])
	})
	return categories
}

// distinctDSASubcategories collects the second-level picker entries shown
// after selecting the top-level DSA group — the real NeetCode roadmap
// categories (Arrays & Hashing, Two Pointers, ...), sorted by
// catalog.CategoryRank to match the roadmap's own sequence.
func distinctDSASubcategories(problems []catalog.ProblemStatus) []string {
	var categories []string
	seen := make(map[string]bool)
	for _, p := range problems {
		if !catalog.IsGroupedCategory(p.Category) {
			continue
		}
		if !seen[p.Category] {
			seen[p.Category] = true
			categories = append(categories, p.Category)
		}
	}
	sort.Slice(categories, func(i, j int) bool {
		return catalog.CategoryRank(categories[i]) < catalog.CategoryRank(categories[j])
	})
	return categories
}

// groupCounts aggregates solved/total across every problem whose
// top-level group matches group — used for the top-level picker's rows,
// where DSA's count needs to sum across all its subcategories while an
// ungrouped category (debug, ...) is its own group of one.
func groupCounts(problems []catalog.ProblemStatus, group string) (solved, total int) {
	for _, p := range problems {
		if catalog.TopLevelGroup(p.Category) == group {
			total++
			if p.Solved {
				solved++
			}
		}
	}
	return solved, total
}

// categoryCounts aggregates solved/total for an exact category match —
// used by the second-level DSA subcategory picker, where each row is a
// real leaf category rather than a top-level group.
func categoryCounts(problems []catalog.ProblemStatus, category string) (solved, total int) {
	for _, p := range problems {
		if p.Category == category {
			total++
			if p.Solved {
				solved++
			}
		}
	}
	return solved, total
}

// filterByCategory builds the picker's problem list for one category:
// the category's problems in List's order, except problems that are
// due (catalog.Due: the mock second pass, or date-based review) float
// to the top -- a due marker is useless on a problem sitting below the
// fold of an alphabetical list. Every categoryProblems assignment goes
// through here, so the ordering can't diverge between the category,
// DSA-subcategory, and resume paths.
func filterByCategory(problems []catalog.ProblemStatus, category string) []catalog.ProblemStatus {
	var out []catalog.ProblemStatus
	for _, p := range problems {
		if p.Category == category {
			out = append(out, p)
		}
	}
	now := time.Now()
	return catalog.SortDueFirst(out, func(p catalog.ProblemStatus) bool { return catalog.Due(p, now) })
}

// --- stageCategories ---

func (m appModel) updateCategories(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.categoryCursor > 0 {
			m.categoryCursor--
		}
	case "down", "j":
		if m.categoryCursor < len(m.categories)-1 {
			m.categoryCursor++
		}
	case "enter":
		if len(m.categories) == 0 {
			return m, nil
		}
		selected := m.categories[m.categoryCursor]
		if selected == exercise.CategoryDSA {
			m.dsaCategories = distinctDSASubcategories(m.problems)
			m.dsaCategoryCursor = 0
			m.stage = stageDSACategories
		} else {
			m.category = selected
			m.categoryProblems = filterByCategory(m.problems, selected)
			m.problemCursor = 0
			m.problemFilter = ""
			m.stage = stageProblems
		}
	case "q", "esc", "ctrl+c":
		m.stage = stageMain
	}
	return m, nil
}

// --- stageDSACategories ---

func (m appModel) updateDSACategories(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.dsaCategoryCursor > 0 {
			m.dsaCategoryCursor--
		}
	case "down", "j":
		if m.dsaCategoryCursor < len(m.dsaCategories)-1 {
			m.dsaCategoryCursor++
		}
	case "enter":
		if len(m.dsaCategories) == 0 {
			return m, nil
		}
		m.category = m.dsaCategories[m.dsaCategoryCursor]
		m.categoryProblems = filterByCategory(m.problems, m.category)
		m.problemCursor = 0
		m.problemFilter = ""
		m.stage = stageProblems
	case "q", "esc", "ctrl+c":
		m.stage = stageCategories
	}
	return m, nil
}

// --- stageProblems ---

// filterProblems returns the problems whose title contains filter
// (case-insensitive), preserving order — the picker-list sibling of
// filterModels, so due-first ordering survives filtering.
func filterProblems(problems []catalog.ProblemStatus, filter string) []catalog.ProblemStatus {
	if filter == "" {
		return problems
	}
	lower := strings.ToLower(filter)
	var out []catalog.ProblemStatus
	for _, p := range problems {
		if strings.Contains(strings.ToLower(p.Title), lower) {
			out = append(out, p)
		}
	}
	return out
}

// visibleProblems is what the picker actually shows and navigates:
// categoryProblems through the live filter.
func (m appModel) visibleProblems() []catalog.ProblemStatus {
	return filterProblems(m.categoryProblems, m.problemFilter)
}

func (m appModel) leaveProblems() (tea.Model, tea.Cmd) {
	if catalog.IsGroupedCategory(m.category) {
		m.stage = stageDSACategories
	} else {
		m.stage = stageCategories
	}
	return m, nil
}

func (m appModel) updateProblems(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visible := m.visibleProblems()
	switch msg.Type {
	case tea.KeyUp:
		if m.problemCursor > 0 {
			m.problemCursor--
		}
	case tea.KeyDown:
		if m.problemCursor < len(visible)-1 {
			m.problemCursor++
		}
	case tea.KeyBackspace:
		if len(m.problemFilter) > 0 {
			m.problemFilter = m.problemFilter[:len(m.problemFilter)-1]
			m.problemCursor = 0
		}
	case tea.KeyEnter:
		if len(visible) == 0 || m.problemCursor >= len(visible) {
			return m, nil
		}
		m.selectedProblem = visible[m.problemCursor]
		m.langCursor = 0
		m.stage = stageLanguage
		return m.resolveLanguageStage()
	case tea.KeyEsc, tea.KeyCtrlC:
		return m.leaveProblems()
	case tea.KeyRunes:
		// "q" with nothing typed yet backs out, matching every other
		// stage — once the user has started typing, every rune
		// (including "q") feeds the filter, since it might be part of
		// a real title. Same contract as the model picker.
		if m.problemFilter == "" && string(msg.Runes) == "q" {
			return m.leaveProblems()
		}
		m.problemFilter += string(msg.Runes)
		m.problemCursor = 0
	}
	return m, nil
}

// --- stageLanguage ---

// resolveLanguageStage honors cfg.DefaultLanguage wherever the flow
// just landed on the language picker: a variant matching the default
// runs immediately, anything else keeps the picker. Design/behavioral
// problems ride coach/interviewer in the language slot, so they never
// match a python/go/cpp default and always ask — same for a problem
// that simply lacks the default language. Guarded on stage so error
// paths (e.g. loadDaily with no problems) pass through untouched.
func (m appModel) resolveLanguageStage() (tea.Model, tea.Cmd) {
	if m.stage != stageLanguage || m.cfg.DefaultLanguage == "" {
		return m, nil
	}
	for _, v := range m.selectedProblem.Variants {
		if v.Exercise.Language == m.cfg.DefaultLanguage {
			return m.launchExercise(v.Exercise)
		}
	}
	return m, nil
}

func (m appModel) updateLanguage(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.langCursor > 0 {
			m.langCursor--
		}
	case "down", "j":
		if m.langCursor < len(m.selectedProblem.Variants)-1 {
			m.langCursor++
		}
	case "enter":
		if len(m.selectedProblem.Variants) == 0 {
			return m, nil
		}
		return m.launchExercise(m.selectedProblem.Variants[m.langCursor].Exercise)
	case "q", "esc", "ctrl+c":
		m.stage = stageProblems
	}
	return m, nil
}

// --- stageStats ---

// updateStats mirrors the pre-merge statsModel: any keypress goes back.
func (m appModel) updateStats(tea.KeyMsg) (tea.Model, tea.Cmd) {
	m.stage = stageMain
	return m, nil
}

// --- stageSettings ---

// updateSettings handles the top-level Worker model / Orchestrator model
// role choice — a plain 2-item list, same up/down/enter shape as every
// other short list in this program. Enter records which Config field is
// being edited (settingsEditing) and moves to stageProviderChoice, which
// asks Local vs API for that role.
func (m appModel) updateSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.settingsCursor > 0 {
			m.settingsCursor--
		}
		return m, nil
	case tea.KeyDown:
		if m.settingsCursor < len(settingsRoleLabels)-1 {
			m.settingsCursor++
		}
		return m, nil
	case tea.KeyEsc, tea.KeyCtrlC:
		m.stage = stageMain
		return m, nil
	case tea.KeyEnter:
		switch m.settingsCursor {
		case 0:
			m.settingsEditing = settingsTargetWorker
		case 1:
			m.settingsEditing = settingsTargetOrchestrator
		case 2:
			return m.cycleDefaultLanguage()
		case 3:
			return m.toggleTutorNotes()
		}
		m.stage = stageProviderChoice
		m.settingsCursor = 0
		return m, nil
	}
	return m, nil
}

// cycleDefaultLanguage advances the default-language preference
// ask → python → go → cpp → ask and persists immediately — a 4-value
// cycle doesn't warrant a sub-stage, and showing the new value in the
// header line is feedback enough.
func (m appModel) cycleDefaultLanguage() (tea.Model, tea.Cmd) {
	next := map[string]string{"": "python", "python": "go", "go": "cpp", "cpp": ""}
	m.cfg.DefaultLanguage = next[m.cfg.DefaultLanguage]
	if err := config.SaveSettings(m.cfg.SettingsPath(), settingsFromConfig(m.cfg)); err != nil {
		m.err = err
		return m, nil
	}
	m.err = nil
	return m, nil
}

// toggleTutorNotes flips whether sessions get the tutor's editor
// highlight/note tool at all (issue #25's durable variant; M-h in a
// session still toggles rendering of notes that do exist).
func (m appModel) toggleTutorNotes() (tea.Model, tea.Cmd) {
	m.cfg.DisableTutorNotes = !m.cfg.DisableTutorNotes
	if err := config.SaveSettings(m.cfg.SettingsPath(), settingsFromConfig(m.cfg)); err != nil {
		m.err = err
		return m, nil
	}
	m.err = nil
	return m, nil
}

// --- stageProviderChoice ---

// updateProviderChoice handles the Local (Ollama) / API (OpenRouter)
// provider choice for whichever role updateSettings just picked
// (m.settingsEditing) — a plain list, 2 items normally, or 3 when editing
// the orchestrator (a trailing "None (disable routing)" entry, since the
// worker model can never be unset).
func (m appModel) updateProviderChoice(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	maxCursor := 1
	if m.settingsEditing == settingsTargetOrchestrator {
		maxCursor = 2
	}
	switch msg.Type {
	case tea.KeyUp:
		if m.settingsCursor > 0 {
			m.settingsCursor--
		}
		return m, nil
	case tea.KeyDown:
		if m.settingsCursor < maxCursor {
			m.settingsCursor++
		}
		return m, nil
	case tea.KeyEsc, tea.KeyCtrlC:
		m.stage = stageSettings
		return m, nil
	case tea.KeyEnter:
		switch m.settingsCursor {
		case 0:
			return m.loadModelPicker()
		case 1:
			return m.loadAPIModelEntry(), nil
		case 2:
			return m.clearOrchestratorModel()
		}
	}
	return m, nil
}

// clearOrchestratorModel disables routing by clearing cfg.OrchestratorModel
// — the "None (disable routing)" entry in stageProviderChoice, only
// reachable when editing the orchestrator role.
func (m appModel) clearOrchestratorModel() (tea.Model, tea.Cmd) {
	m.cfg.OrchestratorModel = ""
	if err := config.SaveSettings(m.cfg.SettingsPath(), settingsFromConfig(m.cfg)); err != nil {
		m.err = err
		return m, nil
	}
	m.err = nil
	m.stage = stageSettings
	return m, nil
}

// --- stageAPIModelEntry ---

// updateAPIModelEntry handles a single-line unmasked text input for a
// bare OpenRouter model slug — same rune/backspace shape as
// updateOpenRouterKeyEntry's key input, unmasked since a model slug
// isn't a secret. Enter delegates to selectModelOrPromptForKey (already
// built for handleModelEnter's directly-typed openrouter: path), which
// handles both "key already available" and "no key yet" correctly
// without any new logic needed here.
func (m appModel) updateAPIModelEntry(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyUp:
		if m.apiModelCursor > 0 {
			m.apiModelCursor--
		}
		return m, nil
	case tea.KeyDown:
		if m.apiModelCursor < len(m.apiModelFiltered)-1 {
			m.apiModelCursor++
		}
		return m, nil
	case tea.KeyEsc, tea.KeyCtrlC:
		m.stage = stageProviderChoice
		m.apiModelInput = ""
		m.apiModelFiltered = nil
		return m, nil
	case tea.KeyBackspace:
		if len(m.apiModelInput) > 0 {
			m.apiModelInput = m.apiModelInput[:len(m.apiModelInput)-1]
			m.apiModelFiltered = filterModels(suggestedOpenRouterModels, m.apiModelInput)
			m.apiModelCursor = 0
		}
		return m, nil
	case tea.KeyEnter:
		if len(m.apiModelFiltered) > 0 {
			sel := m.apiModelFiltered[m.apiModelCursor]
			m.apiModelInput = ""
			m.apiModelFiltered = nil
			return m.selectModelOrPromptForKey(tutor.OpenRouterModelPrefix + sel)
		}
		slug := strings.TrimSpace(m.apiModelInput)
		if slug == "" {
			return m, nil
		}
		m.apiModelInput = ""
		return m.selectModelOrPromptForKey(tutor.OpenRouterModelPrefix + slug)
	case tea.KeyRunes:
		// "q" with nothing typed yet backs out, matching stageModelPicker's
		// identical convention (see handleModelEnter's own comment on the
		// same trade-off for tags that start with "q").
		if m.apiModelInput == "" && string(msg.Runes) == "q" {
			m.stage = stageProviderChoice
			m.apiModelFiltered = nil
			return m, nil
		}
		m.apiModelInput += string(msg.Runes)
		m.apiModelFiltered = filterModels(suggestedOpenRouterModels, m.apiModelInput)
		m.apiModelCursor = 0
		return m, nil
	}
	return m, nil
}

// --- stageModelPicker ---

func (m appModel) updateModelPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.modelDownloading {
		// Nothing to do with input mid-download — it isn't
		// interruptible, matching boot.go's live build/pull panels.
		return m, nil
	}
	if m.modelPendingDownloadTag != "" {
		return m.handleModelDownloadPrompt(msg)
	}
	switch msg.Type {
	case tea.KeyUp:
		if m.modelCursor > 0 {
			m.modelCursor--
		}
		return m, nil
	case tea.KeyDown:
		if m.modelCursor < len(m.modelFiltered)-1 {
			m.modelCursor++
		}
		return m, nil
	case tea.KeyBackspace:
		if len(m.modelFilter) > 0 {
			m.modelFilter = m.modelFilter[:len(m.modelFilter)-1]
			m.modelFiltered = filterModels(m.models, m.modelFilter)
			m.modelCursor = 0
			m.modelWarning = ""
		}
		return m, nil
	case tea.KeyEsc, tea.KeyCtrlC:
		// Back to the provider choice, not stageMain — stageModelPicker
		// is only reachable via stageProviderChoice -> Local now.
		m.stage = stageProviderChoice
		return m, nil
	case tea.KeyEnter:
		return m.handleModelEnter()
	case tea.KeyRunes:
		// "q" with nothing typed yet backs out, matching every other
		// stage in this program — once the user has started typing,
		// every rune (including "q") feeds the filter/custom tag
		// instead, since it might be part of a real model name.
		if m.modelFilter == "" && string(msg.Runes) == "q" {
			m.stage = stageProviderChoice
			return m, nil
		}
		m.modelFilter += string(msg.Runes)
		m.modelFiltered = filterModels(m.models, m.modelFilter)
		m.modelCursor = 0
		m.modelWarning = ""
		return m, nil
	}
	return m, nil
}

// isLocalModel reports whether name is actually pulled locally (as
// opposed to merely suggested) — selecting a local entry can skip
// straight to selection; anything else needs the not-pulled check first.
func (m appModel) isLocalModel(name string) bool {
	for _, local := range m.localModels {
		if local == name {
			return true
		}
	}
	return false
}

// handleModelEnter selects the highlighted entry if it's already pulled
// locally, or — for a highlighted-but-unpulled suggested entry, or when
// the typed filter matches nothing — treats the tag as a candidate:
// checkModelTag checks it against Ollama and, if it isn't pulled, asks
// whether to download it (see handleModelDownloadPrompt).
//
// An OpenRouterModelPrefix-prefixed tag is intercepted before either
// path: there's no "pulled locally" concept for a hosted model, and
// checkModelTag's checkModelFn call is an Ollama /api/tags lookup that
// would just misreport it as not-pulled.
func (m appModel) handleModelEnter() (tea.Model, tea.Cmd) {
	if len(m.modelFiltered) > 0 {
		sel := m.modelFiltered[m.modelCursor]
		if strings.HasPrefix(sel, tutor.OpenRouterModelPrefix) {
			return m.selectModelOrPromptForKey(sel)
		}
		if m.isLocalModel(sel) {
			return m.selectModel(sel)
		}
		return m.checkModelTag(sel)
	}

	tag := strings.TrimSpace(m.modelFilter)
	if tag == "" {
		return m, nil
	}
	if strings.HasPrefix(tag, tutor.OpenRouterModelPrefix) {
		return m.selectModelOrPromptForKey(tag)
	}
	return m.checkModelTag(tag)
}

// selectModelOrPromptForKey selects name immediately if an OpenRouter
// API key is already available (settings.json or OPENROUTER_API_KEY,
// resolved into cfg by config.Load), otherwise asks for one first via
// stageOpenRouterKeyEntry — entering it there persists it to
// settings.json (see updateOpenRouterKeyEntry), so this only ever asks
// once across sessions, not once per pick.
func (m appModel) selectModelOrPromptForKey(name string) (tea.Model, tea.Cmd) {
	if m.cfg.OpenRouterAPIKey != "" {
		return m.selectModel(name)
	}
	m.stage = stageOpenRouterKeyEntry
	m.openRouterPendingModel = name
	m.openRouterKeyInput = ""
	m.modelWarning = ""
	return m, nil
}

// checkModelTag checks tag against Ollama: if it's actually already
// pulled (e.g. a tag Ollama resolves some other way, or a race with a
// pull that just finished outside this picker), select it immediately;
// otherwise show why, and ask whether to download it.
func (m appModel) checkModelTag(tag string) (tea.Model, tea.Cmd) {
	check := checkModelFn(ollamaHost, tag)
	if check.OK {
		return m.selectModel(tag)
	}
	m.modelWarning = check.Detail
	m.modelPendingDownloadTag = tag
	return m, nil
}

// handleModelDownloadPrompt handles the y/n answer to "download <tag>?"
// — y starts a live preflight.PullModel stream (see the pullLineMsg/
// pullDoneMsg cases in Update), n cancels back to the picker with
// nothing selected, and esc/ctrl+c back out to the main menu entirely,
// matching stageModelPicker's normal esc/ctrl+c behavior.
func (m appModel) handleModelDownloadPrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc, tea.KeyCtrlC:
		m.stage = stageMain
		return m, nil
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "y", "Y":
			tag := m.modelPendingDownloadTag
			m.modelPendingDownloadTag = ""
			m.modelDownloading = true
			m.modelDownloadTarget = tag
			lineCh, errCh := pullModelFn(ollamaHost, tag)
			m.modelDownloadLineCh = lineCh
			m.modelDownloadErrCh = errCh
			return m, waitForPullLine(lineCh, errCh)
		case "n", "N":
			m.modelPendingDownloadTag = ""
			m.modelWarning = ""
			return m, nil
		}
	}
	return m, nil
}

// selectModel persists the pick immediately (same call the pre-merge
// runModelPicker made in run.go) and updates cfg in place so any
// exercise/sandbox launched later in this same process uses it right
// away, without waiting for a fresh Config.Load. Selection itself is
// never blocked on whether the model actually supports real tool
// calling — checkToolCallingCmd runs that check in the background and
// only surfaces a warning (toolCallingCheckMsg in Update) if it fails,
// so a slow/unreachable check can't stall the picker.
//
// Which Config field gets written depends on m.settingsEditing (set by
// updateSettings when the role choice was made) — either branch must
// carry forward every other field, or the one this call isn't touching
// would silently get wiped from settings.json on save.
func (m appModel) selectModel(name string) (tea.Model, tea.Cmd) {
	if m.settingsEditing == settingsTargetOrchestrator {
		m.cfg.OrchestratorModel = name
	} else {
		m.cfg.TutorModel = name
	}
	if err := config.SaveSettings(m.cfg.SettingsPath(), settingsFromConfig(m.cfg)); err != nil {
		m.err = err
		return m, nil
	}
	m.err = nil
	m.stage = stageMain
	m.toolCallingWarning = ""
	return m, checkToolCallingCmd(name, m.cfg.OpenRouterAPIKey)
}

// --- stageOpenRouterKeyEntry ---

// updateOpenRouterKeyEntry handles a single-line masked text input for
// the OpenRouter API key — same rune/backspace shape as modelFilter's
// handling in updateModelPicker, kept separate rather than shared since
// this one never filters a list and needs its own Enter/Esc behavior.
func (m appModel) updateOpenRouterKeyEntry(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc, tea.KeyCtrlC:
		// Back to the provider choice, not stageModelPicker — a single
		// consistent "cancel returns to the provider choice" behavior
		// regardless of whether this stage was reached via typing
		// openrouter: directly in the local picker or via the Settings ->
		// API path.
		m.stage = stageProviderChoice
		m.openRouterPendingModel = ""
		m.openRouterKeyInput = ""
		return m, nil
	case tea.KeyBackspace:
		if len(m.openRouterKeyInput) > 0 {
			m.openRouterKeyInput = m.openRouterKeyInput[:len(m.openRouterKeyInput)-1]
		}
		return m, nil
	case tea.KeyEnter:
		key := strings.TrimSpace(m.openRouterKeyInput)
		if key == "" {
			return m, nil
		}
		m.cfg.OpenRouterAPIKey = key
		pending := m.openRouterPendingModel
		m.openRouterPendingModel = ""
		m.openRouterKeyInput = ""
		return m.selectModel(pending)
	case tea.KeyRunes:
		m.openRouterKeyInput += string(msg.Runes)
		return m, nil
	}
	return m, nil
}

// --- View ---

func (m appModel) View() string {
	right := m.renderRight()
	if m.width == 0 || m.height == 0 {
		return right
	}
	// Only the main menu gets the panel footer -- its key hints used to
	// float inside the body. Every other stage still carries its own
	// in-body hint line.
	footer := ""
	if m.stage == stageMain {
		footer = "↑/↓ or j/k move · 1-5 jump · enter select · q quit"
	}
	panel := renderDashboardPanel(m.width, m.height, m.phase, right, layoutCentered, footer)
	return placeBlock(m.width, m.height, panel)
}

func (m appModel) renderRight() string {
	switch m.stage {
	case stageCategories:
		return m.renderCategories()
	case stageDSACategories:
		return m.renderDSACategories()
	case stageProblems:
		return m.renderProblems()
	case stageLanguage:
		return m.renderLanguage()
	case stageResumeDraft:
		return m.renderResumeDraft()
	case stageStats:
		return m.renderStats()
	case stageSettings:
		return m.renderSettings()
	case stageProviderChoice:
		return m.renderProviderChoice()
	case stageModelPicker:
		return m.renderModelPicker()
	case stageAPIModelEntry:
		return m.renderAPIModelEntry()
	case stageOpenRouterKeyEntry:
		return m.renderOpenRouterKeyEntry()
	default:
		return m.renderMain()
	}
}

func (m appModel) renderMain() string {
	var b strings.Builder
	for i, label := range menuLabels {
		numLabel := fmt.Sprintf("%d. %s", i+1, label)
		if i == m.cursor {
			row := fmt.Sprintf("❯ %-*s", menuRightColWidth-2, numLabel)
			b.WriteString(menuRowHighlight.Render(row))
			b.WriteString("\n  " + menuSubtitleStyle.Render(menuDescriptions[i]))
		} else {
			b.WriteString("  " + numLabel)
		}
		b.WriteString("\n\n\n")
	}

	if m.err != nil {
		b.WriteString(renderFriendlyError("something went wrong", m.err))
		b.WriteString("\n\n")
	}

	if m.toolCallingWarning != "" {
		b.WriteString(hintStyle.Render("  " + m.toolCallingWarning))
		b.WriteString("\n\n")
	}

	// The dashboard block; the key hints that used to float here live
	// in the panel footer now (see View's footer for stageMain).
	if board := renderHomeboard(m.problems, m.homeAttempts, time.Now()); board != "" {
		b.WriteString("\n")
		b.WriteString(board)
	}
	if m.homeLoadErr != nil {
		b.WriteString("\n")
		b.WriteString(renderFriendlyError("couldn't load your progress", m.homeLoadErr))
	}
	return b.String()
}

// renderFriendlyError is the one shape every screen-level failure
// renders as: a short human headline in the fail color, with the raw
// error kept visible on a dim detail line underneath — readable at a
// glance without swallowing what actually happened.
func renderFriendlyError(headline string, err error) string {
	return failStyle.Render("  "+headline) + "\n" + checkDimStyle.Render("    "+err.Error())
}

func (m appModel) renderCategories() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("Practice"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("choose a category"))
	b.WriteString("\n\n")

	for i, cat := range m.categories {
		solved, total := groupCounts(m.problems, cat)
		label := fmt.Sprintf("%-16s", catalog.DisplayCategory(cat))
		status := fmt.Sprintf("%d/%d solved", solved, total)
		if i == m.categoryCursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s %s", label, status)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s", categoryStyle.Render(label), checkDimStyle.Render(status)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · q back"))
	return b.String()
}

func (m appModel) renderDSACategories() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("DSA"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("choose a topic"))
	b.WriteString("\n\n")

	for i, cat := range m.dsaCategories {
		solved, total := categoryCounts(m.problems, cat)
		label := fmt.Sprintf("%-26s", catalog.DisplayCategory(cat))
		status := fmt.Sprintf("%d/%d solved", solved, total)
		if i == m.dsaCategoryCursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s %s", label, status)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s", categoryStyle.Render(label), checkDimStyle.Render(status)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · q back"))
	return b.String()
}

func (m appModel) renderProblems() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render(catalog.DisplayCategory(m.category)))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("choose a problem"))
	b.WriteString("\n\n")

	// Set by RunApp when the session just picked here failed to launch
	// (e.g. Docker went down) -- shown right where the pick happened,
	// on the same screen it resumes to, instead of a stderr line the
	// very next alt-screen program immediately wiped (issue #230).
	if m.err != nil {
		b.WriteString(renderFriendlyError("couldn't start that session", m.err))
		b.WriteString("\n\n")
	}

	b.WriteString(fmt.Sprintf("%s%s", checkDimStyle.Render("› "), m.problemFilter))
	b.WriteString("\n\n")

	visible := m.visibleProblems()
	if len(visible) == 0 {
		b.WriteString(checkDimStyle.Render("no matches"))
		b.WriteString("\n")
	}
	cursor := min(m.problemCursor, max(len(visible)-1, 0))
	start, end := problemWindow(cursor, len(visible))
	if start > 0 {
		b.WriteString(checkDimStyle.Render(fmt.Sprintf("  ↑ %d more", start)))
		b.WriteString("\n")
	}
	for i := start; i < end; i++ {
		p := visible[i]
		label := fmt.Sprintf("%-30s", truncateTitle(p.Title, 30))
		badge := fmt.Sprintf("%-5s", renderDifficultyBadge(p))
		status := "not attempted"
		statusStyle := checkDimStyle
		if p.Attempts > 0 {
			plural := "s"
			if p.Attempts == 1 {
				plural = ""
			}
			status = fmt.Sprintf("%d attempt%s", p.Attempts, plural)
			statusStyle = failStyle
			if p.Solved {
				status = "solved (" + status + ")"
				statusStyle = passStyle
			}
		}
		marker := ""
		if catalog.MockDue(p) {
			// The roadmap's second pass: coach passed, interviewer mock
			// still untouched -- nudge without nagging.
			marker = "  · mock due"
		} else if catalog.ReviewDue(p, time.Now()) {
			// Date-based resurfacing: a stale failure or a solved
			// problem gone untouched for a month.
			marker = "  · review due"
		}
		if i == cursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s %s %s%s", label, stripBadgeStyle(p), status, marker)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s %s%s", label, badge, statusStyle.Render(status), dueMarkerStyle.Render(marker)))
		}
		b.WriteString("\n")
	}
	if end < len(visible) {
		b.WriteString(checkDimStyle.Render(fmt.Sprintf("  ↓ %d more", len(visible)-end)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("type to search · ↑/↓ move · enter select · esc back"))
	return b.String()
}

// maxVisibleProblems bounds the picker window so long categories
// (trees has 15 problems) never overflow a short terminal; the window
// slides to keep the cursor visible with ↑/↓ "N more" indicators.
const maxVisibleProblems = 12

// problemWindow returns the [start, end) slice of an n-problem list to
// draw with the cursor kept in view — centered when there's room,
// clamped at the edges. Pure so the windowing math is testable without
// rendering.
func problemWindow(cursor, n int) (int, int) {
	if n <= maxVisibleProblems {
		return 0, n
	}
	start := cursor - maxVisibleProblems/2
	start = max(0, min(start, n-maxVisibleProblems))
	return start, start + maxVisibleProblems
}

// difficultyBadge is the letter(s) for a problem's difficulty: one of
// E/M/H when its variants agree (coding variants always do — the fill
// rates the problem, not the language), or a range like "M/H" when
// they differ (system design rates coach sessions medium and
// interviewer mocks hard). Empty for unrated.
func difficultyBadge(p catalog.ProblemStatus) string {
	letters := map[string]string{
		exercise.DifficultyEasy:   "E",
		exercise.DifficultyMedium: "M",
		exercise.DifficultyHard:   "H",
	}
	var parts []string
	for _, d := range []string{exercise.DifficultyEasy, exercise.DifficultyMedium, exercise.DifficultyHard} {
		for _, v := range p.Variants {
			if v.Exercise.Difficulty == d {
				parts = append(parts, letters[d])
				break
			}
		}
	}
	return strings.Join(parts, "/")
}

// difficultyBadgeStyles maps each letter to the app palette: teal for
// easy (the palette's "good" color — no green exists, same reasoning
// as passStyle), gold for medium, red for hard.
var difficultyBadgeStyles = map[string]lipgloss.Style{
	"E": passStyle,
	"M": dueMarkerStyle,
	"H": failStyle,
}

// renderDifficultyBadge wraps the badge letters in [ ] with their
// per-letter colors; empty difficulty renders as empty (no brackets).
func renderDifficultyBadge(p catalog.ProblemStatus) string {
	badge := difficultyBadge(p)
	if badge == "" {
		return ""
	}
	var styled []string
	for _, letter := range strings.Split(badge, "/") {
		styled = append(styled, difficultyBadgeStyles[letter].Render(letter))
	}
	return checkDimStyle.Render("[") + strings.Join(styled, checkDimStyle.Render("/")) + checkDimStyle.Render("]")
}

// stripBadgeStyle is the cursor row's badge: cursorRowStyle re-renders
// the whole row, so the badge goes in plain to keep the row's single
// highlight style from fighting per-letter colors.
func stripBadgeStyle(p catalog.ProblemStatus) string {
	badge := difficultyBadge(p)
	if badge == "" {
		return fmt.Sprintf("%-5s", "")
	}
	return fmt.Sprintf("%-5s", "["+badge+"]")
}

func (m appModel) renderLanguage() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render(m.selectedProblem.Title))
	b.WriteString("\n")
	// A design problem's variants are session styles (coach/interviewer
	// riding the language slot -- see exercise.LanguageCoach's doc
	// comment), not languages; say so. Kind is uniform across a
	// problem's variants, so checking the first is enough.
	subtitle := "choose a language"
	if len(m.selectedProblem.Variants) > 0 && m.selectedProblem.Variants[0].Exercise.Kind == exercise.KindDesign {
		subtitle = "choose a session style"
	}
	b.WriteString(checkDimStyle.Render(subtitle))
	b.WriteString("\n\n")

	for i, v := range m.selectedProblem.Variants {
		lang := fmt.Sprintf("%-8s", v.Exercise.Language)
		status := "not attempted"
		statusStyle := checkDimStyle
		if v.LastResult != "" {
			status = v.LastResult
			statusStyle = failStyle
			if v.LastResult == tracker.ResultPass {
				statusStyle = passStyle
			}
		}
		if i == m.langCursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s %s", lang, status)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s", langStyle.Render(lang), statusStyle.Render(status)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · q back"))
	return b.String()
}

// codingWeakSpotMinAttempts is the evidence bar for calling a coding
// category weak — under this many attempts, one bad day would brand a
// whole topic.
const codingWeakSpotMinAttempts = 3

func (m appModel) renderStats() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("Stats"))
	b.WriteString("\n\n")

	total, attempted, solved := 0, 0, 0
	for _, s := range m.statsStatuses {
		total++
		if s.Attempts > 0 {
			attempted++
		}
		if s.LastResult == tracker.ResultPass {
			solved++
		}
	}
	fmt.Fprintf(&b, "%s solved · %s attempted · %d total exercises\n\n",
		passStyle.Render(fmt.Sprintf("%d", solved)),
		checkDimStyle.Render(fmt.Sprintf("%d", attempted)),
		total)

	b.WriteString(catalog.FormatSummary(m.statsStatuses) + "\n\n")

	if len(m.statsWeakDims) > 0 {
		b.WriteString(hintStyle.Render("Rubric weak spots (design grading)"))
		b.WriteString("\n")
		shown := m.statsWeakDims
		if len(shown) > 5 {
			shown = shown[:5]
		}
		for _, d := range shown {
			line := fmt.Sprintf("  %-32s missing %d/%d · adequate %d/%d", d.Name, d.Missing, d.Total(), d.Adequate, d.Total())
			style := checkDimStyle
			if d.Missing > 0 {
				style = failStyle
			}
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if len(m.statsCodingWeak) > 0 {
		b.WriteString(hintStyle.Render("Coding weak spots"))
		b.WriteString("\n")
		shown := m.statsCodingWeak
		if len(shown) > 5 {
			shown = shown[:5]
		}
		for _, c := range shown {
			line := fmt.Sprintf("  %-32s failed %d/%d", catalog.DisplayCategory(c.Category), c.Fails, c.Attempts)
			// Half-or-worse is a real problem area; anything listed at
			// all is still worth a dim line (it has fails by
			// construction).
			style := checkDimStyle
			if c.FailRatio() >= 0.5 {
				style = failStyle
			}
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if len(m.statsRecent) == 0 {
		b.WriteString(checkDimStyle.Render("No attempts yet — go practice something!"))
		b.WriteString("\n")
	} else {
		b.WriteString(hintStyle.Render("Recent activity"))
		b.WriteString("\n")
		for _, a := range m.statsRecent {
			resultStyle := failStyle
			if a.Result == tracker.ResultPass {
				resultStyle = passStyle
			}
			fmt.Fprintf(&b, "%s  %-28s %s\n", a.Date, a.ExerciseID, resultStyle.Render(a.Result))
		}
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("press any key to go back"))
	return b.String()
}

// settingsOptionLabels/settingsOptionDescriptions back updateProviderChoice's
// base 2-item list and renderProviderChoice's rendering of it (a 3rd,
// orchestrator-only "None" entry is appended there, not here) — kept
// together so the list and its descriptions can't drift out of sync.
var (
	settingsOptionLabels = []string{"Local (Ollama)", "API (OpenRouter)"}
	settingsOptionDescs  = []string{
		"Pick from models pulled on this machine",
		"Use a hosted model via an OpenRouter API key",
	}
)

// settingsRoleLabels/settingsRoleDescs back updateSettings' top-level
// list and renderSettings' rendering of it: the two model roles (enter
// opens the provider choice) plus the two in-place preferences (enter
// cycles/toggles and saves immediately — no sub-stage for a 4-value
// cycle or a bool).
var (
	settingsRoleLabels = []string{"Worker model", "Orchestrator model", "Default language", "Tutor editor notes"}
	settingsRoleDescs  = []string{
		"Answers coding questions — always required",
		"Routes turns to the worker, or answers directly — optional",
		"Skip the language picker for problems available in it — enter cycles",
		"Let the tutor highlight lines and leave notes in the editor — enter toggles",
	}
)

// settingsFromConfig snapshots every persisted preference from cfg —
// the single source for "carry everything forward" saves, so adding a
// new Settings field can't silently wipe it from one save site.
func settingsFromConfig(cfg config.Config) config.Settings {
	return config.Settings{
		TutorModel:        cfg.TutorModel,
		OrchestratorModel: cfg.OrchestratorModel,
		GraderModel:       cfg.GraderModel,
		OpenRouterAPIKey:  cfg.OpenRouterAPIKey,
		DefaultLanguage:   cfg.DefaultLanguage,
		DisableTutorNotes: cfg.DisableTutorNotes,
	}
}

func (m appModel) renderSettings() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("Settings"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render(fmt.Sprintf("Worker model: %s", m.cfg.TutorModel)))
	b.WriteString("\n")
	orchestratorStatus := m.cfg.OrchestratorModel
	if orchestratorStatus == "" {
		orchestratorStatus = "none (routing off)"
	}
	b.WriteString(checkDimStyle.Render(fmt.Sprintf("Orchestrator model: %s", orchestratorStatus)))
	b.WriteString("\n")
	langStatus := m.cfg.DefaultLanguage
	if langStatus == "" {
		langStatus = "ask every time"
	}
	b.WriteString(checkDimStyle.Render(fmt.Sprintf("Default language: %s", langStatus)))
	b.WriteString("\n")
	notesStatus := "on"
	if m.cfg.DisableTutorNotes {
		notesStatus = "off"
	}
	b.WriteString(checkDimStyle.Render(fmt.Sprintf("Tutor editor notes: %s", notesStatus)))
	b.WriteString("\n\n")

	for i, label := range settingsRoleLabels {
		if i == m.settingsCursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s", label)))
			b.WriteString("\n  " + menuSubtitleStyle.Render(settingsRoleDescs[i]))
		} else {
			b.WriteString(fmt.Sprintf("  %s", label))
		}
		b.WriteString("\n\n")
	}

	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · esc back"))
	return b.String()
}

// renderProviderChoice renders the Local (Ollama) / API (OpenRouter)
// choice for whichever role updateSettings picked (m.settingsEditing) —
// a 3rd "None (disable routing)" entry is appended only when editing the
// orchestrator. Copies the shared label/desc slices before appending so
// the orchestrator-only 3rd item never leaks into the worker's list.
func (m appModel) renderProviderChoice() string {
	labels := append([]string{}, settingsOptionLabels...)
	descs := append([]string{}, settingsOptionDescs...)
	if m.settingsEditing == settingsTargetOrchestrator {
		labels = append(labels, "None (disable routing)")
		descs = append(descs, "Answer every turn with the worker model only")
	}

	var b strings.Builder
	title := "Worker model"
	if m.settingsEditing == settingsTargetOrchestrator {
		title = "Orchestrator model"
	}
	b.WriteString(hintStyle.Render(title))
	b.WriteString("\n\n")

	for i, label := range labels {
		if i == m.settingsCursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s", label)))
			b.WriteString("\n  " + menuSubtitleStyle.Render(descs[i]))
		} else {
			b.WriteString(fmt.Sprintf("  %s", label))
		}
		b.WriteString("\n\n")
	}

	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · esc back"))
	return b.String()
}

// renderAPIModelEntry's input is unmasked (unlike
// renderOpenRouterKeyEntry's) — a model slug isn't a secret.
func (m appModel) renderAPIModelEntry() string {
	currentModel := m.cfg.TutorModel
	if m.settingsEditing == settingsTargetOrchestrator {
		currentModel = m.cfg.OrchestratorModel
	}

	var b strings.Builder
	b.WriteString(hintStyle.Render("OpenRouter model"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("pick a suggestion, or type any model slug, e.g. anthropic/claude-3.5-sonnet"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("%s%s", checkDimStyle.Render("› "), m.apiModelInput))
	b.WriteString("\n\n")

	if len(m.apiModelFiltered) == 0 {
		b.WriteString(checkDimStyle.Render("no suggested matches -- enter confirms the typed slug above"))
		b.WriteString("\n")
	} else {
		for i, slug := range m.apiModelFiltered {
			label := slug
			if tutor.OpenRouterModelPrefix+slug == currentModel {
				label += "  " + hintStyle.Render("(current)")
			}
			if i == m.apiModelCursor {
				b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s", label)))
			} else {
				b.WriteString(fmt.Sprintf("  %s", label))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · esc back"))
	return b.String()
}

func (m appModel) renderModelPicker() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("Tutor model"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("type to search, or enter a custom tag"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("%s%s", checkDimStyle.Render("› "), m.modelFilter))
	b.WriteString("\n\n")

	switch {
	case m.modelLoading:
		b.WriteString(checkDimStyle.Render("loading models from " + ollamaHost + "..."))
		b.WriteString("\n")
	case m.modelLoadErr != nil:
		b.WriteString(renderFriendlyError("couldn't reach Ollama", m.modelLoadErr))
		b.WriteString("\n")
		b.WriteString(checkDimStyle.Render("you can still type a model tag directly"))
		b.WriteString("\n")
	case len(m.modelFiltered) == 0:
		b.WriteString(checkDimStyle.Render("no matches"))
		b.WriteString("\n")
	default:
		for i, name := range m.modelFiltered {
			label := name
			if name == m.cfg.TutorModel {
				label += "  " + hintStyle.Render("(current)")
			} else if !m.isLocalModel(name) {
				label += "  " + checkDimStyle.Render("(not pulled)")
			}
			if m.modelNoTools[name] {
				label += "  " + failStyle.Render("(no tools)")
			}
			if i == m.modelCursor {
				b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s", label)))
			} else {
				b.WriteString(fmt.Sprintf("  %s", label))
			}
			b.WriteString("\n")
		}
	}

	switch {
	case m.modelDownloading:
		b.WriteString("\n")
		b.WriteString(hintStyle.Render(fmt.Sprintf("downloading %s...", m.modelDownloadTarget)))
		b.WriteString("\n")
		for _, line := range m.modelDownloadLines {
			fmt.Fprintf(&b, "    %s\n", buildLogStyle.Render(line))
		}
	case m.modelPendingDownloadTag != "":
		b.WriteString("\n")
		b.WriteString(hintStyle.Render(m.modelWarning))
		b.WriteString("\n")
		b.WriteString(checkDimStyle.Render(fmt.Sprintf("download %s? (y/n)", m.modelPendingDownloadTag)))
		b.WriteString("\n")
	case m.modelWarning != "":
		b.WriteString("\n")
		b.WriteString(hintStyle.Render(m.modelWarning))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · esc back"))
	return b.String()
}

// renderOpenRouterKeyEntry shows asterisks in place of the typed key —
// this is a real secret on screen, unlike everything else the picker
// handles, so it's masked even though this is a local, single-user TUI.
func (m appModel) renderOpenRouterKeyEntry() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("OpenRouter API key"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render(fmt.Sprintf("needed once, to use %s", m.openRouterPendingModel)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("%s%s", checkDimStyle.Render("› "), strings.Repeat("*", len(m.openRouterKeyInput))))
	b.WriteString("\n\n")

	b.WriteString(checkDimStyle.Render("saved to settings.json — you won't be asked again"))
	b.WriteString("\n\n")
	b.WriteString(checkDimStyle.Render("enter confirm · esc cancel"))
	return b.String()
}
