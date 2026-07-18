#include <string>
#include <vector>

struct Record {
    std::string key;
    int value;

    bool operator==(const Record& other) const {
        return key == other.key && value == other.value;
    }
};

// Removes duplicate records -- two records with the same key and
// value are duplicates. Currently keeps both -- find and fix the bug.
std::vector<Record*> dedupe(const std::vector<Record*>& records) {
    std::vector<Record*> result;
    for (Record* r : records) {
        bool found = false;
        for (Record* seen : result) {
            if (seen == r) {
                found = true;
                break;
            }
        }
        if (!found) result.push_back(r);
    }
    return result;
}
