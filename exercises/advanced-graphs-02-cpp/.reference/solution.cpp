#include <cstdlib>
#include <vector>

// MinCostConnectPoints returns the minimum total cost to connect all
// points, where the cost between two points is their Manhattan distance
// (the minimum spanning tree over the complete graph of points).
int MinCostConnectPoints(std::vector<std::vector<int>>& points) {
    int n = static_cast<int>(points.size());
    if (n <= 1) return 0;

    std::vector<bool> inTree(n, false);
    std::vector<int> minDist(n, 1 << 30);
    minDist[0] = 0;

    int total = 0;
    for (int count = 0; count < n; count++) {
        int u = -1;
        for (int v = 0; v < n; v++) {
            if (!inTree[v] && (u == -1 || minDist[v] < minDist[u])) {
                u = v;
            }
        }
        inTree[u] = true;
        total += minDist[u];

        for (int v = 0; v < n; v++) {
            if (!inTree[v]) {
                int dist = std::abs(points[u][0] - points[v][0]) + std::abs(points[u][1] - points[v][1]);
                if (dist < minDist[v]) {
                    minDist[v] = dist;
                }
            }
        }
    }
    return total;
}
