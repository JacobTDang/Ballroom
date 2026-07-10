from solution import group_anagrams


def normalize(groups):
    return sorted(sorted(g) for g in groups)


def test_group_anagrams():
    got = normalize(group_anagrams(["eat", "tea", "tan", "ate", "nat", "bat"]))
    want = normalize([["bat"], ["nat", "tan"], ["ate", "eat", "tea"]])
    assert got == want

    assert normalize(group_anagrams([""])) == normalize([[""]])
    assert normalize(group_anagrams(["a"])) == normalize([["a"]])
    assert normalize(group_anagrams(["abc", "bca", "cab", "xyz"])) == normalize(
        [["abc", "bca", "cab"], ["xyz"]]
    )
