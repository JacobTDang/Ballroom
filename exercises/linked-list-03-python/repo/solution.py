from __future__ import annotations


class ListNode:
    """Singly linked list node."""

    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


def reorder_list(head: ListNode | None) -> None:
    """Reorder head in place from L0->L1->...->Ln into
    L0->Ln->L1->Ln-1->..."""
    pass
