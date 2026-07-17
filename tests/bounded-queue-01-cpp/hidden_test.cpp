#include "solution.cpp"

#include <atomic>
#include <chrono>
#include <cstdio>
#include <map>
#include <mutex>
#include <thread>
#include <vector>

int main() {
    // 1. Every item arrives exactly once under contention.
    {
        BoundedQueue q(4);
        const int producers = 3, per = 200, consumers = 3;
        const int total = producers * per;
        std::vector<std::thread> threads;
        for (int p = 0; p < producers; p++) {
            threads.emplace_back([&q, p, per] {
                for (int i = 0; i < per; i++) q.Put(p * per + i);
            });
        }
        std::mutex mu;
        std::map<int, int> seen;
        for (int c = 0; c < consumers; c++) {
            threads.emplace_back([&] {
                for (int i = 0; i < total / consumers; i++) {
                    int v = q.Get();
                    std::lock_guard<std::mutex> g(mu);
                    seen[v]++;
                }
            });
        }
        for (auto& t : threads) t.join();
        if ((int)seen.size() != total) {
            fprintf(stderr, "saw %zu distinct items, want %d\n", seen.size(), total);
            return 1;
        }
        for (auto& [v, n] : seen) {
            if (n != 1) {
                fprintf(stderr, "item %d consumed %d times, want once\n", v, n);
                return 1;
            }
        }
    }

    // 2. Get blocks until a Put arrives.
    {
        BoundedQueue q(2);
        std::atomic<bool> done{false};
        std::atomic<int> got{-1};
        std::thread t([&] { got = q.Get(); done = true; });
        std::this_thread::sleep_for(std::chrono::milliseconds(50));
        if (done) {
            fprintf(stderr, "Get on an empty queue returned immediately (%d), want it to block\n", got.load());
            t.detach();
            return 1;
        }
        q.Put(7);
        for (int i = 0; i < 100 && !done; i++) std::this_thread::sleep_for(std::chrono::milliseconds(10));
        if (!done || got != 7) {
            fprintf(stderr, "Get never returned 7 after Put(7)\n");
            t.detach();
            return 1;
        }
        t.join();
    }

    // 3. Put blocks at capacity until a Get frees a slot.
    {
        BoundedQueue q(2);
        q.Put(1);
        q.Put(2);
        std::atomic<bool> done{false};
        std::thread t([&] { q.Put(3); done = true; });
        std::this_thread::sleep_for(std::chrono::milliseconds(50));
        if (done) {
            fprintf(stderr, "Put on a full queue returned immediately, want it to block\n");
            t.join();
            return 1;
        }
        if (int v = q.Get(); v != 1) {
            fprintf(stderr, "Get = %d, want FIFO order (1 first)\n", v);
            t.detach();
            return 1;
        }
        for (int i = 0; i < 100 && !done; i++) std::this_thread::sleep_for(std::chrono::milliseconds(10));
        if (!done) {
            fprintf(stderr, "blocked Put never completed after a Get\n");
            t.detach();
            return 1;
        }
        t.join();
    }

    printf("all assertions passed\n");
    return 0;
}
