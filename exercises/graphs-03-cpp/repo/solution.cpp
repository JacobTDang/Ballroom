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

// CloneGraph returns a deep copy of the connected graph reachable
// from node -- every node (including neighbor references) is a
// brand new node, never shared with the input.
Node* CloneGraph(Node* node) {
    return nullptr;
}
