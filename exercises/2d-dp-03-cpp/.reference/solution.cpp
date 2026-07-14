#include <algorithm>
#include <vector>

// MaxProfit returns the maximum profit achievable buying and selling
// prices with unlimited transactions, subject to a mandatory one-day
// cooldown after selling before buying again.
int MaxProfit(std::vector<int>& prices) {
    if (prices.empty()) return 0;
    int hold = -prices[0];
    int sold = 0;
    int rest = 0;
    for (size_t i = 1; i < prices.size(); i++) {
        int prevHold = hold, prevSold = sold, prevRest = rest;
        hold = std::max(prevHold, prevRest - prices[i]);
        sold = prevHold + prices[i];
        rest = std::max(prevRest, prevSold);
    }
    return std::max(sold, rest);
}
