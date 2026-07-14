#include <vector>

// MinCostClimbingStairs returns the minimum cost to reach the top of a
// staircase where cost[i] is the cost of stepping on stair i, and you
// may start from step 0 or step 1 for free, climbing 1 or 2 steps at a
// time.
int MinCostClimbingStairs(std::vector<int>& cost) {
    int n = cost.size();
    int prev = 0, curr = 0;
    for (int i = 2; i <= n; i++) {
        int next = curr + cost[i - 1];
        int alt = prev + cost[i - 2];
        if (alt < next) next = alt;
        prev = curr;
        curr = next;
    }
    return curr;
}
