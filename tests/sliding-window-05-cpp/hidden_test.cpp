#include <cassert>
#include <cstdio>
#include <string>

std::string MinWindow(const std::string& s, const std::string& t);

int main() {
    assert(MinWindow("ADOBECODEBANC", "ABC") == "BANC");
    assert(MinWindow("a", "a") == "a");
    assert(MinWindow("a", "aa") == "");
    assert(MinWindow("ab", "b") == "b");
    assert(MinWindow("bba", "ab") == "ba");
    printf("all assertions passed\n");
    return 0;
}
