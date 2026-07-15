from solution import TreeNode, is_balanced


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


def test_is_balanced_case_1():
    assert is_balanced(build_tree([3, 9, 20, None, None, 15, 7])) is True


def test_is_balanced_case_2():
    assert is_balanced(build_tree([1, 2, 2, 3, 3, None, None, 4, 4])) is False


def test_is_balanced_case_3():
    assert is_balanced(build_tree([])) is True


def test_is_balanced_case_4():
    assert is_balanced(build_tree([1])) is True


def test_is_balanced_case_5():
    assert is_balanced(build_tree([1, 2, 3, 4, 5, 6, 7])) is True


def test_is_balanced_case_6():
    assert is_balanced(build_tree([1, 2, 3, 4, None, None, None, 5])) is False
