#include "solution.cpp"

#include <cstdio>
#include <vector>

int main() {
    // First-try success: one call, zero sleeps.
    {
        std::vector<long> slept;
        int calls = 0;
        bool ok = Retry([&] { calls++; return true; }, 5, 100, 1000,
                        [&](long ms) { slept.push_back(ms); });
        if (!ok || calls != 1 || !slept.empty()) {
            fprintf(stderr, "first-try: ok=%d calls=%d sleeps=%zu, want true/1/0\n", ok, calls, slept.size());
            return 1;
        }
    }
    // Exponential delays, exact.
    {
        std::vector<long> slept;
        int calls = 0;
        bool ok = Retry([&] { return ++calls > 3; }, 5, 100, 10000,
                        [&](long ms) { slept.push_back(ms); });
        std::vector<long> want{100, 200, 400};
        if (!ok || calls != 4 || slept != want) {
            fprintf(stderr, "exponential: ok=%d calls=%d, delays wrong\n", ok, calls);
            return 1;
        }
    }
    // Cap flattens.
    {
        std::vector<long> slept;
        int calls = 0;
        bool ok = Retry([&] { return ++calls > 4; }, 6, 100, 250,
                        [&](long ms) { slept.push_back(ms); });
        std::vector<long> want{100, 200, 250, 250};
        if (!ok || slept != want) {
            fprintf(stderr, "cap: delays wrong\n");
            return 1;
        }
    }
    // Exhaustion: false, exact attempts, no trailing sleep.
    {
        std::vector<long> slept;
        int calls = 0;
        bool ok = Retry([&] { calls++; return false; }, 3, 100, 1000,
                        [&](long ms) { slept.push_back(ms); });
        if (ok || calls != 3 || slept.size() != 2) {
            fprintf(stderr, "exhaustion: ok=%d calls=%d sleeps=%zu, want false/3/2\n", ok, calls, slept.size());
            return 1;
        }
    }
    printf("all assertions passed\n");
    return 0;
}
