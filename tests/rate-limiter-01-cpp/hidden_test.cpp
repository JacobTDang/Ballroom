#include <cassert>
#include <chrono>
#include <cstdio>
#include <thread>

#include "solution.hpp"

int main() {
    {
        RateLimiter rl(3, std::chrono::milliseconds(10000));
        assert(rl.allow());
        assert(rl.allow());
        assert(rl.allow());
        assert(!rl.allow());
    }
    {
        RateLimiter rl(1, std::chrono::milliseconds(50));
        assert(rl.allow());
        assert(!rl.allow());
        std::this_thread::sleep_for(std::chrono::milliseconds(60));
        assert(rl.allow());
    }
    printf("all assertions passed\n");
    return 0;
}
