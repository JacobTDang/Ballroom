#include <cassert>
#include <cstdio>

int ClimbStairs(int n);

void testOne() { assert(ClimbStairs(1) == 1); }
void testTwo() { assert(ClimbStairs(2) == 2); }
void testThree() { assert(ClimbStairs(3) == 3); }
void testFive() { assert(ClimbStairs(5) == 8); }
void testFour() { assert(ClimbStairs(4) == 5); }
void testTen() { assert(ClimbStairs(10) == 89); }
void testBoundaryMax() { assert(ClimbStairs(45) == 1836311903); }

int main() {
    testOne();
    testTwo();
    testThree();
    testFive();
    testFour();
    testTen();
    testBoundaryMax();
    std::printf("all tests passed\n");
    return 0;
}
