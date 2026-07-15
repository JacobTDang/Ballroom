_FITS = {
    "motorcycle": ["motorcycle", "compact", "large"],
    "car": ["compact", "large"],
    "bus": ["large"],
}


class ParkingLot:
    """Parking lot with motorcycle, compact, and large spots. Vehicles
    park in the smallest spot type that fits and has space."""

    def __init__(self, motorcycle_spots: int, compact_spots: int, large_spots: int):
        self._free = {
            "motorcycle": motorcycle_spots,
            "compact": compact_spots,
            "large": large_spots,
        }
        self._tickets = {}  # ticket -> spot type occupied
        self._next_ticket = 1

    def park(self, vehicle: str) -> int:
        for spot in _FITS[vehicle]:
            if self._free[spot] > 0:
                self._free[spot] -= 1
                ticket = self._next_ticket
                self._next_ticket += 1
                self._tickets[ticket] = spot
                return ticket
        return -1

    def leave(self, ticket: int) -> bool:
        spot = self._tickets.pop(ticket, None)
        if spot is None:
            return False
        self._free[spot] += 1
        return True

    def available(self, spot_type: str) -> int:
        return self._free[spot_type]
