-- Minimum viable nvim config for the practice image (spec Section 6 risk item).
-- Deliberately no plugin manager / LSP for MVP — bare editor with sane defaults,
-- but not an unstyled one: a built-in colorscheme and a real statusline so it
-- doesn't feel like a blank terminal you got dropped into.
vim.opt.number = true
vim.opt.relativenumber = true
vim.opt.expandtab = true
vim.opt.shiftwidth = 2
vim.opt.tabstop = 2
vim.opt.smartindent = true
vim.opt.mouse = "a"
vim.opt.ignorecase = true
vim.opt.smartcase = true
vim.opt.termguicolors = true
vim.opt.scrolloff = 4
vim.opt.signcolumn = "yes"
vim.opt.clipboard = "unnamedplus"

-- Autowrite (issue #226): the tutor's tools (internal/tutor/filecontext.go,
-- tools.go) and `ballroom submit` (internal/session/submit.go) both read
-- the solution straight off disk, so an unsaved buffer is invisible to
-- either -- forget :w and the tutor reviews stale code while submit
-- grades the wrong version. This built-in option covers Vim's own
-- write-before-leaving cases (:next, :make, ...); the autocmd block
-- further down covers the case that actually matters here: staying in
-- the same buffer and just switching to the tutor pane or pausing to
-- think.
vim.opt.autowrite = true

vim.opt.cursorline = true
vim.cmd.colorscheme("habamax")

vim.opt.laststatus = 3
vim.opt.statusline = " %{toupper(mode())} │ %f %m%r%h%w%=%y  ln %l/%L col %c │ press <space>? for help "

-- The problem statement opens as pre-rendered plain text (problem.txt,
-- see docker/entrypoint.sh + orchestrator.PrepareWorkspace) -- no
-- markdown conceal tricks needed here. linebreak just keeps any line
-- wider than the pane wrapping at word boundaries instead of mid-word.
vim.api.nvim_create_autocmd({ "BufRead", "BufNewFile" }, {
  pattern = "problem.txt",
  callback = function()
    vim.opt_local.wrap = true
    vim.opt_local.linebreak = true
  end,
})

-- Debounced write-on-change (issue #226), guarded to modifiable ordinary
-- buffers only (buftype == "" excludes terminals/quickfix/etc.;
-- expand("%") ~= "" excludes unnamed scratch buffers with nowhere to
-- write) so the read-only problem statement split (opened with `set
-- readonly nomodifiable`, see docker/tmux.conf's nvim invocation) is
-- never written to.
--
-- InsertLeave/BufLeave/FocusLost write immediately -- they're already
-- discrete, low-frequency events, not per-keystroke -- and FocusLost is
-- what catches switching to the tutor pane (M-2) mid-edit: tmux's
-- focus-events forwards that as a real FocusLost even though the outer
-- terminal window itself never lost OS focus (see tmux.conf's
-- focus-events comment). TextChanged fires far more often (any buffered
-- change, not every keystroke, but still often enough to matter), so
-- it's debounced with a short timer instead of writing on every change.
local function can_autowrite()
  return vim.bo.modifiable and vim.bo.buftype == "" and vim.fn.expand("%") ~= ""
end

local pending_write_timer = nil

local function debounced_write()
  if not can_autowrite() then
    return
  end
  if pending_write_timer then
    pending_write_timer:stop()
    pending_write_timer:close()
  end
  pending_write_timer = vim.defer_fn(function()
    pending_write_timer = nil
    if can_autowrite() then
      vim.cmd("silent! write")
    end
  end, 300)
end

vim.api.nvim_create_autocmd({ "InsertLeave", "BufLeave", "FocusLost" }, {
  pattern = "*",
  callback = function()
    if can_autowrite() then
      vim.cmd("silent! write")
    end
  end,
})

vim.api.nvim_create_autocmd("TextChanged", {
  pattern = "*",
  callback = debounced_write,
})

-- Registers highlight groups + the add_highlight/clear_all RPC targets used
-- by tutor/chat.sh (issue #24: tutor-driven highlights/notes in the editor
-- pane). See lua/ballroom_highlight.lua for the implementation.
require("ballroom_highlight")

-- A floating keybinding cheatsheet (<leader>? / :BallroomHelp) plus
-- <leader>w/<leader>s shortcuts, for someone who doesn't know vim
-- (issue #243). See lua/ballroom_help.lua for the implementation.
require("ballroom_help")
