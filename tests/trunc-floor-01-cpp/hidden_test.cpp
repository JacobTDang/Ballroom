#include <cassert>
#include <cstdio>

int align(int t, int k);

int main() {
    assert(align(7, 4) == 4);
    assert(align(8, 4) == 8);
    assert(align(0, 4) == 0);
    assert(align(-1, 4) == -4);
    assert(align(-7, 4) == -8);
    assert(align(-8, 4) == -8);
    assert(align(-9, 4) == -12);
    assert(align(-10, 3) == -12);
    assert(align(-12, 3) == -12);
    printf("all assertions passed\n");
    return 0;
}
