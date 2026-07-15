#include <cassert>
#include <cstdio>

#include "solution.hpp"

int main() {
    {
        MyHashMap m;
        m.put(1, 100);
        assert(m.get(1) == 100);                 // put/get roundtrip
    }
    {
        MyHashMap m;
        assert(m.get(42) == -1);                 // missing key
    }
    {
        MyHashMap m;
        m.put(1, 100);
        m.put(1, 200);
        assert(m.get(1) == 200);                 // overwrite
    }
    {
        MyHashMap m;
        m.put(7, 70);
        m.remove(7);
        assert(m.get(7) == -1);                  // remove then get
    }
    {
        MyHashMap m;
        m.put(1, 10);
        m.remove(99);
        assert(m.get(1) == 10);                  // remove missing = no-op
    }
    {
        MyHashMap m;                             // colliding keys chain
        int keys[] = {1, 1025, 2049, 1001, 2001};
        for (int k : keys) m.put(k, k * 3);
        for (int k : keys) assert(m.get(k) == k * 3);
    }
    {
        MyHashMap m;                             // remove one of a chain
        m.put(1, 11);
        m.put(1025, 22);
        m.put(2049, 33);
        m.remove(1025);
        assert(m.get(1) == 11);
        assert(m.get(1025) == -1);
        assert(m.get(2049) == 33);
    }
    {
        MyHashMap m;
        m.put(0, 0);
        assert(m.get(0) == 0);                   // zero key and value
    }
    {
        MyHashMap m;
        m.put(1000000, 999);
        assert(m.get(1000000) == 999);           // upper bound key
    }
    {
        MyHashMap m;                             // many keys, no aliasing
        for (int k = 0; k < 500; k++) m.put(k, k * 2);
        for (int k = 0; k < 500; k++) assert(m.get(k) == k * 2);
    }
    {
        MyHashMap m;                             // interleaved sequence
        m.put(5, 1);
        m.put(6, 2);
        m.remove(5);
        m.put(6, 3);
        m.put(5, 4);
        assert(m.get(5) == 4);
        assert(m.get(6) == 3);
    }
    printf("all assertions passed\n");
    return 0;
}
