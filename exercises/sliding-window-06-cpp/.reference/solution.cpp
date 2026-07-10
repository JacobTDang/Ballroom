#include <deque>
#include <vector>

// MaxSlidingWindow returns the maximum of each contiguous window of
// size k as it slides from the start of nums to the end.
std::vector<int> MaxSlidingWindow(const std::vector<int>& nums, int k) {
    std::deque<int> dq;  // indices into nums, values strictly decreasing
    std::vector<int> res;
    for (int i = 0; i < static_cast<int>(nums.size()); i++) {
        while (!dq.empty() && nums[dq.back()] < nums[i]) {
            dq.pop_back();
        }
        dq.push_back(i);
        if (dq.front() <= i - k) {
            dq.pop_front();
        }
        if (i >= k - 1) {
            res.push_back(nums[dq.front()]);
        }
    }
    return res;
}
