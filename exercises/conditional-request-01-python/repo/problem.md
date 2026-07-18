# Conditional Requests

A versioned key-value store: every write gets a fresh **etag**, and
reads and writes can be made *conditional* on it — the mechanism
behind HTTP's `If-Match` / `If-None-Match` and the fix for the classic
lost-update bug (two clients read the same resource, both write back,
the second write silently erases the first).

There's no injected clock here — versions are sequence-based, not
time-based.

## `Get(key, if_none_match)`

- Key doesn't exist: **404**.
- Key exists and `if_none_match` matches its current etag: **304**,
  the caller's cached copy is still good.
- Otherwise: **200** with the current `(etag, body)`.

## `Put(key, body, if_match)`

- Key doesn't exist yet (a create):
  - `if_match` unset: **200**, the key is created with a fresh etag.
  - `if_match` set to anything: **412** — you claimed a version of
    something that doesn't exist.
- Key already exists (an update):
  - `if_match` unset: **428** (precondition required) — blind
    overwrites of an existing resource are refused, on purpose. This
    is the whole point: no client accidentally clobbers another
    client's write.
  - `if_match` set but stale (doesn't match the current etag): **412**
    — someone else wrote first. The store's state is **left exactly
    as it was** — a failed conditional write must never partially
    apply.
  - `if_match` set and current: **200**, the value updates and a new
    etag is issued. The old etag never matches again.

## `Delete(key)`

Removes the key. Deleting a key that doesn't exist is a no-op.

## The invariant the tests enforce

The full status-code matrix above, a 412 leaving state byte-for-byte
unchanged (checked with a follow-up `Get`), and no etag resurrection:
deleting a key and recreating it must never hand out an etag that
matches what a client cached from before the delete — not on the read
side (a stale `if_none_match` must not 304 against the new resource)
and not on the write side (a stale `if_match` must not succeed against
it).

API: `ConditionalStore()`, `.get(key, if_none_match=None) -> (status, etag, body)`, `.put(key, body, if_match=None) -> (status, etag)`, `.delete(key)`.
