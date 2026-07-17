#include "solution.cpp"

#include <algorithm>
#include <chrono>
#include <cstdio>

int main() {
    for (int run = 0; run < 3; run++) {
        std::vector<int> inputs(600);
        for (int i = 0; i < 600; i++) inputs[i] = i % 250;

        // Deterministic per-item jitter -- a shared RNG called from the
        // worker threads would be the test racing with itself.
        auto got = FanOutIn(inputs, 8, [](int v) {
            std::this_thread::sleep_for(std::chrono::microseconds((v * 37) % 2000));
            return v * 3;
        });

        if (got.size() != inputs.size()) {
            fprintf(stderr, "run %d: got %zu results, want %zu -- fan-in dropped work\n", run, got.size(), inputs.size());
            return 1;
        }
        std::vector<int> want;
        for (int v : inputs) want.push_back(v * 3);
        std::sort(want.begin(), want.end());
        std::sort(got.begin(), got.end());
        if (got != want) {
            fprintf(stderr, "run %d: result multiset differs\n", run);
            return 1;
        }
    }
    printf("all assertions passed\n");
    return 0;
}
