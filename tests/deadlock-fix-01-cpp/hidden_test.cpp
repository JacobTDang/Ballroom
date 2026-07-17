#include "solution.cpp"

#include <atomic>
#include <chrono>
#include <cstdio>
#include <cstdlib>
#include <thread>
#include <vector>

int main() {
    Account a(1, 10000), b(2, 10000);
    std::atomic<bool> finished{false};

    std::thread watchdog([&finished] {
        for (int i = 0; i < 1000 && !finished; i++)
            std::this_thread::sleep_for(std::chrono::milliseconds(10));
        if (!finished) {
            fprintf(stderr, "crossed transfers deadlocked\n");
            _Exit(1);
        }
    });

    std::vector<std::thread> threads;
    for (int i = 0; i < 4; i++) {
        threads.emplace_back([&, i] {
            for (int j = 0; j < 50; j++) {
                if (i % 2 == 0) Transfer(a, b, 1);
                else Transfer(b, a, 1);
            }
        });
    }
    for (auto& t : threads) t.join();
    finished = true;
    watchdog.join();

    if (a.Balance() + b.Balance() != 20000) {
        fprintf(stderr, "total balance %d, want 20000 conserved\n", a.Balance() + b.Balance());
        return 1;
    }

    Account c(3, 5), d(4, 0);
    if (Transfer(c, d, 10)) {
        fprintf(stderr, "Transfer succeeded with insufficient funds\n");
        return 1;
    }
    if (c.Balance() != 5 || d.Balance() != 0) {
        fprintf(stderr, "failed transfer moved money\n");
        return 1;
    }

    printf("all assertions passed\n");
    return 0;
}
