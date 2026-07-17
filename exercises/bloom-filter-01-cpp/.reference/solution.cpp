#include <cstdint>
#include <string>
#include <vector>

// Double hashing: two FNV-1a variants generate all k probe positions
// as h1 + i*h2 (Kirsch-Mitzenmacher) -- k hashes as good as k
// independent ones without computing k real hashes.
class BloomFilter {
public:
    BloomFilter(int bits, int hashes) : bits_(bits, false), hashes_(hashes) {}

    void Add(const std::string& key) {
        for (auto p : Positions(key)) bits_[p] = true;
    }

    bool MightContain(const std::string& key) {
        for (auto p : Positions(key)) {
            if (!bits_[p]) return false;
        }
        return true;
    }

private:
    static uint64_t Fnv1a(const std::string& data) {
        uint64_t h = 0xCBF29CE484222325ULL;
        for (unsigned char c : data) {
            h ^= c;
            h *= 0x100000001B3ULL;
        }
        return h;
    }

    // splitmix64 finalizer: FNV alone clusters on similar keys over
    // power-of-two table sizes; the avalanche makes probes independent.
    static uint64_t Mix(uint64_t h) {
        h ^= h >> 30;
        h *= 0xBF58476D1CE4E5B9ULL;
        h ^= h >> 27;
        h *= 0x94D049BB133111EBULL;
        h ^= h >> 31;
        return h;
    }

    std::vector<size_t> Positions(const std::string& key) {
        uint64_t h1 = Mix(Fnv1a(key));
        uint64_t h2 = Mix(h1 ^ 0x9E3779B97F4A7C15ULL) | 1;
        std::vector<size_t> out(hashes_);
        for (int i = 0; i < hashes_; i++) {
            out[i] = (h1 + (uint64_t)i * h2) % bits_.size();
        }
        return out;
    }

    std::vector<bool> bits_;
    int hashes_;
};
