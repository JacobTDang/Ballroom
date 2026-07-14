from solution import ListNode, remove_nth_from_end


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


def test_remove_nth_from_end():
    assert to_list(remove_nth_from_end(build_list([1, 2, 3, 4, 5]), 2)) == [1, 2, 3, 5]
    assert to_list(remove_nth_from_end(build_list([1]), 1)) == []
    assert to_list(remove_nth_from_end(build_list([1, 2]), 1)) == [1]
    assert to_list(remove_nth_from_end(build_list([1, 2]), 2)) == [2]
    assert to_list(remove_nth_from_end(build_list([1, 2, 3, 4, 5]), 5)) == [2, 3, 4, 5]
    assert to_list(remove_nth_from_end(build_list([1, 2, 3, 4, 5]), 1)) == [1, 2, 3, 4]
