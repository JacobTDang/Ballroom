#include "solution.cpp"

#include <atomic>
#include <chrono>
#include <cstdio>
#include <thread>
#include <vector>

int main() {
    std::atomic<int> init_calls{0};
    Lazy lazy([&init_calls] {
        init_calls.fetch_add(1);
        std::this_thread::sleep_for(std::chrono::milliseconds(10));
        return 42;
    });

    const int callers = 50;
    std::vector<int> results(callers, -1);
    std::vector<std::thread> threads;
    for (int i = 0; i < callers; i++) {
        threads.emplace_back([&, i] { results[i] = lazy.Get(); });
    }
    for (auto& t : threads) t.join();

    if (init_calls != 1) {
        fprintf(stderr, "init ran %d times under contention, want exactly once\n", init_calls.load());
        return 1;
    }
    for (int i = 0; i < callers; i++) {
        if (results[i] != 42) {
            fprintf(stderr, "caller %d got %d, want 42\n", i, results[i]);
            return 1;
        }
    }

    std::atomic<int> calls2{0};
    Lazy l2([&calls2] { calls2.fetch_add(1); return 7; });
    for (int i = 0; i < 5; i++) {
        if (l2.Get() != 7) {
            fprintf(stderr, "sequential Get wrong value\n");
            return 1;
        }
    }
    if (calls2 != 1) {
        fprintf(stderr, "init ran %d times sequentially, want once\n", calls2.load());
        return 1;
    }

    printf("all assertions passed\n");
    return 0;
}
