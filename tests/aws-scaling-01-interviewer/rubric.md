# Grading Rubric — Scaling to Millions of Users on AWS

Grade each dimension: strong / adequate / missing. A passing design is
adequate-or-better on every dimension with at least two strongs.

## 1. Use cases & constraints

- The app and its read/write ratio pinned down; user-growth stages laid
  out as the structure of the answer.

## 2. Back-of-envelope estimates

- Requests/sec estimated at the key milestones and actually used to
  call the next bottleneck ("at ~N req/s one box saturates because…").

## 3. Bottleneck-driven narrative

- Every addition is preceded by the specific bottleneck that forced it
  — no cargo-culting components in before they're needed.
- The early steps are right: split web from DB, then LB + second web
  server (removing the single point of failure), before anything
  exotic.

## 4. Static content & network

- Static assets moved to object storage with a CDN in front, justified
  by bandwidth/latency rather than fashion.

## 5. Database scaling sequence

- Read replicas for read load, then caching (hot queries, sessions
  externalized so web tier stays stateless), and only at the write
  ceiling: federation/sharding or NoSQL — in an order the numbers
  justify.

## 6. Elasticity & async work

- Autoscaling tied to diurnal/peak patterns; slow work moved off the
  request path into queues + worker pools with an example.

## 7. Operations

- Monitoring/alerting framed as how you find the next bottleneck;
  backups and multi-AZ failover present; a note on what 100M users
  would force next.

## 8. Communication & trade-offs

- Each stage's trade-off stated (cost, complexity, consistency);
  AWS/generic naming used consistently; the candidate drove the
  staged structure.
