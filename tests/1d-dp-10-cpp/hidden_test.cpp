#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

bool WordBreak(std::string s, std::vector<std::string>& wordDict);

void testClassic() {
    std::vector<std::string> wordDict = {"leet", "code"};
    assert(WordBreak("leetcode", wordDict) == true);
}

void testReusedWord() {
    std::vector<std::string> wordDict = {"apple", "pen"};
    assert(WordBreak("applepenapple", wordDict) == true);
}

void testImpossible() {
    std::vector<std::string> wordDict = {"cats", "dog", "sand", "and", "cat"};
    assert(WordBreak("catsandog", wordDict) == false);
}

void testSingleWord() {
    std::vector<std::string> wordDict = {"a"};
    assert(WordBreak("a", wordDict) == true);
}

void testLeftoverCharUnmatched() {
    std::vector<std::string> wordDict = {"a"};
    assert(WordBreak("ab", wordDict) == false);
}

void testTrailingCharNeverMatches() {
    std::vector<std::string> wordDict = {"a", "aa"};
    assert(WordBreak("aaaaaaaaaaaaaaaaaaaab", wordDict) == false);
}

void testMultiplePaths() {
    std::vector<std::string> wordDict = {"car", "ca", "rs"};
    assert(WordBreak("cars", wordDict) == true);
}

void testSimpleConcatenation() {
    std::vector<std::string> wordDict = {"goal", "special"};
    assert(WordBreak("goalspecial", wordDict) == true);
}

int main() {
    testClassic();
    testReusedWord();
    testImpossible();
    testSingleWord();
    testLeftoverCharUnmatched();
    testTrailingCharNeverMatches();
    testMultiplePaths();
    testSimpleConcatenation();
    std::printf("all tests passed\n");
    return 0;
}
