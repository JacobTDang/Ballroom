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

ListNode* ReverseKGroup(ListNode* head, int k);

ListNode* buildList(const std::vector<int>& vals) {
    ListNode dummy;
    ListNode* cur = &dummy;
    for (int v : vals) {
        cur->next = new ListNode(v);
        cur = cur->next;
    }
    return dummy.next;
}

std::vector<int> toVector(ListNode* head) {
    std::vector<int> out;
    for (ListNode* n = head; n != nullptr; n = n->next) {
        out.push_back(n->val);
    }
    return out;
}

int main() {
    assert(toVector(ReverseKGroup(buildList({1, 2, 3, 4, 5}), 2)) ==
           std::vector<int>({2, 1, 4, 3, 5}));
    assert(toVector(ReverseKGroup(buildList({1, 2, 3, 4, 5}), 3)) ==
           std::vector<int>({3, 2, 1, 4, 5}));
    assert(toVector(ReverseKGroup(buildList({1, 2, 3, 4, 5}), 1)) ==
           std::vector<int>({1, 2, 3, 4, 5}));
    assert(toVector(ReverseKGroup(buildList({1, 2, 3, 4, 5, 6}), 6)) ==
           std::vector<int>({6, 5, 4, 3, 2, 1}));
    assert(toVector(ReverseKGroup(buildList({1}), 1)) == std::vector<int>({1}));
    assert(toVector(ReverseKGroup(buildList({1, 2, 3, 4, 5}), 4)) ==
           std::vector<int>({4, 3, 2, 1, 5}));
    assert(toVector(ReverseKGroup(buildList({1, 2, 3, 4, 5, 6, 7, 8}), 2)) ==
           std::vector<int>({2, 1, 4, 3, 6, 5, 8, 7}));
    printf("all assertions passed\n");
    return 0;
}
