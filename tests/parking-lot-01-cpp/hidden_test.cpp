#include <cassert>
#include <cstdio>

#include "solution.hpp"

int main() {
    {
        ParkingLot lot(2, 3, 4);                   // initial availability
        assert(lot.available("motorcycle") == 2);
        assert(lot.available("compact") == 3);
        assert(lot.available("large") == 4);
    }
    {
        ParkingLot lot(1, 1, 1);                   // smallest fit first
        assert(lot.park("motorcycle") > 0);
        assert(lot.available("motorcycle") == 0);
        assert(lot.available("compact") == 1);
    }
    {
        ParkingLot lot(1, 1, 1);                   // motorcycle overflow
        lot.park("motorcycle");
        lot.park("motorcycle");
        assert(lot.available("compact") == 0);
        lot.park("motorcycle");
        assert(lot.available("large") == 0);
    }
    {
        ParkingLot lot(1, 1, 1);                   // car: compact first
        assert(lot.park("car") > 0);
        assert(lot.available("compact") == 0);
        assert(lot.available("large") == 1);
    }
    {
        ParkingLot lot(5, 0, 0);                   // car never fits motorcycle spot
        assert(lot.park("car") == -1);
        assert(lot.available("motorcycle") == 5);
    }
    {
        ParkingLot a(5, 5, 0);                     // bus needs large
        assert(a.park("bus") == -1);
        ParkingLot b(0, 0, 1);
        assert(b.park("bus") > 0);
    }
    {
        ParkingLot lot(0, 1, 0);                   // full lot
        assert(lot.park("car") > 0);
        assert(lot.park("car") == -1);
    }
    {
        ParkingLot lot(0, 3, 0);                   // tickets 1,2,3
        assert(lot.park("car") == 1);
        assert(lot.park("car") == 2);
        assert(lot.park("car") == 3);
    }
    {
        ParkingLot lot(0, 0, 1);                   // leave frees the spot
        int ticket = lot.park("bus");
        assert(lot.available("large") == 0);
        assert(lot.leave(ticket) == true);
        assert(lot.available("large") == 1);
        assert(lot.park("bus") > 0);
    }
    {
        ParkingLot lot(1, 1, 1);                   // unknown ticket
        assert(lot.leave(99) == false);
    }
    {
        ParkingLot lot(0, 1, 0);                   // double leave
        int ticket = lot.park("car");
        assert(lot.leave(ticket) == true);
        assert(lot.leave(ticket) == false);
        assert(lot.available("compact") == 1);
    }
    {
        ParkingLot lot(0, 1, 0);                   // tickets not reused
        int t1 = lot.park("car");
        lot.leave(t1);
        assert(lot.park("car") != t1);
    }
    printf("all assertions passed\n");
    return 0;
}
