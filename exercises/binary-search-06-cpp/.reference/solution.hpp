#pragma once

#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

// TimeMap stores multiple values per key, each tagged with the
// timestamp it was set at.
class TimeMap {
public:
    // set appends value for key — timestamps arrive strictly
    // increasing, so the per-key vector stays sorted without needing
    // to insert.
    void set(const std::string& key, const std::string& value, int timestamp) {
        store_[key].emplace_back(timestamp, value);
    }

    // get binary-searches for the entry with the largest timestamp <=
    // the query timestamp.
    std::string get(const std::string& key, int timestamp) {
        auto it = store_.find(key);
        if (it == store_.end()) return "";
        const auto& entries = it->second;
        int lo = 0, hi = static_cast<int>(entries.size()) - 1;
        std::string res;
        while (lo <= hi) {
            int mid = lo + (hi - lo) / 2;
            if (entries[mid].first <= timestamp) {
                res = entries[mid].second;
                lo = mid + 1;
            } else {
                hi = mid - 1;
            }
        }
        return res;
    }

private:
    std::unordered_map<std::string, std::vector<std::pair<int, std::string>>> store_;
};
