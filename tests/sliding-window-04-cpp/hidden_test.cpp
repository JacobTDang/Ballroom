#include <cassert>
#include <cstdio>
#include <string>

bool CheckInclusion(const std::string& s1, const std::string& s2);

int main() {
    assert(CheckInclusion("ab", "eidbaooo") == true);
    assert(CheckInclusion("ab", "eidboaoo") == false);
    assert(CheckInclusion("adc", "dcda") == true);
    assert(CheckInclusion("hello", "ooolleoooleh") == false);
    assert(CheckInclusion("a", "a") == true);
    assert(CheckInclusion("abc", "ab") == false);
    assert(CheckInclusion("aa", "ab") == false);
    assert(CheckInclusion("abcd", "dcba") == true);
    assert(CheckInclusion("ab", "a") == false);
    printf("all assertions passed\n");
    return 0;
}
