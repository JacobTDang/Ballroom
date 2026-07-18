#include <cassert>
#include <cstdio>
#include <vector>

bool settles_bill(const std::vector<double>& amounts, double bill);

int main() {
    assert(settles_bill({0.1, 0.1, 0.1}, 0.3) == true);
    assert(settles_bill({10.0, 20.0, 70.0}, 100.0) == true);
    assert(settles_bill({0.7, 0.1}, 0.8) == true);
    assert(settles_bill({1.1, 2.2}, 3.3) == true);
    assert(settles_bill({10.00, 10.00}, 20.01) == false);
    assert(settles_bill({5.00, 5.00}, 11.00) == false);
    assert(settles_bill({}, 0.0) == true);
    printf("all assertions passed\n");
    return 0;
}
