#include <algorithm>
#include <vector>

// CanAttendMeetings returns true if a person could attend every
// meeting in intervals without any two of them overlapping.
bool CanAttendMeetings(std::vector<std::vector<int>>& intervals) {
    std::vector<std::vector<int>> sorted = intervals;
    std::sort(sorted.begin(), sorted.end(), [](const std::vector<int>& a, const std::vector<int>& b) {
        return a[0] < b[0];
    });

    for (size_t i = 1; i < sorted.size(); i++) {
        if (sorted[i][0] < sorted[i - 1][1]) {
            return false;
        }
    }
    return true;
}
