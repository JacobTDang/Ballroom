# Scale a System to Millions of Users on AWS

A simple web app — users, some data, a dashboard — starts on a single
box and grows to millions of users. Tell the scaling story: at each
stage, name the bottleneck, then fix exactly that with the right
building block.

This one is a narrative, not a single snapshot: the deliverable is the
progression, with each step justified by the load that broke the
previous step.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go — here step 3 IS the progression.

## Scope to establish in step 1

- What does the app do? Keep it generic (read-mostly web app with a
  database) unless you want to pick something concrete.
- The stages to cover: 1 user → 10K → 100K → 1M → 10M+.
- Put numbers on it at each stage: requests/sec, DB reads vs writes,
  data size — enough to point at the next bottleneck.

## Suggested defaults (if you want a starting point)

- Read-heavy (~40:1), mostly-static assets plus a user dashboard
- Start: everything on one box (web server + database together)
- You may use AWS names (EC2, ELB, RDS, ElastiCache, S3, CloudFront,
  SQS, Auto Scaling) or generic names — the reasoning is what's graded

## What good looks like

By the end you should have: a staged narrative where every addition is
preceded by the bottleneck that forced it — separate web/DB boxes,
load balancer + horizontal web tier, static assets to object storage +
CDN, DB read replicas then caching, autoscaling, async workers via
queues, and finally sharding/federation or NoSQL for write scale —
plus monitoring as the thing that tells you which bottleneck is next.
