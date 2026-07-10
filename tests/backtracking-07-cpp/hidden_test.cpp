#include <algorithm>
#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::vector<std::vector<std::string>> Partition(std::string s);

std::vector<std::vector<std::string>> normalize(std::vector<std::vector<std::string>> lists) {
    std::sort(lists.begin(), lists.end());
    return lists;
}

void check(std::string s, std::vector<std::vector<std::string>> want) {
    auto got = normalize(Partition(s));
    assert(got == normalize(want));
}

int main() {
    check("aab", {{"a", "a", "b"}, {"aa", "b"}});
    check("a", {{"a"}});
    check("aba", {{"a", "b", "a"}, {"aba"}});
    printf("all assertions passed\n");
    return 0;
}
