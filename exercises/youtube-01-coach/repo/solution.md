# My Design: YouTube

## Step 1 — Use cases, constraints, estimates

<!-- Which use cases are in scope (upload, watch, view counts)? What's
     out (recommendations, comments, live, ads)? Estimate: hours
     uploaded/minute -> storage/day across encodings, and whether
     storage or egress bandwidth dominates the cost -- show the
     arithmetic, don't just assert it. -->

## Step 2 — High-level design

<!-- Draw the upload pipeline (upload service -> queue -> transcoding
     workers -> object storage) separately from the watch path (CDN ->
     player). Where does the metadata DB sit relative to both? -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Transcoding: how does one upload become multiple
       resolutions/bitrates for adaptive streaming? What's the video's
       state machine from upload to watchable?
     - Serving: why put a CDN in front of object storage here
       specifically -- what's the traffic-shape argument for it?
     - View counts: how are they incremented without a synchronous
       write on every single view? -->

## Step 4 — Scale it

<!-- Where does this break at 10x? How does the transcoding fleet scale
     with upload volume? What's cached versus tiered to cold storage
     for old, rarely-watched videos? -->
