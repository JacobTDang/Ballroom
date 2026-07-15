class ParkingLot:
    """Parking lot with motorcycle, compact, and large spots. Vehicles
    park in the smallest spot type that fits and has space."""

    def __init__(self, motorcycle_spots: int, compact_spots: int, large_spots: int):
        pass

    def park(self, vehicle: str) -> int:
        """Park a "motorcycle", "car", or "bus". Return a unique positive
        ticket number, or -1 if no fitting spot is free."""
        raise NotImplementedError

    def leave(self, ticket: int) -> bool:
        """Free the spot held by ticket. Return False for unknown or
        already-freed tickets."""
        raise NotImplementedError

    def available(self, spot_type: str) -> int:
        """Number of free spots of the given type."""
        raise NotImplementedError
