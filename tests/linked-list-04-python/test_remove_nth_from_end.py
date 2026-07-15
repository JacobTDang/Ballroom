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


def test_remove_nth_from_end_case_1():
    assert to_list(remove_nth_from_end(build_list([1, 2, 3, 4, 5]), 2)) == [1, 2, 3, 5]


def test_remove_nth_from_end_case_2():
    assert to_list(remove_nth_from_end(build_list([1]), 1)) == []


def test_remove_nth_from_end_case_3():
    assert to_list(remove_nth_from_end(build_list([1, 2]), 1)) == [1]


def test_remove_nth_from_end_case_4():
    assert to_list(remove_nth_from_end(build_list([1, 2]), 2)) == [2]


def test_remove_nth_from_end_case_5():
    assert to_list(remove_nth_from_end(build_list([1, 2, 3, 4, 5]), 5)) == [2, 3, 4, 5]


def test_remove_nth_from_end_case_6():
    assert to_list(remove_nth_from_end(build_list([1, 2, 3, 4, 5]), 1)) == [1, 2, 3, 4]
