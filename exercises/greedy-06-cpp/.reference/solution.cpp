#include <vector>

// MergeTriplets returns whether target can be formed by taking the
// elementwise max of some subset of triplets.
bool MergeTriplets(std::vector<std::vector<int>>& triplets, std::vector<int>& target) {
    bool matched[3] = {false, false, false};
    for (auto& t : triplets) {
        if (t[0] > target[0] || t[1] > target[1] || t[2] > target[2]) continue;
        for (int i = 0; i < 3; i++) {
            if (t[i] == target[i]) matched[i] = true;
        }
    }
    return matched[0] && matched[1] && matched[2];
}
