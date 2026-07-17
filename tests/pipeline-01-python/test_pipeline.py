import random
import time

from solution import fan_out_in


def test_nothing_dropped_nothing_duplicated():
    for _ in range(3):
        inputs = [i % 250 for i in range(600)]
        got = fan_out_in(inputs, 8, lambda v: (time.sleep(random.random() * 0.002), v * 3)[1])
        assert len(got) == len(inputs), f"fan-in dropped work: {len(got)} of {len(inputs)}"
        assert sorted(got) == sorted(v * 3 for v in inputs)
