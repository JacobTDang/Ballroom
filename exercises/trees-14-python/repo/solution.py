from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def max_path_sum(root: TreeNode) -> int:
    """Return the maximum sum of any non-empty path between two
    nodes in root's tree (the path need not pass through root)."""
    return 0
