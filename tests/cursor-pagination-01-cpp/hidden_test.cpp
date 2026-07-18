#include "solution.cpp"

#include <cstdio>

#define CHECK(cond, msg)                       \
    if (!(cond)) {                              \
        fprintf(stderr, "FAILED: %s\n", msg);   \
        return 1;                               \
    }

static std::string Tamper(const std::string& token) {
    std::string t = token;
    size_t mid = t.size() / 2;
    t[mid] = (t[mid] != 'X') ? 'X' : 'Y';
    return t;
}

static std::vector<Record> MakeRecords(const std::vector<long>& ids) {
    std::vector<Record> out;
    for (long id : ids) out.push_back({id, "r"});
    return out;
}

int main() {
    std::vector<Record> items;
    std::string token, err;

    // exactly-once full walk
    {
        std::vector<long> ids;
        for (long i = 1; i <= 23; i++) ids.push_back(i);
        CursorStore store(MakeRecords(ids));

        std::map<long, int> seen;
        token = "";
        while (true) {
            CHECK(store.List(5, token, &items, &token, &err), "List failed on full walk");
            for (auto& r : items) seen[r.id]++;
            if (token.empty()) break;
        }
        CHECK(seen.size() == ids.size(), "walk didn't visit every id");
        for (long id : ids) {
            CHECK(seen[id] == 1, "an id wasn't seen exactly once on the full walk");
        }
    }

    // empty final token / non-empty when more remain
    {
        CursorStore store(MakeRecords({1, 2, 3}));
        CHECK(store.List(10, "", &items, &token, &err), "List failed");
        CHECK(items.size() == 3, "expected all 3 records in one page");
        CHECK(token.empty(), "nothing left, but next_page_token wasn't empty");
    }
    {
        CursorStore store(MakeRecords({1, 2, 3, 4, 5}));
        CHECK(store.List(2, "", &items, &token, &err), "List failed");
        CHECK(items.size() == 2, "expected a 2-record page");
        CHECK(!token.empty(), "more records remain, but next_page_token was empty");
    }

    // tampered token errors
    {
        std::vector<long> ids;
        for (long i = 1; i <= 9; i++) ids.push_back(i);
        CursorStore store(MakeRecords(ids));
        CHECK(store.List(3, "", &items, &token, &err), "List failed");
        CHECK(!store.List(3, Tamper(token), &items, &token, &err),
              "a tampered page_token was silently accepted");
    }

    // page_size clamps
    {
        std::vector<long> ids;
        for (long i = 1; i <= 60; i++) ids.push_back(i);
        CursorStore store(MakeRecords(ids));

        CHECK(store.List(0, "", &items, &token, &err), "List failed");
        CHECK(items.size() == 10, "page_size<=0 didn't fall back to the 10-record default");

        CHECK(store.List(-5, "", &items, &token, &err), "List failed");
        CHECK(items.size() == 10, "negative page_size didn't fall back to the 10-record default");

        CHECK(store.List(10000, "", &items, &token, &err), "List failed");
        CHECK(items.size() == 50, "oversized page_size wasn't clamped to the 50-record max");
    }

    // parameter change invalidates token
    {
        std::vector<long> ids;
        for (long i = 1; i <= 29; i++) ids.push_back(i);
        CursorStore store(MakeRecords(ids));
        CHECK(store.List(5, "", &items, &token, &err), "List failed");
        CHECK(!store.List(7, token, &items, &token, &err),
              "resuming with a different page_size was silently honored");
    }

    // insert mid-walk never duplicates or skips
    {
        std::vector<long> seedIDs;
        for (long i = 1; i <= 20; i++) seedIDs.push_back(i * 10); // 10, 20, ..., 200
        CursorStore store(MakeRecords(seedIDs));

        CHECK(store.List(5, "", &items, &token, &err), "List failed");
        CHECK(items.size() == 5 && items[0].id == 10 && items[4].id == 50, "unexpected first page");

        std::map<long, int> seen;
        for (auto& r : items) seen[r.id]++;

        CHECK(store.Insert({5, "new-before-cursor"}, &err), "Insert failed");
        CHECK(store.Insert({999, "new-after-cursor"}, &err), "Insert failed");

        while (true) {
            CHECK(store.List(5, token, &items, &token, &err), "List failed mid-walk");
            for (auto& r : items) seen[r.id]++;
            if (token.empty()) break;
        }

        for (long id : seedIDs) {
            CHECK(seen[id] == 1, "an original id was skipped or duplicated after a mid-walk insert");
        }
    }

    printf("all assertions passed\n");
    return 0;
}
