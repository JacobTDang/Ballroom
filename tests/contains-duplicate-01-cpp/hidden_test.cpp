#include <cassert>
#include <cstdio>
#include <vector>

bool contains_duplicate(const std::vector<int>& nums);

int main() {
    assert(contains_duplicate({1, 2, 3, 1}) == true);
    assert(contains_duplicate({1, 2, 3, 4}) == false);
    assert(contains_duplicate({1, 1, 1, 3, 3, 4, 3, 2, 4, 2}) == true);
    assert(contains_duplicate({1}) == false);
    assert(contains_duplicate({1, 1}) == true);
    assert(contains_duplicate({1, 2}) == false);
    assert(contains_duplicate({-1, -1}) == true);
    assert(contains_duplicate({-5, -3, -1, 1, 3, 5}) == false);
    assert(contains_duplicate({0, 4, 5, 0, 3, 6}) == true);
    assert(contains_duplicate({7, 7, 7, 7, 7}) == true);
    assert(contains_duplicate({1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1}) == true);
    assert(contains_duplicate({-1000000000, 1000000000}) == false);
    assert(contains_duplicate({1000000000, 1000000000}) == true);
    printf("all assertions passed\n");
    return 0;
}
