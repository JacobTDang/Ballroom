#include <functional>
#include <vector>

// CountComponents returns the number of connected components in the
// undirected graph of n nodes described by edges.
int CountComponents(int n, std::vector<std::vector<int>>& edges) {
    std::vector<int> parent(n);
    for (int i = 0; i < n; i++) parent[i] = i;

    std::function<int(int)> find = [&](int x) {
        while (parent[x] != x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    };

    int components = n;
    for (auto& e : edges) {
        int rootA = find(e[0]), rootB = find(e[1]);
        if (rootA != rootB) {
            parent[rootA] = rootB;
            components--;
        }
    }
    return components;
}
