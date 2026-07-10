from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def good_nodes(root: TreeNode | None) -> int:
    """Count nodes X in root's tree where no node on the path from
    root to X has a value greater than X."""
    return 0
