#include <cassert>
#include <cstdio>
#include <vector>

int max_of(const std::vector<int>& v);

int main() {
    assert(max_of({3, 1, 4, 1, 5, 9, 2, 6}) == 9);
    assert(max_of({-5, -1, -10}) == -1);
    assert(max_of({42}) == 42);
    printf("all assertions passed\n");
    return 0;
}
