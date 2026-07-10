#include <cassert>
#include <cstdio>
#include <vector>

bool contains_duplicate(const std::vector<int>& nums);

int main() {
    assert(contains_duplicate({1, 2, 3, 1}) == true);
    assert(contains_duplicate({1, 2, 3, 4}) == false);
    assert(contains_duplicate({1, 1, 1, 3, 3, 4, 3, 2, 4, 2}) == true);
    assert(contains_duplicate({1}) == false);
    assert(contains_duplicate({-1, -1}) == true);
    assert(contains_duplicate({0, 4, 5, 0, 3, 6}) == true);
    printf("all assertions passed\n");
    return 0;
}
