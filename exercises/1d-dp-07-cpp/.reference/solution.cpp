#include <string>
#include <vector>

// NumDecodings returns the number of ways to decode the digit string s
// using the 'A'-'Z' -> "1"-"26" mapping.
int NumDecodings(std::string s) {
    int n = s.size();
    if (n == 0) return 0;

    std::vector<int> dp(n + 1, 0);
    dp[0] = 1;
    if (s[0] != '0') dp[1] = 1;

    for (int i = 2; i <= n; i++) {
        if (s[i - 1] != '0') dp[i] += dp[i - 1];
        int twoDigit = (s[i - 2] - '0') * 10 + (s[i - 1] - '0');
        if (twoDigit >= 10 && twoDigit <= 26) dp[i] += dp[i - 2];
    }

    return dp[n];
}
