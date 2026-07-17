from solution import BloomFilter


def test_zero_false_negatives():
    b = BloomFilter(16384, 4)
    for i in range(500):
        b.add(f"present-{i}")
    for i in range(500):
        assert b.might_contain(f"present-{i}"), \
            f"added key present-{i} reported absent -- bloom filters must never false-negative"


def test_false_positive_rate_bounded():
    b = BloomFilter(16384, 4)
    for i in range(500):
        b.add(f"present-{i}")
    probes = 10000
    fps = sum(1 for i in range(probes) if b.might_contain(f"absent-{i}"))
    assert fps < probes * 2 // 100, \
        f"{fps}/{probes} absent keys reported present -- false-positive rate must stay under 2%"


def test_empty_filter_contains_nothing():
    b = BloomFilter(1024, 3)
    assert not any(b.might_contain(f"anything-{i}") for i in range(100))
