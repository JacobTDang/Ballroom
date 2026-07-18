#include <vector>

// Removes every occurrence of target from v, in place, and returns it.
std::vector<int> remove_value(std::vector<int> v, int target) {
    size_t i = 0;
    while (i < v.size()) {
        if (v[i] == target) {
            v.erase(v.begin() + i);
        } else {
            i++;
        }
    }
    return v;
}
