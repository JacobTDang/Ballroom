#include <cassert>
#include <cstdio>
#include <vector>

int LastStoneWeight(std::vector<int>& stones);

int main() {
    std::vector<int> a = {2, 7, 4, 1, 8, 1};
    assert(LastStoneWeight(a) == 1);
    std::vector<int> b = {1};
    assert(LastStoneWeight(b) == 1);
    std::vector<int> c = {1, 1};
    assert(LastStoneWeight(c) == 0);
    std::vector<int> d = {1, 3};
    assert(LastStoneWeight(d) == 2);
    std::vector<int> e = {2, 2};
    assert(LastStoneWeight(e) == 0);
    std::vector<int> f = {10, 4, 2, 10};
    assert(LastStoneWeight(f) == 2);
    std::vector<int> g = {1, 1, 1, 1};
    assert(LastStoneWeight(g) == 0);
    printf("all assertions passed\n");
    return 0;
}
