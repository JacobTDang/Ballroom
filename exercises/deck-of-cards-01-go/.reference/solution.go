package main

import "math/rand"

// Card is one playing card: a suit ("clubs", "diamonds", "hearts",
// "spades") and a rank 1 (ace) through 13 (king).
type Card struct {
	Suit string
	Rank int
}

var suits = []string{"clubs", "diamonds", "hearts", "spades"}

// Deck is a standard 52-card deck, dealt from the front of the
// canonical clubs -> diamonds -> hearts -> spades order.
type Deck struct {
	cards []Card
}

func NewDeck() *Deck {
	d := &Deck{}
	d.Reset()
	return d
}

// Size returns the number of cards remaining.
func (d *Deck) Size() int {
	return len(d.cards)
}

// Deal removes and returns the next n cards (all remaining if fewer).
func (d *Deck) Deal(n int) []Card {
	if n > len(d.cards) {
		n = len(d.cards)
	}
	dealt := d.cards[:n]
	d.cards = d.cards[n:]
	return dealt
}

// Shuffle randomly reorders the remaining cards.
func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

// Reset restores the full 52-card deck in canonical order.
func (d *Deck) Reset() {
	d.cards = d.cards[:0]
	for _, s := range suits {
		for r := 1; r <= 13; r++ {
			d.cards = append(d.cards, Card{Suit: s, Rank: r})
		}
	}
}
