from __future__ import annotations


class ListNode:
    """Singly linked list node."""

    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


def reorder_list(head: ListNode | None) -> None:
    """Reorder head in place from L0->L1->...->Ln into
    L0->Ln->L1->Ln-1->..."""
    if head is None or head.next is None:
        return

    slow, fast = head, head
    while fast.next is not None and fast.next.next is not None:
        slow = slow.next
        fast = fast.next.next

    second = slow.next
    slow.next = None
    prev = None
    while second is not None:
        nxt = second.next
        second.next = prev
        prev = second
        second = nxt

    first = head
    while prev is not None:
        first_next = first.next
        second_next = prev.next
        first.next = prev
        if first_next is not None:
            prev.next = first_next
        first = first_next
        prev = second_next
