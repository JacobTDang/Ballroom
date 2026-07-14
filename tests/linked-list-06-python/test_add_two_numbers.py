from solution import ListNode, add_two_numbers


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


def test_add_two_numbers():
    assert to_list(add_two_numbers(build_list([2, 4, 3]), build_list([5, 6, 4]))) == [7, 0, 8]
    assert to_list(add_two_numbers(build_list([0]), build_list([0]))) == [0]
    assert to_list(
        add_two_numbers(build_list([9, 9, 9, 9, 9, 9, 9]), build_list([9, 9, 9, 9]))
    ) == [8, 9, 9, 9, 0, 0, 0, 1]
    assert to_list(add_two_numbers(build_list([5]), build_list([5]))) == [0, 1]
    assert to_list(add_two_numbers(build_list([1, 8]), build_list([0]))) == [1, 8]
    assert to_list(add_two_numbers(build_list([2, 4, 3]), build_list([5, 6, 4, 9]))) == [
        7,
        0,
        8,
        9,
    ]
