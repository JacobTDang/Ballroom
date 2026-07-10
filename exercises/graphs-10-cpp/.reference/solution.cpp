#include <functional>
#include <vector>

// ValidTree reports whether the n nodes and given undirected edges
// form a valid tree (connected, no cycles).
bool ValidTree(int n, std::vector<std::vector<int>>& edges) {
    if (static_cast<int>(edges.size()) != n - 1) return false;

    std::vector<int> parent(n);
    for (int i = 0; i < n; i++) parent[i] = i;

    std::function<int(int)> find = [&](int x) {
        while (parent[x] != x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    };

    for (auto& e : edges) {
        int rootA = find(e[0]), rootB = find(e[1]);
        if (rootA == rootB) return false;
        parent[rootA] = rootB;
    }
    return true;
}
