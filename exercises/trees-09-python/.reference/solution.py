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
    if root is None:
        return []
    res = []
    queue = [root]
    while queue:
        next_queue = []
        for i, node in enumerate(queue):
            if i == len(queue) - 1:
                res.append(node.val)
            if node.left is not None:
                next_queue.append(node.left)
            if node.right is not None:
                next_queue.append(node.right)
        queue = next_queue
    return res
