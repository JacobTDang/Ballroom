#include <cassert>
#include <cstdio>
#include <vector>

int FindKthLargest(std::vector<int>& nums, int k);

int main() {
    std::vector<int> a = {3, 2, 1, 5, 6, 4};
    assert(FindKthLargest(a, 2) == 5);
    std::vector<int> b = {3, 2, 3, 1, 2, 4, 5, 5, 6};
    assert(FindKthLargest(b, 4) == 4);
    std::vector<int> c = {1};
    assert(FindKthLargest(c, 1) == 1);
    std::vector<int> d = {2, 1};
    assert(FindKthLargest(d, 2) == 1);
    std::vector<int> e = {5, 5, 5, 5};
    assert(FindKthLargest(e, 2) == 5);
    std::vector<int> f = {-1, -5, -3, -2, -4};
    assert(FindKthLargest(f, 3) == -3);
    printf("all assertions passed\n");
    return 0;
}
