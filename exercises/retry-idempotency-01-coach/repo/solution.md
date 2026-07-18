# My Design: Charge Exactly Once

## Step 1 — Scope & resource model

<!-- The Charge resource: fields, states, resource name (AIP-121/122).
     Who calls this, how often — and how often do they retry?
     Put numbers on charges/day and the retry tail. -->

## Step 2 — Methods & wire shapes

<!-- Create + get/list (AIP-131/133): request/response bodies with
     real field names, and status codes for success, validation
     failure, and a declined card (an application outcome — is that a
     transport error or a successful call?). -->

## Step 3 — The hard part: idempotency keys

<!-- Start from the failure: a timeout after the server committed.
     Why can't the client distinguish lost-request from lost-response?
     Then the mechanism (AIP-155 / Stripe canon): who generates the
     key and its scope; what the server stores; what a replay returns;
     same key + different body; and a second request arriving while
     the first is still in flight. Work one example end to end. -->

## Step 4 — Evolution & operations

<!-- TTL and storage: how long keys live, what expiry means for a
     very late retry, what the store costs at your scale. Which other
     methods need keys — and which are inherently idempotent
     already? -->
