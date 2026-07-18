#include <cassert>
#include <cstdio>
#include <string>
#include <vector>

std::vector<std::string> generate_report(const std::vector<std::string>& items);

int main() {
    // One sequence of three calls exercises all three traps: a lone
    // call passes even with the bug (nothing to leak from yet), a
    // second call with different input is polluted by the first, and
    // a third call repeating the first call's input shows that
    // pollution as an outright duplicate.
    std::vector<std::string> first = generate_report({"apples"});
    assert(first == std::vector<std::string>({"LOW STOCK: apples"}));

    std::vector<std::string> second = generate_report({"bananas"});
    assert(second == std::vector<std::string>({"LOW STOCK: bananas"}));

    std::vector<std::string> third = generate_report({"apples"});
    assert(third == std::vector<std::string>({"LOW STOCK: apples"}));

    printf("all assertions passed\n");
    return 0;
}
