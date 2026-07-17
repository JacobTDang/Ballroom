#include "solution.cpp"

#include <cstdio>

int main() {
    struct Case {
        const char* pattern;
        const char* s;
        bool want;
    } cases[] = {
        {"abc", "abc", true},   {"abc", "abd", false},
        {"a?c", "abc", true},   {"a?c", "ac", false},
        {"?", "", false},       {"*", "", true},
        {"*", "anything", true},{"a*", "a", true},
        {"a*b*c", "aXXbYYc", true}, {"a*b*c", "aXXbYY", false},
        {"*.go", "main.go", true},  {"*.go", "main.gox", false},
        {"a*a", "aa", true},    {"a*a", "aba", true},
        {"a*a", "ab", false},   {"**", "x", true},
        {"[a-c]x", "bx", true}, {"[a-c]x", "dx", false},
        {"[xyz]", "y", true},   {"[xyz]", "w", false},
        {"file[0-9].txt", "file7.txt", true},
        {"file[0-9].txt", "fileX.txt", false},
        {"[abc", "a", false},   {"[", "[", false},
    };
    for (const auto& c : cases) {
        if (Match(c.pattern, c.s) != c.want) {
            fprintf(stderr, "Match(%s, %s) != %d\n", c.pattern, c.s, c.want);
            return 1;
        }
    }
    printf("all assertions passed\n");
    return 0;
}
