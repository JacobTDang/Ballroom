#include "solution.cpp"

#include <atomic>
#include <chrono>
#include <cstdio>

static std::atomic<int> in_flight{0}, high_water{0}, ran{0};

static void reset() { in_flight = 0; high_water = 0; ran = 0; }

static std::vector<std::function<void()>> instrumented(int n) {
    std::vector<std::function<void()>> tasks;
    for (int i = 0; i < n; i++) {
        tasks.push_back([] {
            int cur = in_flight.fetch_add(1) + 1;
            int hw = high_water.load();
            while (cur > hw && !high_water.compare_exchange_weak(hw, cur)) {}
            std::this_thread::sleep_for(std::chrono::milliseconds(15));
            in_flight.fetch_sub(1);
            ran.fetch_add(1);
        });
    }
    return tasks;
}

int main() {
    reset();
    RunLimited(instrumented(32), 4);
    if (ran != 32) {
        fprintf(stderr, "%d tasks ran, want 32\n", ran.load());
        return 1;
    }
    if (high_water > 4) {
        fprintf(stderr, "high-water %d exceeded limit 4\n", high_water.load());
        return 1;
    }
    if (high_water < 2) {
        fprintf(stderr, "high-water %d: no real parallelism under limit 4\n", high_water.load());
        return 1;
    }

    reset();
    RunLimited(instrumented(6), 1);
    if (ran != 6 || high_water != 1) {
        fprintf(stderr, "limit 1: ran=%d high_water=%d, want 6 and exactly 1\n", ran.load(), high_water.load());
        return 1;
    }

    reset();
    RunLimited(instrumented(3), 10);
    if (ran != 3) {
        fprintf(stderr, "limit>n: %d tasks ran, want 3\n", ran.load());
        return 1;
    }

    printf("all assertions passed\n");
    return 0;
}
