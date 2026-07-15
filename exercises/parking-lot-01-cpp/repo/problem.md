# Design a Parking Lot

Design a parking lot with three spot sizes — `"motorcycle"`,
`"compact"`, and `"large"` — created with a fixed count of each.

Vehicles: a `"motorcycle"` fits in any spot; a `"car"` fits in compact
or large; a `"bus"` needs a large spot. A vehicle always parks in the
**smallest** spot type that fits and has space.

- `park(vehicle)` — assign a spot; return a unique positive ticket
  number, or `-1` if nothing fits
- `leave(ticket)` — free that vehicle's spot; return whether the
  ticket was valid and active
- `available(spotType)` — free spots of that type

## Examples

```
lot = new ParkingLot(1 motorcycle, 1 compact, 1 large)
park("car")        -> ticket 1   (takes the compact spot)
available("compact") -> 0
park("car")        -> ticket 2   (compact full -> takes large)
park("bus")        -> -1         (no large spot left)
leave(2)           -> true
park("bus")        -> ticket 3
```

## Constraints

- Ticket numbers increase from 1 and are never reused.
- `leave` on an unknown or already-freed ticket returns false.
