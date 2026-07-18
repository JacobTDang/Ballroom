# My Design: Payments Error Surface

## Step 1 — Scope & resource model

<!-- Who consumes errors (SDKs, dashboards, alerting), and what must
     each be able to decide from the error alone? List the failure
     families a payments API actually produces. -->

## Step 2 — Methods & wire shapes

<!-- The canonical code set (AIP-193 / google.rpc flavor): each code,
     its HTTP mapping, one example trigger. The envelope every
     endpoint shares — write it as JSON with real field names, and
     one fully worked example error response. -->

## Step 3 — The hard part: the retry contract

<!-- Classify every code retryable / not / conditionally. Where does
     Retry-After apply? The timeout-after-charge ambiguity: what do
     you tell the client to do (this seeds the idempotency question
     later in the track)? How do machine-readable details
     (reason/domain/metadata) keep clients off message-string
     parsing? -->

## Step 4 — Evolution & operations

<!-- What may change in an error later (new codes? new detail
     fields?) without breaking SDKs? The no-leak policy: what never
     appears in any error, and your 404-vs-403 disclosure call.
     Partial failure: the shape a batch endpoint uses per item. -->
