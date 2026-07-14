#include <cassert>
#include <cstdio>
#include <vector>

int Search(const std::vector<int>& nums, int target);

int main() {
    assert(Search({-1, 0, 3, 5, 9, 12}, 9) == 4);
    assert(Search({-1, 0, 3, 5, 9, 12}, 2) == -1);
    assert(Search({5}, 5) == 0);
    assert(Search({2, 5}, 5) == 1);
    assert(Search({2, 5}, 1) == -1);
    assert(Search({1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1) == 0);
    assert(Search({1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10) == 9);
    assert(Search({-10, -5, 0, 5, 10}, -5) == 1);
    printf("all assertions passed\n");
    return 0;
}
