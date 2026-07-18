#include <cassert>
#include <cstdio>
#include <stdexcept>
#include <vector>

int max_adjacent_diff(const std::vector<int>& v);

int main() {
    assert(max_adjacent_diff({3, 1, 4, 1, 5, 9, 2, 6}) == 7);
    assert(max_adjacent_diff({5, 5}) == 0);
    assert(max_adjacent_diff({-5, -1, -10}) == 9);
    assert(max_adjacent_diff({1, 100}) == 99);

    bool empty_threw = false;
    try {
        max_adjacent_diff({});
    } catch (const std::invalid_argument&) {
        empty_threw = true;
    }
    assert(empty_threw);

    bool single_threw = false;
    try {
        max_adjacent_diff({42});
    } catch (const std::invalid_argument&) {
        single_threw = true;
    }
    assert(single_threw);

    printf("all assertions passed\n");
    return 0;
}
