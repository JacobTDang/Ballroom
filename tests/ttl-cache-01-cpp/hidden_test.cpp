#include "solution.cpp"

#include <cstdio>

#define CHECK(cond, msg)                          \
    if (!(cond)) {                                \
        fprintf(stderr, "FAILED: %s\n", msg);    \
        return 1;                                 \
    }

int main() {
    int v;
    {
        TTLCache c(2, 1000);
        c.PutAt("a", 1, 0);
        CHECK(c.GetAt("a", 10, &v) && v == 1, "basic get");
        CHECK(!c.GetAt("missing", 10, &v), "missing key hit");
    }
    {
        TTLCache c(2, 100000);
        c.PutAt("a", 1, 0);
        c.PutAt("b", 2, 1);
        c.GetAt("a", 2, &v);
        c.PutAt("c", 3, 3);
        CHECK(!c.GetAt("b", 4, &v), "b survived despite being LRU");
        CHECK(c.GetAt("a", 4, &v) && v == 1, "a evicted despite recent use");
        CHECK(c.GetAt("c", 4, &v) && v == 3, "c missing after insert");
    }
    {
        TTLCache c(2, 100);
        c.PutAt("a", 1, 0);
        CHECK(c.GetAt("a", 99, &v), "expired early");
        CHECK(!c.GetAt("a", 100, &v), "alive at exactly ttl");
    }
    {
        TTLCache c(2, 100);
        c.PutAt("a", 1, 0);
        CHECK(c.GetAt("a", 99, &v), "setup");
        CHECK(!c.GetAt("a", 100, &v), "get must not extend the TTL");
    }
    {
        TTLCache c(2, 100);
        c.PutAt("a", 1, 0);
        c.PutAt("b", 2, 0);
        c.PutAt("x", 10, 200);
        c.PutAt("y", 20, 201);
        CHECK(c.GetAt("x", 202, &v) && v == 10, "corpse counted against capacity (x)");
        CHECK(c.GetAt("y", 202, &v) && v == 20, "corpse counted against capacity (y)");
    }
    {
        TTLCache c(2, 100);
        c.PutAt("a", 1, 0);
        c.PutAt("a", 2, 50);
        CHECK(c.GetAt("a", 149, &v) && v == 2, "rewrite value/expiry");
        CHECK(!c.GetAt("a", 150, &v), "rewrite expiry end");
    }
    printf("all assertions passed\n");
    return 0;
}
