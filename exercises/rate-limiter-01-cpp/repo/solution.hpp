#pragma once

#include <chrono>

class RateLimiter {
public:
    RateLimiter(int limit, std::chrono::milliseconds window)
        : limit_(limit), window_(window) {}

    // Return true if a new request should be allowed right now.
    bool allow() {
        // TODO: implement
        return false;
    }

private:
    int limit_;
    std::chrono::milliseconds window_;
    std::chrono::steady_clock::time_point window_start_{};
    int count_ = 0;
};
