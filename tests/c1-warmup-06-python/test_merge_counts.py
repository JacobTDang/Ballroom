from solution import merge_counts


def test_example_overlap():
    assert merge_counts({"cafe": 2, "gas": 1}, {"cafe": 3}) == {"cafe": 5, "gas": 1}


def test_example_one_empty():
    assert merge_counts({}, {"gym": 4}) == {"gym": 4}


def test_example_both_empty():
    assert merge_counts({}, {}) == {}


def test_disjoint_keys():
    assert merge_counts({"a": 1}, {"b": 2}) == {"a": 1, "b": 2}


def test_keys_sorted():
    result = merge_counts({"zoo": 1, "bar": 2}, {"map": 3, "ant": 4})
    assert list(result.keys()) == ["ant", "bar", "map", "zoo"]


def test_all_keys_shared():
    assert merge_counts({"x": 1, "y": 2}, {"x": 10, "y": 20}) == {"x": 11, "y": 22}


def test_inputs_not_mutated():
    a = {"cafe": 2}
    b = {"cafe": 3}
    merge_counts(a, b)
    assert a == {"cafe": 2} and b == {"cafe": 3}


def test_first_map_empty():
    assert merge_counts({"deli": 7}, {}) == {"deli": 7}
