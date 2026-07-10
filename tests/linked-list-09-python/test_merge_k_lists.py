from solution import ListNode, merge_k_lists


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


def check(lists_vals, want):
    lists = [build_list(v) for v in lists_vals]
    assert to_list(merge_k_lists(lists)) == want


def test_merge_k_lists():
    check([[1, 4, 5], [1, 3, 4], [2, 6]], [1, 1, 2, 3, 4, 4, 5, 6])
    check([], [])
    check([[]], [])
    check([[1], [], [2]], [1, 2])
