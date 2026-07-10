package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JacobTDang/Ballroom/internal/catalog"
	"github.com/JacobTDang/Ballroom/internal/config"
	"github.com/JacobTDang/Ballroom/internal/exercise"
	"github.com/JacobTDang/Ballroom/internal/preflight"
	"github.com/JacobTDang/Ballroom/internal/tracker"
)

// listModelsFn and checkModelFn are vars (not direct calls) so tests can
// substitute fakes instead of making real HTTP calls to Ollama — same
// indirection pattern as buildImageFn in boot.go.
var listModelsFn = preflight.ListModels
var checkModelFn = preflight.CheckModel

// suggestedModels are known-good tutor model tags surfaced in the picker
// even before they're pulled locally, so they're discoverable without
// already knowing their exact Ollama tag. Selecting one that isn't
// pulled yet goes through the same not-pulled warn-then-confirm flow as
// typing an arbitrary tag (see handleModelEnter).
var suggestedModels = []string{
	config.DeepSeekCoderV2LiteModel,
	config.Qwen25Coder14BModel,
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
	stageStats
	stageModelPicker
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

// recentAttemptsLimit caps how many recent attempts Stats shows.
const recentAttemptsLimit = 10

// menuChoice is one of the main menu options.
type menuChoice int

const (
	menuPractice menuChoice = iota
	menuSandbox
	menuStats
	menuModelPicker
)

var menuLabels = []string{"Practice", "Sandbox", "Stats", "Model"}

var menuDescriptions = []string{
	"Pick a category and work through exercises",
	"Free practice, no grading, persists across sessions",
	"See your progress across categories",
	"Choose which Ollama model tutors your sessions",
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
type appResume struct {
	stage    appStage
	category string
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

	// stageCategories / source data for stageProblems
	problems       []catalog.ProblemStatus
	categories     []string
	categoryCursor int

	// stageDSACategories — second-level picker shown when the top-level
	// selection is the grouped DSA entry, listing the NeetCode roadmap
	// subcategories (Arrays & Hashing, Two Pointers, ...) it collapses.
	dsaCategories     []string
	dsaCategoryCursor int

	// stageProblems
	category         string
	categoryProblems []catalog.ProblemStatus
	problemCursor    int

	// stageLanguage
	selectedProblem catalog.ProblemStatus
	langCursor      int

	// stageStats
	statsStatuses []catalog.ExerciseStatus
	statsRecent   []tracker.Attempt

	// stageModelPicker
	modelLoading  bool
	modelLoadErr  error
	localModels   []string // exactly what Ollama reports pulled, for isLocalModel
	models        []string // localModels + suggestedModels, deduplicated
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
	if resume.stage == stageProblems {
		m = m.loadPractice()
		if m.err == nil {
			m.category = resume.category
			m.categoryProblems = filterByCategory(m.problems, resume.category)
			m.problemCursor = 0
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
			m.modelWarning = fmt.Sprintf("download failed: %v", msg.err)
			m.modelDownloadLines = nil
			return m, nil
		}
		return m.selectModel(m.modelDownloadTarget)

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
		case stageStats:
			return m.updateStats(msg)
		case stageModelPicker:
			return m.updateModelPicker(msg)
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
	case "1", "2", "3", "4":
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
	case menuSandbox:
		m.outcome = outcomeRunSandbox
		return m, tea.Quit
	case menuStats:
		return m.loadStats(), nil
	case menuModelPicker:
		return m.loadModelPicker()
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
	m.err = nil
	m.statsStatuses = statuses
	m.statsRecent = recent
	m.stage = stageStats
	return m
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

func filterByCategory(problems []catalog.ProblemStatus, category string) []catalog.ProblemStatus {
	var out []catalog.ProblemStatus
	for _, p := range problems {
		if p.Category == category {
			out = append(out, p)
		}
	}
	return out
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
		m.stage = stageProblems
	case "q", "esc", "ctrl+c":
		m.stage = stageCategories
	}
	return m, nil
}

// --- stageProblems ---

func (m appModel) updateProblems(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.problemCursor > 0 {
			m.problemCursor--
		}
	case "down", "j":
		if m.problemCursor < len(m.categoryProblems)-1 {
			m.problemCursor++
		}
	case "enter":
		if len(m.categoryProblems) == 0 {
			return m, nil
		}
		m.selectedProblem = m.categoryProblems[m.problemCursor]
		m.langCursor = 0
		m.stage = stageLanguage
	case "q", "esc", "ctrl+c":
		if catalog.IsGroupedCategory(m.category) {
			m.stage = stageDSACategories
		} else {
			m.stage = stageCategories
		}
	}
	return m, nil
}

// --- stageLanguage ---

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
		m.exerciseToRun = m.selectedProblem.Variants[m.langCursor].Exercise
		m.outcome = outcomeRunExercise
		return m, tea.Quit
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
		m.stage = stageMain
		return m, nil
	case tea.KeyEnter:
		return m.handleModelEnter()
	case tea.KeyRunes:
		// "q" with nothing typed yet backs out, matching every other
		// stage in this program — once the user has started typing,
		// every rune (including "q") feeds the filter/custom tag
		// instead, since it might be part of a real model name.
		if m.modelFilter == "" && string(msg.Runes) == "q" {
			m.stage = stageMain
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
func (m appModel) handleModelEnter() (tea.Model, tea.Cmd) {
	if len(m.modelFiltered) > 0 {
		sel := m.modelFiltered[m.modelCursor]
		if m.isLocalModel(sel) {
			return m.selectModel(sel)
		}
		return m.checkModelTag(sel)
	}

	tag := strings.TrimSpace(m.modelFilter)
	if tag == "" {
		return m, nil
	}
	return m.checkModelTag(tag)
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
// away, without waiting for a fresh Config.Load.
func (m appModel) selectModel(name string) (tea.Model, tea.Cmd) {
	m.cfg.TutorModel = name
	if err := config.SaveSettings(m.cfg.SettingsPath(), config.Settings{TutorModel: name}); err != nil {
		m.err = err
		return m, nil
	}
	m.err = nil
	m.stage = stageMain
	return m, nil
}

// --- View ---

func (m appModel) View() string {
	right := m.renderRight()
	if m.width == 0 || m.height == 0 {
		return right
	}
	panel := renderDashboardPanel(m.width, m.height, m.phase, right)
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
	case stageStats:
		return m.renderStats()
	case stageModelPicker:
		return m.renderModelPicker()
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
		b.WriteString(failStyle.Render("  " + m.err.Error()))
		b.WriteString("\n\n")
	}

	b.WriteString("\n")
	b.WriteString(menuSubtitleStyle.Render("↑/↓ or j/k move · 1-4 jump · enter select · q quit"))
	return b.String()
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

	for i, p := range m.categoryProblems {
		label := fmt.Sprintf("%-30s", truncateTitle(p.Title, 30))
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
		if i == m.problemCursor {
			b.WriteString(cursorRowStyle.Render(fmt.Sprintf("❯ %s %s", label, status)))
		} else {
			b.WriteString(fmt.Sprintf("  %s %s", label, statusStyle.Render(status)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("↑/↓ move · enter select · q back"))
	return b.String()
}

func (m appModel) renderLanguage() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render(m.selectedProblem.Title))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("choose a language"))
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
		b.WriteString(failStyle.Render("couldn't reach Ollama: " + m.modelLoadErr.Error()))
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
