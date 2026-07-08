#include <cstdio>

int Counter(int n);

int main() {
    int got = Counter(200);
    if (got != 200) {
        fprintf(stderr, "Counter(200) = %d, want 200\n", got);
        return 1;
    }
    printf("all assertions passed\n");
    return 0;
}
