#include <vector>

// DailyTemperatures returns, for each day, how many days until a
// warmer temperature, or 0 if there isn't one.
std::vector<int> DailyTemperatures(const std::vector<int>& temperatures) {
    std::vector<int> res(temperatures.size(), 0);
    std::vector<int> stack;  // indices, decreasing temperature
    for (int i = 0; i < static_cast<int>(temperatures.size()); i++) {
        while (!stack.empty() && temperatures[stack.back()] < temperatures[i]) {
            int top = stack.back();
            stack.pop_back();
            res[top] = i - top;
        }
        stack.push_back(i);
    }
    return res;
}
