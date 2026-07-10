#include <cassert>
#include <cstdio>
#include <vector>

#include "solution.hpp"

int main() {
    {
        std::vector<int> nums = {4, 5, 8, 2};
        KthLargest kl(3, nums);
        assert(kl.add(3) == 4);
        assert(kl.add(5) == 5);
        assert(kl.add(10) == 5);
        assert(kl.add(9) == 8);
        assert(kl.add(4) == 8);
    }
    {
        std::vector<int> nums = {};
        KthLargest kl(1, nums);
        assert(kl.add(-3) == -3);
        assert(kl.add(-2) == -2);
    }
    printf("all assertions passed\n");
    return 0;
}
