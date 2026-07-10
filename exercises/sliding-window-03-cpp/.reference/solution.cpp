#include <algorithm>
#include <array>
#include <string>

// CharacterReplacement returns the length of the longest substring of
// s that can be made to contain only one repeating letter after at
// most k character replacements.
int CharacterReplacement(const std::string& s, int k) {
    std::array<int, 26> count{};
    int left = 0, maxFreq = 0, best = 0;
    for (int right = 0; right < static_cast<int>(s.size()); right++) {
        count[s[right] - 'A']++;
        maxFreq = std::max(maxFreq, count[s[right] - 'A']);
        while (right - left + 1 - maxFreq > k) {
            count[s[left] - 'A']--;
            left++;
        }
        best = std::max(best, right - left + 1);
    }
    return best;
}
