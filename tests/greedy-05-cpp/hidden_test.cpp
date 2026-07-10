#include <cassert>
#include <cstdio>
#include <vector>

bool IsNStraightHand(std::vector<int>& hand, int groupSize);

void testClassic() {
    std::vector<int> hand = {1, 2, 3, 6, 2, 3, 4, 7, 8};
    assert(IsNStraightHand(hand, 3) == true);
}

void testNotDivisible() {
    std::vector<int> hand = {1, 2, 3, 4, 5};
    assert(IsNStraightHand(hand, 4) == false);
}

void testMissingCard() {
    std::vector<int> hand = {1, 2, 3, 4, 5, 7};
    assert(IsNStraightHand(hand, 3) == false);
}

void testGroupSizeOne() {
    std::vector<int> hand = {5, 5, 5};
    assert(IsNStraightHand(hand, 1) == true);
}

void testExactSingleGroup() {
    std::vector<int> hand = {1, 2, 3};
    assert(IsNStraightHand(hand, 3) == true);
}

int main() {
    testClassic();
    testNotDivisible();
    testMissingCard();
    testGroupSizeOne();
    testExactSingleGroup();
    std::printf("all tests passed\n");
    return 0;
}
