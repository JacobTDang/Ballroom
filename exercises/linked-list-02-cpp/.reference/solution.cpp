// ListNode is a singly linked list node.
struct ListNode {
    int val;
    ListNode* next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode* next) : val(x), next(next) {}
};

// MergeTwoLists merges two sorted linked lists into one sorted list.
ListNode* MergeTwoLists(ListNode* list1, ListNode* list2) {
    ListNode dummy;
    ListNode* cur = &dummy;
    while (list1 != nullptr && list2 != nullptr) {
        if (list1->val <= list2->val) {
            cur->next = list1;
            list1 = list1->next;
        } else {
            cur->next = list2;
            list2 = list2->next;
        }
        cur = cur->next;
    }
    cur->next = (list1 != nullptr) ? list1 : list2;
    return dummy.next;
}
