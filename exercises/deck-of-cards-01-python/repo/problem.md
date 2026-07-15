# Design a Deck of Cards

Design a standard 52-card deck.

A card has a suit (`"clubs"`, `"diamonds"`, `"hearts"`, `"spades"`)
and a rank (`1` = ace through `13` = king). A new deck holds all 52
cards in canonical order: all clubs ace→king, then diamonds, then
hearts, then spades — and cards are dealt from the front of that
order.

- `size()` — cards remaining
- `deal(n)` — remove and return the next `n` cards, in order; if fewer
  than `n` remain, return all that are left
- `shuffle()` — randomly reorder the remaining cards
- `reset()` — restore the full 52-card deck in canonical order

## Examples

```
d = new Deck
d.size()    -> 52
d.deal(1)   -> [(clubs, 1)]
d.deal(2)   -> [(clubs, 2), (clubs, 3)]
d.size()    -> 49
d.reset()
d.size()    -> 52
```

## Constraints

- Dealing past an empty deck returns an empty result, not an error.
- `shuffle()` must keep exactly the remaining cards — same cards, new
  order.
