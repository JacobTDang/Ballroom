#include <cassert>
#include <cstdio>
#include <vector>

int max_below_limit(const std::vector<int>& v, int limit);

int main() {
    assert(max_below_limit({3, 7, 2, 9, 5}, 7) == 7);
    assert(max_below_limit({3, 7, 2, 9, 5}, 6) == 5);
    assert(max_below_limit({10, 20, 30}, 5) == -1);
    assert(max_below_limit({-5, -1, -10}, -2) == -5);
    assert(max_below_limit({5}, 10) == 5);
    assert(max_below_limit({15}, 10) == -1);
    printf("all assertions passed\n");
    return 0;
}
