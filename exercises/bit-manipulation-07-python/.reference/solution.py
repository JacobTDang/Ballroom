INT_MIN = -(2**31)
INT_MAX = 2**31 - 1


def reverse(x: int) -> int:
    sign = -1 if x < 0 else 1
    x = abs(x)

    result = 0
    while x != 0:
        digit = x % 10
        x //= 10
        result = result * 10 + digit
        if result > INT_MAX:
            return 0

    result *= sign
    if result < INT_MIN or result > INT_MAX:
        return 0
    return result
