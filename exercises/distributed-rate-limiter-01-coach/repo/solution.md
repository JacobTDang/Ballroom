# My Design: Distributed Rate Limiter

## Step 1 — Use cases, constraints, estimates

<!-- What's being limited, per what key (client/user/IP), across how
     many servers? What does the client see when it's over the limit?
     Estimate: total decision QPS the limiter must sustain, number of
     distinct clients and how many are concurrently active, memory per
     client's counter state. -->

## Step 2 — High-level design

<!-- Where does the check happen (gateway/middleware vs each service)?
     What's the shared state behind it, and why does it need to be
     shared at all across a fleet of stateless API servers? -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Algorithm: token bucket vs fixed window vs sliding-window
       log/counter. Demonstrate the fixed window's boundary-burst flaw,
       don't just name it. Which one, and why?
     - Distributed correctness: the check-then-increment race across
       concurrent requests hitting different API servers — how is it
       made atomic (Lua script, INCR+EXPIRE, or equivalent)?
     - Config: how do per-client tiers and per-endpoint limits get to
       every API server? -->

## Step 4 — Scale it

<!-- Per-request latency budget for the check itself — can it be local-
     first with a global backstop? Fail-open or fail-closed when the
     counter store is unavailable, and what that choice costs. How does
     the counter store scale (sharding by client key), and what happens
     when one client is a hot key? Anything different across regions? -->
