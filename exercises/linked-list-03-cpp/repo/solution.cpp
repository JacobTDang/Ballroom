// ListNode is a singly linked list node.
struct ListNode {
    int val;
    ListNode* next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode* next) : val(x), next(next) {}
};

// ReorderList reorders head in place from L0->L1->...->Ln into
// L0->Ln->L1->Ln-1->...
void ReorderList(ListNode* head) {
    // TODO: implement
}
