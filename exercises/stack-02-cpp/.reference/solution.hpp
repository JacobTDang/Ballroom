#pragma once

#include <vector>

// MinStack is a stack that also tracks its minimum element in O(1).
class MinStack {
public:
    void push(int val) {
        stack_.push_back(val);
        if (minStack_.empty() || val < minStack_.back()) {
            minStack_.push_back(val);
        } else {
            minStack_.push_back(minStack_.back());
        }
    }

    void pop() {
        stack_.pop_back();
        minStack_.pop_back();
    }

    int top() {
        return stack_.back();
    }

    int getMin() {
        return minStack_.back();
    }

private:
    std::vector<int> stack_;
    std::vector<int> minStack_;
};
