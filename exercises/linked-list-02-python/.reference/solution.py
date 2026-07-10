from __future__ import annotations


class ListNode:
    """Singly linked list node."""

    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


def merge_two_lists(list1: ListNode | None, list2: ListNode | None) -> ListNode | None:
    """Merge two sorted linked lists into one sorted list."""
    dummy = ListNode()
    cur = dummy
    while list1 is not None and list2 is not None:
        if list1.val <= list2.val:
            cur.next = list1
            list1 = list1.next
        else:
            cur.next = list2
            list2 = list2.next
        cur = cur.next
    cur.next = list1 if list1 is not None else list2
    return dummy.next
