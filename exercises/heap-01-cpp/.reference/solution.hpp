#pragma once

#include <queue>
#include <vector>

// KthLargest tracks the kth largest value seen so far in a stream of
// integers, using a min-heap capped at size k -- the heap's smallest
// element (the top) is always the kth largest overall.
class KthLargest {
public:
    KthLargest(int k, std::vector<int>& nums) : k_(k) {
        for (int n : nums) add(n);
    }

    int add(int val) {
        heap_.push(val);
        if (static_cast<int>(heap_.size()) > k_) heap_.pop();
        return heap_.top();
    }

private:
    int k_;
    std::priority_queue<int, std::vector<int>, std::greater<int>> heap_;
};
