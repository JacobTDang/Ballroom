#include <cassert>
#include <cstdio>

#include "solution.hpp"

int main() {
    LRUCache cache(2);
    cache.put(1, 100);
    cache.put(2, 200);
    assert(cache.get(1) == 100);   // access 1 -> 1 is now most recently used

    cache.put(3, 300);             // evicts 2 (least recently used)
    assert(cache.get(2) == -1);
    assert(cache.get(3) == 300);

    cache.put(4, 400);             // evicts 1 (3 and then 4 are more recent)
    assert(cache.get(1) == -1);
    assert(cache.get(3) == 300);
    assert(cache.get(4) == 400);

    {
        LRUCache c(2);
        c.put(1, 1);
        c.put(2, 2);
        c.put(1, 10);              // update, not a new insert -- must not evict 2
        assert(c.get(2) == 2);
        assert(c.get(1) == 10);
    }
    {
        LRUCache c(2);
        c.put(1, 1);
        c.put(2, 2);
        c.get(1);                  // 1 is now most recently used
        c.put(3, 3);                // should evict 2, not 1
        assert(c.get(2) == -1);
        assert(c.get(1) == 1);
    }
    {
        LRUCache c(1);
        c.put(1, 1);
        c.put(2, 2);
        assert(c.get(1) == -1);
        assert(c.get(2) == 2);
    }
    {
        LRUCache c(2);
        assert(c.get(999) == -1);
    }

    printf("all assertions passed\n");
    return 0;
}
