from solution import ListNode, merge_two_lists


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


def test_merge_two_lists_case_1():
    assert to_list(merge_two_lists(build_list([1, 2, 4]), build_list([1, 3, 4]))) == [
        1,
        1,
        2,
        3,
        4,
        4,
    ]


def test_merge_two_lists_case_2():
    assert to_list(merge_two_lists(build_list([]), build_list([]))) == []


def test_merge_two_lists_case_3():
    assert to_list(merge_two_lists(build_list([]), build_list([0]))) == [0]


def test_merge_two_lists_case_4():
    assert to_list(merge_two_lists(build_list([5]), build_list([1, 2, 4]))) == [1, 2, 4, 5]


def test_merge_two_lists_case_5():
    assert to_list(merge_two_lists(build_list([1, 1, 1]), build_list([1, 1, 1]))) == [
        1,
        1,
        1,
        1,
        1,
        1,
    ]


def test_merge_two_lists_case_6():
    assert to_list(merge_two_lists(build_list([-3, -1, 2]), build_list([-2, 0, 5]))) == [
        -3,
        -2,
        -1,
        0,
        2,
        5,
    ]
