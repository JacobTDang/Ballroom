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

vim.opt.cursorline = true
vim.cmd.colorscheme("habamax")

vim.opt.laststatus = 3
vim.opt.statusline = " %{toupper(mode())} │ %f %m%r%h%w%=%y  ln %l/%L col %c "

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

-- Registers highlight groups + the add_highlight/clear_all RPC targets used
-- by tutor/chat.sh (issue #24: tutor-driven highlights/notes in the editor
-- pane). See lua/ballroom_highlight.lua for the implementation.
require("ballroom_highlight")
