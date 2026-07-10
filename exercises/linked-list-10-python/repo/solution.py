from __future__ import annotations


class ListNode:
    """Singly linked list node."""

    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


def reverse_k_group(head: ListNode | None, k: int) -> ListNode | None:
    """Reverse head k nodes at a time, leaving any remaining group
    shorter than k untouched."""
    return None
