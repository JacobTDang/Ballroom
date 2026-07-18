# Design Instagram

Design a photo-sharing service like Instagram: users post photos and
scroll a feed of accounts they follow.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Core use cases: post a photo, view the home feed, like a post.
  What's out of scope — stories, reels, DMs, an explore/discovery
  page?
- The defining tension: most accounts have a normal-sized follower
  list, but a celebrity account has millions. How should that shape
  how a feed gets built?
- Put numbers on it: photos posted/day, storage per photo across the
  renditions you'd actually serve, feed reads/day, and the read:write
  ratio.

## Suggested defaults (if you want a starting point)

- 500 million daily active users, 100 million photos posted per day
- Feed reads massively outnumber posts — assume roughly 1,000:1
- Each photo is transcoded into a handful of renditions (thumbnail,
  feed-sized, full) rather than served as the original upload
- Average follower count is a few hundred; a small fraction of accounts
  have millions

## What good looks like

By the end you should have: stated assumptions with the storage and
feed-read arithmetic, a high-level diagram separating the media path
(upload → object storage + CDN) from the metadata/feed path (post
metadata, follow graph, feed service), a concrete fanout decision for
building the feed (including a deliberate answer for celebrity
accounts), where materialized feeds live and how they're paginated,
and a non-synchronous plan for like/comment counters at this write
volume.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
