from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def build_tree(preorder: list[int], inorder: list[int]) -> TreeNode | None:
    """Reconstruct the unique binary tree whose preorder and inorder
    traversals are preorder and inorder."""
    inorder_idx = {v: i for i, v in enumerate(inorder)}
    pre = 0

    def build(in_lo: int, in_hi: int) -> TreeNode | None:
        nonlocal pre
        if in_lo > in_hi:
            return None
        root_val = preorder[pre]
        pre += 1
        root = TreeNode(root_val)
        mid = inorder_idx[root_val]
        root.left = build(in_lo, mid - 1)
        root.right = build(mid + 1, in_hi)
        return root

    return build(0, len(inorder) - 1)
