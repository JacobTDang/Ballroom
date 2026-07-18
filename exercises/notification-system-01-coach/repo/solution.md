# My Design: Notification System

## Step 1 — Use cases, constraints, estimates

<!-- What's the entry point (one internal send API), and what's out of
     scope (in-app inbox, campaign tooling)? Estimate: notifications/day
     -> peak sends/sec, and what differs per channel (push/SMS/email)
     in cost and latency. -->

## Step 2 — High-level design

<!-- Boxes and arrows: ingest API -> preference/opt-out check -> queue
     -> channel workers -> third-party providers. Why async end to
     end? -->

## Step 3 — Core components

<!-- Go deep where this question hinges:
     - Preferences & capping: where are opt-outs and per-channel
       preferences checked? What stops one user from being spammed?
     - Reliability: how does one triggering event avoid sending twice
       when a retry happens? What happens to a notification that keeps
       failing?
     - Prioritization: how does a one-time code get delivered ahead of
       a digest email queued behind a million others? -->

## Step 4 — Scale it

<!-- Where does this break at 10x? Do channels scale independently of
     each other? What happens when one provider (e.g. the SMS gateway)
     is slow or down -- can it stall push and email too? -->
