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
    ListNode* node = head;
    int count = 0;
    while (node != nullptr && count < k) {
        node = node->next;
        count++;
    }
    if (count < k) return head;

    // node is now the head of the rest of the list, already
    // recursively reversed in groups of k.
    ListNode* newHead = ReverseKGroup(node, k);

    ListNode* cur = head;
    ListNode* prev = newHead;
    for (int i = 0; i < k; i++) {
        ListNode* next = cur->next;
        cur->next = prev;
        prev = cur;
        cur = next;
    }
    return prev;
}
