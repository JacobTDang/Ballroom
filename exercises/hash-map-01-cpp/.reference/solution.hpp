#pragma once

#include <utility>
#include <vector>

// Hash map for non-negative integer keys and values, built without
// std::unordered_map / std::map.
class MyHashMap {
public:
    MyHashMap() : buckets_(kBuckets) {}

    // Insert key with value, or update it if key already exists.
    void put(int key, int value) {
        auto& bucket = buckets_[key % kBuckets];
        for (auto& kv : bucket) {
            if (kv.first == key) {
                kv.second = value;
                return;
            }
        }
        bucket.emplace_back(key, value);
    }

    // Return the value for key, or -1 if key is absent.
    int get(int key) {
        for (const auto& kv : buckets_[key % kBuckets]) {
            if (kv.first == key) return kv.second;
        }
        return -1;
    }

    // Delete key if present; do nothing otherwise.
    void remove(int key) {
        auto& bucket = buckets_[key % kBuckets];
        for (size_t i = 0; i < bucket.size(); i++) {
            if (bucket[i].first == key) {
                bucket.erase(bucket.begin() + i);
                return;
            }
        }
    }

private:
    static const int kBuckets = 1024;
    std::vector<std::vector<std::pair<int, int>>> buckets_;
};
