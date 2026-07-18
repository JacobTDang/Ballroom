# Upload 10 GB, Resumably

Design the upload surface for a media service: clients upload files
from under a megabyte (thumbnails) to ten gigabytes (raw video), often
over connections that drop mid-transfer. A failed 9-GB upload that
starts over from zero is the failure mode you are designing away.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go. The canon here is the resumable
upload protocol used by the big storage APIs (session URI, offset
negotiation, 308) — derive what you can from the requirements first.

## Scope to establish in step 1

Talk through and pin down with your coach:

- The resource model: what is an upload *session*, and how does it
  relate to the final media resource? Who may touch a session?
- Client mix: mobile on flaky links, servers on fat pipes — do they
  deserve the same protocol?
- What failure data do you assume? Drop rates, typical file sizes.

## Suggested defaults (if you want a starting point)

- File sizes: p50 40 MB, p95 2 GB, max 10 GB; thumbnails under 1 MB
- A multi-hour mobile upload sees a connection drop roughly hourly
- Metadata (title, tags) accompanies every upload

## What good looks like

By the end you should have: a tiered protocol (single-shot for small,
resumable for large, with thresholds); a session model with initiation,
expiry, and auth scope; the offset dance — chunk uploads, the
query-committed-offset recovery step, server-authoritative offsets,
out-of-order rejection; integrity checks and what a mismatch does; a
finalization step that produces the durable resource; and a retry
matrix saying which failures retry the chunk vs restart the session.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
