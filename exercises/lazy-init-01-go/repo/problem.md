# Once-Only Lazy Initialization

An expensive initializer must run **exactly once**, no matter how many
threads ask first — and every caller gets its result.

The starter's `if not initialized: initialize()` is the textbook
broken double-init: two threads pass the check together and the
expensive work runs twice (or worse, a caller sees a half-built value).

## The invariant the tests enforce

With many concurrent `Get()` calls against a deliberately slow
initializer, the initializer runs exactly once and every caller
returns its value.

API: `NewLazy(init func() int) *Lazy`, `(*Lazy).Get() int`. Tests run with `-race`.

Think: your language ships a primitive built for exactly this — or do
it by hand with the double-checked pattern, done right.
