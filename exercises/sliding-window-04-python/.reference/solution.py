from collections import Counter


def check_inclusion(s1: str, s2: str) -> bool:
    """Return whether s2 contains a permutation of s1 as a contiguous
    substring."""
    if len(s1) > len(s2):
        return False
    need = Counter(s1)
    window = Counter(s2[: len(s1)])
    if need == window:
        return True
    for i in range(len(s1), len(s2)):
        window[s2[i]] += 1
        left_char = s2[i - len(s1)]
        window[left_char] -= 1
        if window[left_char] == 0:
            del window[left_char]
        if need == window:
            return True
    return False
