#include <cassert>
#include <cstdio>
#include <vector>

int Search(const std::vector<int>& nums, int target);

int main() {
    assert(Search({4, 5, 6, 7, 0, 1, 2}, 0) == 4);
    assert(Search({4, 5, 6, 7, 0, 1, 2}, 3) == -1);
    assert(Search({1}, 0) == -1);
    assert(Search({5, 1, 3}, 5) == 0);
    assert(Search({1, 3}, 3) == 1);
    assert(Search({9, 10, 1, 2, 3, 4, 5, 6, 7, 8}, 8) == 9);
    assert(Search({9, 10, 1, 2, 3, 4, 5, 6, 7, 8}, 100) == -1);
    printf("all assertions passed\n");
    return 0;
}
