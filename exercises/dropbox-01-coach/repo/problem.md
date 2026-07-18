# Design Dropbox

Design a file-hosting service like Dropbox: files sync across a user's
devices and can be shared with others.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Core use cases: upload/download, sync across a user's own devices,
  sharing a file or folder with another user. What's out of scope —
  real-time collaborative editing, in-browser previews?
- A file can be huge and can change a little at a time. What does that
  imply for what gets re-uploaded on an edit?
- Put numbers on it: users × average storage each, daily active
  uploaders, upload bandwidth, and where dedup could cut the storage
  bill.

## Suggested defaults (if you want a starting point)

- 500 million users, ~2 GB of storage used on average → measured in
  exabytes, not terabytes
- ~10% of users sync at least one change on a given day
- A reasonable chunk size for splitting files (a few MB) — pick one and
  justify it
- Sharing is link-based or invite-based; both are fine to assume

## What good looks like

By the end you should have: stated assumptions with the storage and
bandwidth arithmetic, a high-level diagram separating the metadata
service (file tree, versions, block lists) from block storage (the
bytes themselves), a chunking scheme with content-addressed blocks that
enables delta sync and dedup, the sync protocol narrated end to end
(edit on one device → other devices notified → blocks pulled), a
deliberate answer for concurrent edits to the same file, and a scaling
story for the metadata database and block storage separately.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
