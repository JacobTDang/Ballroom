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
    parts: list[str] = []

    def walk(node: TreeNode | None) -> None:
        if node is None:
            parts.append("#")
            return
        parts.append(str(node.val))
        walk(node.left)
        walk(node.right)

    walk(root)
    return ",".join(parts)


def deserialize(data: str) -> TreeNode | None:
    """Decode a string produced by serialize back into the original
    tree."""
    tokens = data.split(",")
    idx = 0

    def walk() -> TreeNode | None:
        nonlocal idx
        if idx >= len(tokens) or tokens[idx] == "#":
            idx += 1
            return None
        node = TreeNode(int(tokens[idx]))
        idx += 1
        node.left = walk()
        node.right = walk()
        return node

    return walk()
