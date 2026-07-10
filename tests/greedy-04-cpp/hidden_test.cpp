#include <cassert>
#include <cstdio>
#include <vector>

int CanCompleteCircuit(std::vector<int>& gas, std::vector<int>& cost);

void testClassic() {
    std::vector<int> gas = {1, 2, 3, 4, 5};
    std::vector<int> cost = {3, 4, 5, 1, 2};
    assert(CanCompleteCircuit(gas, cost) == 3);
}

void testImpossible() {
    std::vector<int> gas = {2, 3, 4};
    std::vector<int> cost = {3, 4, 3};
    assert(CanCompleteCircuit(gas, cost) == -1);
}

void testSingleExact() {
    std::vector<int> gas = {5};
    std::vector<int> cost = {4};
    assert(CanCompleteCircuit(gas, cost) == 0);
}

void testSingleInsufficient() {
    std::vector<int> gas = {3};
    std::vector<int> cost = {4};
    assert(CanCompleteCircuit(gas, cost) == -1);
}

void testStartAtLastIndex() {
    std::vector<int> gas = {5, 1, 2, 3, 4};
    std::vector<int> cost = {4, 4, 1, 5, 1};
    assert(CanCompleteCircuit(gas, cost) == 4);
}

int main() {
    testClassic();
    testImpossible();
    testSingleExact();
    testSingleInsufficient();
    testStartAtLastIndex();
    std::printf("all tests passed\n");
    return 0;
}
