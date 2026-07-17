# Bounded Blocking Queue

Implement a fixed-capacity FIFO queue shared between producer and
consumer threads:

- `Put(v)` adds an item, **blocking while the queue is full**.
- `Get()` removes and returns the oldest item, **blocking while the
  queue is empty**.

The starter code has neither bound nor blocking — it "works" until two
threads actually overlap.

## The invariant the tests enforce

Every item put is got exactly once (nothing lost, nothing duplicated,
under many producers and consumers at once), a `Get()` on an empty
queue waits for the next `Put`, and a `Put()` on a full queue waits for
the next `Get`.

API: `NewBoundedQueue(capacity int) *BoundedQueue`, `(*BoundedQueue).Put(v int)`, `(*BoundedQueue).Get() int`. Tests run with `-race`.

Think: what state is shared, and who needs to be woken when it changes?
