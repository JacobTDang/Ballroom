// Rounds t down to the start of its k-wide bucket. Currently wrong for
// negative t -- find and fix the bug.
int align(int t, int k) {
    return t / k * k;
}
