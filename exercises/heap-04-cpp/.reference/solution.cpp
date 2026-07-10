#include <queue>
#include <vector>

// FindKthLargest returns the kth largest element of nums (1st
// largest is the maximum), via a min-heap capped at size k.
int FindKthLargest(std::vector<int>& nums, int k) {
    std::priority_queue<int, std::vector<int>, std::greater<int>> heap;
    for (int n : nums) {
        heap.push(n);
        if (static_cast<int>(heap.size()) > k) heap.pop();
    }
    return heap.top();
}
