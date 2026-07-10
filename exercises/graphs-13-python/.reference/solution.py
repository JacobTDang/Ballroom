from collections import defaultdict, deque


def ladder_length(begin_word: str, end_word: str, word_list: list[str]) -> int:
    word_set = set(word_list)
    if end_word not in word_set:
        return 0

    patterns = defaultdict(list)

    def add_patterns(word: str) -> None:
        for i in range(len(word)):
            pattern = word[:i] + "*" + word[i + 1 :]
            patterns[pattern].append(word)

    for w in word_set:
        add_patterns(w)
    add_patterns(begin_word)

    visited = {begin_word}
    queue = deque([begin_word])
    steps = 1

    while queue:
        for _ in range(len(queue)):
            word = queue.popleft()
            if word == end_word:
                return steps
            for i in range(len(word)):
                pattern = word[:i] + "*" + word[i + 1 :]
                for neighbor in patterns[pattern]:
                    if neighbor not in visited:
                        visited.add(neighbor)
                        queue.append(neighbor)
        steps += 1
    return 0
