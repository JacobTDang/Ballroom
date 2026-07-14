#include <cassert>
#include <cstdio>
#include <vector>

#include "solution.hpp"

int main() {
    {
        DetectSquares ds;
        ds.add({3, 10});
        ds.add({11, 2});
        ds.add({3, 2});

        assert(ds.count({11, 10}) == 1);
        assert(ds.count({14, 8}) == 0);

        ds.add({11, 2});
        assert(ds.count({11, 10}) == 2);
    }
    {
        DetectSquares ds;
        assert(ds.count({0, 0}) == 0);
    }
    {
        DetectSquares ds;
        ds.add({0, 2});
        ds.add({2, 0});
        ds.add({2, 2});
        ds.add({-2, 0});
        ds.add({-2, 2});

        assert(ds.count({0, 0}) == 2);
    }
    {
        DetectSquares ds;
        ds.add({1, 4});
        ds.add({4, 1});
        ds.add({4, 4});

        assert(ds.count({1, 1}) == 1);
    }
    {
        DetectSquares ds;
        ds.add({1, 4});
        ds.add({1, 4});
        ds.add({1, 4});
        ds.add({4, 1});
        ds.add({4, 1});
        ds.add({4, 4});

        assert(ds.count({1, 1}) == 6);
    }
    {
        DetectSquares ds;
        ds.add({5, 5});
        ds.add({5, 9});
        ds.add({9, 5});
        ds.add({9, 9});

        assert(ds.count({100, 100}) == 0);
    }
    {
        DetectSquares ds;
        ds.add({0, 2});
        ds.add({2, 0});
        ds.add({2, 2});
        ds.add({-2, 0});
        ds.add({-2, 2});

        assert(ds.count({0, 0}) == 2);
        assert(ds.count({0, 0}) == 2);

        ds.add({0, 2});
        assert(ds.count({0, 0}) == 4);
    }
    std::printf("all assertions passed\n");
    return 0;
}
