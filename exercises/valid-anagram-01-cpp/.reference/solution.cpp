#include <string>
#include <unordered_map>

// Returns true if t is an anagram of s.
bool is_anagram(const std::string& s, const std::string& t) {
    if (s.size() != t.size()) return false;
    std::unordered_map<char, int> counts;
    for (char c : s) counts[c]++;
    for (char c : t) {
        if (--counts[c] < 0) return false;
    }
    return true;
}
