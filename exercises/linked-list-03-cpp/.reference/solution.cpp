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
    if (head == nullptr || head->next == nullptr) return;

    ListNode* slow = head;
    ListNode* fast = head;
    while (fast->next != nullptr && fast->next->next != nullptr) {
        slow = slow->next;
        fast = fast->next->next;
    }

    ListNode* second = slow->next;
    slow->next = nullptr;
    ListNode* prev = nullptr;
    while (second != nullptr) {
        ListNode* next = second->next;
        second->next = prev;
        prev = second;
        second = next;
    }

    ListNode* first = head;
    while (prev != nullptr) {
        ListNode* firstNext = first->next;
        ListNode* secondNext = prev->next;
        first->next = prev;
        if (firstNext != nullptr) {
            prev->next = firstNext;
        }
        first = firstNext;
        prev = secondNext;
    }
}
