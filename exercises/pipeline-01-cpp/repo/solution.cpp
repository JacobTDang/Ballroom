#include <atomic>
#include <functional>
#include <mutex>
#include <thread>
#include <vector>

// FanOutIn fans inputs out to `workers` threads running stage and
// collects every result (order doesn't matter).
//
// TODO: this version detaches the workers and returns immediately --
// results go missing, and the vector races with the return.
std::vector<int> FanOutIn(const std::vector<int>& inputs, int workers,
                          std::function<int(int)> stage) {
    static std::vector<int> results;
    results.clear();
    static std::mutex mu;
    static std::atomic<size_t> next{0};
    next = 0;

    for (int w = 0; w < workers; w++) {
        std::thread([&] {
            while (true) {
                size_t i = next.fetch_add(1);
                if (i >= inputs.size()) return;
                int v = stage(inputs[i]);
                std::lock_guard<std::mutex> g(mu);
                results.push_back(v);
            }
        }).detach();
    }
    return results;
}
