#include <cassert>
#include <cstdio>
#include <string>

int LengthOfLongestSubstring(const std::string& s);

int main() {
    assert(LengthOfLongestSubstring("abcabcbb") == 3);
    assert(LengthOfLongestSubstring("bbbbb") == 1);
    assert(LengthOfLongestSubstring("pwwkew") == 3);
    assert(LengthOfLongestSubstring("") == 0);
    assert(LengthOfLongestSubstring(" ") == 1);
    assert(LengthOfLongestSubstring("au") == 2);
    assert(LengthOfLongestSubstring("dvdf") == 3);
    assert(LengthOfLongestSubstring("abba") == 2);
    printf("all assertions passed\n");
    return 0;
}
