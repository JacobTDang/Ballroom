#include <algorithm>
#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::vector<std::string> FindWords(std::vector<std::vector<char>>& board,
                                    std::vector<std::string>& words);

void check(std::vector<std::vector<char>> board, std::vector<std::string> words,
           std::vector<std::string> want) {
    auto got = FindWords(board, words);
    std::sort(got.begin(), got.end());
    std::sort(want.begin(), want.end());
    assert(got == want);
}

int main() {
    check({{'o', 'a', 'a', 'n'}, {'e', 't', 'a', 'e'}, {'i', 'h', 'k', 'r'}, {'i', 'f', 'l', 'v'}},
          {"oath", "pea", "eat", "rain"}, {"eat", "oath"});
    check({{'a', 'b'}, {'c', 'd'}}, {"abcb"}, {});
    check({{'a'}}, {"a"}, {"a"});
    check({{'a', 'a'}}, {"aaa"}, {});
    check({{'a', 'b'}, {'c', 'd'}}, {"abdc"}, {"abdc"});
    printf("all assertions passed\n");
    return 0;
}
