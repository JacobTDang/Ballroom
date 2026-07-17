from solution import match


def test_match_table():
    cases = [
        ("abc", "abc", True),
        ("abc", "abd", False),
        ("a?c", "abc", True),
        ("a?c", "ac", False),
        ("?", "", False),
        ("*", "", True),
        ("*", "anything", True),
        ("a*", "a", True),
        ("a*b*c", "aXXbYYc", True),
        ("a*b*c", "aXXbYY", False),
        ("*.go", "main.go", True),
        ("*.go", "main.gox", False),
        ("a*a", "aa", True),
        ("a*a", "aba", True),
        ("a*a", "ab", False),
        ("**", "x", True),
        ("[a-c]x", "bx", True),
        ("[a-c]x", "dx", False),
        ("[xyz]", "y", True),
        ("[xyz]", "w", False),
        ("file[0-9].txt", "file7.txt", True),
        ("file[0-9].txt", "fileX.txt", False),
        ("[abc", "a", False),
        ("[", "[", False),
    ]
    for pattern, s, want in cases:
        assert match(pattern, s) == want, f"match({pattern!r}, {s!r}) != {want}"
