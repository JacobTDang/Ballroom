# My Design: Resumable Upload API

## Step 1 — Scope & resource model

<!-- The nouns: media resources, upload sessions. Session ownership,
     lifetime, and auth. Client mix and failure assumptions — put
     numbers on sizes and drop rates; they decide the thresholds. -->

## Step 2 — Methods & wire shapes

<!-- Initiation: what starts a session and what comes back? The
     chunk upload call with its headers (ranges), the status query,
     the small-file single-shot path. One worked example of a chunk
     round-trip. -->

## Step 3 — The hard part: resume & integrity

<!-- The offset protocol: who is authoritative about committed
     bytes? What exactly does a client do after a mid-chunk drop?
     Out-of-order and overlapping chunks. Checksums: per chunk,
     final, or both — and what a mismatch does to the session. -->

## Step 4 — Evolution & operations

<!-- Finalization and metadata: when does the durable resource
     exist? Abandoned-session GC. The retry matrix: which errors
     retry the chunk, which restart the session. Chunk-size
     trade-offs at your numbers. -->
