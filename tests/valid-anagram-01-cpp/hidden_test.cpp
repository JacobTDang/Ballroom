#include <cassert>
#include <cstdio>
#include <string>

bool is_anagram(const std::string& s, const std::string& t);

int main() {
    assert(is_anagram("anagram", "nagaram") == true);
    assert(is_anagram("rat", "car") == false);
    assert(is_anagram("ab", "a") == false);
    assert(is_anagram("aacc", "ccac") == false);
    assert(is_anagram("a", "a") == true);
    assert(is_anagram("aabbcc", "abcabc") == true);
    printf("all assertions passed\n");
    return 0;
}
