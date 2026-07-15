# Design Mint.com

Design a personal-finance service: users connect their bank accounts,
the service pulls in transactions, categorizes them, and alerts users
when they blow past a budget.

Work through it with your coach using the 4-step method, writing each
step into `solution.md` as you go.

## Scope to establish in step 1

- Core use cases: connect a financial account, pull transactions,
  categorize each transaction, monthly spending overview by category,
  budget alerts. Out of scope? (Bill pay, credit scores, investments.)
- The shape of the workload: is this read-heavy or write-heavy? How
  fresh do transactions need to be?
- Put numbers on it: users, accounts per user, transactions per month,
  bytes per transaction, storage over 5 years.

## Suggested defaults (if you want a starting point)

- 10 million users, ~3 connected accounts each
- ~30 transactions per account per month; data pulled from financial
  institutions on a daily batch cadence (no real-time requirement)
- Users check the app a few times a month — write-heavy, unlike most
  consumer apps
- Budgets: default per-category budgets, user can override; alert when
  a category crosses its budget

## What good looks like

By the end you should have: estimates that recognize the write-heavy,
batch-friendly shape, a high-level design with the transaction
extraction pipeline separate from the user-facing read path, a
categorization approach (rules first, with overrides), a budget/alert
mechanism that doesn't recompute everything on every read, and a
scaling story built on async processing and precomputed monthly
analyses.
