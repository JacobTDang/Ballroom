from solution import find_order


def is_valid_order(num_courses, prerequisites, order):
    if len(order) != num_courses or len(set(order)) != num_courses:
        return False
    pos = {course: i for i, course in enumerate(order)}
    for course, pre in prerequisites:
        if pos[pre] >= pos[course]:
            return False
    return True


def test_valid():
    prereqs = [[1, 0], [2, 0], [3, 1], [3, 2]]
    order = find_order(4, prereqs)
    assert is_valid_order(4, prereqs, order)


def test_cycle():
    order = find_order(2, [[1, 0], [0, 1]])
    assert order == []


def test_no_prerequisites():
    order = find_order(3, [])
    assert is_valid_order(3, [], order)
