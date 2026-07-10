from solution import TreeNode, max_path_sum


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


def test_max_path_sum():
    assert max_path_sum(build_tree([1, 2, 3])) == 6
    assert max_path_sum(build_tree([-10, 9, 20, None, None, 15, 7])) == 42
    assert max_path_sum(build_tree([-3])) == -3
    assert max_path_sum(build_tree([2, -1])) == 2
