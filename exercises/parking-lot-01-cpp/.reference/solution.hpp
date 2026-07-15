#pragma once

#include <map>
#include <string>
#include <vector>

// Parking lot with motorcycle, compact, and large spots. Vehicles park
// in the smallest spot type that fits and has space.
class ParkingLot {
public:
    ParkingLot(int motorcycle_spots, int compact_spots, int large_spots) {
        free_["motorcycle"] = motorcycle_spots;
        free_["compact"] = compact_spots;
        free_["large"] = large_spots;
    }

    // Park a "motorcycle", "car", or "bus". Return a unique positive
    // ticket number, or -1 if no fitting spot is free.
    int park(const std::string& vehicle) {
        for (const auto& spot : fits(vehicle)) {
            if (free_[spot] > 0) {
                free_[spot]--;
                int ticket = next_ticket_++;
                tickets_[ticket] = spot;
                return ticket;
            }
        }
        return -1;
    }

    // Free the spot held by ticket. Return false for unknown or
    // already-freed tickets.
    bool leave(int ticket) {
        auto it = tickets_.find(ticket);
        if (it == tickets_.end()) return false;
        free_[it->second]++;
        tickets_.erase(it);
        return true;
    }

    // Number of free spots of the given type.
    int available(const std::string& spot_type) { return free_[spot_type]; }

private:
    static std::vector<std::string> fits(const std::string& vehicle) {
        if (vehicle == "motorcycle") return {"motorcycle", "compact", "large"};
        if (vehicle == "car") return {"compact", "large"};
        return {"large"};  // bus
    }

    std::map<std::string, int> free_;
    std::map<int, std::string> tickets_;
    int next_ticket_ = 1;
};
