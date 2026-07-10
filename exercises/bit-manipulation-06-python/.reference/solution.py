def get_sum(a: int, b: int) -> int:
    mask = 0xFFFFFFFF
    a &= mask
    b &= mask
    while b != 0:
        carry = ((a & b) << 1) & mask
        a = (a ^ b) & mask
        b = carry
    if a >= 0x80000000:
        return a - 0x100000000
    return a
