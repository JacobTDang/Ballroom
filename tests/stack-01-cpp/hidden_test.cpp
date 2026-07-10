#include <cassert>
#include <cstdio>
#include <string>

bool IsValid(const std::string& s);

int main() {
    assert(IsValid("()") == true);
    assert(IsValid("()[]{}") == true);
    assert(IsValid("(]") == false);
    assert(IsValid("([)]") == false);
    assert(IsValid("{[]}") == true);
    assert(IsValid("(") == false);
    assert(IsValid("]") == false);
    printf("all assertions passed\n");
    return 0;
}
