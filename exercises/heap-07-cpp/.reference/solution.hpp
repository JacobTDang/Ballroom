#pragma once

#include <queue>
#include <vector>

// MedianFinder tracks the running median of a stream of integers,
// using two heaps that split the stream in half: small holds the
// lower half (max-heap, so its top is the largest of the low half)
// and large holds the upper half (min-heap, so its top is the
// smallest of the high half). Kept balanced within 1 of each other
// after every insert, so the median is always at the top of one (or
// both) heaps.
class MedianFinder {
public:
    void addNum(int num) {
        small_.push(num);
        large_.push(small_.top());
        small_.pop();
        if (large_.size() > small_.size()) {
            small_.push(large_.top());
            large_.pop();
        }
    }

    double findMedian() {
        if (small_.size() > large_.size()) return small_.top();
        return (small_.top() + large_.top()) / 2.0;
    }

private:
    std::priority_queue<int> small_;                                        // max-heap
    std::priority_queue<int, std::vector<int>, std::greater<int>> large_;   // min-heap
};
