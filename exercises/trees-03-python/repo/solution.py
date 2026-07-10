from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def diameter_of_binary_tree(root: TreeNode | None) -> int:
    """Return the number of edges on the longest path between any
    two nodes in root's tree."""
    return 0
