# Transcode Video: Long-Running Operations

Design the API for a video platform's transcoding service: a client
submits an uploaded video for transcoding and the work takes anywhere
from thirty seconds to twenty minutes. Holding the connection open is
not an answer — design how callers start the work, track it, cancel
it, and collect the result.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go. The canonical pattern here is the
long-running operation (AIP-151) — try to derive its shape from the
requirements before you lean on it.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Resources: the videos, the transcode jobs — is a job its own
  resource? Who owns it, and what is its name?
- Callers: server-side services, a web console, mobile apps — which
  of them poll, and which can receive a push?
- What does the client actually need while work runs: progress? an
  ETA? partial results?

## Suggested defaults (if you want a starting point)

- 50,000 transcodes per day, p50 two minutes, p99 twenty minutes
- Callers: a web console (a person watching a progress bar) and
  server-side batch pipelines (thousands of jobs in flight)
- Results: a set of output renditions attached to the video

## What good looks like

By the end you should have: an operation resource with a name a client
can GET like anything else; a clean split between in-flight metadata
(progress, stage) and the terminal result (response or error, never
both); distinct terminal states including cancelled; a cancel method
whose semantics survive being called twice; a polling contract with
backoff plus an honest paragraph on the push alternative; and a
retention story for finished operations.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
