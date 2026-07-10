#include <algorithm>
#include <array>
#include <vector>

// LeastInterval returns the minimum number of CPU intervals needed
// to run every task, with identical tasks separated by at least n
// intervals.
int LeastInterval(std::vector<char>& tasks, int n) {
    std::array<int, 26> freq{};
    for (char t : tasks) freq[t - 'A']++;
    int maxFreq = *std::max_element(freq.begin(), freq.end());
    int maxCount = 0;
    for (int f : freq) {
        if (f == maxFreq) maxCount++;
    }

    int frameSize = (maxFreq - 1) * (n + 1) + maxCount;
    return std::max(static_cast<int>(tasks.size()), frameSize);
}
