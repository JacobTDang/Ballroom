#include <mutex>

// TokenBucket: one mutex makes check-and-take a single atomic step --
// the whole fix. Refill clamps under the same lock.
class TokenBucket {
public:
    explicit TokenBucket(int capacity) : capacity_(capacity), tokens_(capacity) {}

    bool Allow() {
        std::lock_guard<std::mutex> g(mu_);
        if (tokens_ > 0) {
            tokens_--;
            return true;
        }
        return false;
    }

    void Refill(int n) {
        std::lock_guard<std::mutex> g(mu_);
        tokens_ += n;
        if (tokens_ > capacity_) tokens_ = capacity_;
    }

private:
    int capacity_;
    int tokens_;
    std::mutex mu_;
};
