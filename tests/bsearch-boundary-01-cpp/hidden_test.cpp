#include <cassert>
#include <cstdio>
#include <vector>

int first_at_least(const std::vector<int>& v, int target);

int main() {
    assert(first_at_least({1, 3, 5, 7, 9}, 6) == 3);
    assert(first_at_least({1, 3, 5, 7, 9}, 1) == 0);
    assert(first_at_least({1, 3, 5, 7, 9}, 0) == 0);
    assert(first_at_least({1, 3, 5, 7, 9}, 9) == 4);
    assert(first_at_least({1, 3, 5, 7, 9}, 10) == 5);
    assert(first_at_least({5}, 10) == 1);
    assert(first_at_least({2, 2, 2, 2}, 2) == 0);
    printf("all assertions passed\n");
    return 0;
}
