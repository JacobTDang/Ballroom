from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def max_depth(root: TreeNode | None) -> int:
    """Return the number of nodes along the longest path from root
    down to the farthest leaf."""
    if root is None:
        return 0
    return 1 + max(max_depth(root.left), max_depth(root.right))
