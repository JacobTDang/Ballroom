#include <algorithm>
#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::vector<std::string> LetterCombinations(std::string digits);

void check(std::string digits, std::vector<std::string> want) {
    auto got = LetterCombinations(digits);
    std::sort(got.begin(), got.end());
    std::sort(want.begin(), want.end());
    assert(got == want);
}

int main() {
    check("23", {"ad", "ae", "af", "bd", "be", "bf", "cd", "ce", "cf"});
    check("", {});
    check("2", {"a", "b", "c"});
    printf("all assertions passed\n");
    return 0;
}
