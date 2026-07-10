#include <unordered_set>
#include <vector>

// Returns the length of the longest run of consecutive integers present
// in nums (order doesn't matter, duplicates don't count extra).
int longest_consecutive(const std::vector<int>& nums) {
    std::unordered_set<int> set(nums.begin(), nums.end());

    int longest = 0;
    for (int n : set) {
        if (set.count(n - 1)) continue;  // n isn't the start of a sequence
        int length = 1;
        while (set.count(n + length)) length++;
        if (length > longest) longest = length;
    }
    return longest;
}
