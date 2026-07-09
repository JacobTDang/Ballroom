#pragma once

#include <chrono>

class RateLimiter {
public:
    RateLimiter(int limit, std::chrono::milliseconds window)
        : limit_(limit), window_(window) {}

    // Return true if a new request should be allowed right now.
    bool allow() {
        auto now = std::chrono::steady_clock::now();
        if (window_start_.time_since_epoch().count() == 0 ||
            now - window_start_ >= window_) {
            window_start_ = now;
            count_ = 0;
        }
        if (count_ >= limit_) {
            return false;
        }
        count_++;
        return true;
    }

private:
    int limit_;
    std::chrono::milliseconds window_;
    std::chrono::steady_clock::time_point window_start_{};
    int count_ = 0;
};
