#include <functional>
#include <vector>

// FindRedundantConnection returns the edge that can be removed to
// turn the graph back into a tree, using union-find to detect the
// first edge that connects two already-connected nodes.
std::vector<int> FindRedundantConnection(std::vector<std::vector<int>>& edges) {
    int n = static_cast<int>(edges.size());
    std::vector<int> parent(n + 1);
    for (int i = 0; i <= n; i++) parent[i] = i;

    std::function<int(int)> find = [&](int x) {
        while (parent[x] != x) {
            parent[x] = parent[parent[x]];
            x = parent[x];
        }
        return x;
    };

    for (auto& e : edges) {
        int rootA = find(e[0]), rootB = find(e[1]);
        if (rootA == rootB) return e;
        parent[rootA] = rootB;
    }
    return {};
}
