#include <algorithm>
#include <cassert>
#include <cstdio>
#include <vector>

std::vector<std::vector<int>> ThreeSum(std::vector<int>& nums);

// normalize sorts each triplet ascending, then sorts the list of
// triplets — 3Sum's valid outputs aren't uniquely ordered, so tests
// compare as sets rather than asserting an exact sequence.
std::vector<std::vector<int>> normalize(std::vector<std::vector<int>> triplets) {
    for (auto& t : triplets) std::sort(t.begin(), t.end());
    std::sort(triplets.begin(), triplets.end());
    return triplets;
}

void check(std::vector<int> nums, std::vector<std::vector<int>> want) {
    auto got = normalize(ThreeSum(nums));
    assert(got == normalize(want));
}

int main() {
    check({-1, 0, 1, 2, -1, -4}, {{-1, -1, 2}, {-1, 0, 1}});
    check({0, 1, 1}, {});
    check({0, 0, 0}, {{0, 0, 0}});
    check({}, {});
    check({0, 0, 0, 0}, {{0, 0, 0}});
    check({-2, 0, 1, 1, 2}, {{-2, 0, 2}, {-2, 1, 1}});
    check({-3, -2, -1}, {});
    check({1, 2, 3}, {});
    check({1, -1}, {});
    check({3, -2, 1, 0, -1, -3, 2, -2, 0},
          {{-3, 0, 3}, {-3, 1, 2}, {-2, -1, 3}, {-2, 0, 2}, {-1, 0, 1}});
    printf("all assertions passed\n");
    return 0;
}
