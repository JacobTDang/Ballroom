def retry(op, max_attempts, base_ms, cap_ms, sleep):
    """Call op until it succeeds (returns), backing off exponentially
    (base * 2^attempt, capped) between tries; after max_attempts
    failures, re-raise the last exception.

    TODO: this retries with a FIXED delay every time -- no exponential
    growth, no cap, and it even sleeps after the final failure.
    """
    last = None
    for _ in range(max_attempts):
        try:
            return op()
        except Exception as e:  # noqa: BLE001 -- retrying anything is the point here
            last = e
            sleep(base_ms)
    raise last
