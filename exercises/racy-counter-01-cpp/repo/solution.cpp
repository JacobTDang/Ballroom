#include <thread>
#include <vector>

// Counter increments a shared counter n times, once per thread, and
// returns the final count. It should always return exactly n.
int Counter(int n) {
    int count = 0;
    std::vector<std::thread> threads;
    for (int i = 0; i < n; i++) {
        threads.emplace_back([&count]() {
            count++;
        });
    }
    for (auto& t : threads) t.join();
    return count;
}
