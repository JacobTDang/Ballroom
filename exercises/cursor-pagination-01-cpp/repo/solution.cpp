#include <algorithm>
#include <cstdlib>
#include <map>
#include <string>
#include <vector>

// Record is a single row in the store, ordered by id.
struct Record {
    long id;
    std::string name;
};

// CursorStore: keyset pagination over id-sorted records.
//
// TODO: this paginates by raw position offset -- it looks fine until
// the data changes between calls. No tamper check, no page-size
// clamp, no page-size binding in the token. Every rule in the problem
// statement is still yours to build.
class CursorStore {
public:
    explicit CursorStore(std::vector<Record> records) {
        for (auto& r : records) records_[r.id] = r;
    }

    bool Insert(const Record& r, std::string* err) {
        records_[r.id] = r;
        return true;
    }

    bool List(int page_size, const std::string& page_token,
              std::vector<Record>* items, std::string* next_page_token,
              std::string* err) {
        long offset = page_token.empty() ? 0 : std::strtol(page_token.c_str(), nullptr, 10);

        std::vector<Record> ordered;
        ordered.reserve(records_.size());
        for (auto& [id, r] : records_) ordered.push_back(r);
        std::sort(ordered.begin(), ordered.end(),
                  [](const Record& a, const Record& b) { return a.id < b.id; });

        long end = offset + page_size;
        if (end > (long)ordered.size()) end = ordered.size();
        if (offset > (long)ordered.size()) offset = ordered.size();
        if (offset < 0) offset = 0;

        items->assign(ordered.begin() + offset, ordered.begin() + end);
        *next_page_token = (end < (long)ordered.size()) ? std::to_string(end) : "";
        return true;
    }

private:
    std::map<long, Record> records_;
};
