#include <deque>

// SlidingWindow keeps the timestamps of allowed requests (a deque:
// they arrive in order, so eviction pops from the front) and evicts
// entries a full window old before deciding. Denied requests are never
// recorded.
class SlidingWindow {
public:
    SlidingWindow(int limit, long window_ms)
        : limit_(limit), window_ms_(window_ms) {}

    bool AllowAt(long now_ms) {
        long cutoff = now_ms - window_ms_;
        while (!allowed_.empty() && allowed_.front() <= cutoff) {
            allowed_.pop_front();
        }
        if ((int)allowed_.size() < limit_) {
            allowed_.push_back(now_ms);
            return true;
        }
        return false;
    }

private:
    int limit_;
    long window_ms_;
    std::deque<long> allowed_;
};
