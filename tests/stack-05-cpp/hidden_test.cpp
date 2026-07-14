#include <cassert>
#include <cstdio>
#include <vector>

std::vector<int> DailyTemperatures(const std::vector<int>& temperatures);

int main() {
    assert((DailyTemperatures({73, 74, 75, 71, 69, 72, 76, 73}) ==
            std::vector<int>{1, 1, 4, 2, 1, 1, 0, 0}));
    assert((DailyTemperatures({30, 40, 50, 60}) == std::vector<int>{1, 1, 1, 0}));
    assert((DailyTemperatures({30, 60, 90}) == std::vector<int>{1, 1, 0}));
    assert((DailyTemperatures({80, 79, 78}) == std::vector<int>{0, 0, 0}));
    assert((DailyTemperatures({75}) == std::vector<int>{0}));
    assert((DailyTemperatures({55, 55, 55, 60}) == std::vector<int>{3, 2, 1, 0}));
    printf("all assertions passed\n");
    return 0;
}
