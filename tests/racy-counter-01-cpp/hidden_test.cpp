#include <cstdio>

int Counter(int n);

int main() {
    int got = Counter(200);
    if (got != 200) {
        fprintf(stderr, "Counter(200) = %d, want 200\n", got);
        return 1;
    }

    int got_small = Counter(1);
    if (got_small != 1) {
        fprintf(stderr, "Counter(1) = %d, want 1\n", got_small);
        return 1;
    }

    int got_large = Counter(2000);
    if (got_large != 2000) {
        fprintf(stderr, "Counter(2000) = %d, want 2000\n", got_large);
        return 1;
    }

    printf("all assertions passed\n");
    return 0;
}
