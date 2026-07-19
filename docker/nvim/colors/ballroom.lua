-- The editor's own colorscheme, replacing the stock habamax nvim ships
-- with -- this is where the user spends the most time in a session, so
-- it gets the app's real palette instead of an unrelated default theme.
--
-- The hex values below are hand-copied from internal/palette (Lua can't
-- import Go, and this container has no build step over its config --
-- see internal/palette's own doc comment), named to match that
-- package's constants so the two stay easy to eyeball against each
-- other. internal/palette/docker_drift_test.go greps this file (along
-- with tmux.conf and clock.sh) and asserts every hex it finds is a real
-- palette color, so a typo or an invented color here fails a Go test
-- instead of just looking wrong in a session.
--
-- No treesitter parsers ship in this image (init.lua's own comment:
-- "Deliberately no plugin manager / LSP for MVP"), so highlighting
-- comes entirely from Neovim's bundled classic regex syntax files
-- (syntax/go.vim, python.vim, cpp.vim, ...) -- which is exactly the
-- classic group set defined below (Statement, Type, Identifier, ...),
-- not treesitter's @-prefixed captures.

local teal = "#2FA6A6" -- palette.Teal -- "good": passes, code, the tutor's own voice
local gold = "#E8A93C" -- palette.Gold -- attention: streaks, due markers, hints
local red = "#F03C3C" -- palette.Red -- failures
local purple = "#9B5FB0" -- palette.Purple -- selection, language column
local blue = "#3C7DC4" -- palette.Blue -- category labels
local orange = "#F0862E" -- palette.Orange -- banner mid-tone
local cyan = "#3ED6D6" -- palette.Cyan -- disco-ball sparkle

local cream = "#F2EBDD" -- palette.Cream -- brightest text
local pale_gray = "#D9D3C4" -- palette.PaleGray -- ordinary dim text
local warm_gray = "#96918B" -- palette.WarmGray -- pane metadata
local dim_gray = "#6B6B6B" -- palette.DimGray -- disco-ball body

local rule = "#3A3D4D" -- palette.Rule -- borders, rules, card frames
local ink = "#000000" -- palette.Ink -- text on a colored background

local card_bg = "#14151C" -- palette.CardBg -- editor-card body
local card_header_bg = "#1E2029" -- palette.CardHeaderBg -- card header bar, status bar row
local gutter_fg = "#5C5852" -- palette.GutterFg -- line-number gutter

-- Well-behaved colorscheme boilerplate (:help colorscheme): reset
-- anything a previous scheme left behind before defining this one, so
-- re-sourcing (or an earlier stock scheme having already run) can never
-- leave a stray group behind.
vim.cmd("highlight clear")
if vim.fn.exists("syntax_on") == 1 then
  vim.cmd("syntax reset")
end
vim.o.background = "dark"
vim.g.colors_name = "ballroom"

-- palette.Pink (#E0468C) is deliberately never used below: it's already
-- the tutor's note/highlight color in this same editor
-- (ballroom_highlight.lua's BallroomTutorNote/Sign), and reusing it for
-- ordinary syntax or UI chrome would make a real tutor note harder to
-- spot against everyday pink keywords or identifiers instead of easier.

local groups = {
  -- Base text: the brightest tone (Cream) on the same near-black
  -- surface the tutor pane's own editor cards use (CardBg) -- the
  -- session's two editors (nvim here, the tutor's read-only code cards)
  -- share one "this is code" background.
  Normal = { fg = cream, bg = card_bg },

  -- Identifiers (variable/parameter names) stay close to Normal rather
  -- than getting their own accent -- they're the single most frequent
  -- token in real code, and giving every one of them a loud color reads
  -- as noise, not signal. PaleGray ("ordinary dim text") is the
  -- palette's own name for exactly this role.
  Identifier = { fg = pale_gray },

  -- Comments: dim and italic, deliberately close to the muted tone
  -- Monokai (the tutor pane's own chroma style, markdown.go) uses for
  -- its comments -- the editor and the tutor's code cards read as
  -- related, not two unrelated syntax themes.
  Comment = { fg = dim_gray, italic = true },

  String = { fg = cyan },
  Function = { fg = teal, bold = true },
  Number = { fg = orange },
  Constant = { fg = blue },
  -- Type doubles as palette.Purple's "language column" role -- a type
  -- name and a language name (Go, Python, C++) are the same idea, "what
  -- kind of thing is this."
  Type = { fg = purple },
  -- Keyword/Statement share Gold's "attention" role (control flow is
  -- exactly what should draw the eye first) -- Keyword bold for the
  -- word itself (if/for/return/import), Statement plain so a bare
  -- control-flow region doesn't out-bold the keyword that names it.
  Keyword = { fg = gold, bold = true },
  Statement = { fg = gold },

  -- Line numbers: the ordinary gutter is palette's own dedicated
  -- gutter tone; the cursor's own line uses Gold ("attention") bold so
  -- "where am I" is answerable at a glance the same way the tutor
  -- pane's own mode pill or the host's due-markers use Gold for a nudge.
  LineNr = { fg = gutter_fg },
  CursorLineNr = { fg = gold, bold = true },
  -- CursorLine only sets bg, deliberately -- the same "layer a tint
  -- over what's already there" contract Visual below relies on, so
  -- syntax colors on the current line stay exactly as colored. The
  -- header-bar tone is a level brighter than CardBg by design
  -- elsewhere in this app (card.go's card header vs. body), so reusing
  -- it here keeps "a subtly raised surface" meaning one thing.
  CursorLine = { bg = card_header_bg },

  -- Visual selection: bg only (no fg), so every token keeps its own
  -- syntax color under the highlight instead of the selection
  -- flattening everything to one color. Rule -- "borders, rules, card
  -- frames" -- reads as structure being drawn over the text, and never
  -- collides with Type's own Purple the way reusing "selection" here
  -- literally would.
  Visual = { bg = rule },
  -- Search matches get a full fg+bg override (unlike Visual) because a
  -- match should be unmistakable regardless of what it's a match
  -- inside of -- Ink-on-Gold guarantees contrast even over Gold's own
  -- Keyword/Statement text.
  Search = { fg = ink, bg = gold, bold = true },

  StatusLine = { fg = cream, bg = card_header_bg, bold = true },
  StatusLineNC = { fg = warm_gray, bg = card_header_bg },

  Pmenu = { fg = pale_gray, bg = card_header_bg },
  PmenuSel = { fg = ink, bg = teal, bold = true },

  -- Not in the requested group list, but left alone they're the
  -- loudest thing on screen: `highlight clear` (above) resets every
  -- group this file doesn't touch back to Vim's own literal defaults,
  -- and a few of those are hardcoded named colors nobody's dark-mode
  -- palette overrides automatically (verified live -- SignColumn came
  -- back "Cyan on Grey", a real vim default with roots in an
  -- assumed-light terminal, not this colorscheme choosing cyan).
  --
  -- SignColumn/FoldColumn: init.lua sets signcolumn=yes, so this two
  -- cell strip runs down every line of every file, always -- it reads
  -- as part of the gutter, so it gets the gutter's own tone.
  SignColumn = { fg = gutter_fg, bg = card_bg },
  FoldColumn = { fg = gutter_fg, bg = card_bg },
  -- EndOfBuffer/NonText: the `~` past-EOF markers -- also the gutter
  -- family, not code, so the same dim tone rather than Vim's default
  -- bold blue.
  EndOfBuffer = { fg = gutter_fg },
  NonText = { fg = gutter_fg },
  -- MatchParen: the bracket a cursor sits on/near is common enough in
  -- real code to matter -- the same "raised surface" CursorLine and
  -- Pmenu already use, so a matched pair reads as a highlight, not an
  -- alarm.
  MatchParen = { bg = card_header_bg, bold = true },
  -- Special/SpecialKey: escape sequences and format specifiers inside
  -- otherwise-plain text (e.g. \n in a string) -- String's own Cyan,
  -- since that's almost always where they actually appear.
  Special = { fg = cyan },
  SpecialKey = { fg = cyan },
  -- Directory: netrw's own listing color (entrypoint.sh opens netrw
  -- when a session has no solution.* file to open directly) -- Type's
  -- own Purple, the same "what kind of thing is this" role a
  -- directory entry plays next to a file one.
  Directory = { fg = purple },

  -- No LSP ships in this image (init.lua), so nothing populates these
  -- today -- defined anyway for completeness and so a future LSP
  -- addition inherits the palette instead of Neovim's own defaults.
  DiagnosticError = { fg = red },
  DiagnosticWarn = { fg = gold },
  DiagnosticInfo = { fg = blue },
  DiagnosticHint = { fg = cyan },
}

for group, opts in pairs(groups) do
  vim.api.nvim_set_hl(0, group, opts)
end
