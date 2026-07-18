#include <cassert>
#include <chrono>
#include <cstdio>
#include <string>
#include <vector>

std::string build_log(const std::vector<std::string>& chunks);

int main() {
    // Tier 1: exact small-n correctness.
    assert(build_log({"c", "b", "a"}) == "abc");
    assert(build_log({"single"}) == "single");
    assert(build_log({}) == "");
    assert(build_log({"d", "c", "b", "a"}) == "abcd");
    assert(build_log({"world", "hello"}) == "helloworld");
    assert(build_log({"x", "x", "x"}) == "xxx");

    // Tier 2: an in-test stopwatch, not a harness timeout. n=200,000
    // chunks of 15 bytes each (~3MB total copying for a correct,
    // linear-time build):
    //   fixed:     ~3,000,000 bytes copied total -- well under 100ms
    //              even at -O0 under AddressSanitizer -- >=100x
    //              headroom under the 10s bound.
    //   quadratic: ~15 * 200,000^2 / 2 =~ 3*10^11 bytes (~300GB) of
    //              copying -- tens of seconds (20-40s observed) --
    //              >=2x OVER the 10s bound.
    // The 10s cutoff sits comfortably between those two, so this
    // cannot flake based on machine speed; it only distinguishes
    // O(n) from O(n^2). std::chrono::steady_clock is monotonic by
    // definition.
    const int n = 200000;
    std::vector<std::string> chunks;
    chunks.reserve(n);
    for (int i = 0; i < n; i++) {
        char buf[16];
        std::snprintf(buf, sizeof(buf), "%014d;", i);
        chunks.emplace_back(buf);
    }

    auto start = std::chrono::steady_clock::now();
    std::string result = build_log(chunks);
    auto elapsed = std::chrono::steady_clock::now() - start;
    double elapsed_s = std::chrono::duration<double>(elapsed).count();

    assert(result.size() == static_cast<size_t>(n) * 15);
    assert(result.compare(0, 15, chunks.back()) == 0);
    assert(result.compare(result.size() - 15, 15, chunks.front()) == 0);
    assert(elapsed_s < 10.0);

    printf("all assertions passed\n");
    return 0;
}
