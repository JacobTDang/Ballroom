package main

// Retry: delay_i = min(cap, base << i) -- doubling until the cap
// flattens it. The sleep happens only BETWEEN attempts (never after
// the last failure), and the final error is the operation's own, not
// a wrapper that hides it.
func Retry(op func() error, maxAttempts int, baseMillis, capMillis int64, sleep func(int64)) error {
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		err = op()
		if err == nil {
			return nil
		}
		if attempt == maxAttempts-1 {
			break // out of budget: no pointless sleep after the last try
		}
		delay := baseMillis << uint(attempt)
		if delay > capMillis || delay <= 0 { // <= 0 guards shift overflow
			delay = capMillis
		}
		sleep(delay)
	}
	return err
}
