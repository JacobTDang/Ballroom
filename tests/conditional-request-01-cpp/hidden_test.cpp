#include "solution.cpp"

#include <cstdio>

#define CHECK(cond, msg)                       \
    if (!(cond)) {                              \
        fprintf(stderr, "FAILED: %s\n", msg);   \
        return 1;                               \
    }

int main() {
    // Get on a missing key.
    {
        ConditionalStore store;
        GetResult r = store.Get("missing", "");
        CHECK(r.status == 404, "Get on a missing key wasn't 404");
    }

    // Put create flows.
    {
        ConditionalStore store;
        PutResult r = store.Put("k1", "v1", "");
        CHECK(r.status == 200 && !r.etag.empty(), "fresh create wasn't 200 with a fresh etag");

        ConditionalStore store2;
        PutResult r2 = store2.Put("k1", "v1", "some-version");
        CHECK(r2.status == 412, "create with if_match set wasn't 412");
    }

    // Update requires If-Match (428), and a 428 leaves state unchanged.
    {
        ConditionalStore store;
        PutResult created = store.Put("k1", "v1", "");

        PutResult r = store.Put("k1", "v2", "");
        CHECK(r.status == 428, "update without if_match wasn't 428");

        GetResult g = store.Get("k1", "");
        CHECK(g.status == 200 && g.etag == created.etag && g.body == "v1",
              "state changed after a 428");
    }

    // Stale If-Match -> 412, state unchanged.
    {
        ConditionalStore store;
        PutResult created = store.Put("k1", "v1", "");

        PutResult r = store.Put("k1", "v2", "not-" + created.etag);
        CHECK(r.status == 412, "stale if_match wasn't 412");

        GetResult g = store.Get("k1", "");
        CHECK(g.status == 200 && g.etag == created.etag && g.body == "v1",
              "a failed conditional write changed state");
    }

    // Correct If-Match succeeds and rotates the etag.
    {
        ConditionalStore store;
        PutResult created = store.Put("k1", "v1", "");

        PutResult r = store.Put("k1", "v2", created.etag);
        CHECK(r.status == 200 && !r.etag.empty() && r.etag != created.etag,
              "update with correct if_match didn't rotate the etag");

        GetResult g = store.Get("k1", "");
        CHECK(g.status == 200 && g.etag == r.etag && g.body == "v2", "update didn't take effect");
    }

    // Get's If-None-Match matrix.
    {
        ConditionalStore store;
        PutResult created = store.Put("k1", "v1", "");

        GetResult r = store.Get("k1", created.etag);
        CHECK(r.status == 304, "matching if_none_match wasn't 304");

        GetResult r2 = store.Get("k1", "stale-etag");
        CHECK(r2.status == 200 && r2.etag == created.etag && r2.body == "v1",
              "stale if_none_match didn't return a fresh 200");
    }

    // No etag resurrection after delete + recreate.
    {
        ConditionalStore store;
        PutResult first = store.Put("b", "first", "");
        store.Delete("b");
        PutResult second = store.Put("b", "first", "");

        CHECK(second.etag != first.etag, "a recreated resource reused its old etag");

        GetResult g = store.Get("b", first.etag);
        CHECK(g.status == 200, "a stale pre-delete etag falsely matched the recreated resource");

        PutResult p = store.Put("b", "second", first.etag);
        CHECK(p.status == 412, "a stale pre-delete etag falsely satisfied If-Match");
    }

    printf("all assertions passed\n");
    return 0;
}
