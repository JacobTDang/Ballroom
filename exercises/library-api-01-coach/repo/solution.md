# My Design: Library API

## Step 1 — Scope & resource model

<!-- Which nouns become resources, and why? Write the resource names
     concretely (AIP-122 style): publishers/{publisher},
     publishers/{publisher}/books/{book}, ... Flat vs hierarchical —
     what did you choose for copies, members, loans, and what does
     each choice cost? -->

## Step 2 — Methods & wire shapes

<!-- The five standard methods (AIP-131..135) for your core resources:
     HTTP verb, path, request body, response body. Which are
     idempotent? What does Create return? Sketch List with
     page_size/page_token. Show one full request/response example. -->

## Step 3 — The hard part: checkout & return

<!-- Checkout and return don't fit CRUD. Custom methods
     (POST .../copies/{copy}:checkout, AIP-136)? A loans collection
     you Create/Delete in? Argue the trade-off and specify the one
     you chose: request shape, response, and every failure case
     (already checked out, member over limit, copy lost). -->

## Step 4 — Evolution & operations

<!-- What changes safely later (new fields, new methods)? What would
     break clients? Error consistency: which HTTP codes and error
     shapes do all endpoints share? Rate limits or quotas for the
     public catalog callers? -->
