#include <cassert>
#include <cstdio>
#include <vector>

struct ListNode {
    int val;
    ListNode* next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode* next) : val(x), next(next) {}
};

bool HasCycle(ListNode* head);

ListNode* buildCycleList(const std::vector<int>& vals, int pos) {
    if (vals.empty()) return nullptr;
    std::vector<ListNode*> nodes;
    for (int v : vals) nodes.push_back(new ListNode(v));
    for (size_t i = 0; i + 1 < nodes.size(); i++) nodes[i]->next = nodes[i + 1];
    if (pos >= 0) nodes.back()->next = nodes[pos];
    return nodes[0];
}

int main() {
    assert(HasCycle(buildCycleList({3, 2, 0, -4}, 1)) == true);
    assert(HasCycle(buildCycleList({1, 2}, 0)) == true);
    assert(HasCycle(buildCycleList({1}, -1)) == false);
    assert(HasCycle(buildCycleList({}, -1)) == false);
    assert(HasCycle(buildCycleList({1, 2, 3}, -1)) == false);
    printf("all assertions passed\n");
    return 0;
}
