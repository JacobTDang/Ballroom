// Rounds t down to the start of its k-wide bucket.
int align(int t, int k) {
    int q = t / k;
    if (t % k != 0 && t < 0) q--;
    return q * k;
}
