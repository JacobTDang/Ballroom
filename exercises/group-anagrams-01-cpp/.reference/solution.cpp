#include <array>
#include <string>
#include <unordered_map>
#include <vector>

// Groups the strings in strs into vectors of anagrams of each other, in
// any order (both between groups and within a group).
std::vector<std::vector<std::string>> group_anagrams(const std::vector<std::string>& strs) {
    std::unordered_map<std::string, std::vector<std::string>> groups;
    for (const auto& s : strs) {
        std::array<int, 26> counts{};
        for (char c : s) counts[c - 'a']++;
        std::string key;
        key.reserve(26);
        for (int c : counts) key.push_back(static_cast<char>(c));
        groups[key].push_back(s);
    }
    std::vector<std::vector<std::string>> result;
    result.reserve(groups.size());
    for (auto& [key, group] : groups) {
        result.push_back(std::move(group));
    }
    return result;
}
