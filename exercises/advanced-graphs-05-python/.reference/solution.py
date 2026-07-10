def alien_order(words: list[str]) -> str:
    adj: dict[str, set[str]] = {}
    in_degree: dict[str, int] = {}

    for word in words:
        for c in word:
            if c not in adj:
                adj[c] = set()
                in_degree[c] = 0

    for w1, w2 in zip(words, words[1:]):
        min_len = min(len(w1), len(w2))
        if len(w1) > len(w2) and w1[:min_len] == w2[:min_len]:
            return ""
        for c1, c2 in zip(w1, w2):
            if c1 != c2:
                if c2 not in adj[c1]:
                    adj[c1].add(c2)
                    in_degree[c2] += 1
                break

    queue = sorted(c for c, d in in_degree.items() if d == 0)
    order: list[str] = []
    while queue:
        c = queue.pop(0)
        order.append(c)
        for n in sorted(adj[c]):
            in_degree[n] -= 1
            if in_degree[n] == 0:
                queue.append(n)

    if len(order) != len(in_degree):
        return ""
    return "".join(order)
