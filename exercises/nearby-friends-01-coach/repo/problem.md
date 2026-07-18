# Design Nearby Friends

Design a nearby-friends feature: opted-in users see which of their
friends are currently close by, updating live as people move.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Core use case: two friends who've both opted in see each other on a
  map when nearby, and it updates as they move. What's out of scope —
  discovering strangers, turn-by-turn navigation?
- This is a location-*writing* problem before it's anything else: every
  opted-in user's device pings its position regularly, whether or not
  anyone is currently looking. What does that do to the write volume
  compared to a typical read-heavy app?
- Put numbers on it: active opted-in users × ping interval → location
  writes/sec, and how much state one user's location record needs.

## Suggested defaults (if you want a starting point)

- 100 million users have location sharing on; ~10% concurrently active
- Devices ping location roughly every 30 seconds while moving (less
  often while stationary — decide if that's in scope)
- Average friend list size a few hundred, but only the opted-in subset
  matters for this feature

## What good looks like

By the end you should have: stated assumptions with the location-write
QPS arithmetic (this should come out as the defining, write-heavy
number), a high-level diagram (clients ping a location service backed
by an in-memory store with a spatial index), a concrete geo-indexing
scheme with its precision/candidate-count trade-off and the
neighbor-cell edge case handled, a deliberate push-vs-pull answer for
telling a user their friend is nearby, a freshness/TTL story so stale
locations don't show ghosts, and how per-friend-pair opt-in is actually
enforced rather than just hidden in the UI.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
