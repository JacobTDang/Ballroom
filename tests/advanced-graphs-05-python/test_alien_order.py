from solution import alien_order


def is_valid_alien_order(words, order):
    """Checks the ORDERING PROPERTY rather than an exact string, since the
    topological order implied by words is not unique: every distinct
    character appearing in words must appear exactly once in order, and
    every adjacent word pair's first differing character must respect
    that order."""
    if not order:
        return False

    pos = {}
    for i, c in enumerate(order):
        if c in pos:
            return False  # duplicate character in order
        pos[c] = i

    seen = set()
    for word in words:
        seen.update(word)
    if seen != set(pos.keys()):
        return False

    for w1, w2 in zip(words, words[1:]):
        min_len = min(len(w1), len(w2))
        if len(w1) > len(w2) and w1[:min_len] == w2[:min_len]:
            return False
        for c1, c2 in zip(w1, w2):
            if c1 != c2:
                if pos[c1] >= pos[c2]:
                    return False
                break

    return True


def test_valid():
    words = ["wrt", "wrf", "er", "ett", "rftt"]
    order = alien_order(words)
    assert is_valid_alien_order(words, order)


def test_invalid_prefix():
    assert alien_order(["abc", "ab"]) == ""


def test_cycle():
    assert alien_order(["z", "x", "z"]) == ""
