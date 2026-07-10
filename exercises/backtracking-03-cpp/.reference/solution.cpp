#include <functional>
#include <vector>

// Permute returns every permutation of nums.
std::vector<std::vector<int>> Permute(std::vector<int>& nums) {
    std::vector<std::vector<int>> res;
    std::vector<int> cur;
    std::vector<bool> used(nums.size(), false);
    std::function<void()> backtrack = [&]() {
        if (cur.size() == nums.size()) {
            res.push_back(cur);
            return;
        }
        for (size_t i = 0; i < nums.size(); i++) {
            if (used[i]) continue;
            used[i] = true;
            cur.push_back(nums[i]);
            backtrack();
            cur.pop_back();
            used[i] = false;
        }
    };
    backtrack();
    return res;
}
