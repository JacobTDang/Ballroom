def is_n_straight_hand(hand: list[int], group_size: int) -> bool:
    if len(hand) % group_size != 0:
        return False

    count: dict[int, int] = {}
    for c in hand:
        count[c] = count.get(c, 0) + 1

    for k in sorted(count.keys()):
        need = count.get(k, 0)
        if need == 0:
            continue
        for i in range(group_size):
            c = k + i
            if count.get(c, 0) < need:
                return False
            count[c] = count.get(c, 0) - need
    return True
