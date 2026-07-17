def retry(op, max_attempts, base_ms, cap_ms, sleep):
    """delay_i = min(cap, base * 2^i) -- doubling until the cap
    flattens it. Sleeps only BETWEEN attempts (never after the final
    failure); the re-raised exception is the operation's own last."""
    last = None
    for attempt in range(max_attempts):
        try:
            return op()
        except Exception as e:  # noqa: BLE001 -- retrying anything is the point here
            last = e
            if attempt == max_attempts - 1:
                break  # out of budget: no pointless final sleep
            sleep(min(cap_ms, base_ms * (2 ** attempt)))
    raise last
