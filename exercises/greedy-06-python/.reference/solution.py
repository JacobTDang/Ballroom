def merge_triplets(triplets: list[list[int]], target: list[int]) -> bool:
    matched: set[int] = set()
    for t in triplets:
        if t[0] > target[0] or t[1] > target[1] or t[2] > target[2]:
            continue
        for i in range(3):
            if t[i] == target[i]:
                matched.add(i)
    return len(matched) == 3
