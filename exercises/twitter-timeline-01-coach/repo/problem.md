# Design the Twitter Timeline & Search

Design the core of Twitter: users post tweets, view a home timeline of
tweets from people they follow, and search across all tweets.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

- Core use cases: post a tweet, view home timeline, view a user
  timeline, search tweets. What's out of scope? (DMs, notifications,
  ads, trending.)
- The defining tension: most users follow a few hundred people, but a
  celebrity has 100M followers. How does that shape delivery?
- Put numbers on it: tweets/day, timeline reads/day, fanout writes per
  tweet, storage per tweet, search index size.

## Suggested defaults (if you want a starting point)

- 100 million active users, 500 million tweets per day
- Each user follows 200 people on average; reads dominate writes
  heavily on the timeline (~100:1)
- Timeline should feel instant (< 200 ms); search can be slower
- A tweet is ~140 bytes of text plus metadata; media lives elsewhere

## What good looks like

By the end you should have: stated assumptions with fanout arithmetic,
a high-level diagram separating the write path (tweet → fanout) from
the read path (timeline fetch), a concrete answer to the celebrity
fanout problem, the timeline storage choice (and why it's in memory),
how search is indexed and queried, and a scaling story for the
read-heavy load.
