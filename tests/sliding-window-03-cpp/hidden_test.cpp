#include <cassert>
#include <cstdio>
#include <string>

int CharacterReplacement(const std::string& s, int k);

int main() {
    assert(CharacterReplacement("ABAB", 2) == 4);
    assert(CharacterReplacement("AABABBA", 1) == 4);
    assert(CharacterReplacement("ABCDE", 1) == 2);
    assert(CharacterReplacement("AAAA", 0) == 4);
    assert(CharacterReplacement("A", 0) == 1);
    assert(CharacterReplacement("ABBB", 2) == 4);
    assert(CharacterReplacement("", 2) == 0);
    assert(CharacterReplacement("AAAA", 4) == 4);
    assert(CharacterReplacement("ABABABAB", 3) == 7);
    printf("all assertions passed\n");
    return 0;
}
