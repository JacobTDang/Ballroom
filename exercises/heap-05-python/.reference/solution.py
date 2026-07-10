from collections import Counter


def least_interval(tasks: list[str], n: int) -> int:
    """Return the minimum number of CPU intervals needed to run
    every task, with identical tasks separated by at least n
    intervals."""
    freq = Counter(tasks)
    max_freq = max(freq.values())
    max_count = sum(1 for f in freq.values() if f == max_freq)

    frame_size = (max_freq - 1) * (n + 1) + max_count
    return max(len(tasks), frame_size)
