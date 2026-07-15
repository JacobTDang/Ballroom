from solution import TreeNode, deserialize, serialize


def build_tree(vals):
    if not vals or vals[0] is None:
        return None
    root = TreeNode(vals[0])
    queue = [root]
    i = 1
    while queue and i < len(vals):
        node = queue.pop(0)
        if i < len(vals):
            if vals[i] is not None:
                node.left = TreeNode(vals[i])
                queue.append(node.left)
            i += 1
        if i < len(vals):
            if vals[i] is not None:
                node.right = TreeNode(vals[i])
                queue.append(node.right)
            i += 1
    return root


def to_level_order(root):
    if root is None:
        return []
    out = []
    queue = [root]
    while queue:
        node = queue.pop(0)
        if node is None:
            out.append(None)
            continue
        out.append(node.val)
        queue.append(node.left)
        queue.append(node.right)
    while out and out[-1] is None:
        out.pop()
    return out


def check(vals):
    original = build_tree(vals)
    round_tripped = deserialize(serialize(original))
    assert to_level_order(round_tripped) == to_level_order(original)


def test_serialize_deserialize_round_trip_case_1():
    check([1, 2, 3, None, None, 4, 5])


def test_serialize_deserialize_round_trip_case_2():
    check([])


def test_serialize_deserialize_round_trip_case_3():
    check([1])


def test_serialize_deserialize_round_trip_case_4():
    check([-1, -2, -3])


def test_serialize_deserialize_round_trip_case_5():
    check([5, 4, 7, 3, None, 2, None, -1, None, 9])


def test_serialize_deserialize_round_trip_case_6():
    check([0])


def test_serialize_deserialize_round_trip_case_7():
    check([100, None, 200, None, None, None, 300])
