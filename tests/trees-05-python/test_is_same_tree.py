from solution import TreeNode, is_same_tree


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


def test_is_same_tree_case_1():
    assert is_same_tree(build_tree([1, 2, 3]), build_tree([1, 2, 3])) is True


def test_is_same_tree_case_2():
    assert is_same_tree(build_tree([1, 2]), build_tree([1, None, 2])) is False


def test_is_same_tree_case_3():
    assert is_same_tree(build_tree([1, 2, 1]), build_tree([1, 1, 2])) is False


def test_is_same_tree_case_4():
    assert is_same_tree(build_tree([]), build_tree([])) is True


def test_is_same_tree_case_5():
    assert is_same_tree(build_tree([1]), build_tree([])) is False


def test_is_same_tree_case_6():
    assert is_same_tree(build_tree([-5]), build_tree([-5])) is True


def test_is_same_tree_case_7():
    assert is_same_tree(build_tree([1]), build_tree([2])) is False
