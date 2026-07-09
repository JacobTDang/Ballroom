-- Tutor-driven code highlighting/notes (issue #24), wired over nvim's RPC
-- server socket (see docker/entrypoint.sh --listen and tutor/chat.sh, which
-- calls in via `nvim --server <sock> --remote-expr`).
--
-- Everything here is purely buffer-side: nvim_buf_add_highlight (ephemeral
-- highlight, namespaced) and a sign-column marker + nvim_buf_set_extmark
-- virtual text for the note. None of this ever calls nvim_buf_set_lines or
-- otherwise edits buffer content, so the highlighted/annotated file is byte
-- -for-byte identical to what gets saved/submitted.

local M = {}

local NS = vim.api.nvim_create_namespace("ballroom_tutor")
local SIGN_GROUP = "BallroomTutor"

vim.api.nvim_set_hl(0, "BallroomTutorHighlight", { bg = "#3C7DC4", fg = "#F2EBDD" })
vim.api.nvim_set_hl(0, "BallroomTutorNote", { fg = "#E0468C", italic = true })
vim.api.nvim_set_hl(0, "BallroomTutorSign", { fg = "#E0468C", bold = true })

vim.fn.sign_define("BallroomTutorSign", { text = "»", texthl = "BallroomTutorSign" })

-- Find a loaded buffer whose file matches `file` (exact path or basename).
-- Falls back to the current buffer so a relative/differently-rooted path
-- from the tutor still lands somewhere sane instead of silently no-op'ing.
local function resolve_buf(file)
  if file and file ~= "" then
    for _, buf in ipairs(vim.api.nvim_list_bufs()) do
      if vim.api.nvim_buf_is_loaded(buf) then
        local name = vim.api.nvim_buf_get_name(buf)
        if name == file or vim.fn.fnamemodify(name, ":t") == vim.fn.fnamemodify(file, ":t") then
          return buf
        end
      end
    end
  end
  return vim.api.nvim_get_current_buf()
end

--- Add a highlighted range + note. 1-indexed, inclusive line numbers (matches
--- how a human, and the tutor's prose, refers to lines).
---
--- Returns "ok" or an "ballroom_highlight error: ..." string — this is
--- reached over RPC from a shell script parsing model output, so a bad
--- range/args must degrade gracefully (report, don't throw) rather than
--- crash the caller or the editor.
function M.add_highlight(file, line_start, line_end, note)
  local ok, err = pcall(function()
    local buf = resolve_buf(file)
    local last_line = vim.api.nvim_buf_line_count(buf)
    line_start = math.max(1, math.floor(tonumber(line_start) or 1))
    line_end = math.min(last_line, math.floor(tonumber(line_end) or line_start))
    if line_start > last_line or line_end < line_start then
      error("line range out of bounds (buffer has " .. last_line .. " lines)")
    end

    for lnum = line_start, line_end do
      vim.api.nvim_buf_add_highlight(buf, NS, "BallroomTutorHighlight", lnum - 1, 0, -1)
    end
    vim.fn.sign_place(0, SIGN_GROUP, "BallroomTutorSign", buf, { lnum = line_start, priority = 20 })

    if note and note ~= "" then
      vim.api.nvim_buf_set_extmark(buf, NS, line_start - 1, 0, {
        virt_text = { { "  <- tutor: " .. note, "BallroomTutorNote" } },
        virt_text_pos = "eol",
      })
    end
  end)
  if not ok then
    return "ballroom_highlight error: " .. tostring(err)
  end
  return "ok"
end

--- Clear every tutor highlight/note in every loaded buffer. Session-scoped
--- only (extmarks/signs live in memory) — there is nothing on disk to undo.
function M.clear_all()
  for _, buf in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_is_loaded(buf) then
      vim.api.nvim_buf_clear_namespace(buf, NS, 0, -1)
      vim.fn.sign_unplace(SIGN_GROUP, { buffer = buf })
    end
  end
  return "ok"
end

return M
