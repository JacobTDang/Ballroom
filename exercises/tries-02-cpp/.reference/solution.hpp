#pragma once

#include <string>
#include <unordered_map>

// WordDictionary supports adding words and searching, where a search
// query may use '.' to match any single character.
class WordDictionary {
public:
    void addWord(const std::string& word) {
        WordDictionary* node = this;
        for (char c : word) {
            if (node->children_.find(c) == node->children_.end()) {
                node->children_[c] = new WordDictionary();
            }
            node = node->children_[c];
        }
        node->isEnd_ = true;
    }

    bool search(const std::string& word) {
        return searchFrom(word, 0);
    }

private:
    std::unordered_map<char, WordDictionary*> children_;
    bool isEnd_ = false;

    bool searchFrom(const std::string& word, size_t idx) {
        WordDictionary* node = this;
        for (size_t i = idx; i < word.size(); i++) {
            char c = word[i];
            if (c == '.') {
                for (auto& [ch, child] : node->children_) {
                    if (child->searchFrom(word, i + 1)) return true;
                }
                return false;
            }
            auto it = node->children_.find(c);
            if (it == node->children_.end()) return false;
            node = it->second;
        }
        return node->isEnd_;
    }
};
