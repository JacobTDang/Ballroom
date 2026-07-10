from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def is_balanced(root: TreeNode | None) -> bool:
    """Return whether every node's left and right subtrees differ in
    height by no more than 1."""
    return False
