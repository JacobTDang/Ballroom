#include <cassert>
#include <cstdio>

bool IsHappy(int n);

int main() {
    assert(IsHappy(19) == true);
    assert(IsHappy(2) == false);
    assert(IsHappy(1) == true);
    assert(IsHappy(7) == true);
    assert(IsHappy(4) == false);
    assert(IsHappy(100) == true);
    assert(IsHappy(3) == false);
    assert(IsHappy(986) == false);
    assert(IsHappy(2147483647) == false);
    std::printf("all assertions passed\n");
    return 0;
}
