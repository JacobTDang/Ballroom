#include <cassert>
#include <cstdio>

#include "solution.hpp"

int main() {
    {
        MinStack s;
        s.push(-2);
        s.push(0);
        s.push(-3);
        assert(s.getMin() == -3);
        s.pop();
        assert(s.top() == 0);
        assert(s.getMin() == -2);
    }
    {
        MinStack s;
        s.push(1);
        s.push(1);
        s.push(1);
        assert(s.getMin() == 1);
        s.pop();
        assert(s.getMin() == 1);
    }
    {
        MinStack s;
        s.push(5);
        s.push(3);
        s.push(7);
        s.push(1);
        assert(s.getMin() == 1);
        s.pop();
        assert(s.getMin() == 3);
    }
    printf("all assertions passed\n");
    return 0;
}
