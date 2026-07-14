#include <vector>

// PlusOne returns digits+1 as a new digits array.
std::vector<int> PlusOne(std::vector<int>& digits) {
    std::vector<int> result = digits;

    for (int i = static_cast<int>(result.size()) - 1; i >= 0; i--) {
        if (result[i] < 9) {
            result[i]++;
            return result;
        }
        result[i] = 0;
    }

    result.insert(result.begin(), 1);
    return result;
}
