from solution import TreeNode, right_side_view


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


def test_right_side_view():
    assert right_side_view(build_tree([1, 2, 3, None, 5, None, 4])) == [1, 3, 4]
    assert right_side_view(build_tree([1, None, 3])) == [1, 3]
    assert right_side_view(build_tree([])) == []
    assert right_side_view(build_tree([1, 2, 3, 4])) == [1, 3, 4]
    assert right_side_view(build_tree([1, 2, 3, 4, 5, 6, 7])) == [1, 3, 7]
    assert right_side_view(build_tree([1, 2, None, 3])) == [1, 2, 3]
