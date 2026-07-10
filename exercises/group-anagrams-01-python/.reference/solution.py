from collections import defaultdict


def group_anagrams(strs: list[str]) -> list[list[str]]:
    """Group the strings in strs into lists of anagrams of each other, in
    any order (both between groups and within a group)."""
    groups: dict[tuple[int, ...], list[str]] = defaultdict(list)
    for s in strs:
        counts = [0] * 26
        for c in s:
            counts[ord(c) - ord("a")] += 1
        groups[tuple(counts)].append(s)
    return list(groups.values())
