from solution import group_anagrams


def normalize(groups):
    return sorted(sorted(g) for g in groups)


def test_group_anagrams_case_1():
    got = normalize(group_anagrams(["eat", "tea", "tan", "ate", "nat", "bat"]))
    want = normalize([["bat"], ["nat", "tan"], ["ate", "eat", "tea"]])
    assert got == want


def test_group_anagrams_case_2():
    assert normalize(group_anagrams([""])) == normalize([[""]])
    assert normalize(group_anagrams(["a"])) == normalize([["a"]])
    assert normalize(group_anagrams(["abc", "bca", "cab", "xyz"])) == normalize(
        [["abc", "bca", "cab"], ["xyz"]]
    )
    assert normalize(group_anagrams(["cat", "dog", "bird"])) == normalize(
        [["cat"], ["dog"], ["bird"]]
    )
    assert normalize(group_anagrams(["abc", "bca", "cab", "acb"])) == normalize(
        [["abc", "bca", "cab", "acb"]]
    )
    assert normalize(group_anagrams(["", "", ""])) == normalize([["", "", ""]])
    assert normalize(
        group_anagrams(["bat", "tab", "cat", "act", "dog", "god", "xyz"])
    ) == normalize([["bat", "tab"], ["cat", "act"], ["dog", "god"], ["xyz"]])
    assert normalize(group_anagrams(["a", "b", "a", "c", "b"])) == normalize(
        [["a", "a"], ["b", "b"], ["c"]]
    )
