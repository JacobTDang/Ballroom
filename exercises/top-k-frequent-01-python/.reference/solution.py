from collections import Counter


def top_k_frequent(nums: list[int], k: int) -> list[int]:
    """Return the k most frequent elements in nums, in any order."""
    counts = Counter(nums)
    buckets: list[list[int]] = [[] for _ in range(len(nums) + 1)]
    for n, c in counts.items():
        buckets[c].append(n)

    result: list[int] = []
    for c in range(len(buckets) - 1, -1, -1):
        for n in buckets[c]:
            result.append(n)
            if len(result) == k:
                return result
    return result
