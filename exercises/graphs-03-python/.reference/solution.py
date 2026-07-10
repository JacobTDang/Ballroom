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
    if node is None:
        return None
    visited: dict[Node, Node] = {}

    def dfs(n: Node) -> Node:
        if n in visited:
            return visited[n]
        copy_node = Node(n.val)
        visited[n] = copy_node
        for nb in n.neighbors:
            copy_node.neighbors.append(dfs(nb))
        return copy_node

    return dfs(node)
