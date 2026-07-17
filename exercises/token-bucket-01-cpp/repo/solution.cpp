#include <chrono>
#include <thread>

// TokenBucket: a rate limiter shared by many threads. Allow takes a
// token if available; Refill (called by an external ticker) adds
// tokens, clamped at capacity.
//
// TODO: check-then-decrement below is two separate steps -- under
// contention this hands out more tokens than exist.
class TokenBucket {
public:
    explicit TokenBucket(int capacity) : capacity_(capacity), tokens_(capacity) {}

    bool Allow() {
        if (tokens_ > 0) {
            std::this_thread::sleep_for(std::chrono::microseconds(1));
            tokens_--;
            return true;
        }
        return false;
    }

    void Refill(int n) {
        tokens_ += n;
        if (tokens_ > capacity_) tokens_ = capacity_;
    }

private:
    int capacity_;
    int tokens_;
};
