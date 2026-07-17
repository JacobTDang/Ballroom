#include "solution.cpp"

#include <atomic>
#include <chrono>
#include <cstdio>
#include <cstdlib>
#include <thread>
#include <vector>

int main() {
    const int n = 4, rounds = 5;
    Barrier b(n);
    std::atomic<int> arrivals[rounds];
    for (int r = 0; r < rounds; r++) arrivals[r] = 0;
    std::atomic<bool> failed{false};
    std::atomic<bool> finished{false};

    std::thread watchdog([&finished] {
        for (int i = 0; i < 1000 && !finished; i++)
            std::this_thread::sleep_for(std::chrono::milliseconds(10));
        if (!finished) {
            fprintf(stderr, "barrier deadlocked (participants never all released)\n");
            _Exit(1);
        }
    });

    std::vector<std::thread> threads;
    for (int p = 0; p < n; p++) {
        threads.emplace_back([&, p] {
            for (int r = 0; r < rounds; r++) {
                std::this_thread::sleep_for(std::chrono::milliseconds(p * 3));
                arrivals[r].fetch_add(1);
                b.Wait();
                if (arrivals[r].load() != n) {
                    failed = true;
                    return;
                }
            }
        });
    }
    for (auto& t : threads) t.join();
    finished = true;
    watchdog.join();

    if (failed) {
        fprintf(stderr, "a participant proceeded before all arrived\n");
        return 1;
    }
    printf("all assertions passed\n");
    return 0;
}
