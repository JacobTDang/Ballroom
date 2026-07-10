#include <algorithm>
#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::vector<std::string> GenerateParenthesis(int n);

void check(int n, std::vector<std::string> want) {
    auto got = GenerateParenthesis(n);
    std::sort(got.begin(), got.end());
    std::sort(want.begin(), want.end());
    assert(got == want);
}

int main() {
    check(1, {"()"});
    check(2, {"(())", "()()"});
    check(3, {"((()))", "(()())", "(())()", "()(())", "()()()"});
    printf("all assertions passed\n");
    return 0;
}
