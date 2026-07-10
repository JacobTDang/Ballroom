#include <cassert>
#include <cstdio>
#include <string>

bool IsPalindrome(const std::string& s);

int main() {
    assert(IsPalindrome("A man, a plan, a canal: Panama") == true);
    assert(IsPalindrome("race a car") == false);
    assert(IsPalindrome(" ") == true);
    assert(IsPalindrome("0P") == false);
    assert(IsPalindrome("Was it a car or a cat I saw?") == true);
    assert(IsPalindrome(".,") == true);
    assert(IsPalindrome("a_b") == false);
    printf("all assertions passed\n");
    return 0;
}
