def _match_class(pattern, p, c):
    """Match c against the class starting at pattern[p] == '['.
    Returns (matched, index_past_bracket, valid)."""
    q = p + 1
    matched = False
    n = len(pattern)
    while q < n and pattern[q] != "]":
        if q + 2 < n and pattern[q + 1] == "-" and pattern[q + 2] != "]":
            if pattern[q] <= c <= pattern[q + 2]:
                matched = True
            q += 3
        else:
            if pattern[q] == c:
                matched = True
            q += 1
    if q >= n:
        return False, 0, False  # unclosed class
    return matched, q + 1, True


def match(pattern, s):
    """The classic two-pointer loop: on '*' remember both positions;
    on a later mismatch, back up to just after the star and let it
    swallow one more character. That remembered pair IS the
    backtracking state -- no recursion needed."""
    p = i = 0
    star_p = -1
    star_i = 0

    while i < len(s):
        if p < len(pattern):
            ch = pattern[p]
            if ch == "*":
                star_p, star_i = p, i
                p += 1
                continue
            # The class branch must run before the literal branch: a
            # pattern "[" against the string "[" would otherwise match
            # itself literally instead of being an invalid class.
            if ch == "[":
                ok, nxt, valid = _match_class(pattern, p, s[i])
                if not valid:
                    return False
                if ok:
                    p = nxt
                    i += 1
                    continue
            elif ch == "?" or ch == s[i]:
                p += 1
                i += 1
                continue
        if star_p == -1:
            return False
        star_i += 1
        p = star_p + 1
        i = star_i

    while p < len(pattern) and pattern[p] == "*":
        p += 1
    return p == len(pattern)
