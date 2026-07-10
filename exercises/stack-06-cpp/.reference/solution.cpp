#include <algorithm>
#include <vector>

// CarFleet returns the number of car fleets that will arrive at target.
int CarFleet(int target, const std::vector<int>& position, const std::vector<int>& speed) {
    int n = static_cast<int>(position.size());
    std::vector<int> idx(n);
    for (int i = 0; i < n; i++) idx[i] = i;
    std::sort(idx.begin(), idx.end(),
              [&](int a, int b) { return position[a] > position[b]; });

    std::vector<double> stack;  // arrival times of fleets found so far
    for (int i : idx) {
        double t = static_cast<double>(target - position[i]) / speed[i];
        if (stack.empty() || t > stack.back()) {
            stack.push_back(t);
        }
    }
    return static_cast<int>(stack.size());
}
