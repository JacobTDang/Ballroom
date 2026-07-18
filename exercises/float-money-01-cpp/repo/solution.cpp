#include <vector>

// Returns whether amounts sums to bill. Currently unreliable -- find
// and fix the bug.
bool settles_bill(const std::vector<double>& amounts, double bill) {
    double total = 0.0;
    for (double amount : amounts) {
        total += amount;
    }
    return total == bill;
}
