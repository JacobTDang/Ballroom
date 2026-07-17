#include <cstdint>
#include <string>
#include <unordered_map>

// Each entry remembers its write time (expiry never moves) and a
// monotonic touch sequence (recency does). Puts at capacity purge
// expired corpses first, then evict the least recently used live
// entry. O(n) scans are fine at exercise scale; a production cache
// pairs the map with a linked list.
class TTLCache {
public:
    TTLCache(int capacity, long ttl_ms) : capacity_(capacity), ttl_ms_(ttl_ms) {}

    void PutAt(const std::string& key, int value, long now_ms) {
        seq_++;
        auto it = items_.find(key);
        if (it != items_.end()) {
            it->second = {value, now_ms, seq_};
            return;
        }
        for (auto e = items_.begin(); e != items_.end();) {
            if (Expired(e->second, now_ms)) e = items_.erase(e);
            else ++e;
        }
        if ((int)items_.size() >= capacity_) {
            auto lru = items_.begin();
            for (auto e = items_.begin(); e != items_.end(); ++e) {
                if (e->second.touched < lru->second.touched) lru = e;
            }
            items_.erase(lru);
        }
        items_[key] = {value, now_ms, seq_};
    }

    bool GetAt(const std::string& key, long now_ms, int* out) {
        auto it = items_.find(key);
        if (it == items_.end() || Expired(it->second, now_ms)) return false;
        seq_++;
        it->second.touched = seq_; // recency refreshed; written_at not
        *out = it->second.value;
        return true;
    }

private:
    struct Entry {
        int value;
        long written_at;
        int64_t touched;
    };

    bool Expired(const Entry& e, long now_ms) const {
        return now_ms - e.written_at >= ttl_ms_;
    }

    int capacity_;
    long ttl_ms_;
    std::unordered_map<std::string, Entry> items_;
    int64_t seq_ = 0;
};
