#include <string>
#include <vector>

// BloomFilter: a bit array + k hashes. Add sets k bits; MightContain
// checks them -- "definitely not" or "probably yes".
//
// TODO: one weak hash into a fixed 64-slot table, ignoring both
// parameters -- no false negatives, but the table saturates instantly
// and almost every absent key collides.
class BloomFilter {
public:
    BloomFilter(int bits, int hashes) : table_(64, false) {}

    void Add(const std::string& key) {
        table_[Hash(key)] = true;
    }

    bool MightContain(const std::string& key) {
        return table_[Hash(key)];
    }

private:
    int Hash(const std::string& key) {
        int h = 0;
        for (char c : key) h += c;
        return h % 64;
    }

    std::vector<bool> table_;
};
