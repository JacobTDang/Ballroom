#pragma once

#include <string>
#include <unordered_map>

// Trie is a prefix tree over lowercase English letters.
class Trie {
public:
    void insert(const std::string& word) {
        // TODO: implement
    }

    bool search(const std::string& word) {
        // TODO: implement
        return false;
    }

    bool startsWith(const std::string& prefix) {
        // TODO: implement
        return false;
    }

private:
    std::unordered_map<char, Trie*> children_;
    bool isEnd_ = false;
};
