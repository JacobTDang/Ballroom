#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> remove_value(std::vector<int> v, int target);

int main() {
    assert((remove_value({1, 2, 2, 2, 3}, 2) == std::vector<int>{1, 3}));
    assert((remove_value({2, 2, 5, 2}, 2) == std::vector<int>{5}));
    assert((remove_value({1, 3, 5}, 9) == std::vector<int>{1, 3, 5}));
    assert((remove_value({2, 2, 2, 2}, 2) == std::vector<int>{}));
    assert((remove_value({5, 2, 2, 5}, 2) == std::vector<int>{5, 5}));
    assert((remove_value({2}, 2) == std::vector<int>{}));
    assert((remove_value({}, 2) == std::vector<int>{}));
    printf("all assertions passed\n");
    return 0;
}
