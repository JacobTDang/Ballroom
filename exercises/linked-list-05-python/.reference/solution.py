from __future__ import annotations


class Node:
    """Linked list node with an extra random pointer that can point
    to any node in the list, or None."""

    def __init__(self, x: int, next: Node | None = None, random: Node | None = None):
        self.val = x
        self.next = next
        self.random = random


def copy_random_list(head: Node | None) -> Node | None:
    """Return a deep copy of head — every node (including random
    targets) is a brand new node, never shared with the input."""
    if head is None:
        return None
    copies: dict[Node, Node] = {}
    cur = head
    while cur is not None:
        copies[cur] = Node(cur.val)
        cur = cur.next
    cur = head
    while cur is not None:
        copies[cur].next = copies[cur.next] if cur.next is not None else None
        copies[cur].random = copies[cur.random] if cur.random is not None else None
        cur = cur.next
    return copies[head]
