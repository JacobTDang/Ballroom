#include <algorithm>
#include <functional>
#include <vector>

// CombinationSum2 returns every unique combination of candidates
// (each usable at most once) that sums to target.
std::vector<std::vector<int>> CombinationSum2(std::vector<int>& candidates, int target) {
    std::vector<int> sorted = candidates;
    std::sort(sorted.begin(), sorted.end());

    std::vector<std::vector<int>> res;
    std::vector<int> cur;
    std::function<void(int, int)> backtrack = [&](int start, int remain) {
        if (remain == 0) {
            res.push_back(cur);
            return;
        }
        for (int i = start; i < static_cast<int>(sorted.size()); i++) {
            if (i > start && sorted[i] == sorted[i - 1]) continue;
            if (sorted[i] > remain) break;
            cur.push_back(sorted[i]);
            backtrack(i + 1, remain - sorted[i]);
            cur.pop_back();
        }
    };
    backtrack(0, target);
    return res;
}
