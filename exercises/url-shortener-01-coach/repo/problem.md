# Design Pastebin / Bit.ly

Design a URL-shortening service: a user submits a long URL and gets back
a short link; anyone who visits the short link is redirected to the
original URL.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

Talk through and pin down with your coach:

- Core use cases: shorten, redirect — what else is in scope? Expiration?
  Click analytics? Custom aliases? A user account system?
- What can you assume? Traffic volume, read/write ratio, URL length,
  how long links live.
- Put numbers on it: writes/sec, reads/sec, storage over 3 years,
  bandwidth.

## Suggested defaults (if you want a starting point)

- 100 million new links per month
- 10:1 read-to-write ratio
- Links expire after 3 years by default
- Analytics: click counts per link, viewable by the creator

## What good looks like

By the end you should have: stated assumptions with estimates, a
high-level diagram, the short-code generation scheme (and why it
doesn't collide), your data model and database choice with reasoning,
the redirect flow (including the HTTP status code choice), and a
scaling story for the read-heavy load.

After you submit with M-q, a distilled reference design
(`reference.md`) appears alongside the rubric — compare your design
against it before moving on.
