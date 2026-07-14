#include <cassert>
#include <cstdio>

#include "solution.hpp"

int main() {
    {
        Trie trie;
        trie.insert("apple");
        assert(trie.search("apple") == true);
        assert(trie.search("app") == false);
        assert(trie.startsWith("app") == true);
        trie.insert("app");
        assert(trie.search("app") == true);
    }
    {
        Trie trie;
        trie.insert("banana");
        assert(trie.startsWith("ban") == true);
        assert(trie.search("ban") == false);
    }
    {
        Trie trie;
        assert(trie.search("anything") == false);
        assert(trie.startsWith("a") == false);
    }
    {
        Trie trie;
        trie.insert("app");
        trie.insert("apple");
        trie.insert("application");
        assert(trie.search("app") == true);
        assert(trie.search("apple") == true);
        assert(trie.search("appl") == false);
        assert(trie.startsWith("appl") == true);
    }
    printf("all assertions passed\n");
    return 0;
}
