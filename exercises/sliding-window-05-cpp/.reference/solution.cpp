#include <string>
#include <unordered_map>

// MinWindow returns the shortest substring of s containing every
// character of t (with duplicates), or "" if no such substring exists.
std::string MinWindow(const std::string& s, const std::string& t) {
    if (t.empty() || s.size() < t.size()) return "";
    std::unordered_map<char, int> need, window;
    for (char c : t) need[c]++;
    int required = static_cast<int>(need.size());
    int have = 0;

    int bestLen = -1, bestStart = 0;
    int left = 0;
    for (int right = 0; right < static_cast<int>(s.size()); right++) {
        char c = s[right];
        window[c]++;
        auto it = need.find(c);
        if (it != need.end() && window[c] == it->second) have++;

        while (have == required) {
            if (bestLen == -1 || right - left + 1 < bestLen) {
                bestLen = right - left + 1;
                bestStart = left;
            }
            char leftChar = s[left];
            window[leftChar]--;
            auto lit = need.find(leftChar);
            if (lit != need.end() && window[leftChar] < lit->second) have--;
            left++;
        }
    }
    if (bestLen == -1) return "";
    return s.substr(bestStart, bestLen);
}
