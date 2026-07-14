from solution import TreeNode, lowest_common_ancestor


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


def find_node(root, val):
    while root is not None:
        if val == root.val:
            return root
        root = root.left if val < root.val else root.right
    return None


TREE = [6, 2, 8, 0, 4, 7, 9, None, None, 3, 5]


def test_lowest_common_ancestor():
    root = build_tree(TREE)
    assert lowest_common_ancestor(root, find_node(root, 2), find_node(root, 8)).val == 6

    root = build_tree(TREE)
    assert lowest_common_ancestor(root, find_node(root, 2), find_node(root, 4)).val == 2

    root = build_tree(TREE)
    assert lowest_common_ancestor(root, find_node(root, 0), find_node(root, 5)).val == 2

    root = build_tree(TREE)
    assert lowest_common_ancestor(root, find_node(root, 7), find_node(root, 9)).val == 8

    root = build_tree(TREE)
    assert lowest_common_ancestor(root, find_node(root, 6), find_node(root, 6)).val == 6

    root = build_tree(TREE)
    assert lowest_common_ancestor(root, find_node(root, 0), find_node(root, 3)).val == 2

    root = build_tree([2, 1])
    assert lowest_common_ancestor(root, find_node(root, 2), find_node(root, 1)).val == 2
