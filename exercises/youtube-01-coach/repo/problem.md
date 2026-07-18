# Design YouTube

Design a video-sharing service like YouTube: users upload videos,
anyone can watch them.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Core use cases: upload a video, watch a video, view counts. What's
  out of scope — recommendations, comments, live streaming, ads?
- A watch has to work well on a phone on a bad connection and a TV on
  a great one. What does that imply about how a single upload turns
  into what actually gets served?
- Put numbers on it: hours of video uploaded per minute, the storage
  that implies once you account for multiple encodings, and which of
  storage or bandwidth actually dominates the cost.

## Suggested defaults (if you want a starting point)

- ~500 hours of video uploaded per minute
- Each uploaded video is transcoded into several resolutions/bitrates
  for adaptive streaming, not served as the raw upload
- Billions of watch views per day
- Assume egress bandwidth, not storage, is the number that should worry
  you — confirm that with arithmetic rather than taking it on faith

## What good looks like

By the end you should have: stated assumptions with the storage and
bandwidth arithmetic (and the explicit comparison between them), a
high-level diagram separating the upload pipeline (upload → queue →
transcoding workers → object storage) from the watch path (CDN →
player), an async transcoding design producing multiple
resolutions/bitrates with a clear video-state machine, a serving story
that puts a CDN in front of object storage and explains why given the
popularity skew of videos, and a non-synchronous plan for view counts
at this write volume.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
