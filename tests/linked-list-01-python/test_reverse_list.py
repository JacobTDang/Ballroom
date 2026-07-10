from solution import ListNode, reverse_list


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


def test_reverse_list():
    assert to_list(reverse_list(build_list([1, 2, 3, 4, 5]))) == [5, 4, 3, 2, 1]
    assert to_list(reverse_list(build_list([1, 2]))) == [2, 1]
    assert to_list(reverse_list(build_list([]))) == []
    assert to_list(reverse_list(build_list([7]))) == [7]
