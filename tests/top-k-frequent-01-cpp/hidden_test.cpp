#include <algorithm>
#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> top_k_frequent(const std::vector<int>& nums, int k);

bool same_set(std::vector<int> a, std::vector<int> b) {
    std::sort(a.begin(), a.end());
    std::sort(b.begin(), b.end());
    return a == b;
}

int main() {
    assert(same_set(top_k_frequent({1, 1, 1, 2, 2, 3}, 2), {1, 2}));
    assert(same_set(top_k_frequent({1}, 1), {1}));
    assert(same_set(top_k_frequent({4, 1, -1, 2, -1, 2, 3}, 2), {-1, 2}));
    assert(same_set(top_k_frequent({5, 5, 5, 5, 3, 3, 1}, 1), {5}));
    printf("all assertions passed\n");
    return 0;
}
