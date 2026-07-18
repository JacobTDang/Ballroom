#include <cmath>
#include <vector>

// Returns whether amounts sums to bill, to the nearest cent.
bool settles_bill(const std::vector<double>& amounts, double bill) {
    double total = 0.0;
    for (double amount : amounts) {
        total += amount;
    }
    return std::round(total * 100) == std::round(bill * 100);
}
