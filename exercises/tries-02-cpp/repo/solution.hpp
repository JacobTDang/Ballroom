#pragma once

#include <string>
#include <unordered_map>

// WordDictionary supports adding words and searching, where a search
// query may use '.' to match any single character.
class WordDictionary {
public:
    void addWord(const std::string& word) {
        // TODO: implement
    }

    bool search(const std::string& word) {
        // TODO: implement
        return false;
    }

private:
    std::unordered_map<char, WordDictionary*> children_;
    bool isEnd_ = false;
};
