from solution import ListNode, has_cycle


def build_cycle_list(vals, pos):
    if not vals:
        return None
    nodes = [ListNode(v) for v in vals]
    for i in range(len(nodes) - 1):
        nodes[i].next = nodes[i + 1]
    if pos >= 0:
        nodes[-1].next = nodes[pos]
    return nodes[0]


def test_has_cycle():
    assert has_cycle(build_cycle_list([3, 2, 0, -4], 1)) is True
    assert has_cycle(build_cycle_list([1, 2], 0)) is True
    assert has_cycle(build_cycle_list([1], -1)) is False
    assert has_cycle(build_cycle_list([], -1)) is False
    assert has_cycle(build_cycle_list([1, 2, 3], -1)) is False
