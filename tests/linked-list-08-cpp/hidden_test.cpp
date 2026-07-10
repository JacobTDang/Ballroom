#include <cassert>
#include <cstdio>
#include <vector>

int FindDuplicate(const std::vector<int>& nums);

int main() {
    assert(FindDuplicate({1, 3, 4, 2, 2}) == 2);
    assert(FindDuplicate({3, 1, 3, 4, 2}) == 3);
    assert(FindDuplicate({1, 1}) == 1);
    assert(FindDuplicate({1, 1, 2}) == 1);
    assert(FindDuplicate({2, 2, 2, 2, 2}) == 2);
    printf("all assertions passed\n");
    return 0;
}
