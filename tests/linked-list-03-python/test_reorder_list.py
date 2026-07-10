from solution import ListNode, reorder_list


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


def check(in_vals, want):
    head = build_list(in_vals)
    reorder_list(head)
    assert to_list(head) == want


def test_reorder_list():
    check([1, 2, 3, 4], [1, 4, 2, 3])
    check([1, 2, 3, 4, 5], [1, 5, 2, 4, 3])
    check([1], [1])
    check([1, 2], [1, 2])
