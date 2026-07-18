#include <map>
#include <string>

// GetResult is what Get hands back. etag/body are only meaningful
// when status == 200.
struct GetResult {
    int status;
    std::string etag;
    std::string body;
};

// PutResult is what Put hands back. etag is only meaningful when
// status == 200.
struct PutResult {
    int status;
    std::string etag;
};

namespace {
struct Entry {
    std::string etag;
    std::string body;
};
}  // namespace

// ConditionalStore: every successful write (create or update) draws
// its etag from one store-wide monotonic sequence -- so deleting a
// key and recreating it never reissues an old etag, and a stale
// cached etag can never falsely match again, on either the read or
// write side.
class ConditionalStore {
public:
    GetResult Get(const std::string& key, const std::string& if_none_match) {
        auto it = items_.find(key);
        if (it == items_.end()) return GetResult{404, "", ""};
        const Entry& e = it->second;
        if (!if_none_match.empty() && if_none_match == e.etag) {
            return GetResult{304, "", ""};
        }
        return GetResult{200, e.etag, e.body};
    }

    PutResult Put(const std::string& key, const std::string& body, const std::string& if_match) {
        auto it = items_.find(key);
        if (it == items_.end()) {
            if (!if_match.empty()) {
                // Claimed a version of something that doesn't exist.
                return PutResult{412, ""};
            }
            std::string etag = NextETag();
            items_[key] = Entry{etag, body};
            return PutResult{200, etag};
        }

        Entry& e = it->second;
        if (if_match.empty()) {
            // Blind overwrite of an existing resource: refused on purpose.
            return PutResult{428, ""};
        }
        if (if_match != e.etag) {
            // Stale precondition -- state is left exactly as it was.
            return PutResult{412, ""};
        }

        std::string etag = NextETag();
        e.etag = etag;
        e.body = body;
        return PutResult{200, etag};
    }

    void Delete(const std::string& key) { items_.erase(key); }

private:
    std::string NextETag() { return std::to_string(++seq_); }

    std::map<std::string, Entry> items_;
    long seq_ = 0;
};
