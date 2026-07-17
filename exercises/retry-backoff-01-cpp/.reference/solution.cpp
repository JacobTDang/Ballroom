#include <functional>

// delay_i = min(cap, base << i) -- doubling until the cap flattens it.
// Sleeps only BETWEEN attempts, never after the final failure.
bool Retry(std::function<bool()> op, int max_attempts, long base_ms,
           long cap_ms, std::function<void(long)> sleep) {
    for (int attempt = 0; attempt < max_attempts; attempt++) {
        if (op()) return true;
        if (attempt == max_attempts - 1) break; // out of budget
        long delay = base_ms << attempt;
        if (delay > cap_ms || delay <= 0) delay = cap_ms; // <=0 guards overflow
        sleep(delay);
    }
    return false;
}
