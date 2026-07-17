# Reusable Barrier

A barrier for N participants: each calls `Wait()`, blocks until all N
have arrived, then everyone proceeds together — and the barrier resets
for the next round, over and over.

The starter releases everyone with a closed channel/flag — it works for
one round and falls apart on reuse. Reuse is the actual problem: a
participant from round 2 must not slip through round 1's release.

## The invariant the tests enforce

Across several rounds, no participant proceeds past `Wait()` for round
R until all N have arrived at round R — checked by counting arrivals
before each release.

API: `class Barrier { Barrier(int n); void Wait(); }` — mutex + condition variable + generation counter. Compiled with `-fsanitize=thread`.

Think: generations. How does a waiter know the wake-up it saw belongs
to *its* round?
