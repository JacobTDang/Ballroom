#include <cassert>
#include <cstdio>
#include <vector>

int FindMin(const std::vector<int>& nums);

int main() {
    assert(FindMin({3, 4, 5, 1, 2}) == 1);
    assert(FindMin({4, 5, 6, 7, 0, 1, 2}) == 0);
    assert(FindMin({11, 13, 15, 17}) == 11);
    assert(FindMin({2, 1}) == 1);
    assert(FindMin({1}) == 1);
    assert(FindMin({1, 2, 3, 4, 5}) == 1);
    assert(FindMin({15, 18, 2, 3, 6, 12}) == 2);
    printf("all assertions passed\n");
    return 0;
}
