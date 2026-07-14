#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::string encode(const std::vector<std::string>& strs);
std::vector<std::string> decode(const std::string& s);

int main() {
    std::vector<std::vector<std::string>> cases = {
        {"neet", "code", "love", "you"},
        {},
        {""},
        {"", "", ""},
        {"a#b", "c##d", "5#hello"},
        {"hello world", "foo,bar", "123"},
        {"4#abcd", "hello"},
        {"#####"},
        {"xyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxyxy"},
        {"123", "456", "0"},
        {"", "a", "", "b"},
    };
    for (const auto& strs : cases) {
        std::string encoded = encode(strs);
        auto got = decode(encoded);
        assert(got == strs);
    }
    printf("all assertions passed\n");
    return 0;
}
