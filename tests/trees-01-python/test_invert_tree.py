from solution import TreeNode, invert_tree


def build_tree(vals):
    """Build a binary tree from vals in LeetCode's level-order array
    format (None entries are missing children)."""
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


def to_level_order(root):
    """Serialize a tree back to the same None-padded level-order
    format build_tree consumes, trimming only the trailing run of
    Nones."""
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


def test_invert_tree():
    assert to_level_order(invert_tree(build_tree([4, 2, 7, 1, 3, 6, 9]))) == [
        4,
        7,
        2,
        9,
        6,
        3,
        1,
    ]
    assert to_level_order(invert_tree(build_tree([2, 1, 3]))) == [2, 3, 1]
    assert to_level_order(invert_tree(build_tree([]))) == []
