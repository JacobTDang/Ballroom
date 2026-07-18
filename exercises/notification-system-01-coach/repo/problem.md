# Design a Notification System

Design a notification system that delivers push, SMS, and email
notifications to users at scale, on behalf of many other internal
services.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Core use case: one internal API that other services call to send a
  notification, fanning out to push/SMS/email as appropriate. What's
  out of scope — an in-app notification inbox, marketing-campaign
  tooling?
- The three channels aren't interchangeable: they differ wildly in
  cost, latency, and how they fail. How should that shape the design?
- Put numbers on it: notifications/day → peak sends/sec, and per-channel
  cost/latency differences worth naming.

## Suggested defaults (if you want a starting point)

- An internal platform used by thousands of other services, sending
  ~500 million notifications/day, peaking at several times the average
- Three channels: push (APNs/FCM), SMS, email
- Some notifications are time-critical (a one-time login code); others
  can wait minutes to hours (a weekly digest) — decide how that's
  represented

## What good looks like

By the end you should have: stated assumptions with the peak-sends-per-
second arithmetic, a high-level diagram (ingest API → preference check
→ per-channel queues → channel workers → third-party providers, fully
async), how user preferences/opt-outs and per-user rate caps are
enforced before a send goes out, a reliability story built on
idempotency keys and retries with backoff (with an explicit
at-least-once choice), a way for a one-time code to jump the queue
ahead of a digest, and a scaling story where one slow channel or
provider can't stall the others.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
