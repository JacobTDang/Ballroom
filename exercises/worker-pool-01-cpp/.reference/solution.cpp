#include <atomic>
#include <functional>
#include <thread>
#include <vector>

// ProcessAll: an atomic index cursor feeds `workers` threads; each
// writes results[i] for the i it claimed -- distinct slots, no lock on
// the results. Position, not completion time, decides placement, so
// input order is free.
std::vector<int> ProcessAll(const std::vector<int>& jobs, int workers,
                            std::function<int(int)> fn) {
    std::vector<int> results(jobs.size());
    std::atomic<size_t> next{0};

    std::vector<std::thread> threads;
    for (int w = 0; w < workers; w++) {
        threads.emplace_back([&] {
            while (true) {
                size_t i = next.fetch_add(1);
                if (i >= jobs.size()) return;
                results[i] = fn(jobs[i]);
            }
        });
    }
    for (auto& t : threads) t.join();
    return results;
}
