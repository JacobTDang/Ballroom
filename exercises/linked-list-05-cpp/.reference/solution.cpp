#include <unordered_map>

// Node is a linked list node with an extra random pointer that can
// point to any node in the list, or nullptr.
class Node {
public:
    int val;
    Node* next;
    Node* random;

    Node(int _val) {
        val = _val;
        next = nullptr;
        random = nullptr;
    }
};

// CopyRandomList returns a deep copy of head — every node (including
// random targets) is a brand new node, never shared with the input.
Node* CopyRandomList(Node* head) {
    if (head == nullptr) return nullptr;
    std::unordered_map<Node*, Node*> copies;
    for (Node* cur = head; cur != nullptr; cur = cur->next) {
        copies[cur] = new Node(cur->val);
    }
    for (Node* cur = head; cur != nullptr; cur = cur->next) {
        copies[cur]->next = cur->next ? copies[cur->next] : nullptr;
        copies[cur]->random = cur->random ? copies[cur->random] : nullptr;
    }
    return copies[head];
}
