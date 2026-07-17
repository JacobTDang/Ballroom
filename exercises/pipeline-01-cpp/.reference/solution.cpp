#include <atomic>
#include <functional>
#include <mutex>
#include <thread>
#include <vector>

// FanOutIn: plain locals (no static state leaking across calls), and
// the caller joins every worker before returning -- completeness is
// exactly "all workers finished before the results left the function".
std::vector<int> FanOutIn(const std::vector<int>& inputs, int workers,
                          std::function<int(int)> stage) {
    std::vector<int> results;
    std::mutex mu;
    std::atomic<size_t> next{0};

    std::vector<std::thread> threads;
    for (int w = 0; w < workers; w++) {
        threads.emplace_back([&] {
            while (true) {
                size_t i = next.fetch_add(1);
                if (i >= inputs.size()) return;
                int v = stage(inputs[i]);
                std::lock_guard<std::mutex> g(mu);
                results.push_back(v);
            }
        });
    }
    for (auto& t : threads) t.join();
    return results;
}
