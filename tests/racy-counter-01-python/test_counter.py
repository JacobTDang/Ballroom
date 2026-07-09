from solution import run


def test_counter_is_exact():
    assert run(50) == 50


def test_counter_single_thread():
    assert run(1) == 1


def test_counter_is_exact_larger_n():
    assert run(200) == 200
