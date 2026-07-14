#include <cassert>
#include <cmath>
#include <cstdio>
#include <vector>

double FindMedianSortedArrays(std::vector<int> nums1, std::vector<int> nums2);

void check(std::vector<int> a, std::vector<int> b, double want) {
    double got = FindMedianSortedArrays(a, b);
    assert(std::fabs(got - want) < 1e-5);
}

int main() {
    check({1, 3}, {2}, 2.0);
    check({1, 2}, {3, 4}, 2.5);
    check({}, {1}, 1.0);
    check({2}, {}, 2.0);
    check({1, 2, 3, 4, 5}, {6, 7, 8, 9, 10}, 5.5);
    check({1, 2, 3}, {4, 5, 6, 7, 8, 9}, 5.0);
    check({1}, {2, 3, 4, 5}, 3.0);
    check({100, 200}, {1, 2, 3}, 3.0);
    printf("all assertions passed\n");
    return 0;
}
