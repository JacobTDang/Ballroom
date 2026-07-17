#include "solution.cpp"

#include <atomic>
#include <cstdio>
#include <thread>
#include <vector>

static int hammer(TokenBucket& b, int callers) {
    std::atomic<int> allowed{0};
    std::vector<std::thread> threads;
    for (int i = 0; i < callers; i++) {
        threads.emplace_back([&] {
            if (b.Allow()) allowed.fetch_add(1);
        });
    }
    for (auto& t : threads) t.join();
    return allowed.load();
}

int main() {
    {
        TokenBucket b(100);
        int got = hammer(b, 300);
        if (got != 100) {
            fprintf(stderr, "%d of 300 concurrent Allow calls succeeded, want exactly 100\n", got);
            return 1;
        }
    }
    {
        TokenBucket b(100);
        hammer(b, 300);
        b.Refill(40);
        int got = hammer(b, 200);
        if (got != 40) {
            fprintf(stderr, "%d allowed after Refill(40), want exactly 40\n", got);
            return 1;
        }
    }
    {
        TokenBucket b(50);
        b.Refill(1000);
        int got = hammer(b, 200);
        if (got != 50) {
            fprintf(stderr, "%d allowed after over-refill, want capacity 50\n", got);
            return 1;
        }
    }
    printf("all assertions passed\n");
    return 0;
}
