#include <functional>
#include <thread>
#include <vector>

// RunLimited runs every task, with at most `limit` executing
// concurrently.
//
// TODO: this version launches everything at once -- the limit is
// ignored entirely.
void RunLimited(const std::vector<std::function<void()>>& tasks, int limit) {
    std::vector<std::thread> threads;
    for (const auto& task : tasks) {
        threads.emplace_back(task);
    }
    for (auto& t : threads) t.join();
}
