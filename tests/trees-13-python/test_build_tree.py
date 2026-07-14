from solution import build_tree


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


def test_build_tree():
    assert to_level_order(build_tree([3, 9, 20, 15, 7], [9, 3, 15, 20, 7])) == [
        3,
        9,
        20,
        None,
        None,
        15,
        7,
    ]
    assert to_level_order(build_tree([-1], [-1])) == [-1]
    assert to_level_order(build_tree([1, 2, 3], [3, 2, 1])) == [1, 2, None, 3]
    assert to_level_order(build_tree([1, 2, 4, 5, 3, 6, 7], [4, 2, 5, 1, 6, 3, 7])) == [
        1,
        2,
        3,
        4,
        5,
        6,
        7,
    ]
    assert to_level_order(build_tree([1, 2], [1, 2])) == [1, None, 2]
