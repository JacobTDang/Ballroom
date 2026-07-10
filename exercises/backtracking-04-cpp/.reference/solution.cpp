#include <algorithm>
#include <functional>
#include <vector>

// SubsetsWithDup returns every unique subset of nums, which may
// contain duplicate values.
std::vector<std::vector<int>> SubsetsWithDup(std::vector<int>& nums) {
    std::vector<int> sorted = nums;
    std::sort(sorted.begin(), sorted.end());

    std::vector<std::vector<int>> res;
    std::vector<int> cur;
    std::function<void(int)> backtrack = [&](int start) {
        res.push_back(cur);
        for (int i = start; i < static_cast<int>(sorted.size()); i++) {
            if (i > start && sorted[i] == sorted[i - 1]) continue;
            cur.push_back(sorted[i]);
            backtrack(i + 1);
            cur.pop_back();
        }
    };
    backtrack(0);
    return res;
}
