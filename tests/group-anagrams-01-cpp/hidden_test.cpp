#include <algorithm>
#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::vector<std::vector<std::string>> group_anagrams(const std::vector<std::string>& strs);

// Sorts each group's strings and then sorts the list of groups, so
// results can be compared regardless of ordering — any correct grouping
// is valid no matter which order it comes back in.
std::vector<std::vector<std::string>> normalize(std::vector<std::vector<std::string>> groups) {
    for (auto& g : groups) {
        std::sort(g.begin(), g.end());
    }
    std::sort(groups.begin(), groups.end());
    return groups;
}

int main() {
    {
        auto got = normalize(group_anagrams({"eat", "tea", "tan", "ate", "nat", "bat"}));
        auto want = normalize({{"bat"}, {"nat", "tan"}, {"ate", "eat", "tea"}});
        assert(got == want);
    }
    {
        auto got = normalize(group_anagrams({""}));
        auto want = normalize({{""}});
        assert(got == want);
    }
    {
        auto got = normalize(group_anagrams({"a"}));
        auto want = normalize({{"a"}});
        assert(got == want);
    }
    {
        auto got = normalize(group_anagrams({"abc", "bca", "cab", "xyz"}));
        auto want = normalize({{"abc", "bca", "cab"}, {"xyz"}});
        assert(got == want);
    }
    {
        auto got = normalize(group_anagrams({"cat", "dog", "bird"}));
        auto want = normalize({{"cat"}, {"dog"}, {"bird"}});
        assert(got == want);
    }
    {
        auto got = normalize(group_anagrams({"abc", "bca", "cab", "acb"}));
        auto want = normalize({{"abc", "bca", "cab", "acb"}});
        assert(got == want);
    }
    {
        auto got = normalize(group_anagrams({"", "", ""}));
        auto want = normalize({{"", "", ""}});
        assert(got == want);
    }
    {
        auto got = normalize(group_anagrams({"bat", "tab", "cat", "act", "dog", "god", "xyz"}));
        auto want = normalize({{"bat", "tab"}, {"cat", "act"}, {"dog", "god"}, {"xyz"}});
        assert(got == want);
    }
    {
        auto got = normalize(group_anagrams({"a", "b", "a", "c", "b"}));
        auto want = normalize({{"a", "a"}, {"b", "b"}, {"c"}});
        assert(got == want);
    }
    printf("all assertions passed\n");
    return 0;
}
