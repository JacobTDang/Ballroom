from __future__ import annotations


class Node:
    """Undirected graph node."""

    def __init__(self, val=0, neighbors=None):
        self.val = val
        self.neighbors = neighbors if neighbors is not None else []


def clone_graph(node: Node | None) -> Node | None:
    """Return a deep copy of the connected graph reachable from
    node -- every node (including neighbor references) is a brand
    new node, never shared with the input."""
    return None
