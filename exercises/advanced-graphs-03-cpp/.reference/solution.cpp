#include <vector>

// NetworkDelayTime returns the minimum time for a signal starting at node
// k to reach every node in a directed weighted graph of n nodes described
// by times[i] = [u, v, w], or -1 if some node is unreachable.
int NetworkDelayTime(std::vector<std::vector<int>>& times, int n, int k) {
    const int inf = 1 << 30;

    std::vector<std::vector<std::pair<int, int>>> adj(n + 1);
    for (auto& t : times) {
        adj[t[0]].push_back({t[1], t[2]});
    }

    std::vector<int> dist(n + 1, inf);
    dist[k] = 0;

    std::vector<bool> visited(n + 1, false);
    for (int count = 0; count < n; count++) {
        int u = -1;
        for (int v = 1; v <= n; v++) {
            if (!visited[v] && (u == -1 || dist[v] < dist[u])) {
                u = v;
            }
        }
        if (u == -1 || dist[u] == inf) break;
        visited[u] = true;
        for (auto& edge : adj[u]) {
            int v = edge.first, w = edge.second;
            if (dist[u] + w < dist[v]) {
                dist[v] = dist[u] + w;
            }
        }
    }

    int maxDist = 0;
    for (int v = 1; v <= n; v++) {
        if (dist[v] == inf) return -1;
        if (dist[v] > maxDist) maxDist = dist[v];
    }
    return maxDist;
}
