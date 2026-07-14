from solution import Trie


def test_trie():
    trie = Trie()
    trie.insert("apple")
    assert trie.search("apple") is True
    assert trie.search("app") is False
    assert trie.starts_with("app") is True
    trie.insert("app")
    assert trie.search("app") is True


def test_starts_with_false_for_unrelated_prefix():
    trie = Trie()
    trie.insert("banana")
    assert trie.starts_with("ban") is True
    assert trie.search("ban") is False


def test_empty_trie_has_no_matches():
    trie = Trie()
    assert trie.search("anything") is False
    assert trie.starts_with("a") is False


def test_multiple_words_share_common_prefix():
    trie = Trie()
    trie.insert("app")
    trie.insert("apple")
    trie.insert("application")

    assert trie.search("app") is True
    assert trie.search("apple") is True
    assert trie.search("appl") is False
    assert trie.starts_with("appl") is True
