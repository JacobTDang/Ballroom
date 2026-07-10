#include <algorithm>
#include <cassert>
#include <cstdio>
#include <map>
#include <queue>
#include <set>
#include <vector>

// Node is an undirected graph node.
class Node {
public:
    int val;
    std::vector<Node*> neighbors;
    Node() : val(0) {}
    Node(int _val) : val(_val) {}
    Node(int _val, std::vector<Node*> _neighbors) : val(_val), neighbors(_neighbors) {}
};

Node* CloneGraph(Node* node);

Node* buildGraph(const std::vector<std::vector<int>>& adjList) {
    if (adjList.empty()) return nullptr;
    std::vector<Node*> nodes(adjList.size() + 1, nullptr);  // 1-indexed
    for (size_t v = 1; v <= adjList.size(); v++) nodes[v] = new Node(static_cast<int>(v));
    for (size_t v = 0; v < adjList.size(); v++) {
        for (int nv : adjList[v]) nodes[v + 1]->neighbors.push_back(nodes[nv]);
    }
    return nodes[1];
}

std::map<int, std::vector<int>> adjacency(Node* start) {
    std::map<int, std::vector<int>> out;
    if (start == nullptr) return out;
    std::set<Node*> visited = {start};
    std::queue<Node*> q;
    q.push(start);
    while (!q.empty()) {
        Node* n = q.front();
        q.pop();
        std::vector<int> nbVals;
        for (Node* nb : n->neighbors) {
            nbVals.push_back(nb->val);
            if (visited.find(nb) == visited.end()) {
                visited.insert(nb);
                q.push(nb);
            }
        }
        std::sort(nbVals.begin(), nbVals.end());
        out[n->val] = nbVals;
    }
    return out;
}

std::vector<Node*> allNodes(Node* start) {
    std::vector<Node*> all;
    if (start == nullptr) return all;
    std::set<Node*> visited = {start};
    std::queue<Node*> q;
    q.push(start);
    while (!q.empty()) {
        Node* n = q.front();
        q.pop();
        all.push_back(n);
        for (Node* nb : n->neighbors) {
            if (visited.find(nb) == visited.end()) {
                visited.insert(nb);
                q.push(nb);
            }
        }
    }
    return all;
}

int main() {
    {
        Node* original = buildGraph({{2, 4}, {1, 3}, {2, 4}, {1, 3}});
        std::set<Node*> originalSet;
        for (Node* n : allNodes(original)) originalSet.insert(n);

        Node* clone = CloneGraph(original);
        assert(adjacency(original) == adjacency(clone));

        for (Node* n : allNodes(clone)) {
            assert(originalSet.find(n) == originalSet.end());
        }
    }
    {
        assert(CloneGraph(nullptr) == nullptr);
    }
    {
        Node* original = new Node(1);
        Node* clone = CloneGraph(original);
        assert(clone != original);
        assert(clone->val == 1);
        assert(clone->neighbors.empty());
    }
    printf("all assertions passed\n");
    return 0;
}
