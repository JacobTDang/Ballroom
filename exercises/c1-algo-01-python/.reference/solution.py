from collections import Counter


def top_k(events: list[str], k: int) -> list[str]:
    """Return the k most active account ids, count desc, ties lexicographic."""
    counts = Counter(events)
    ranked = sorted(counts, key=lambda account: (-counts[account], account))
    return ranked[:k]
