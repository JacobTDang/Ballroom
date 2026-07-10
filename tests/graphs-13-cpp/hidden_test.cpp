#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

int LadderLength(std::string beginWord, std::string endWord, std::vector<std::string>& wordList);

void testClassic() {
    std::vector<std::string> wordList = {"hot", "dot", "dog", "lot", "log", "cog"};
    assert(LadderLength("hit", "cog", wordList) == 5);
}

void testEndWordNotInList() {
    std::vector<std::string> wordList = {"hot", "dot", "dog", "lot", "log"};
    assert(LadderLength("hit", "cog", wordList) == 0);
}

void testDirectNeighbor() {
    std::vector<std::string> wordList = {"hot"};
    assert(LadderLength("hit", "hot", wordList) == 2);
}

int main() {
    testClassic();
    testEndWordNotInList();
    testDirectNeighbor();
    std::printf("all tests passed\n");
    return 0;
}
