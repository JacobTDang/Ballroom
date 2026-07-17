# Worker Pool

Run a function over a batch of jobs with a fixed number of concurrent
workers, returning the results **in input order**.

The starter processes jobs one at a time — correct output, zero
parallelism. Your job is to make N workers share the work without
scrambling the result order or racing on the results collection.

## The invariant the tests enforce

- `results[i]` corresponds to `jobs[i]` — always.
- With 8 workers and slow jobs, more than one job is genuinely in
  flight at once (the tests measure the high-water mark).
- Never more than `workers` jobs in flight.

API: `ProcessAll(jobs []int, workers int, fn func(int) int) []int`. Tests run with `-race`.

Think: which index does each worker write, and why does that make the
results collection race-free?
