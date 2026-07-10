// ListNode is a singly linked list node.
struct ListNode {
    int val;
    ListNode* next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode* next) : val(x), next(next) {}
};

// ReverseKGroup reverses head k nodes at a time, leaving any
// remaining group shorter than k untouched.
ListNode* ReverseKGroup(ListNode* head, int k) {
    return nullptr;
}
