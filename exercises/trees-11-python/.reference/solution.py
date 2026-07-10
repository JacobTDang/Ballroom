from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def is_valid_bst(root: TreeNode | None) -> bool:
    """Return whether root is a valid binary search tree."""

    def valid(node: TreeNode | None, lo: float, hi: float) -> bool:
        if node is None:
            return True
        if node.val <= lo or node.val >= hi:
            return False
        return valid(node.left, lo, node.val) and valid(node.right, node.val, hi)

    return valid(root, float("-inf"), float("inf"))
