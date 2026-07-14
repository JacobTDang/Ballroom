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

ListNode* MergeTwoLists(ListNode* list1, ListNode* list2);

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
    assert(toVector(MergeTwoLists(buildList({1, 2, 4}), buildList({1, 3, 4}))) ==
           std::vector<int>({1, 1, 2, 3, 4, 4}));
    assert(toVector(MergeTwoLists(buildList({}), buildList({}))) == std::vector<int>({}));
    assert(toVector(MergeTwoLists(buildList({}), buildList({0}))) == std::vector<int>({0}));
    assert(toVector(MergeTwoLists(buildList({5}), buildList({1, 2, 4}))) ==
           std::vector<int>({1, 2, 4, 5}));
    assert(toVector(MergeTwoLists(buildList({1, 1, 1}), buildList({1, 1, 1}))) ==
           std::vector<int>({1, 1, 1, 1, 1, 1}));
    assert(toVector(MergeTwoLists(buildList({-3, -1, 2}), buildList({-2, 0, 5}))) ==
           std::vector<int>({-3, -2, -1, 0, 2, 5}));
    printf("all assertions passed\n");
    return 0;
}
