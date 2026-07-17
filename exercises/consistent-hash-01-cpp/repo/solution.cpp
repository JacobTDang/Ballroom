#include <algorithm>
#include <string>
#include <vector>

// Ring maps keys to nodes. Adding/removing a node should only remap
// the keys in its neighborhood.
//
// TODO: hash(key) % len(nodes) remaps almost EVERY key whenever the
// node count changes -- the exact failure consistent hashing exists to
// fix. (vnodes is ignored here too.)
class Ring {
public:
    explicit Ring(int vnodes) {}

    void AddNode(const std::string& name) {
        nodes_.push_back(name);
        std::sort(nodes_.begin(), nodes_.end());
    }

    void RemoveNode(const std::string& name) {
        nodes_.erase(std::remove(nodes_.begin(), nodes_.end(), name), nodes_.end());
    }

    std::string Lookup(const std::string& key) {
        if (nodes_.empty()) return "";
        return nodes_[Hash(key) % nodes_.size()];
    }

private:
    static size_t Hash(const std::string& s) {
        uint64_t h = 0xCBF29CE484222325ULL;
        for (unsigned char c : s) {
            h ^= c;
            h *= 0x100000001B3ULL;
        }
        return (size_t)h;
    }

    std::vector<std::string> nodes_;
};
