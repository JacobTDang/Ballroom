-- A way in and out of the editor for someone who doesn't know vim
-- (issue #243): a floating cheatsheet (<leader>? / :BallroomHelp), plus
-- <leader>w to write and <leader>s to jump to the terminal pane and
-- pre-type the submit command -- the two things a novice most needs a
-- shortcut for. mapleader is set here (before any of this file's own
-- <leader> mappings, which is the one thing nvim requires ordering-wise)
-- rather than in init.lua, so this whole feature -- state, keys, and
-- command -- lives in one file, the same way ballroom_highlight.lua is
-- self-contained and just required from init.lua.

vim.g.mapleader = " "

local M = {}

-- win/buf track the open float so toggle() can close it instead of
-- stacking a second one on top -- nvim_win_is_valid guards against a
-- stale handle from a float the user already closed some other way
-- (e.g. a plain :q inside it).
M.win = nil
M.buf = nil

-- essentials is the single table this module renders from -- one place
-- to edit, same reasoning as internal/tui/help.go's helpSections on the
-- host side. Every vim row below is a default Neovim command untouched
-- by init.lua's options; every session row was checked directly against
-- docker/tmux.conf's `bind -n` lines rather than copied from memory.
local essentials = {
  { header = "vim" },
  { "i", "insert mode" },
  { "Esc", "back to normal mode" },
  { ":w", "write (save)" },
  { ":q", "close this split/window" },
  { "u", "undo" },
  { "dd", "delete (cut) the current line" },
  { "gg", "go to the top" },
  { "G", "go to the bottom" },
  { "/", "search forward" },
  { ":noh", "clear search highlighting" },
  { header = "session (works from the editor too)" },
  { "M-1 / M-2 / M-3", "jump to editor / tutor / terminal pane" },
  { "M-q", "submit (pre-types the command; enter confirms)" },
  { "M-0", "leave the session (y/n confirm)" },
  { "M-h", "toggle tutor highlights/notes" },
  { header = "this window" },
  { "<leader>?", "toggle this help" },
  { "<leader>w", "write (save)" },
  { "<leader>s", "jump to the terminal, pre-type submit" },
  { "q / Esc", "close this window" },
}

-- keyColWidth is the fixed key-column width the description hangs off
-- of -- wide enough for the longest key string above ("M-1 / M-2 / M-3").
local keyColWidth = 17

-- render lays essentials out as plain text lines, and returns the width
-- of the widest one -- the float is sized to fit this exactly (capped by
-- the available pane size, see open() below) rather than a guessed fixed
-- size, so it stays readable instead of clipping at small pane sizes.
local function render()
  local lines = { "Ballroom help", "" }
  local width = 0
  for _, row in ipairs(essentials) do
    local line
    if row.header then
      line = row.header
    else
      line = string.format("  %-" .. keyColWidth .. "s %s", row[1], row[2])
    end
    table.insert(lines, line)
    width = math.max(width, #line)
  end
  table.insert(lines, "")
  table.insert(lines, "q or Esc closes this window")
  width = math.max(width, #"q or Esc closes this window")
  return lines, width
end

-- close() tears the float down if it's still open -- safe to call
-- unconditionally (toggle() and the buffer-local q/Esc maps both do).
function M.close()
  if M.win and vim.api.nvim_win_is_valid(M.win) then
    vim.api.nvim_win_close(M.win, true)
  end
  M.win = nil
  M.buf = nil
end

-- open() builds a fresh scratch buffer and floats it over the editor,
-- centered and sized to the content (see render()), never wider/taller
-- than the pane itself has room for -- the editor pane runs at roughly
-- half the outer terminal's height, so this has to shrink to fit there,
-- not just on a full-size terminal.
function M.open()
  local lines, contentWidth = render()
  local width = math.min(contentWidth + 2, math.max(vim.o.columns - 4, 20))
  local height = math.min(#lines, math.max(vim.o.lines - 4, 5))

  M.buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_buf_set_lines(M.buf, 0, -1, false, lines)
  vim.bo[M.buf].modifiable = false
  vim.bo[M.buf].buftype = "nofile"
  vim.bo[M.buf].filetype = "ballroom-help"

  M.win = vim.api.nvim_open_win(M.buf, true, {
    relative = "editor",
    width = width,
    height = height,
    row = math.floor((vim.o.lines - height) / 2),
    col = math.floor((vim.o.columns - width) / 2),
    border = "rounded",
    style = "minimal",
    title = " Ballroom help ",
    title_pos = "center",
  })

  local opts = { buffer = M.buf, nowait = true, silent = true }
  vim.keymap.set("n", "q", M.close, opts)
  vim.keymap.set("n", "<Esc>", M.close, opts)
end

-- toggle() is what <leader>? and :BallroomHelp both call -- open if
-- closed, close if already open, so mashing the same key never stacks
-- a second float on top of the first.
function M.toggle()
  if M.win and vim.api.nvim_win_is_valid(M.win) then
    M.close()
  else
    M.open()
  end
end

vim.api.nvim_create_user_command("BallroomHelp", M.toggle, {
  desc = "Toggle the Ballroom keybinding cheatsheet",
})

vim.keymap.set("n", "<leader>?", M.toggle, { desc = "Toggle the Ballroom keybinding cheatsheet" })
vim.keymap.set("n", "<leader>w", "<cmd>write<CR>", { desc = "Write the current buffer" })

-- <leader>s mirrors docker/tmux.conf's M-q binding exactly (same target
-- pane, same tmux command chain, same "pre-type but don't run" contract
-- -- send-keys types the command into pane 2 without a trailing Enter, so
-- a stray press can't burn an attempt) rather than inventing a second way
-- to reach the terminal pane. Shelled out to the tmux CLI -- same as
-- tmux.conf's own run-shell bind for M-h -- since nvim has no direct way
-- to drive a sibling tmux pane.
vim.keymap.set("n", "<leader>s", function()
  vim.fn.system([[tmux select-pane -t :0.2 \; send-keys "ballroom submit"]])
end, { desc = "Jump to the terminal pane and pre-type ballroom submit" })

return M
