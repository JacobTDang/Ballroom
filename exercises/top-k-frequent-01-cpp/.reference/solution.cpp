#include <unordered_map>
#include <vector>

// Returns the k most frequent elements in nums, in any order.
std::vector<int> top_k_frequent(const std::vector<int>& nums, int k) {
    std::unordered_map<int, int> counts;
    for (int n : nums) counts[n]++;

    std::vector<std::vector<int>> buckets(nums.size() + 1);
    for (auto& [n, c] : counts) {
        buckets[c].push_back(n);
    }

    std::vector<int> result;
    for (int i = static_cast<int>(buckets.size()) - 1; i >= 0 && static_cast<int>(result.size()) < k; i--) {
        for (int n : buckets[i]) {
            result.push_back(n);
            if (static_cast<int>(result.size()) == k) break;
        }
    }
    return result;
}
