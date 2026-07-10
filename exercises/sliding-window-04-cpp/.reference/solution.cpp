#include <array>
#include <string>

// CheckInclusion reports whether s2 contains a permutation of s1 as a
// contiguous substring.
bool CheckInclusion(const std::string& s1, const std::string& s2) {
    if (s1.size() > s2.size()) return false;
    std::array<int, 26> need{}, window{};
    for (size_t i = 0; i < s1.size(); i++) {
        need[s1[i] - 'a']++;
        window[s2[i] - 'a']++;
    }
    if (need == window) return true;
    for (size_t i = s1.size(); i < s2.size(); i++) {
        window[s2[i] - 'a']++;
        window[s2[i - s1.size()] - 'a']--;
        if (need == window) return true;
    }
    return false;
}
