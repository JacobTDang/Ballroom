#include <atomic>
#include <thread>
#include <vector>

// Counter increments a shared counter n times, once per thread, and
// returns the final count. It always returns exactly n.
int Counter(int n) {
    std::atomic<int> count{0};
    std::vector<std::thread> threads;
    for (int i = 0; i < n; i++) {
        threads.emplace_back([&count]() {
            count.fetch_add(1, std::memory_order_relaxed);
        });
    }
    for (auto& t : threads) t.join();
    return count.load();
}
