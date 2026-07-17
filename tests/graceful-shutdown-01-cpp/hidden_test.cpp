#include "solution.cpp"

#include <atomic>
#include <chrono>
#include <cstdio>
#include <cstdlib>

int main() {
    // 1. Stop drains everything accepted.
    {
        std::atomic<int> handled{0};
        Server s(4, [&handled](int) {
            std::this_thread::sleep_for(std::chrono::milliseconds(2));
            handled.fetch_add(1);
        });
        int accepted = 0;
        for (int i = 0; i < 200; i++) {
            if (s.Submit(i)) accepted++;
        }

        std::atomic<bool> stopped{false};
        std::thread stopper([&] { s.Stop(); stopped = true; });
        for (int i = 0; i < 1000 && !stopped; i++)
            std::this_thread::sleep_for(std::chrono::milliseconds(10));
        if (!stopped) {
            fprintf(stderr, "Stop never returned -- deadlock or never drained\n");
            _Exit(1);
        }
        stopper.join();

        if (handled != accepted) {
            fprintf(stderr, "%d handled after Stop, want every accepted job (%d)\n", handled.load(), accepted);
            return 1;
        }
    }

    // 2. Submit refused after Stop; nothing new handled.
    {
        std::atomic<int> handled{0};
        Server s(2, [&handled](int) { handled.fetch_add(1); });
        s.Submit(1);
        s.Stop();
        int before = handled.load();
        if (s.Submit(2)) {
            fprintf(stderr, "Submit accepted a job after Stop returned\n");
            return 1;
        }
        std::this_thread::sleep_for(std::chrono::milliseconds(20));
        if (handled != before) {
            fprintf(stderr, "handled grew from %d to %d after Stop\n", before, handled.load());
            return 1;
        }
    }

    printf("all assertions passed\n");
    return 0;
}
