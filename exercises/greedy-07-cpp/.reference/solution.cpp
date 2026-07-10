#include <string>
#include <unordered_map>
#include <vector>

// PartitionLabels returns the sizes of the parts formed by splitting s
// so that each letter appears in at most one part, in order.
std::vector<int> PartitionLabels(std::string s) {
    std::unordered_map<char, int> last;
    for (size_t i = 0; i < s.size(); i++) {
        last[s[i]] = (int)i;
    }

    std::vector<int> result;
    int start = 0, end = 0;
    for (size_t i = 0; i < s.size(); i++) {
        if (last[s[i]] > end) end = last[s[i]];
        if ((int)i == end) {
            result.push_back(end - start + 1);
            start = (int)i + 1;
        }
    }
    return result;
}
