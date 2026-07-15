#pragma once

#include <string>

// Parking lot with motorcycle, compact, and large spots. Vehicles park
// in the smallest spot type that fits and has space.
class ParkingLot {
public:
    ParkingLot(int motorcycle_spots, int compact_spots, int large_spots) {}

    // Park a "motorcycle", "car", or "bus". Return a unique positive
    // ticket number, or -1 if no fitting spot is free.
    int park(const std::string& vehicle) {
        // TODO: implement
        return -1;
    }

    // Free the spot held by ticket. Return false for unknown or
    // already-freed tickets.
    bool leave(int ticket) {
        // TODO: implement
        return false;
    }

    // Number of free spots of the given type.
    int available(const std::string& spot_type) {
        // TODO: implement
        return 0;
    }
};
