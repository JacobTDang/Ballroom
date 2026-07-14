#include <vector>

// CanPartition returns whether nums can be partitioned into two
// subsets with equal sum.
bool CanPartition(std::vector<int>& nums) {
    int sum = 0;
    for (int n : nums) sum += n;
    if (sum % 2 != 0) return false;

    int target = sum / 2;
    std::vector<bool> reachable(target + 1, false);
    reachable[0] = true;

    for (int n : nums) {
        for (int i = target; i >= n; i--) {
            if (reachable[i - n]) reachable[i] = true;
        }
    }

    return reachable[target];
}
