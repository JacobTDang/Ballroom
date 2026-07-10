from __future__ import annotations


class ListNode:
    """Singly linked list node."""

    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


def reverse_k_group(head: ListNode | None, k: int) -> ListNode | None:
    """Reverse head k nodes at a time, leaving any remaining group
    shorter than k untouched."""
    node = head
    count = 0
    while node is not None and count < k:
        node = node.next
        count += 1
    if count < k:
        return head

    # node is now the head of the rest of the list, already
    # recursively reversed in groups of k.
    new_head = reverse_k_group(node, k)

    cur = head
    prev = new_head
    for _ in range(k):
        nxt = cur.next
        cur.next = prev
        prev = cur
        cur = nxt
    return prev
