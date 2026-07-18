# Cursor Pagination

List records in `id` order through an opaque, keyset-based cursor —
not an offset. The token names a position (the last key seen), not a
slot: unlike offset pagination, a walk stays correct while the
underlying data changes underneath it.

The rules:

- `List` returns up to `page_size` records with `id` greater than the
  cursor encoded in `page_token` (or from the beginning, if
  `page_token` is empty), ordered by `id` ascending.
- `page_size <= 0` falls back to a default; above the server's max, it
  clamps — the response is never held hostage to whatever the caller
  asks for.
- `next_page_token` is empty exactly when there's nothing left; a
  non-empty token always has a next page waiting.
- `page_token` is opaque and tamper-checked: a corrupted token is a
  loud error, never a best-effort guess. A token also remembers the
  `page_size` it was issued for — resuming with a different page size
  is rejected, not silently honored.

## The invariant the tests enforce

A full walk (chaining tokens until empty) visits every record that
existed at the start exactly once, in order — even if records are
inserted mid-walk. That's the whole point of a keyset cursor over an
offset: offsets shift under mutation and silently duplicate or skip
records; a keyset cursor can't, because it names a key, not a
position.

API: `NewCursorStore(records []Record) *CursorStore`, `Insert(r Record) error`, `List(pageSize int, pageToken string) (items []Record, nextPageToken string, err error)` — `err` is non-nil on an invalid/tampered token or a page-size change mid-walk.
