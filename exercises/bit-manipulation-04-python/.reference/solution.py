def reverse_bits(n: int) -> int:
    result = 0
    for _ in range(32):
        result = ((result << 1) | (n & 1)) & 0xFFFFFFFF
        n >>= 1
    return result
