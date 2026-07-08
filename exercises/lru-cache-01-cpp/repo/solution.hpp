#pragma once

#include <list>
#include <unordered_map>

class LRUCache {
public:
    explicit LRUCache(int capacity) : capacity_(capacity) {}

    // Return the value for key, or -1 if not present. Accessing a key
    // marks it as most recently used.
    int get(int key) {
        // TODO: implement
        return -1;
    }

    // Insert or update key with value, marking it most recently used. If
    // inserting a new key would exceed capacity, evict the least recently
    // used entry first.
    void put(int key, int value) {
        // TODO: implement
    }

private:
    int capacity_;
};
