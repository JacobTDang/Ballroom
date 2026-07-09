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

-- Registers highlight groups + the add_highlight/clear_all RPC targets used
-- by tutor/chat.sh (issue #24: tutor-driven highlights/notes in the editor
-- pane). See lua/ballroom_highlight.lua for the implementation.
require("ballroom_highlight")
