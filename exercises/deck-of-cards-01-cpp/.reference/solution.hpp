#pragma once

#include <algorithm>
#include <random>
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
    Deck() { reset(); }

    // Number of cards remaining.
    int size() const { return static_cast<int>(cards_.size()); }

    // Remove and return the next n cards (all remaining if fewer).
    std::vector<Card> deal(int n) {
        if (n > size()) n = size();
        std::vector<Card> dealt(cards_.begin(), cards_.begin() + n);
        cards_.erase(cards_.begin(), cards_.begin() + n);
        return dealt;
    }

    // Randomly reorder the remaining cards.
    void shuffle() {
        static std::mt19937 rng(std::random_device{}());
        std::shuffle(cards_.begin(), cards_.end(), rng);
    }

    // Restore the full 52-card deck in canonical order.
    void reset() {
        cards_.clear();
        for (const char* suit : {"clubs", "diamonds", "hearts", "spades"}) {
            for (int rank = 1; rank <= 13; rank++) {
                cards_.push_back({suit, rank});
            }
        }
    }

private:
    std::vector<Card> cards_;
};
