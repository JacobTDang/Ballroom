from solution import Node, copy_random_list


def build_random_list(vals, random_idx):
    nodes = [Node(v) for v in vals]
    for i in range(len(nodes)):
        if i + 1 < len(nodes):
            nodes[i].next = nodes[i + 1]
        if random_idx[i] >= 0:
            nodes[i].random = nodes[random_idx[i]]
    return nodes


def check(vals, random_idx):
    orig = build_random_list(vals, random_idx)
    orig_set = set(orig)

    head = orig[0] if orig else None
    copy_head = copy_random_list(head)

    copies = []
    cur = copy_head
    while cur is not None:
        assert cur not in orig_set, "copied node shares identity with an original node"
        copies.append(cur)
        cur = cur.next
    assert len(copies) == len(vals)

    copy_index = {n: i for i, n in enumerate(copies)}

    for i, n in enumerate(copies):
        assert n.val == vals[i]
        if random_idx[i] == -1:
            assert n.random is None
            continue
        assert n.random is not None
        assert copy_index[n.random] == random_idx[i]


def test_copy_random_list():
    check([7, 13, 11, 10, 1], [-1, 0, 4, 2, 0])
    check([1, 2], [1, 1])
    check([3, 3, 3], [-1, -1, -1])
    check([], [])
