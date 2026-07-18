package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// stageHelp (issue #242) is a help screen reachable with "?" from every
// stage that takes key input -- the keybindings used to live only in
// README.md, outside the app. Esc/q/? all close it and return to
// m.helpOrigin, the stage it was opened from (see openHelp/updateHelp),
// so backing out never dumps you back at the main menu if you opened
// help from three screens deep.

// helpKey is one key/description row in the table below.
type helpKey struct {
	key  string
	desc string
}

// helpSection is one titled group of rows.
type helpSection struct {
	title string
	keys  []helpKey
}

// helpSections is the single place every keybinding this program
// surfaces gets edited -- renderHelp just walks it. Three groups, matching
// where each set of keys is actually read:
//
//   - Host keys: this program's own tea.KeyMsg handling (updateMain,
//     updateCategories, updateDSACategories, updateProblems,
//     updateLanguage, updateSearch, updateStats, updateSettings).
//   - In-session keys: docker/tmux.conf's `bind -n` lines -- the tmux
//     session wrapping the editor/tutor/terminal panes.
//   - Tutor pane keys: internal/tutor/model.go's Update, the tutor
//     chat's own key handling once that pane has focus.
//
// Every row here was checked directly against those sources rather than
// against README.md's own (separately maintained) table, so this can't
// silently document a key that doesn't exist -- or a tmux bind
// (Ctrl-Shift-Tab) that exists but the README table doesn't mention.
var helpSections = []helpSection{
	{
		title: "Host keys",
		keys: []helpKey{
			{"up/down · j/k", "move the cursor (j/k on menus; pickers reserve letters for search)"},
			{"1-5", "jump straight to a main-menu item"},
			{"enter", "select"},
			{"/", "search the whole catalog, from anywhere"},
			{"?", "show this help"},
			{"q / esc", "back a screen; q quits from the main menu"},
		},
	},
	{
		title: "In-session keys",
		keys: []helpKey{
			{"M-1 / M-2 / M-3", "jump to editor / tutor / terminal pane"},
			{"Ctrl-Tab (or M-Tab)", "cycle to the next pane"},
			{"Ctrl-Shift-Tab", "cycle to the previous pane"},
			{"M-q", "submit (pre-types the command; enter confirms)"},
			{"M-0", "leave the session and return to the picker (y/n confirm)"},
			{"M-h", "toggle tutor highlights/notes in the editor"},
		},
	},
	{
		title: "Tutor pane keys",
		keys: []helpKey{
			{"enter", "send your message"},
			{"Ctrl-D", "exit the tutor chat"},
			{"PgUp / PgDn", "scroll the conversation"},
		},
	},
}

// helpKeyColWidth is the fixed key-column width the description hangs
// off of -- wide enough for the longest key string above
// ("Ctrl-Shift-Tab") plus a little breathing room.
const helpKeyColWidth = 20

// openHelp records which stage help was opened from (so Esc/q/? can
// return to it -- see updateHelp) and switches to stageHelp. Every
// stage's update func that takes a "?" calls this the same way.
func (m appModel) openHelp() (tea.Model, tea.Cmd) {
	m.helpOrigin = m.stage
	m.stage = stageHelp
	return m, nil
}

// updateHelp handles stageHelp's input: Esc, q, and ? are all "close
// help" -- ? doubles as both the open and close key so there's no need
// to remember a different one once you're already looking at the list
// that documents it.
func (m appModel) updateHelp(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc, tea.KeyCtrlC:
		m.stage = m.helpOrigin
		return m, nil
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "q", "?":
			m.stage = m.helpOrigin
			return m, nil
		}
	}
	return m, nil
}

// renderHelp draws helpSections inside the same dashboard panel shell
// every other stage uses (see renderStats for the pattern this mirrors).
func (m appModel) renderHelp() string {
	var b strings.Builder
	b.WriteString(hintStyle.Render("Help"))
	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("keys across the host menu, the session, and the tutor"))
	b.WriteString("\n\n")

	for i, section := range helpSections {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(hintStyle.Render(section.title))
		b.WriteString("\n")
		for _, k := range section.keys {
			label := fmt.Sprintf("%-*s", helpKeyColWidth, k.key)
			b.WriteString(fmt.Sprintf("  %s %s\n", categoryStyle.Render(label), checkDimStyle.Render(k.desc)))
		}
	}

	b.WriteString("\n")
	b.WriteString(checkDimStyle.Render("esc/q/? back"))
	return b.String()
}
