# My Design: Dropbox

## Step 1 — Use cases, constraints, estimates

<!-- Which use cases are in scope (upload/download, multi-device sync,
     sharing)? What's explicitly out? Estimate: users x average
     storage, daily active uploaders, upload bandwidth, and where dedup
     saves real space. Show the arithmetic. -->

## Step 2 — High-level design

<!-- Boxes and arrows: metadata service (file tree, versions, block
     lists) vs block storage (the bytes, in object storage). How do
     clients talk to each? -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Chunking & dedup: how is a file split into blocks, and how is a
       block identified? Why does that make delta sync and cross-user
       dedup possible?
     - Sync protocol: how does a device notice a local change, and how
       do other devices find out? Narrate one edit end to end.
     - Conflicts: two devices edit the same file while offline — what
       happens when both come back online? -->

## Step 4 — Scale it

<!-- Where does this break at 10x the load from step 1? How is the
     metadata DB sharded? What gets cached vs tiered to cold storage?
     Does notification fanout scale differently from raw data
     transfer? -->
