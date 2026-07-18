# My Design: Contacts Batch Operations

## Step 1 — Scope & resource model

<!-- Which batch methods, and why those? What does one item in a
     batch request contain? Put numbers on import sizes and payload
     budgets — they drive the limits in step 3. -->

## Step 2 — Methods & wire shapes

<!-- The actual shapes: batchCreate/batchUpdate/batchGet/batchDelete
     requests and responses, mirroring the singular methods. HTTP
     mapping. One worked example: a 3-item batchCreate request and
     its response. -->

## Step 3 — The hard part: atomicity, errors, limits

<!-- Per method: all-or-nothing or partial success? Say why. How do
     per-item errors come back, and how does a client match them to
     inputs? What are the limits (items, bytes) and what does
     exceeding them return? What happens when a half-succeeded
     batch is retried — how do you avoid duplicate contacts? -->

## Step 4 — Evolution & operations

<!-- Batch vs N singular calls: quantify the win. Where does batch
     stop being enough (100k? 1M?) and what takes over? How do the
     embedded singular requests evolve without breaking batchers? -->
