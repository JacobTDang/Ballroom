# My Design: Transcode API with Long-Running Operations

## Step 1 — Scope & resource model

<!-- The nouns: videos, transcode operations. What is the operation's
     resource name? Who can see it? Which callers poll vs receive
     pushes? Put numbers on job volume and duration. -->

## Step 2 — Methods & wire shapes

<!-- How does a transcode start, and what comes back immediately?
     Write the actual shapes: the start call, GET on the operation,
     list. One worked example of an in-flight operation and a
     finished one. HTTP mapping and status codes. -->

## Step 3 — The hard part: operation lifecycle

<!-- Progress/metadata vs terminal response — where does each live?
     Terminal states: success, failure, cancelled — how does a client
     tell them apart? Cancel semantics: best-effort? idempotent?
     What happens on cancel-after-done? Polling contract (interval,
     backoff) and the push option — argue the trade-off. -->

## Step 4 — Evolution & operations

<!-- How long do finished operations live, and what does a client
     with a lost operation id do? What can you add to metadata
     without breaking pollers? Where does this API creak at 10x the
     job volume? -->
