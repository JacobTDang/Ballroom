#include <vector>

// MaxProfit returns the maximum profit from buying on one day and
// selling on a later day, or 0 if no profit is possible.
int MaxProfit(const std::vector<int>& prices) {
    if (prices.empty()) return 0;
    int minPrice = prices[0];
    int best = 0;
    for (size_t i = 1; i < prices.size(); i++) {
        if (prices[i] - minPrice > best) best = prices[i] - minPrice;
        if (prices[i] < minPrice) minPrice = prices[i];
    }
    return best;
}
