#include <algorithm>
#include <vector>

// MinEatingSpeed returns the minimum bananas-per-hour eating speed
// that lets Koko finish every pile within h hours.
int MinEatingSpeed(const std::vector<int>& piles, int h) {
    int lo = 1, hi = 0;
    for (int p : piles) hi = std::max(hi, p);
    while (lo < hi) {
        int mid = lo + (hi - lo) / 2;
        long long hours = 0;
        for (int p : piles) {
            hours += (p + mid - 1) / mid;  // ceil(p / mid)
        }
        if (hours <= h) {
            hi = mid;
        } else {
            lo = mid + 1;
        }
    }
    return lo;
}
