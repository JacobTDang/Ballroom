# Event Emitter

The pub/sub primitive under every UI framework and plugin system:
`On(event, handler)` subscribes (returning an id), `Once` subscribes
for a single delivery, `Off(id)` unsubscribes, `Emit(event, value)`
calls the event's handlers **in registration order**.

The starter can add handlers and that's it — no ids, no removal, and
`Once` never unhooks itself.

## The invariant the tests enforce

- Handlers fire in registration order, only for their own event.
- A `Once` handler fires exactly once, on the first emit.
- `Off(id)` removes exactly that subscription; an id removed **during
  an emit** (by an earlier handler in the same emit) does not fire.
- Emitting an event with no handlers is a no-op, not an error.

API: `NewEmitter() *Emitter`, `On(event string, fn func(int)) int`, `Once(event string, fn func(int)) int`, `Off(id int)`, `Emit(event string, v int)`.

Think: what does an id have to identify, and what does Emit iterate
over while handlers are mutating the subscription list?
