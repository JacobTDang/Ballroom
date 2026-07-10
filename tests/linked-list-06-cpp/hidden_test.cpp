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

ListNode* AddTwoNumbers(ListNode* l1, ListNode* l2);

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
    assert(toVector(AddTwoNumbers(buildList({2, 4, 3}), buildList({5, 6, 4}))) ==
           std::vector<int>({7, 0, 8}));
    assert(toVector(AddTwoNumbers(buildList({0}), buildList({0}))) == std::vector<int>({0}));
    assert(toVector(AddTwoNumbers(buildList({9, 9, 9, 9, 9, 9, 9}), buildList({9, 9, 9, 9}))) ==
           std::vector<int>({8, 9, 9, 9, 0, 0, 0, 1}));
    assert(toVector(AddTwoNumbers(buildList({5}), buildList({5}))) == std::vector<int>({0, 1}));
    printf("all assertions passed\n");
    return 0;
}
