#include <vector>

// ListNode is a singly linked list node.
struct ListNode {
    int val;
    ListNode* next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode* next) : val(x), next(next) {}
};

// MergeKLists merges k sorted linked lists into one sorted list.
ListNode* MergeKLists(std::vector<ListNode*>& lists) {
    return nullptr;
}
