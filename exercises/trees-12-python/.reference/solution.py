from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def kth_smallest(root: TreeNode | None, k: int) -> int:
    """Return the kth smallest value (1-indexed) among all nodes in
    the BST rooted at root."""
    stack: list[TreeNode] = []
    cur = root
    while True:
        while cur is not None:
            stack.append(cur)
            cur = cur.left
        cur = stack.pop()
        k -= 1
        if k == 0:
            return cur.val
        cur = cur.right
