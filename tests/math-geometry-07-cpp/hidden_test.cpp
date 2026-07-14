#include <cassert>
#include <cstdio>
#include <string>

std::string MultiplyStrings(std::string num1, std::string num2);

int main() {
    assert(MultiplyStrings("2", "3") == "6");
    assert(MultiplyStrings("123", "456") == "56088");
    assert(MultiplyStrings("0", "12345") == "0");
    assert(MultiplyStrings("999", "999") == "998001");
    assert(MultiplyStrings("1", "1") == "1");
    assert(MultiplyStrings("0", "0") == "0");
    assert(MultiplyStrings("1", "999") == "999");
    assert(MultiplyStrings("12345", "67890") == "838102050");
    assert(MultiplyStrings("9", "123") == "1107");
    std::printf("all assertions passed\n");
    return 0;
}
