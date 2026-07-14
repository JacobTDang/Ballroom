#include <algorithm>
#include <vector>

// MinMeetingRooms returns the minimum number of conference rooms
// required so that all meetings in intervals can happen without any
// two overlapping meetings sharing a room.
int MinMeetingRooms(std::vector<std::vector<int>>& intervals) {
    int n = static_cast<int>(intervals.size());
    if (n == 0) return 0;

    std::vector<int> starts(n), ends(n);
    for (int i = 0; i < n; i++) {
        starts[i] = intervals[i][0];
        ends[i] = intervals[i][1];
    }
    std::sort(starts.begin(), starts.end());
    std::sort(ends.begin(), ends.end());

    int rooms = 0, maxRooms = 0;
    int i = 0, j = 0;
    while (i < n) {
        if (starts[i] < ends[j]) {
            rooms++;
            i++;
            maxRooms = std::max(maxRooms, rooms);
        } else {
            rooms--;
            j++;
        }
    }

    return maxRooms;
}
