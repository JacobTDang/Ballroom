#include <algorithm>
#include <cassert>
#include <cstdio>
#include <vector>

std::vector<std::vector<int>> Permute(std::vector<int>& nums);

std::vector<std::vector<int>> normalizeExact(std::vector<std::vector<int>> lists) {
    std::sort(lists.begin(), lists.end());
    return lists;
}

void check(std::vector<int> nums, std::vector<std::vector<int>> want) {
    auto got = normalizeExact(Permute(nums));
    assert(got == normalizeExact(want));
}

int main() {
    check({1, 2, 3}, {{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}});
    check({0, 1}, {{0, 1}, {1, 0}});
    check({1}, {{1}});
    check({1, -1}, {{1, -1}, {-1, 1}});
    printf("all assertions passed\n");
    return 0;
}
