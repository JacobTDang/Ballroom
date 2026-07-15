# Reference Design — Scaling to Millions of Users on AWS

The deliverable is the progression: each addition justified by the
bottleneck that forced it.

## The progression
1. **One box** (web + DB together, EC2): fine to ~hundreds of users.
   Breaks: contention for memory/CPU between web and DB.
2. **Split web / DB** (RDS or DB on its own instance): each scales and
   is tuned independently.
3. **LB + second web server** (ELB, 2× EC2): removes the single point
   of failure and doubles capacity. Web tier must go **stateless**
   (sessions out to a shared store) — the prerequisite for everything
   after.
4. **Static assets → S3 + CloudFront**: bandwidth off the web tier;
   images/js/css served from the edge.
5. **DB read replicas** (RDS read replicas): read load fans out;
   replication lag consequence stated.
6. **Cache layer** (ElastiCache): hot queries and sessions;
   cache-aside; relieves replicas.
7. **Autoscaling** (ASG on the web tier): diurnal peaks without paying
   for the peak all day. Needs monitoring to drive it.
8. **Async workers** (SQS + worker fleet): slow work (emails, reports,
   thumbnails) off the request path.
9. **Write ceiling** (last): federation by function, then shard by
   user, or move hot tables to NoSQL (DynamoDB) — chosen by the write
   numbers, not fashion.

## Keep it running
CloudWatch metrics/alarms are how the *next* bottleneck is found;
automated backups + multi-AZ failover; at 100M+: multi-region,
cell-based isolation.
