#pragma once

#include <string>
#include <vector>

// One playing card: a suit ("clubs", "diamonds", "hearts", "spades")
// and a rank 1 (ace) through 13 (king).
struct Card {
    std::string suit;
    int rank;

    bool operator==(const Card& other) const {
        return suit == other.suit && rank == other.rank;
    }
};

// Standard 52-card deck, dealt from the front of the canonical
// clubs -> diamonds -> hearts -> spades order.
class Deck {
public:
    Deck() {}

    // Number of cards remaining.
    int size() const {
        // TODO: implement
        return 0;
    }

    // Remove and return the next n cards (all remaining if fewer).
    std::vector<Card> deal(int n) {
        // TODO: implement
        return {};
    }

    // Randomly reorder the remaining cards.
    void shuffle() {
        // TODO: implement
    }

    // Restore the full 52-card deck in canonical order.
    void reset() {
        // TODO: implement
    }
};
