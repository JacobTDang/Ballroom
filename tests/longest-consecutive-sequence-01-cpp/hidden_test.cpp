#include <cassert>
#include <cstdio>
#include <vector>

int longest_consecutive(const std::vector<int>& nums);

int main() {
    assert(longest_consecutive({100, 4, 200, 1, 3, 2}) == 4);
    assert(longest_consecutive({0, 3, 7, 2, 5, 8, 4, 6, 0, 1}) == 9);
    assert(longest_consecutive({}) == 0);
    assert(longest_consecutive({1, 2, 0, 1}) == 3);
    assert(longest_consecutive({9, 1, 4, 7, 3, -1, 0, 5, 8, -1, 6}) == 7);
    assert(longest_consecutive({5}) == 1);
    printf("all assertions passed\n");
    return 0;
}
