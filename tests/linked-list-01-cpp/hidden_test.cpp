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

ListNode* ReverseList(ListNode* head);

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
    assert(toVector(ReverseList(buildList({1, 2, 3, 4, 5}))) ==
           std::vector<int>({5, 4, 3, 2, 1}));
    assert(toVector(ReverseList(buildList({1, 2}))) == std::vector<int>({2, 1}));
    assert(toVector(ReverseList(buildList({}))) == std::vector<int>({}));
    assert(toVector(ReverseList(buildList({7}))) == std::vector<int>({7}));
    printf("all assertions passed\n");
    return 0;
}
