from solution import Node, clone_graph


def build_graph(adj_list):
    if not adj_list:
        return None
    nodes = [None] + [Node(v) for v in range(1, len(adj_list) + 1)]  # 1-indexed
    for v, neighbors in enumerate(adj_list, start=1):
        for nv in neighbors:
            nodes[v].neighbors.append(nodes[nv])
    return nodes[1]


def adjacency(start):
    if start is None:
        return {}
    out = {}
    visited = {start}
    queue = [start]
    while queue:
        n = queue.pop(0)
        nb_vals = []
        for nb in n.neighbors:
            nb_vals.append(nb.val)
            if nb not in visited:
                visited.add(nb)
                queue.append(nb)
        out[n.val] = sorted(nb_vals)
    return out


def all_nodes(start):
    if start is None:
        return []
    visited = {start}
    queue = [start]
    result = []
    while queue:
        n = queue.pop(0)
        result.append(n)
        for nb in n.neighbors:
            if nb not in visited:
                visited.add(nb)
                queue.append(nb)
    return result


def test_clone_graph():
    original = build_graph([[2, 4], [1, 3], [2, 4], [1, 3]])
    original_set = set(all_nodes(original))

    clone = clone_graph(original)

    assert adjacency(original) == adjacency(clone)

    for n in all_nodes(clone):
        assert n not in original_set, "cloned node shares identity with an original node"


def test_clone_graph_none_input():
    assert clone_graph(None) is None


def test_clone_graph_single_node_no_neighbors():
    original = Node(1)
    clone = clone_graph(original)
    assert clone is not original
    assert clone.val == 1
    assert clone.neighbors == []
