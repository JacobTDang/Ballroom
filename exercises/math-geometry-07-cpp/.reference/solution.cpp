#include <string>
#include <vector>

// MultiplyStrings returns the product of num1 and num2, both
// non-negative integers given as decimal strings, as a decimal string.
std::string MultiplyStrings(std::string num1, std::string num2) {
    if (num1 == "0" || num2 == "0") return "0";

    int m = static_cast<int>(num1.size());
    int n = static_cast<int>(num2.size());
    std::vector<int> digits(m + n, 0);

    for (int i = m - 1; i >= 0; i--) {
        int d1 = num1[i] - '0';
        for (int j = n - 1; j >= 0; j--) {
            int d2 = num2[j] - '0';
            int sum = d1 * d2 + digits[i + j + 1];
            digits[i + j + 1] = sum % 10;
            digits[i + j] += sum / 10;
        }
    }

    int start = 0;
    while (start < static_cast<int>(digits.size()) - 1 && digits[start] == 0) {
        start++;
    }

    std::string result;
    for (int i = start; i < static_cast<int>(digits.size()); i++) {
        result += static_cast<char>('0' + digits[i]);
    }
    return result;
}
