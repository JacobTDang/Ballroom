#include <functional>

// Retry calls op until it succeeds (returns true), backing off
// exponentially (base * 2^attempt, capped) between tries; returns
// false after max_attempts failures.
//
// TODO: this retries with a FIXED delay every time -- no exponential
// growth, no cap, and it even sleeps after the final failure.
bool Retry(std::function<bool()> op, int max_attempts, long base_ms,
           long cap_ms, std::function<void(long)> sleep) {
    for (int i = 0; i < max_attempts; i++) {
        if (op()) return true;
        sleep(base_ms);
    }
    return false;
}
