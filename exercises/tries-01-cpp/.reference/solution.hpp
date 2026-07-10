#pragma once

#include <string>
#include <unordered_map>

// Trie is a prefix tree over lowercase English letters.
class Trie {
public:
    void insert(const std::string& word) {
        Trie* node = this;
        for (char c : word) {
            if (node->children_.find(c) == node->children_.end()) {
                node->children_[c] = new Trie();
            }
            node = node->children_[c];
        }
        node->isEnd_ = true;
    }

    bool search(const std::string& word) {
        Trie* node = find(word);
        return node != nullptr && node->isEnd_;
    }

    bool startsWith(const std::string& prefix) {
        return find(prefix) != nullptr;
    }

private:
    std::unordered_map<char, Trie*> children_;
    bool isEnd_ = false;

    Trie* find(const std::string& s) {
        Trie* node = this;
        for (char c : s) {
            auto it = node->children_.find(c);
            if (it == node->children_.end()) return nullptr;
            node = it->second;
        }
        return node;
    }
};
