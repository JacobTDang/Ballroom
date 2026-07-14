#include <algorithm>
#include <cassert>
#include <cstdio>
#include <string>
#include <unordered_map>
#include <unordered_set>
#include <vector>

std::string AlienOrder(std::vector<std::string>& words);

// isValidAlienOrder checks the ORDERING PROPERTY rather than an exact
// string, since the topological order implied by words is not unique:
// every distinct character appearing in words must appear exactly once
// in order, and every adjacent word pair's first differing character
// must respect that order.
bool isValidAlienOrder(std::vector<std::string>& words, const std::string& order) {
    if (order.empty()) return false;

    std::unordered_map<char, int> pos;
    for (size_t i = 0; i < order.size(); i++) {
        if (pos.find(order[i]) != pos.end()) return false; // duplicate char
        pos[order[i]] = static_cast<int>(i);
    }

    std::unordered_set<char> seen;
    for (auto& w : words) {
        for (char c : w) seen.insert(c);
    }
    if (seen.size() != pos.size()) return false;
    for (char c : seen) {
        if (pos.find(c) == pos.end()) return false;
    }

    for (size_t i = 0; i + 1 < words.size(); i++) {
        std::string& w1 = words[i];
        std::string& w2 = words[i + 1];
        size_t minLen = std::min(w1.size(), w2.size());
        if (w1.size() > w2.size() && w1.substr(0, minLen) == w2.substr(0, minLen)) {
            return false;
        }
        for (size_t j = 0; j < minLen; j++) {
            if (w1[j] != w2[j]) {
                if (pos[w1[j]] >= pos[w2[j]]) return false;
                break;
            }
        }
    }
    return true;
}

void testValid() {
    std::vector<std::string> words = {"wrt", "wrf", "er", "ett", "rftt"};
    std::string order = AlienOrder(words);
    assert(isValidAlienOrder(words, order));
}

void testInvalidPrefix() {
    std::vector<std::string> words = {"abc", "ab"};
    assert(AlienOrder(words) == "");
}

void testCycle() {
    std::vector<std::string> words = {"z", "x", "z"};
    assert(AlienOrder(words) == "");
}

void testSingleWord() {
    std::vector<std::string> words = {"z"};
    assert(AlienOrder(words) == "z");
}

void testTwoDistinctChars() {
    std::vector<std::string> words = {"a", "b"};
    assert(AlienOrder(words) == "ab");
}

int main() {
    testValid();
    testInvalidPrefix();
    testCycle();
    testSingleWord();
    testTwoDistinctChars();
    std::printf("all tests passed\n");
    return 0;
}
