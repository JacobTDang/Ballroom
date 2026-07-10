#include <vector>

// CountBits returns a vector ans of length n+1 where ans[i] is the number
// of set bits in the binary representation of i.
std::vector<int> CountBits(int n) {
    std::vector<int> ans(n + 1, 0);
    for (int i = 1; i <= n; i++) {
        ans[i] = ans[i >> 1] + (i & 1);
    }
    return ans;
}
