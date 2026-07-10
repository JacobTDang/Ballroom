from solution import TreeNode, is_subtree


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


def test_is_subtree():
    assert is_subtree(build_tree([3, 4, 5, 1, 2]), build_tree([4, 1, 2])) is True
    assert (
        is_subtree(
            build_tree([3, 4, 5, 1, 2, None, None, None, None, 0]),
            build_tree([4, 1, 2]),
        )
        is False
    )
    assert is_subtree(build_tree([1, 1]), build_tree([1])) is True
    assert is_subtree(build_tree([1]), build_tree([1])) is True
