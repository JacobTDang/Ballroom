#include <cassert>
#include <cstdio>
#include <vector>

bool MergeTriplets(std::vector<std::vector<int>>& triplets, std::vector<int>& target);

void testClassic() {
    std::vector<std::vector<int>> triplets = {{2, 5, 3}, {1, 8, 4}, {1, 7, 5}};
    std::vector<int> target = {2, 7, 5};
    assert(MergeTriplets(triplets, target) == true);
}

void testClassicFalse() {
    std::vector<std::vector<int>> triplets = {{3, 4, 5}, {4, 5, 6}};
    std::vector<int> target = {3, 2, 5};
    assert(MergeTriplets(triplets, target) == false);
}

void testSingleExact() {
    std::vector<std::vector<int>> triplets = {{5, 5, 5}};
    std::vector<int> target = {5, 5, 5};
    assert(MergeTriplets(triplets, target) == true);
}

void testAllPoisoned() {
    std::vector<std::vector<int>> triplets = {{10, 1, 1}, {1, 10, 1}, {1, 1, 10}};
    std::vector<int> target = {5, 5, 5};
    assert(MergeTriplets(triplets, target) == false);
}

void testPartialMatch() {
    std::vector<std::vector<int>> triplets = {{2, 1, 1}, {1, 2, 1}};
    std::vector<int> target = {2, 2, 2};
    assert(MergeTriplets(triplets, target) == false);
}

int main() {
    testClassic();
    testClassicFalse();
    testSingleExact();
    testAllPoisoned();
    testPartialMatch();
    std::printf("all tests passed\n");
    return 0;
}
