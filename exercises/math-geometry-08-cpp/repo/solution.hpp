#pragma once

#include <map>
#include <utility>
#include <vector>

// DetectSquares tracks added points and counts axis-aligned squares
// formable with a query point.
class DetectSquares {
public:
    void add(std::vector<int> point) {
        // TODO: implement
    }

    int count(std::vector<int> point) {
        // TODO: implement
        return 0;
    }

private:
    std::map<std::pair<int, int>, int> points_;
};
