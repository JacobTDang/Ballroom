#include <string>
#include <unordered_map>

// LengthOfLongestSubstring returns the length of the longest substring
// of s with no repeating characters.
int LengthOfLongestSubstring(const std::string& s) {
    std::unordered_map<char, int> lastSeen;
    int left = 0, best = 0;
    for (int right = 0; right < static_cast<int>(s.size()); right++) {
        char c = s[right];
        auto it = lastSeen.find(c);
        if (it != lastSeen.end() && it->second >= left) {
            left = it->second + 1;
        }
        lastSeen[c] = right;
        int window = right - left + 1;
        if (window > best) best = window;
    }
    return best;
}
