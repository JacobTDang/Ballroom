#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

struct Record {
    std::string key;
    int value;

    bool operator==(const Record& other) const {
        return key == other.key && value == other.value;
    }
};

std::vector<Record*> dedupe(const std::vector<Record*>& records);

static bool kv_match(const std::vector<Record*>& got, const std::vector<Record>& want) {
    if (got.size() != want.size()) return false;
    for (size_t i = 0; i < got.size(); i++) {
        if (!(*got[i] == want[i])) return false;
    }
    return true;
}

int main() {
    {
        Record r1{"a", 1};
        Record r2{"a", 1};  // distinct object, same value
        Record r3{"b", 2};
        auto got = dedupe({&r1, &r2, &r3});
        assert(kv_match(got, {{"a", 1}, {"b", 2}}));
    }
    {
        Record r1{"x", 5}, r2{"x", 5}, r3{"x", 5};
        auto got = dedupe({&r1, &r2, &r3});
        assert(kv_match(got, {{"x", 5}}));
    }
    {
        Record r1{"a", 1}, r2{"b", 2}, r3{"c", 3};
        auto got = dedupe({&r1, &r2, &r3});
        assert(kv_match(got, {{"a", 1}, {"b", 2}, {"c", 3}}));
    }
    {
        Record r1{"a", 1}, r2{"b", 2}, r3{"a", 1}, r4{"c", 3}, r5{"b", 2};
        auto got = dedupe({&r1, &r2, &r3, &r4, &r5});
        assert(kv_match(got, {{"a", 1}, {"b", 2}, {"c", 3}}));
    }
    {
        Record r1{"a", 1};
        Record r2{"a", 2};
        auto got = dedupe({&r1, &r2});
        assert(kv_match(got, {{"a", 1}, {"a", 2}}));
    }
    printf("all assertions passed\n");
    return 0;
}
