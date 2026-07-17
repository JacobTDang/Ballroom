# Fan-Out / Fan-In Pipeline

Fan a batch of inputs out to parallel workers running a stage function,
and fan the results back in to a single collection. Order does **not**
matter — completeness does.

The starter launches the workers and returns without waiting for the
fan-in to finish: results go missing, and the collection itself races.

## The invariant the tests enforce

The output is exactly the multiset of `stage(x)` for every input —
nothing dropped, nothing duplicated — across repeated runs with slow,
jittery stage functions.

API: `std::vector<int> FanOutIn(const std::vector<int>& inputs, int workers, std::function<int(int)> stage)` (any order). Compiled with `-fsanitize=thread`.

Think: who knows when all workers are done, and who is allowed to close
or finalize the output?
