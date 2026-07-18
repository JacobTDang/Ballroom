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

// ConditionalStore: a versioned key-value store.
//
// TODO: no versioning at all -- Put always succeeds and overwrites
// unconditionally (the classic lost-update bug this store exists to
// prevent), and Get hands out the same constant etag forever. Every
// rule in the problem statement is still yours to build.
class ConditionalStore {
public:
    GetResult Get(const std::string& key, const std::string& if_none_match) {
        auto it = items_.find(key);
        if (it == items_.end()) return GetResult{404, "", ""};
        return GetResult{200, "1", it->second};
    }

    PutResult Put(const std::string& key, const std::string& body, const std::string& if_match) {
        items_[key] = body;
        return PutResult{200, "1"};
    }

    void Delete(const std::string& key) { items_.erase(key); }

private:
    std::map<std::string, std::string> items_;
};
