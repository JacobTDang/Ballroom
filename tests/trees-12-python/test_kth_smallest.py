from solution import TreeNode, kth_smallest


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


def test_kth_smallest_case_1():
    assert kth_smallest(build_tree([3, 1, 4, None, 2]), 1) == 1


def test_kth_smallest_case_2():
    assert kth_smallest(build_tree([3, 1, 4, None, 2]), 2) == 2


def test_kth_smallest_case_3():
    assert kth_smallest(build_tree([3, 1, 4, None, 2]), 4) == 4


def test_kth_smallest_case_4():
    assert kth_smallest(build_tree([5, 3, 6, 2, 4, None, None, 1]), 3) == 3


def test_kth_smallest_case_5():
    assert kth_smallest(build_tree([5, 3, 6, 2, 4, None, None, 1]), 5) == 5


def test_kth_smallest_case_6():
    assert kth_smallest(build_tree([1]), 1) == 1
