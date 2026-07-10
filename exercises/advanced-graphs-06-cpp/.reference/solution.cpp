#include <vector>

// FindCheapestPrice returns the cheapest price from src to dst using at
// most k stops (at most k+1 edges), or -1 if impossible.
int FindCheapestPrice(int n, std::vector<std::vector<int>>& flights, int src, int dst, int k) {
    const int inf = 1 << 30;

    std::vector<int> dist(n, inf);
    dist[src] = 0;

    // Bellman-Ford limited to exactly k+1 relaxation rounds. Each round
    // must relax edges using a SNAPSHOT of the previous round's
    // distances, not the array being updated in place during that same
    // round, or a single round could silently chain multiple edges
    // together and violate the stop limit.
    for (int round = 0; round <= k; round++) {
        std::vector<int> prev = dist;

        for (auto& f : flights) {
            int u = f[0], v = f[1], price = f[2];
            if (prev[u] == inf) continue;
            if (prev[u] + price < dist[v]) {
                dist[v] = prev[u] + price;
            }
        }
    }

    if (dist[dst] == inf) return -1;
    return dist[dst];
}
