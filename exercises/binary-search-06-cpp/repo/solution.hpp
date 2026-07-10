#pragma once

#include <string>

// TimeMap stores multiple values per key, each tagged with the
// timestamp it was set at.
class TimeMap {
public:
    void set(const std::string& key, const std::string& value, int timestamp) {
        // TODO: implement
    }

    std::string get(const std::string& key, int timestamp) {
        // TODO: implement
        return "";
    }
};
