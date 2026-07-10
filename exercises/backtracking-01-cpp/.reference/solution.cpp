#include <functional>
#include <vector>

// Subsets returns every subset of nums (the power set).
std::vector<std::vector<int>> Subsets(std::vector<int>& nums) {
    std::vector<std::vector<int>> res;
    std::vector<int> cur;
    std::function<void(int)> backtrack = [&](int start) {
        res.push_back(cur);
        for (int i = start; i < static_cast<int>(nums.size()); i++) {
            cur.push_back(nums[i]);
            backtrack(i + 1);
            cur.pop_back();
        }
    };
    backtrack(0);
    return res;
}
