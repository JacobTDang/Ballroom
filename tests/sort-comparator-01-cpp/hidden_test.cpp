#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

struct Entry {
    std::string name;
    int score;
};

std::vector<Entry> SortLeaderboard(std::vector<Entry> entries);

static bool names_match(const std::vector<Entry>& got, const std::vector<std::string>& want) {
    if (got.size() != want.size()) return false;
    for (size_t i = 0; i < got.size(); i++) {
        if (got[i].name != want[i]) return false;
    }
    return true;
}

int main() {
    {
        std::vector<Entry> entries = {{"bob", 60}, {"dan", 100}, {"cara", 75}, {"amy", 90}};
        assert(names_match(SortLeaderboard(entries), {"dan", "amy", "cara", "bob"}));
    }
    {
        std::vector<Entry> entries = {{"erin", 90}, {"cara", 75}, {"amy", 90}, {"bob", 90}};
        assert(names_match(SortLeaderboard(entries), {"amy", "bob", "erin", "cara"}));
    }
    {
        std::vector<Entry> entries = {{"zoe", 50}, {"amy", 80}, {"erin", 65}, {"dan", 50}, {"cara", 80}, {"bob", 80}};
        assert(names_match(SortLeaderboard(entries), {"amy", "bob", "cara", "erin", "dan", "zoe"}));
    }
    {
        std::vector<Entry> entries = {{"zed", 10}, {"amy", 10}, {"mno", 10}};
        assert(names_match(SortLeaderboard(entries), {"amy", "mno", "zed"}));
    }
    {
        std::vector<Entry> entries = {{"cara", 10}, {"bob", -5}, {"amy", -5}};
        assert(names_match(SortLeaderboard(entries), {"cara", "amy", "bob"}));
    }
    printf("all assertions passed\n");
    return 0;
}
