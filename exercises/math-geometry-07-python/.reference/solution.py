def multiply_strings(num1: str, num2: str) -> str:
    if num1 == "0" or num2 == "0":
        return "0"

    m, n = len(num1), len(num2)
    digits = [0] * (m + n)

    for i in range(m - 1, -1, -1):
        d1 = ord(num1[i]) - ord("0")
        for j in range(n - 1, -1, -1):
            d2 = ord(num2[j]) - ord("0")
            total = d1 * d2 + digits[i + j + 1]
            digits[i + j + 1] = total % 10
            digits[i + j] += total // 10

    start = 0
    while start < len(digits) - 1 and digits[start] == 0:
        start += 1

    return "".join(str(d) for d in digits[start:])
