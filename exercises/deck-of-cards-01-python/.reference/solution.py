import random

SUITS = ["clubs", "diamonds", "hearts", "spades"]


class Deck:
    """Standard 52-card deck. Cards are (suit, rank) tuples with rank
    1 (ace) through 13 (king), dealt from the front of the canonical
    clubs -> diamonds -> hearts -> spades order."""

    def __init__(self):
        self.reset()

    def size(self) -> int:
        return len(self._cards)

    def deal(self, n: int) -> list:
        dealt = self._cards[:n]
        self._cards = self._cards[n:]
        return dealt

    def shuffle(self) -> None:
        random.shuffle(self._cards)

    def reset(self) -> None:
        self._cards = [(suit, rank) for suit in SUITS for rank in range(1, 14)]
