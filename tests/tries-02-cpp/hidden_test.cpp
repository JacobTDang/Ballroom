#include <cassert>
#include <cstdio>

#include "solution.hpp"

int main() {
    {
        WordDictionary d;
        d.addWord("bad");
        d.addWord("dad");
        d.addWord("mad");
        assert(d.search("pad") == false);
        assert(d.search("bad") == true);
        assert(d.search(".ad") == true);
        assert(d.search("b..") == true);
        assert(d.search("...") == true);
        assert(d.search("....") == false);
        assert(d.search("..d") == true);
        assert(d.search("dab") == false);
    }
    {
        WordDictionary d;
        assert(d.search("a") == false);
        assert(d.search(".") == false);
    }
    printf("all assertions passed\n");
    return 0;
}
