#include <algorithm>
#include <cassert>
#include <cstdio>
#include <vector>

std::vector<std::vector<int>> CombinationSum(std::vector<int>& candidates, int target);

std::vector<std::vector<int>> normalize(std::vector<std::vector<int>> lists) {
    for (auto& l : lists) std::sort(l.begin(), l.end());
    std::sort(lists.begin(), lists.end());
    return lists;
}

void check(std::vector<int> candidates, int target, std::vector<std::vector<int>> want) {
    auto got = normalize(CombinationSum(candidates, target));
    assert(got == normalize(want));
}

int main() {
    check({2, 3, 6, 7}, 7, {{2, 2, 3}, {7}});
    check({2, 3, 5}, 8, {{2, 2, 2, 2}, {2, 3, 3}, {3, 5}});
    check({2}, 1, {});
    printf("all assertions passed\n");
    return 0;
}
