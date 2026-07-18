#include <algorithm>
#include <string>
#include <vector>

struct Entry {
    std::string name;
    int score;
};

// Returns entries sorted by score descending; ties break by name
// ascending.
std::vector<Entry> SortLeaderboard(std::vector<Entry> entries) {
    std::sort(entries.begin(), entries.end(), [](const Entry& a, const Entry& b) {
        if (a.score != b.score) return a.score > b.score;
        return a.name < b.name;
    });
    return entries;
}
