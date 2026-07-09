#pragma once

#include <list>
#include <unordered_map>

class LRUCache {
public:
    explicit LRUCache(int capacity) : capacity_(capacity) {}

    // Return the value for key, or -1 if not present. Accessing a key
    // marks it as most recently used.
    int get(int key) {
        auto it = index_.find(key);
        if (it == index_.end()) {
            return -1;
        }
        order_.splice(order_.begin(), order_, it->second);
        return it->second->second;
    }

    // Insert or update key with value, marking it most recently used. If
    // inserting a new key would exceed capacity, evict the least recently
    // used entry first.
    void put(int key, int value) {
        auto it = index_.find(key);
        if (it != index_.end()) {
            it->second->second = value;
            order_.splice(order_.begin(), order_, it->second);
            return;
        }
        if (static_cast<int>(index_.size()) >= capacity_) {
            auto& lru = order_.back();
            index_.erase(lru.first);
            order_.pop_back();
        }
        order_.emplace_front(key, value);
        index_[key] = order_.begin();
    }

private:
    int capacity_;
    // Front = most recently used, back = least recently used.
    std::list<std::pair<int, int>> order_;
    std::unordered_map<int, std::list<std::pair<int, int>>::iterator> index_;
};
