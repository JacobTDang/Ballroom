#include <string>
#include <unordered_map>

// TTLCache: LRU eviction at capacity, plus per-entry expiry ttl_ms
// after the write. GetAt returns false on miss/expired.
//
// TODO: a plain map -- no eviction, no expiry, no recency. Every rule
// in the problem statement is still yours to build.
class TTLCache {
public:
    TTLCache(int capacity, long ttl_ms) : capacity_(capacity), ttl_ms_(ttl_ms) {}

    void PutAt(const std::string& key, int value, long now_ms) {
        items_[key] = value;
    }

    bool GetAt(const std::string& key, long now_ms, int* out) {
        auto it = items_.find(key);
        if (it == items_.end()) return false;
        *out = it->second;
        return true;
    }

private:
    int capacity_;
    long ttl_ms_;
    std::unordered_map<std::string, int> items_;
};
