from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def is_subtree(root: TreeNode | None, sub_root: TreeNode | None) -> bool:
    """Return whether sub_root's tree matches some node in root's
    tree and everything below it."""
    return False
