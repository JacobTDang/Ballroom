from solution import TreeNode, is_valid_bst


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


def test_is_valid_bst_case_1():
    assert is_valid_bst(build_tree([2, 1, 3])) is True


def test_is_valid_bst_case_2():
    assert is_valid_bst(build_tree([5, 1, 4, None, None, 3, 6])) is False


def test_is_valid_bst_case_3():
    assert is_valid_bst(build_tree([1])) is True


def test_is_valid_bst_case_4():
    assert is_valid_bst(build_tree([2, 2, 2])) is False


def test_is_valid_bst_case_5():
    assert is_valid_bst(build_tree([10, 5, 15, None, None, 6, 20])) is False


def test_is_valid_bst_case_6():
    assert is_valid_bst(build_tree([4, 2, 6, 1, 3, 5, 7])) is True


def test_is_valid_bst_case_7():
    assert is_valid_bst(build_tree([2147483647])) is True
