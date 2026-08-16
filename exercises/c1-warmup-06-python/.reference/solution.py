def merge_counts(a: dict[str, int], b: dict[str, int]) -> dict[str, int]:
    """Return the two merchant count maps merged, summed, keys sorted."""
    merged = dict(a)
    for merchant, count in b.items():
        merged[merchant] = merged.get(merchant, 0) + count
    return {merchant: merged[merchant] for merchant in sorted(merged)}
