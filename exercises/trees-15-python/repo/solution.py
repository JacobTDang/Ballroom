from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def serialize(root: TreeNode | None) -> str:
    """Encode root as a string that deserialize can turn back into
    an equivalent tree. The exact format is up to you."""
    return ""


def deserialize(data: str) -> TreeNode | None:
    """Decode a string produced by serialize back into the original
    tree."""
    return None
