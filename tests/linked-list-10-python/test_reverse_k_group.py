from solution import ListNode, reverse_k_group


def build_list(vals):
    dummy = ListNode()
    cur = dummy
    for v in vals:
        cur.next = ListNode(v)
        cur = cur.next
    return dummy.next


def to_list(head):
    out = []
    while head is not None:
        out.append(head.val)
        head = head.next
    return out


def test_reverse_k_group():
    assert to_list(reverse_k_group(build_list([1, 2, 3, 4, 5]), 2)) == [2, 1, 4, 3, 5]
    assert to_list(reverse_k_group(build_list([1, 2, 3, 4, 5]), 3)) == [3, 2, 1, 4, 5]
    assert to_list(reverse_k_group(build_list([1, 2, 3, 4, 5]), 1)) == [1, 2, 3, 4, 5]
    assert to_list(reverse_k_group(build_list([1, 2, 3, 4, 5, 6]), 6)) == [6, 5, 4, 3, 2, 1]
    assert to_list(reverse_k_group(build_list([1]), 1)) == [1]
