# My Design: Evolving Without Breaking

## Step 1 — Scope & resource model

<!-- The surface being evolved and its clients: how many, how they
     integrate (SDK vs raw HTTP), upgrade lag. Define "breaking"
     concretely before sorting anything. -->

## Step 2 — Methods & wire shapes

<!-- The change taxonomy (AIP-180): sort the concrete list — add an
     optional field, remove/rename a field, tighten validation, add
     an enum value, change a default — into safe / breaking / subtle,
     with one sentence of why each. -->

## Step 3 — The hard part: the compatibility contract

<!-- The rules both sides obey: client unknown-field tolerance,
     server additive-only discipline, never reusing a name with a new
     meaning. Where the version lives (AIP-185) and what alpha/beta/
     GA promise (AIP-181). Work one example: the enum-value addition,
     end to end, on an old client. -->

## Step 4 — Evolution & operations

<!-- The v2 rollout: announce → dual-run → sunset, with headers/docs
     and telemetry proving who is still on v1; the server-side shim
     vs dual-stack trade-off; what running v1 forever actually
     costs. -->
