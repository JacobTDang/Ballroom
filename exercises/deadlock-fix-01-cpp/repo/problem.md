# Fix the Deadlock

Debug-style: the code you're given deadlocks. Two bank accounts, a
`Transfer(from, to, amount)` that locks both — and two threads
transferring in opposite directions each grab one lock and wait forever
for the other. The classic lock-ordering inversion.

Make the tests pass **without** serializing everything through one
global lock for the whole transfer table — fix the ordering instead.

## The invariant the tests enforce

- Crossed concurrent transfers complete (no deadlock — the tests carry
  a watchdog).
- Money is conserved: the total across accounts never changes.
- A transfer with insufficient funds returns false and moves nothing;
  balances never go negative.

API: `class Account { Account(int id, int balance); int Balance(); }` and `bool Transfer(Account& from, Account& to, int amount)` (public `id` and `mu`). `std::lock`/`std::scoped_lock` count as fixing the ordering — that's exactly what they exist for. Compiled with `-fsanitize=thread`, which also reports lock-order inversions.

Think: deadlock needs a cycle. What total order on the locks makes a
cycle impossible?
