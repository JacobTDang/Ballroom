#include <algorithm>
#include <vector>

// FindMedianSortedArrays returns the median of the two sorted arrays
// nums1 and nums2 combined.
double FindMedianSortedArrays(std::vector<int> nums1, std::vector<int> nums2) {
    if (nums1.size() > nums2.size()) std::swap(nums1, nums2);
    int m = static_cast<int>(nums1.size());
    int n = static_cast<int>(nums2.size());
    int lo = 0, hi = m;
    int half = (m + n + 1) / 2;
    const long long inf = 1LL << 40;

    while (lo <= hi) {
        int i = lo + (hi - lo) / 2;
        int j = half - i;

        long long maxLeftA = (i > 0) ? nums1[i - 1] : -inf;
        long long minRightA = (i < m) ? nums1[i] : inf;
        long long maxLeftB = (j > 0) ? nums2[j - 1] : -inf;
        long long minRightB = (j < n) ? nums2[j] : inf;

        if (maxLeftA <= minRightB && maxLeftB <= minRightA) {
            if ((m + n) % 2 == 1) {
                return static_cast<double>(std::max(maxLeftA, maxLeftB));
            }
            return static_cast<double>(std::max(maxLeftA, maxLeftB) +
                                        std::min(minRightA, minRightB)) /
                   2.0;
        } else if (maxLeftA > minRightB) {
            hi = i - 1;
        } else {
            lo = i + 1;
        }
    }
    return 0.0;  // unreachable for valid, well-formed input
}
