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

ListNode* RemoveNthFromEnd(ListNode* head, int n);

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
    assert(toVector(RemoveNthFromEnd(buildList({1, 2, 3, 4, 5}), 2)) ==
           std::vector<int>({1, 2, 3, 5}));
    assert(toVector(RemoveNthFromEnd(buildList({1}), 1)) == std::vector<int>({}));
    assert(toVector(RemoveNthFromEnd(buildList({1, 2}), 1)) == std::vector<int>({1}));
    assert(toVector(RemoveNthFromEnd(buildList({1, 2}), 2)) == std::vector<int>({2}));
    printf("all assertions passed\n");
    return 0;
}
