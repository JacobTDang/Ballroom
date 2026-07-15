package main

// Card is one playing card: a suit ("clubs", "diamonds", "hearts",
// "spades") and a rank 1 (ace) through 13 (king).
type Card struct {
	Suit string
	Rank int
}

// Deck is a standard 52-card deck, dealt from the front of the
// canonical clubs -> diamonds -> hearts -> spades order.
type Deck struct {
}

func NewDeck() *Deck {
	return &Deck{}
}

// Size returns the number of cards remaining.
func (d *Deck) Size() int {
	// TODO: implement
	return 0
}

// Deal removes and returns the next n cards (all remaining if fewer).
func (d *Deck) Deal(n int) []Card {
	// TODO: implement
	return nil
}

// Shuffle randomly reorders the remaining cards.
func (d *Deck) Shuffle() {
	// TODO: implement
}

// Reset restores the full 52-card deck in canonical order.
func (d *Deck) Reset() {
	// TODO: implement
}
