from __future__ import annotations


class TreeNode:
    """Binary tree node."""

    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


def level_order(root: TreeNode | None) -> list[list[int]]:
    """Return root's node values grouped by depth, level by level
    from top to bottom, left to right within each level."""
    if root is None:
        return []
    res = []
    queue = [root]
    while queue:
        level = []
        next_queue = []
        for node in queue:
            level.append(node.val)
            if node.left is not None:
                next_queue.append(node.left)
            if node.right is not None:
                next_queue.append(node.right)
        res.append(level)
        queue = next_queue
    return res
