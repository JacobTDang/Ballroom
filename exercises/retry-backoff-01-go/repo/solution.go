package main

// Retry calls op until it succeeds, backing off exponentially (base *
// 2^attempt, capped), giving up after maxAttempts with the last error.
//
// TODO: this retries with a FIXED delay every time -- no exponential
// growth, no cap logic, and it even sleeps after the final failure.
func Retry(op func() error, maxAttempts int, baseMillis, capMillis int64, sleep func(int64)) error {
	var err error
	for i := 0; i < maxAttempts; i++ {
		err = op()
		if err == nil {
			return nil
		}
		sleep(baseMillis)
	}
	return err
}
