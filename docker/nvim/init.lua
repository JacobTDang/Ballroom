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

-- problem.md opens in the top-left pane every session, so make markdown
-- actually render instead of showing raw markers: conceal hides the
-- **/`/_ syntax via the bundled regex markdown syntax (the image's
-- nvim 0.9.5 has no treesitter markdown parser, so vim.treesitter.start
-- would error here -- verified against the real image), leaving bold
-- text bold and code spans highlighted. concealcursor keeps it
-- concealed in normal mode too -- the problem statement is read-only in
-- practice, and insert mode still reveals raw syntax for anyone who
-- does edit.
vim.api.nvim_create_autocmd("FileType", {
  pattern = "markdown",
  callback = function()
    vim.opt_local.conceallevel = 2
    vim.opt_local.concealcursor = "nc"
    vim.opt_local.wrap = true
    vim.opt_local.linebreak = true
  end,
})

-- Registers highlight groups + the add_highlight/clear_all RPC targets used
-- by tutor/chat.sh (issue #24: tutor-driven highlights/notes in the editor
-- pane). See lua/ballroom_highlight.lua for the implementation.
require("ballroom_highlight")
