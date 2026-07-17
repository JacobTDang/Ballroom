#include "solution.cpp"

#include <cstdio>

#define CHECK(cond, msg)                          \
    if (!(cond)) {                                \
        fprintf(stderr, "FAILED: %s\n", msg);    \
        return 1;                                 \
    }

int main() {
    {
        SlidingWindow s(2, 100);
        CHECK(s.AllowAt(90), "AllowAt(90)");
        CHECK(s.AllowAt(95), "AllowAt(95)");
        CHECK(!s.AllowAt(105), "boundary burst at 105 must be denied");
        CHECK(!s.AllowAt(110), "boundary burst at 110 must be denied");
        CHECK(s.AllowAt(191), "AllowAt(191) after aging out");
    }
    {
        SlidingWindow s(1, 100);
        CHECK(s.AllowAt(1000), "first");
        CHECK(!s.AllowAt(1099), "99ms-old request still counts");
        CHECK(s.AllowAt(1100), "exactly-window-old request must not count");
    }
    {
        SlidingWindow s(2, 100);
        s.AllowAt(0);
        s.AllowAt(1);
        for (long i = 2; i < 50; i++) {
            CHECK(!s.AllowAt(i), "third request inside window allowed");
        }
        CHECK(s.AllowAt(101), "denied requests must not consume the budget");
    }
    {
        SlidingWindow s(2, 100);
        for (long at = 0; at < 1000; at += 60) {
            CHECK(s.AllowAt(at), "steady under-limit rate denied");
        }
    }
    printf("all assertions passed\n");
    return 0;
}
