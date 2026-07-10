from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def right_side_view(root: TreeNode | None) -> list[int]:
    """Return the value of the rightmost node at each depth of
    root's tree, top to bottom."""
    return []
