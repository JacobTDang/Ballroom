#include "solution.cpp"

#include <atomic>
#include <chrono>
#include <cstdio>
#include <thread>

int main() {
    // 1. Results in input order.
    {
        std::vector<int> jobs(100);
        for (int i = 0; i < 100; i++) jobs[i] = i;
        auto got = ProcessAll(jobs, 8, [](int v) {
            std::this_thread::sleep_for(std::chrono::milliseconds(2));
            return v * 2;
        });
        if (got.size() != jobs.size()) {
            fprintf(stderr, "got %zu results, want %zu\n", got.size(), jobs.size());
            return 1;
        }
        for (int i = 0; i < 100; i++) {
            if (got[i] != i * 2) {
                fprintf(stderr, "results[%d] = %d, want %d -- input order required\n", i, got[i], i * 2);
                return 1;
            }
        }
    }

    // 2. Actually parallel, within the bound.
    {
        std::atomic<int> in_flight{0}, high_water{0};
        auto fn = [&](int v) {
            int cur = in_flight.fetch_add(1) + 1;
            int hw = high_water.load();
            while (cur > hw && !high_water.compare_exchange_weak(hw, cur)) {}
            std::this_thread::sleep_for(std::chrono::milliseconds(10));
            in_flight.fetch_sub(1);
            return v;
        };
        std::vector<int> jobs(64, 0);
        ProcessAll(jobs, 8, fn);
        if (high_water < 2) {
            fprintf(stderr, "high-water %d: never ran jobs in parallel\n", high_water.load());
            return 1;
        }
        if (high_water > 8) {
            fprintf(stderr, "high-water %d: exceeded the 8 workers requested\n", high_water.load());
            return 1;
        }
    }

    printf("all assertions passed\n");
    return 0;
}
