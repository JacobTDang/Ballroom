#include "solution.cpp"

#include <cstdio>

int main() {
    {
        BloomFilter b(16384, 4);
        char buf[32];
        for (int i = 0; i < 500; i++) {
            snprintf(buf, sizeof buf, "present-%d", i);
            b.Add(buf);
        }
        for (int i = 0; i < 500; i++) {
            snprintf(buf, sizeof buf, "present-%d", i);
            if (!b.MightContain(buf)) {
                fprintf(stderr, "added key %s reported absent -- no false negatives allowed\n", buf);
                return 1;
            }
        }
        int fps = 0;
        const int probes = 10000;
        for (int i = 0; i < probes; i++) {
            snprintf(buf, sizeof buf, "absent-%d", i);
            if (b.MightContain(buf)) fps++;
        }
        if (fps >= probes * 2 / 100) {
            fprintf(stderr, "%d/%d absent keys reported present -- FP rate must stay under 2%%\n", fps, probes);
            return 1;
        }
    }
    {
        BloomFilter b(1024, 3);
        char buf[32];
        for (int i = 0; i < 100; i++) {
            snprintf(buf, sizeof buf, "anything-%d", i);
            if (b.MightContain(buf)) {
                fprintf(stderr, "empty filter reported %s present\n", buf);
                return 1;
            }
        }
    }
    printf("all assertions passed\n");
    return 0;
}
