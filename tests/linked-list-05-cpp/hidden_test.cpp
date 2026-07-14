#include <cassert>
#include <cstdio>
#include <unordered_map>
#include <unordered_set>
#include <vector>

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

Node* CopyRandomList(Node* head);

std::vector<Node*> buildRandomList(const std::vector<int>& vals, const std::vector<int>& randomIdx) {
    std::vector<Node*> nodes;
    for (int v : vals) nodes.push_back(new Node(v));
    for (size_t i = 0; i < nodes.size(); i++) {
        if (i + 1 < nodes.size()) nodes[i]->next = nodes[i + 1];
        if (randomIdx[i] >= 0) nodes[i]->random = nodes[randomIdx[i]];
    }
    return nodes;
}

void check(std::vector<int> vals, std::vector<int> randomIdx) {
    auto orig = buildRandomList(vals, randomIdx);
    std::unordered_set<Node*> origSet(orig.begin(), orig.end());

    Node* head = orig.empty() ? nullptr : orig[0];
    Node* copyHead = CopyRandomList(head);

    std::vector<Node*> copies;
    for (Node* cur = copyHead; cur != nullptr; cur = cur->next) {
        assert(origSet.find(cur) == origSet.end());  // must be a deep copy
        copies.push_back(cur);
    }
    assert(copies.size() == vals.size());

    std::unordered_map<Node*, int> copyIndex;
    for (size_t i = 0; i < copies.size(); i++) copyIndex[copies[i]] = static_cast<int>(i);

    for (size_t i = 0; i < copies.size(); i++) {
        assert(copies[i]->val == vals[i]);
        if (randomIdx[i] == -1) {
            assert(copies[i]->random == nullptr);
            continue;
        }
        assert(copies[i]->random != nullptr);
        assert(copyIndex.count(copies[i]->random) == 1);
        assert(copyIndex.at(copies[i]->random) == randomIdx[i]);
    }
}

int main() {
    check({7, 13, 11, 10, 1}, {-1, 0, 4, 2, 0});
    check({1, 2}, {1, 1});
    check({3, 3, 3}, {-1, -1, -1});
    check({}, {});
    check({5}, {0});
    printf("all assertions passed\n");
    return 0;
}
