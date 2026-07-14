def my_pow(x: float, n: int) -> float:
    exp = n
    if exp < 0:
        x = 1 / x
        exp = -exp

    result = 1.0
    base = x
    while exp > 0:
        if exp % 2 == 1:
            result *= base
        base *= base
        exp //= 2
    return result
