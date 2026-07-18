# Design WhatsApp

Design a messaging service like WhatsApp: one-to-one and group chats
with delivery receipts.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Core use cases: 1:1 messaging, group messaging, and delivery states
  (sent/delivered/read). What's out of scope — voice/video calls,
  stories, payments?
- Messaging is fundamentally different from a request/response API:
  a message has to reach a specific device that may not be connected
  right now. What kind of connection does that call for?
- Put numbers on it: concurrent connections at peak, messages/sec, and
  how many connections a single chat server can realistically hold —
  that last number is what sizes your fleet.

## Suggested defaults (if you want a starting point)

- 2 billion users, ~100 million devices concurrently connected at peak
- ~100 billion messages/day
- A single chat server can hold on the order of a million persistent
  connections (say what technology choice gets you there)
- End-to-end encryption is in place; the server routes ciphertext and
  never sees message contents

## What good looks like

By the end you should have: stated assumptions with the
connections-per-server and messages/sec arithmetic, a high-level design
built on persistent connections to chat servers plus a session registry
mapping user → server, a message narrated end to end (sender → their
server → registry lookup → recipient's server → recipient), how
sent/delivered/read receipts flow back and how offline recipients are
handled without losing messages, a deliberate group-send fanout
strategy, and a scaling story where the session registry's own
availability is treated as the critical shared state it is.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
