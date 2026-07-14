// ClimbStairs returns the number of distinct ways to climb a staircase
// of n steps, taking 1 or 2 steps at a time.
int ClimbStairs(int n) {
    if (n <= 2) return n;
    int prev = 1, curr = 2;
    for (int i = 3; i <= n; i++) {
        int next = prev + curr;
        prev = curr;
        curr = next;
    }
    return curr;
}
