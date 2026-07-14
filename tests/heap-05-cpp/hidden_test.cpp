#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

int LeastInterval(std::vector<char>& tasks, int n);

std::vector<char> toChars(const std::string& s) {
    return std::vector<char>(s.begin(), s.end());
}

int main() {
    auto t1 = toChars("AAABBB");
    assert(LeastInterval(t1, 2) == 8);
    auto t2 = toChars("AAABBB");
    assert(LeastInterval(t2, 0) == 6);
    auto t3 = toChars("AAAAAABCDEFG");
    assert(LeastInterval(t3, 2) == 16);
    auto t4 = toChars("A");
    assert(LeastInterval(t4, 5) == 1);
    auto t5 = toChars("AAAB");
    assert(LeastInterval(t5, 3) == 9);
    auto t6 = toChars("AB");
    assert(LeastInterval(t6, 2) == 2);
    printf("all assertions passed\n");
    return 0;
}
