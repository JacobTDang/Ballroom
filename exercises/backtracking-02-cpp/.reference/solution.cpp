#include <functional>
#include <vector>

// CombinationSum returns every unique combination of candidates
// (each usable unlimited times) that sums to target.
std::vector<std::vector<int>> CombinationSum(std::vector<int>& candidates, int target) {
    std::vector<std::vector<int>> res;
    std::vector<int> cur;
    std::function<void(int, int)> backtrack = [&](int start, int remain) {
        if (remain == 0) {
            res.push_back(cur);
            return;
        }
        if (remain < 0) return;
        for (int i = start; i < static_cast<int>(candidates.size()); i++) {
            cur.push_back(candidates[i]);
            backtrack(i, remain - candidates[i]);
            cur.pop_back();
        }
    };
    backtrack(0, target);
    return res;
}
