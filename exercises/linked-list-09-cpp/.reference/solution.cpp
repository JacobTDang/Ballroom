#include <vector>

// ListNode is a singly linked list node.
struct ListNode {
    int val;
    ListNode* next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode* next) : val(x), next(next) {}
};

ListNode* mergeTwoLists(ListNode* a, ListNode* b) {
    ListNode dummy;
    ListNode* cur = &dummy;
    while (a != nullptr && b != nullptr) {
        if (a->val <= b->val) {
            cur->next = a;
            a = a->next;
        } else {
            cur->next = b;
            b = b->next;
        }
        cur = cur->next;
    }
    cur->next = (a != nullptr) ? a : b;
    return dummy.next;
}

// MergeKLists merges k sorted linked lists into one sorted list.
ListNode* MergeKLists(std::vector<ListNode*>& lists) {
    if (lists.empty()) return nullptr;
    while (lists.size() > 1) {
        std::vector<ListNode*> merged;
        for (size_t i = 0; i < lists.size(); i += 2) {
            if (i + 1 < lists.size()) {
                merged.push_back(mergeTwoLists(lists[i], lists[i + 1]));
            } else {
                merged.push_back(lists[i]);
            }
        }
        lists = merged;
    }
    return lists[0];
}
