#include <unordered_set>
#include <vector>

// Returns true if any value appears at least twice in nums.
bool contains_duplicate(const std::vector<int>& nums) {
    std::unordered_set<int> seen;
    for (int n : nums) {
        if (seen.count(n)) return true;
        seen.insert(n);
    }
    return false;
}
