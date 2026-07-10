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

    def dfs(node: TreeNode | None, max_so_far: float) -> int:
        if node is None:
            return 0
        count = 0
        if node.val >= max_so_far:
            count = 1
            max_so_far = node.val
        count += dfs(node.left, max_so_far)
        count += dfs(node.right, max_so_far)
        return count

    return dfs(root, float("-inf"))
