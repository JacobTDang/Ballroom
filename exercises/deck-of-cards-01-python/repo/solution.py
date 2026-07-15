SUITS = ["clubs", "diamonds", "hearts", "spades"]


class Deck:
    """Standard 52-card deck. Cards are (suit, rank) tuples with rank
    1 (ace) through 13 (king), dealt from the front of the canonical
    clubs -> diamonds -> hearts -> spades order."""

    def __init__(self):
        pass

    def size(self) -> int:
        """Number of cards remaining."""
        raise NotImplementedError

    def deal(self, n: int) -> list:
        """Remove and return the next n cards (all remaining if fewer)."""
        raise NotImplementedError

    def shuffle(self) -> None:
        """Randomly reorder the remaining cards."""
        raise NotImplementedError

    def reset(self) -> None:
        """Restore the full 52-card deck in canonical order."""
        raise NotImplementedError
