from solution import Entry, sort_leaderboard


def _names(entries):
    return [e.name for e in entries]


def test_no_ties_sorts_by_score_only():
    entries = [Entry("bob", 60), Entry("dan", 100), Entry("cara", 75), Entry("amy", 90)]
    assert _names(sort_leaderboard(entries)) == ["dan", "amy", "cara", "bob"]


def test_tie_group_breaks_ascending_by_name():
    entries = [Entry("erin", 90), Entry("cara", 75), Entry("amy", 90), Entry("bob", 90)]
    assert _names(sort_leaderboard(entries)) == ["amy", "bob", "erin", "cara"]


def test_multiple_tie_groups():
    entries = [
        Entry("zoe", 50), Entry("amy", 80), Entry("erin", 65),
        Entry("dan", 50), Entry("cara", 80), Entry("bob", 80),
    ]
    assert _names(sort_leaderboard(entries)) == ["amy", "bob", "cara", "erin", "dan", "zoe"]


def test_fully_tied_list_sorts_by_name():
    entries = [Entry("zed", 10), Entry("amy", 10), Entry("mno", 10)]
    assert _names(sort_leaderboard(entries)) == ["amy", "mno", "zed"]


def test_negative_scores_with_tie():
    entries = [Entry("cara", 10), Entry("bob", -5), Entry("amy", -5)]
    assert _names(sort_leaderboard(entries)) == ["cara", "amy", "bob"]
