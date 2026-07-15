from solution import Deck


def test_new_deck_has_52_cards():
    assert Deck().size() == 52


def test_first_card_is_the_ace_of_clubs():
    assert Deck().deal(1) == [("clubs", 1)]


def test_deal_follows_canonical_order():
    d = Deck()
    d.deal(1)
    assert d.deal(3) == [("clubs", 2), ("clubs", 3), ("clubs", 4)]


def test_suits_change_at_the_right_boundary():
    d = Deck()
    d.deal(13)
    assert d.deal(1) == [("diamonds", 1)]


def test_deal_reduces_size():
    d = Deck()
    d.deal(5)
    assert d.size() == 47


def test_all_52_cards_are_unique():
    cards = Deck().deal(52)
    assert len(set(cards)) == 52


def test_deal_past_empty_returns_remainder():
    d = Deck()
    d.deal(50)
    assert len(d.deal(5)) == 2
    assert d.size() == 0
    assert d.deal(1) == []


def test_reset_restores_canonical_deck():
    d = Deck()
    d.deal(30)
    d.reset()
    assert d.size() == 52
    assert d.deal(1) == [("clubs", 1)]


def test_shuffle_keeps_exactly_the_remaining_cards():
    d = Deck()
    dealt = d.deal(10)
    d.shuffle()
    rest = d.deal(52)
    assert len(rest) == 42
    assert len(set(rest)) == 42
    assert set(rest) | set(dealt) == set(Deck().deal(52))


def test_shuffle_does_not_change_size():
    d = Deck()
    d.deal(7)
    d.shuffle()
    assert d.size() == 45
