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

-- Underlying notes/data store, independent of what's currently rendered
-- (issue #25: toggling visibility off must hide rendering without deleting
-- this). Each entry is {buf, line_start, line_end, note}. M.visible tracks
-- on/off; render() (below) is the only thing that touches the namespace's
-- extmarks/signs, so add_highlight and toggle() both funnel through it.
M.notes = {}
M.visible = true

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

-- Actually paints one stored note's highlight/sign/virtual text into its
-- buffer. The only place that touches the namespace's extmarks/signs, so
-- both add_highlight (new note) and toggle() (replaying stored notes back
-- on) go through here and can't drift apart.
local function render(rec)
  for lnum = rec.line_start, rec.line_end do
    vim.api.nvim_buf_add_highlight(rec.buf, NS, "BallroomTutorHighlight", lnum - 1, 0, -1)
  end
  vim.fn.sign_place(0, SIGN_GROUP, "BallroomTutorSign", rec.buf, { lnum = rec.line_start, priority = 20 })

  if rec.note and rec.note ~= "" then
    vim.api.nvim_buf_set_extmark(rec.buf, NS, rec.line_start - 1, 0, {
      virt_text = { { "  <- tutor: " .. rec.note, "BallroomTutorNote" } },
      virt_text_pos = "eol",
    })
  end
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

    local rec = { buf = buf, line_start = line_start, line_end = line_end, note = note }
    table.insert(M.notes, rec)
    -- Store the note regardless of visibility (issue #25: toggling off must
    -- not drop notes added while hidden), but only paint it if currently on.
    if M.visible then
      render(rec)
    end
  end)
  if not ok then
    return "ballroom_highlight error: " .. tostring(err)
  end
  return "ok"
end

--- Toggle visibility of every tutor highlight/note on/off (issue #25).
--- Turning off only clears the namespace's rendered extmarks/signs — it
--- never touches M.notes, so the underlying data survives untouched.
--- Turning back on replays every stored note through render(), so what
--- reappears is exactly what was hidden (plus anything added while off).
---
--- Returns "shown" or "hidden" so a caller (e.g. the tmux keybind) could
--- surface the new state if it wanted to.
function M.toggle()
  M.visible = not M.visible
  if M.visible then
    for _, rec in ipairs(M.notes) do
      if vim.api.nvim_buf_is_valid(rec.buf) then
        render(rec)
      end
    end
  else
    for _, buf in ipairs(vim.api.nvim_list_bufs()) do
      if vim.api.nvim_buf_is_loaded(buf) then
        vim.api.nvim_buf_clear_namespace(buf, NS, 0, -1)
        vim.fn.sign_unplace(SIGN_GROUP, { buffer = buf })
      end
    end
  end
  return M.visible and "shown" or "hidden"
end

--- Clear every tutor highlight/note in every loaded buffer, and forget the
--- underlying data too (unlike toggle(), this is a real delete — there is
--- nothing on disk to undo, extmarks/signs/M.notes all live in memory only).
function M.clear_all()
  for _, buf in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_is_loaded(buf) then
      vim.api.nvim_buf_clear_namespace(buf, NS, 0, -1)
      vim.fn.sign_unplace(SIGN_GROUP, { buffer = buf })
    end
  end
  M.notes = {}
  return "ok"
end

return M
