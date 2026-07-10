from solution import TreeNode, good_nodes


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


def test_good_nodes():
    assert good_nodes(build_tree([3, 1, 4, 3, None, 1, 5])) == 4
    assert good_nodes(build_tree([3, 3, None, 4, 2])) == 3
    assert good_nodes(build_tree([1])) == 1
