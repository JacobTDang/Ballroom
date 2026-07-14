#include <algorithm>
#include <cassert>
#include <cstdio>
#include <vector>

std::vector<std::vector<int>> CombinationSum2(std::vector<int>& candidates, int target);

std::vector<std::vector<int>> normalize(std::vector<std::vector<int>> lists) {
    for (auto& l : lists) std::sort(l.begin(), l.end());
    std::sort(lists.begin(), lists.end());
    return lists;
}

void check(std::vector<int> candidates, int target, std::vector<std::vector<int>> want) {
    auto got = normalize(CombinationSum2(candidates, target));
    assert(got == normalize(want));
}

int main() {
    check({10, 1, 2, 7, 6, 1, 5}, 8, {{1, 1, 6}, {1, 2, 5}, {1, 7}, {2, 6}});
    check({2, 5, 2, 1, 2}, 5, {{1, 2, 2}, {5}});
    check({1, 1, 1, 2, 2}, 4, {{1, 1, 2}, {2, 2}});
    printf("all assertions passed\n");
    return 0;
}
