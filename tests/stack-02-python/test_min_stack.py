from solution import MinStack


def test_min_stack():
    s = MinStack()
    s.push(-2)
    s.push(0)
    s.push(-3)
    assert s.get_min() == -3
    s.pop()
    assert s.top() == 0
    assert s.get_min() == -2


def test_min_updates_as_equal_values_are_pushed_and_popped():
    s = MinStack()
    s.push(1)
    s.push(1)
    s.push(1)
    assert s.get_min() == 1
    s.pop()
    assert s.get_min() == 1


def test_min_reverts_after_popping_new_min():
    s = MinStack()
    s.push(5)
    s.push(3)
    s.push(7)
    s.push(1)
    assert s.get_min() == 1
    s.pop()
    assert s.get_min() == 3
