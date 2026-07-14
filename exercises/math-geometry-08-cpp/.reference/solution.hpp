#pragma once

#include <map>
#include <utility>
#include <vector>

// DetectSquares tracks added points and counts axis-aligned squares
// formable with a query point, using a frequency map from point to the
// number of times it was added. For a query point, every previously
// added point sharing its x-coordinate forms a candidate vertical edge;
// the two horizontal partner corners at that same side length are then
// checked for existence.
class DetectSquares {
public:
    void add(std::vector<int> point) {
        points_[{point[0], point[1]}]++;
    }

    int count(std::vector<int> point) {
        int qx = point[0], qy = point[1];
        int total = 0;

        for (auto& [p, freq] : points_) {
            if (p.first != qx || p.second == qy) continue;
            int side = p.second - qy;
            for (int cx : {qx + side, qx - side}) {
                auto it1 = points_.find({cx, qy});
                if (it1 == points_.end()) continue;
                auto it2 = points_.find({cx, p.second});
                if (it2 == points_.end()) continue;
                total += freq * it1->second * it2->second;
            }
        }

        return total;
    }

private:
    std::map<std::pair<int, int>, int> points_;
};
