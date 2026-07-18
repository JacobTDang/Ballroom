#include <algorithm>
#include <cstdlib>
#include <map>
#include <sstream>
#include <string>
#include <vector>

// Record is a single row in the store, ordered by id.
struct Record {
    long id;
    std::string name;
};

namespace {

constexpr int kDefaultPageSize = 10;
constexpr int kMaxPageSize = 50;

const std::string kB64Chars =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    "abcdefghijklmnopqrstuvwxyz"
    "0123456789-_";

const std::vector<int>& B64Table() {
    static std::vector<int> table = [] {
        std::vector<int> t(256, -1);
        for (size_t i = 0; i < kB64Chars.size(); i++) t[(unsigned char)kB64Chars[i]] = (int)i;
        return t;
    }();
    return table;
}

std::string Base64Encode(const std::string& input) {
    std::string out;
    int val = 0, bits = -6;
    for (unsigned char c : input) {
        val = (val << 8) + c;
        bits += 8;
        while (bits >= 0) {
            out.push_back(kB64Chars[(val >> bits) & 0x3F]);
            bits -= 6;
        }
    }
    if (bits > -6) out.push_back(kB64Chars[((val << 8) >> (bits + 8)) & 0x3F]);
    return out;
}

bool Base64Decode(const std::string& input, std::string* out) {
    const std::vector<int>& table = B64Table();
    int val = 0, bits = -8;
    out->clear();
    for (unsigned char c : input) {
        if (table[c] == -1) return false;
        val = (val << 6) + table[c];
        bits += 6;
        if (bits >= 0) {
            out->push_back(char((val >> bits) & 0xFF));
            bits -= 8;
        }
    }
    return true;
}

std::string Checksum(const std::string& payload) {
    long long h = 0;
    for (unsigned char c : payload) h = (h * 131 + c) % 1000000007LL;
    std::ostringstream os;
    os << std::hex << h;
    return os.str();
}

int ClampPageSize(int page_size) {
    if (page_size <= 0) return kDefaultPageSize;
    if (page_size > kMaxPageSize) return kMaxPageSize;
    return page_size;
}

std::string EncodeToken(long last_id, int page_size) {
    std::ostringstream payload;
    payload << last_id << ":" << page_size;
    std::string raw = payload.str() + ":" + Checksum(payload.str());
    return Base64Encode(raw);
}

std::vector<std::string> SplitColon(const std::string& s) {
    std::vector<std::string> parts;
    std::string cur;
    for (char c : s) {
        if (c == ':') {
            parts.push_back(cur);
            cur.clear();
        } else {
            cur.push_back(c);
        }
    }
    parts.push_back(cur);
    return parts;
}

bool DecodeToken(const std::string& token, long* last_id, int* page_size, std::string* err) {
    std::string raw;
    if (!Base64Decode(token, &raw)) {
        *err = "invalid page token";
        return false;
    }
    std::vector<std::string> parts = SplitColon(raw);
    if (parts.size() != 3) {
        *err = "invalid page token";
        return false;
    }
    std::string payload = parts[0] + ":" + parts[1];
    if (Checksum(payload) != parts[2]) {
        *err = "invalid page token: checksum mismatch";
        return false;
    }
    char* endptr = nullptr;
    long id = std::strtol(parts[0].c_str(), &endptr, 10);
    if (endptr == parts[0].c_str() || *endptr != '\0') {
        *err = "invalid page token";
        return false;
    }
    int size = (int)std::strtol(parts[1].c_str(), &endptr, 10);
    if (endptr == parts[1].c_str() || *endptr != '\0') {
        *err = "invalid page token";
        return false;
    }
    *last_id = id;
    *page_size = size;
    return true;
}

}  // namespace

// CursorStore: keyset pagination. The token names the last id seen
// (plus the page_size it was issued for), never a raw offset -- so a
// walk can't be thrown off by inserts, and resuming with a different
// page size is a loud error rather than a silent behavior change.
class CursorStore {
public:
    explicit CursorStore(std::vector<Record> records) {
        for (auto& r : records) records_[r.id] = r;
    }

    bool Insert(const Record& r, std::string* err) {
        if (records_.count(r.id)) {
            *err = "duplicate id";
            return false;
        }
        records_[r.id] = r;
        return true;
    }

    bool List(int page_size, const std::string& page_token,
              std::vector<Record>* items, std::string* next_page_token,
              std::string* err) {
        int effective = ClampPageSize(page_size);

        bool has_cursor = false;
        long cursor_id = 0;
        if (!page_token.empty()) {
            int encoded_size = 0;
            if (!DecodeToken(page_token, &cursor_id, &encoded_size, err)) return false;
            if (encoded_size != effective) {
                *err = "page_token was issued for a different page_size";
                return false;
            }
            has_cursor = true;
        }

        std::vector<Record> candidates;
        candidates.reserve(records_.size());
        for (auto& [id, r] : records_) {
            if (!has_cursor || id > cursor_id) candidates.push_back(r);
        }
        std::sort(candidates.begin(), candidates.end(),
                  [](const Record& a, const Record& b) { return a.id < b.id; });

        size_t page_len = std::min((size_t)effective, candidates.size());
        items->assign(candidates.begin(), candidates.begin() + page_len);

        if (candidates.size() > page_len) {
            *next_page_token = EncodeToken(items->back().id, effective);
        } else {
            *next_page_token = "";
        }
        return true;
    }

private:
    std::map<long, Record> records_;
};
