#include <algorithm>
#include <cassert>
#include <cstdio>
#include <vector>

std::vector<std::vector<int>> SubsetsWithDup(std::vector<int>& nums);

std::vector<std::vector<int>> normalize(std::vector<std::vector<int>> lists) {
    for (auto& l : lists) std::sort(l.begin(), l.end());
    std::sort(lists.begin(), lists.end());
    return lists;
}

void check(std::vector<int> nums, std::vector<std::vector<int>> want) {
    auto got = normalize(SubsetsWithDup(nums));
    assert(got == normalize(want));
}

int main() {
    check({1, 2, 2}, {{}, {1}, {1, 2}, {1, 2, 2}, {2}, {2, 2}});
    check({0}, {{}, {0}});
    printf("all assertions passed\n");
    return 0;
}
