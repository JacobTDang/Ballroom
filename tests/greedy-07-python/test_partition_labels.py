from solution import partition_labels


def test_classic():
    assert partition_labels("ababcbacadefegdehijhklij") == [9, 7, 8]


def test_all_unique():
    assert partition_labels("abcde") == [1, 1, 1, 1, 1]


def test_all_same():
    assert partition_labels("aaaa") == [4]


def test_single_char():
    assert partition_labels("a") == [1]


def test_multiple_equal_partitions():
    assert partition_labels("aabbcc") == [2, 2, 2]
