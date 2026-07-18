# My Design: Rate Limits & Quotas

## Step 1 — Scope & resource model

<!-- The two mechanisms and their windows: billable long-window
     quotas vs short-window protective rate limits. The tiers, the
     traffic numbers, and what abuse looks like concretely. -->

## Step 2 — Methods & wire shapes

<!-- What every limited response carries: the 429 error shape
     (consistent with the platform error model), Retry-After, and
     the limit/remaining/reset headers on ordinary responses too.
     One worked example of a limited call. -->

## Step 3 — The hard part: the algorithm and its edges

<!-- Token bucket / sliding window — chosen, with burst behavior
     explained, and the fixed-window boundary flaw named. Cost
     weighting for expensive endpoints. Where the counters live when
     there are three gateways — and what the distributed counting
     trade-off costs you. -->

## Step 4 — Evolution & operations

<!-- Client guidance you document: backoff with jitter, SDK
     behavior, programmatic discovery of one's own limits. Fairness:
     tenant isolation, tier changes, graceful degradation vs hard
     cutoff, and how you change a limit without an incident. -->
