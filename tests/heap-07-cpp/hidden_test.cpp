#include <cassert>
#include <cmath>
#include <cstdio>

#include "solution.hpp"

bool closeEnough(double a, double b) {
    return std::fabs(a - b) < 1e-5;
}

int main() {
    {
        MedianFinder mf;
        mf.addNum(1);
        mf.addNum(2);
        assert(closeEnough(mf.findMedian(), 1.5));
        mf.addNum(3);
        assert(closeEnough(mf.findMedian(), 2.0));
    }
    {
        MedianFinder mf;
        mf.addNum(42);
        assert(closeEnough(mf.findMedian(), 42.0));
    }
    {
        MedianFinder mf;
        for (int n : {5, 1, 9, 3, 7}) mf.addNum(n);
        assert(closeEnough(mf.findMedian(), 5.0));
        mf.addNum(10);
        assert(closeEnough(mf.findMedian(), 6.0));
    }
    {
        MedianFinder mf;
        for (int n : {-5, -1, -3}) mf.addNum(n);
        assert(closeEnough(mf.findMedian(), -3.0));
        mf.addNum(-2);
        assert(closeEnough(mf.findMedian(), -2.5));
    }
    printf("all assertions passed\n");
    return 0;
}
