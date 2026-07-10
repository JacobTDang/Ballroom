from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def _is_same_tree(p: TreeNode | None, q: TreeNode | None) -> bool:
    if p is None and q is None:
        return True
    if p is None or q is None or p.val != q.val:
        return False
    return _is_same_tree(p.left, q.left) and _is_same_tree(p.right, q.right)


def is_subtree(root: TreeNode | None, sub_root: TreeNode | None) -> bool:
    """Return whether sub_root's tree matches some node in root's
    tree and everything below it."""
    if root is None:
        return sub_root is None
    if _is_same_tree(root, sub_root):
        return True
    return is_subtree(root.left, sub_root) or is_subtree(root.right, sub_root)
