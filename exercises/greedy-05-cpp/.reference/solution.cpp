#include <algorithm>
#include <unordered_map>
#include <vector>

// IsNStraightHand returns whether hand can be rearranged into groups of
// groupSize consecutive cards.
bool IsNStraightHand(std::vector<int>& hand, int groupSize) {
    if ((int)hand.size() % groupSize != 0) return false;

    std::unordered_map<int, int> count;
    for (int c : hand) count[c]++;

    std::vector<int> keys;
    keys.reserve(count.size());
    for (auto& kv : count) keys.push_back(kv.first);
    std::sort(keys.begin(), keys.end());

    for (int k : keys) {
        int need = count[k];
        if (need == 0) continue;
        for (int i = 0; i < groupSize; i++) {
            int c = k + i;
            if (count[c] < need) return false;
            count[c] -= need;
        }
    }
    return true;
}
