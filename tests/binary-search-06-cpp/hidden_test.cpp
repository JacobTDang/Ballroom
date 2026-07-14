#include <cassert>
#include <cstdio>

#include "solution.hpp"

int main() {
    {
        TimeMap m;
        m.set("foo", "bar", 1);
        assert(m.get("foo", 1) == "bar");
        assert(m.get("foo", 3) == "bar");
        m.set("foo", "bar2", 4);
        assert(m.get("foo", 4) == "bar2");
        assert(m.get("foo", 5) == "bar2");
    }
    {
        TimeMap m;
        m.set("foo", "bar", 5);
        assert(m.get("foo", 1) == "");
    }
    {
        TimeMap m;
        assert(m.get("missing", 1) == "");
    }
    {
        TimeMap m;
        m.set("foo", "foo-val", 1);
        m.set("bar", "bar-val", 2);
        assert(m.get("foo", 10) == "foo-val");
        assert(m.get("bar", 10) == "bar-val");
    }
    printf("all assertions passed\n");
    return 0;
}
