#include <algorithm>
#include <cstdint>
#include <map>
#include <string>
#include <vector>

// Ring: each node occupies `vnodes` positions (virtual nodes smooth
// the balance); a key belongs to the first position clockwise from its
// hash -- lower_bound over the sorted positions, wrapping to the
// front. Removing a node deletes only its own positions, so every
// other key keeps its owner.
class Ring {
public:
    explicit Ring(int vnodes) : vnodes_(vnodes) {}

    void AddNode(const std::string& name) {
        for (int i = 0; i < vnodes_; i++) {
            uint64_t p = Hash(name + "#" + std::to_string(i));
            if (owner_.count(p)) continue; // vanishing collision odds
            owner_[p] = name;
            positions_.push_back(p);
        }
        std::sort(positions_.begin(), positions_.end());
    }

    void RemoveNode(const std::string& name) {
        positions_.erase(
            std::remove_if(positions_.begin(), positions_.end(),
                           [&](uint64_t p) { return owner_[p] == name; }),
            positions_.end());
        for (auto it = owner_.begin(); it != owner_.end();) {
            if (it->second == name) it = owner_.erase(it);
            else ++it;
        }
    }

    std::string Lookup(const std::string& key) {
        if (positions_.empty()) return "";
        uint64_t h = Hash(key);
        auto it = std::lower_bound(positions_.begin(), positions_.end(), h);
        if (it == positions_.end()) it = positions_.begin(); // wrap
        return owner_[*it];
    }

private:
    static uint64_t Hash(const std::string& s) {
        uint64_t h = 0xCBF29CE484222325ULL;
        for (unsigned char c : s) {
            h ^= c;
            h *= 0x100000001B3ULL;
        }
        return h;
    }

    int vnodes_;
    std::vector<uint64_t> positions_;
    std::map<uint64_t, std::string> owner_;
};
