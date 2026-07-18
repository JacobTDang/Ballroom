# My Design: Settings API with Optimistic Concurrency

## Step 1 — Scope & resource model

<!-- The settings resource and its name. Who reads/writes, at what
     rates, with what gap between their read and their write? Write
     out the exact interleaving that loses an update today. -->

## Step 2 — Methods & wire shapes

<!-- GET and update shapes with the concurrency fields/headers in
     them. Where does the etag appear (header, resource field,
     both)? One worked example: read → concurrent write → stale
     write rejected, with real status codes. -->

## Step 3 — The hard part: validators & recovery

<!-- What is the etag derived from — version counter or content
     hash — and why? Strong vs weak validators: what does each
     guarantee, which do you need? The read side: If-None-Match /
     304. And the recovery protocol: after a 412, what exactly does
     the client do — re-GET and replay? merge? show a diff? -->

## Step 4 — Evolution & operations

<!-- Policy: are unconditional writes allowed at all? What do
     automations do differently from humans? Where would etags be
     overkill, and what does that boundary tell you? -->
